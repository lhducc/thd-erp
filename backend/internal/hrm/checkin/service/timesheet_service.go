package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"erp/backend/pkg/pointer"
	"erp/backend/pkg/timeonly"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

var vietnamLoc = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Printf("Failed to load Asia/Ho_Chi_Minh timezone: %v, falling back to UTC", err)
		// You could also try loading from IANA database or use a fixed offset
		loc = time.FixedZone("ICT", 7*60*60) // UTC+7
	}
	log.Println("Vietnam location loaded:", loc)
	return loc
}()

type timesheetServiceImp struct {
	timesheetRepo        repo_interface.TimeSheetRepoInterface
	tsListRepo           repo_interface.TimesheetListInterface
	attendanceRecordRepo repo_interface.AttendanceRecordRepository
	empWSRepo            repo_interface.EmployeeWorkShiftRepo
	tsDetailsRepo        repo_interface.TimeSheetDetailRepoInterface
}

func NewTimesheetService(
	timesheetRepo repo_interface.TimeSheetRepoInterface,
	tsListRepo repo_interface.TimesheetListInterface,
	attendanceRecordRepo repo_interface.AttendanceRecordRepository,
	empWSRepo repo_interface.EmployeeWorkShiftRepo,
	tsDetailsRepo repo_interface.TimeSheetDetailRepoInterface,
) service_interface.TimeSheetServiceInterface {
	return &timesheetServiceImp{
		timesheetRepo:        timesheetRepo,
		tsListRepo:           tsListRepo,
		attendanceRecordRepo: attendanceRecordRepo,
		empWSRepo:            empWSRepo,
		tsDetailsRepo:        tsDetailsRepo,
	}
}

func (t *timesheetServiceImp) FindByEmployeeIDAndMonth(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error) {
	return t.timesheetRepo.FindByEmployeeAndMonth(ctx, employeeID, month, year)
}

func (t *timesheetServiceImp) CalculatorTimeSheetList(ctx context.Context, timesheetListID string) error {
	tsl, err := t.tsListRepo.GetByID(ctx, timesheetListID)
	if err != nil {
		return errors.New("Lỗi khi lấy danh sách tính toán")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan model.TimeSheet, len(tsl.Timesheets))
	results := make(chan error, len(tsl.Timesheets))

	numWorkers := 5
	if numWorkers > len(tsl.Timesheets) {
		numWorkers = len(tsl.Timesheets)
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for ts := range jobs {
				if ctx.Err() != nil {
					results <- nil
					return
				}

				err := t.calculateForEmployee(ctx, ts, tsl)
				results <- err
			}
		}()
	}

	for _, ts := range tsl.Timesheets {
		jobs <- ts
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var firstErr error
	for i := 0; i < len(tsl.Timesheets); i++ {
		err := <-results
		if err != nil && firstErr == nil {
			firstErr = err
			cancel()
		}
	}

	return firstErr
}

// concurrency
func (t *timesheetServiceImp) calculateForEmployee(
	ctx context.Context,
	ts model.TimeSheet,
	tsl *model.TimeSheetList,
) error {
	log.Printf("vietnam location: %v", vietnamLoc)
	// Check if context is cancelled before expensive operations
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Fetch existing timesheet details to preserve manual adjustments
	existingDetails, err := t.tsDetailsRepo.GetDetailsByTimeSheetID(ctx, ts.TimeSheetID)
	if err != nil {
		log.Printf("Error fetching existing details for timesheet %d: %v", ts.TimeSheetID, err)
		return err
	}

	// CreateTimeSheets map for existing details for quick lookup
	existingDetailMap := make(map[time.Time]model.TimeSheetDetail)
	for _, detail := range existingDetails {
		// Use normalized date as key for consistent comparison
		dateKey := normalizeToDay(detail.Date, vietnamLoc)
		existingDetailMap[dateKey] = detail
	}

	records, err := t.attendanceRecordRepo.ListHistoryRecordApproveByEmpID(ctx, ts.EmployeeID, tsl.StartDate, tsl.EndDate)
	if err != nil {
		return err
	}

	// Check context after I/O operations
	if ctx.Err() != nil {
		return ctx.Err()
	}

	employeeShifts, err := t.empWSRepo.GetAllByEmployeeID(ts.EmployeeID)
	if err != nil {
		return err
	}

	// Map day-shift
	shiftMap := make(map[time.Time]model.EmployeeWorkshift)
	for _, shift := range employeeShifts {
		// Ensure shift.Date has a location before using it
		if shift.Date.Location() == nil {
			shift.Date = shift.Date.In(vietnamLoc)
		}
		dateKey := normalizeToDay(shift.Date, vietnamLoc)
		shiftMap[dateKey] = shift
	}

	details := make([]model.TimeSheetDetail, 0)
	totalWorkDays := 0.0
	totalLateMinutes := 0
	lateShifts := 0

	// Make sure both dates have proper location
	if tsl.StartDate.Location() == nil {
		tsl.StartDate = tsl.StartDate.In(vietnamLoc)
	}
	if tsl.EndDate.Location() == nil {
		tsl.EndDate = tsl.EndDate.In(vietnamLoc)
	}

	currentDate := tsl.StartDate.In(vietnamLoc)
	endDate := tsl.EndDate.In(vietnamLoc)

	currentDate = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, vietnamLoc)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, vietnamLoc)

	for currentDate.Before(tsl.EndDate) || currentDate.Equal(endDate) {
		// Check context before each iteration
		if ctx.Err() != nil {
			return ctx.Err()
		}

		dateKey := normalizeToDay(currentDate, vietnamLoc)

		// Check if this date has a manual adjustment
		if existingDetail, exists := existingDetailMap[dateKey]; exists && existingDetail.IsManuallyAdjusted {
			// Preserve manually adjusted details
			details = append(details, existingDetail)

			// Add to totals
			totalWorkDays += existingDetail.WorkDays
			if existingDetail.IsLate {
				totalLateMinutes += existingDetail.LateMinutes
				lateShifts++
			}

			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		// Process only if shift exists
		if shift, exists := shiftMap[dateKey]; exists {
			// Start with either existing detail or new one
			detail := existingDetailMap[dateKey]
			detail.TimeSheetID = ts.TimeSheetID
			detail.Date = dateKey
			detail.DayOfWeek = int(dateKey.Weekday())
			detail.IsWorkingDay = true
			detail.IsWeekend = isWeekend(dateKey)
			detail.WorkShiftID = &shift.WorkShiftID
			detail.IsManuallyAdjusted = false // Reset manual adjustment flag

			// Get check-in/out time ranges from shift
			checkinFrom, checkinTo, checkoutFrom, checkoutTo, err := t.getShiftTimeRange(
				shift.WorkShift, dateKey)

			if err != nil {
				log.Printf("Error parsing shift times: %v", err)
				currentDate = currentDate.AddDate(0, 0, 1)
				continue
			}

			checkInRecords, checkOutRecords := t.classifyRecords(
				records, checkinFrom, checkinTo, checkoutFrom, checkoutTo)

			if len(checkInRecords) > 0 {
				detail.CheckInRecordID = &checkInRecords[0].AttendanceRecordID
			}

			if len(checkOutRecords) > 0 {
				detail.CheckOutRecordID = &checkOutRecords[0].AttendanceRecordID
			}

			if len(checkInRecords) > 0 && len(checkOutRecords) > 0 {
				firstCheckIn := checkInRecords[0].Timestamp.In(vietnamLoc)
				lastCheckOut := checkOutRecords[0].Timestamp.In(vietnamLoc)

				// Calculate real working hours
				workHours := lastCheckOut.Sub(firstCheckIn).Hours()
				detail.WorkHours = math.Round(workHours*100) / 100

				// Convert TimeOnly to time.Time on the given date
				shiftStart := time.Date(
					dateKey.Year(),
					dateKey.Month(),
					dateKey.Day(),
					shift.WorkShift.StartTime.Hour(),
					shift.WorkShift.StartTime.Minute(),
					0, 0, dateKey.Location(),
				)

				// Check valid checkin/checkout times
				hasValidCheckIn := checkinFrom != nil && checkinTo != nil &&
					firstCheckIn.After(checkinFrom.Add(-1*time.Minute)) &&
					firstCheckIn.Before(checkinTo.Add(1*time.Minute))

				hasValidCheckOut := checkoutFrom != nil && checkoutTo != nil &&
					lastCheckOut.After(checkoutFrom.Add(-1*time.Minute)) &&
					lastCheckOut.Before(checkoutTo.Add(1*time.Minute))

				if hasValidCheckIn && hasValidCheckOut {
					// Save original value before any potential manual adjustment
					if detail.OriginalWorkDays == nil {
						detail.OriginalWorkDays = &detail.WorkDays
					}

					detail.WorkDays = float64(shift.WorkShift.WorkDay)
					totalWorkDays += float64(shift.WorkShift.WorkDay)

					if firstCheckIn.After(shiftStart) {
						lateMinutes := int(firstCheckIn.Sub(shiftStart).Minutes())
						detail.LateMinutes = lateMinutes
						detail.IsLate = true
						totalLateMinutes += lateMinutes
						lateShifts++
					} else {
						detail.IsLate = false
					}
				} else {
					detail.WorkDays = 0
					detail.IsAbsent = true
					detail.AbsentReason = pointer.String("Không chấm công đủ theo quy định")
				}
			} else {
				detail.WorkDays = 0
				detail.IsAbsent = true
				if len(checkInRecords) == 0 && len(checkOutRecords) == 0 {
					detail.AbsentReason = pointer.String("Không có dữ liệu chấm công")
				} else if len(checkInRecords) == 0 {
					detail.AbsentReason = pointer.String("Không có dữ liệu check-in")
				} else {
					detail.AbsentReason = pointer.String("Không có dữ liệu check-out")
				}
			}

			details = append(details, detail)
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	ts.TotalWorkDays = totalWorkDays
	ts.TotalLateMinutes = totalLateMinutes
	ts.LateShifts = lateShifts
	ts.Details = details

	// Check context before saving
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := t.timesheetRepo.UpdateTimesheetAndCreateDetail(ctx, &ts); err != nil {
		return err
	}

	return nil
}

func isWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func (t *timesheetServiceImp) getShiftTimeRange(workShift model.WorkShifts, date time.Time) (checkinFrom, checkinTo, checkoutFrom, checkoutTo *time.Time, err error) {
	location := vietnamLoc

	parseTime := func(ti *timeonly.TimeOnly) (*time.Time, error) {
		if ti == nil || ti.IsZero() {
			return nil, nil
		}
		result := time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			ti.Hour(),
			ti.Minute(),
			0, 0, location,
		)
		return &result, nil
	}

	checkinFrom, err = parseTime(workShift.CheckinFrom)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	checkinTo, err = parseTime(workShift.CheckinTo)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	checkoutFrom, err = parseTime(workShift.CheckoutFrom)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	checkoutTo, err = parseTime(workShift.CheckoutTo)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return checkinFrom, checkinTo, checkoutFrom, checkoutTo, nil
}

func (t *timesheetServiceImp) classifyRecords(records []model.AttendanceRecord,
	checkinFrom, checkinTo,
	checkoutFrom, checkoutTo *time.Time) (
	checkInRecords []*model.AttendanceRecord,
	checkOutRecords []*model.AttendanceRecord) {

	for i := range records {
		record := &records[i]

		// Ensure record timestamp has a location before calling In()
		if record.Timestamp.Location() == nil {
			record.Timestamp = time.Date(
				record.Timestamp.Year(),
				record.Timestamp.Month(),
				record.Timestamp.Day(),
				record.Timestamp.Hour(),
				record.Timestamp.Minute(),
				record.Timestamp.Second(),
				record.Timestamp.Nanosecond(),
				vietnamLoc,
			)
		}

		recordTime := record.Timestamp.In(vietnamLoc)

		if checkinFrom != nil && checkinTo != nil {
			if (recordTime.After(*checkinFrom) || recordTime.Equal(*checkinFrom)) &&
				(recordTime.Before(*checkinTo) || recordTime.Equal(*checkinTo)) {
				checkInRecords = append(checkInRecords, record)
				continue
			}
		}

		if checkoutFrom != nil && checkoutTo != nil {
			if (recordTime.After(*checkoutFrom) || recordTime.Equal(*checkoutFrom)) &&
				(recordTime.Before(*checkoutTo) || recordTime.Equal(*checkoutTo)) {
				checkOutRecords = append(checkOutRecords, record)
			}
		}
	}

	sort.Slice(checkInRecords, func(i, j int) bool {
		return checkInRecords[i].Timestamp.Before(checkInRecords[j].Timestamp)
	})

	sort.Slice(checkOutRecords, func(i, j int) bool {
		return checkOutRecords[i].Timestamp.After(checkOutRecords[j].Timestamp)
	})

	return checkInRecords, checkOutRecords
}

func normalizeToDay(t time.Time, loc *time.Location) time.Time {
	// Check if the time has no location or has a nil location
	if t.Location() == nil {
		// If no location is set, create a new time with the provided location
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}

	localTime := t.In(loc)
	dayStart := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, loc)
	return dayStart
}

func (t *timesheetServiceImp) ResetWorkDayAdjustment(
	ctx context.Context,
	timesheetDetailID int,
) error {
	// Get detail data
	detail, err := t.tsDetailsRepo.GetDetailByID(ctx, timesheetDetailID)
	if err != nil {
		return err
	}

	//get timesheet
	if err != nil {
		return errors.New("Lỗi khi lấy timesheet")
	}

	// Reset
	if detail.OriginalWorkDays != nil {
		detail.WorkDays = *detail.OriginalWorkDays
	}
	detail.OriginalWorkDays = nil
	detail.WorkDaysAdjusted = 0
	detail.IsManuallyAdjusted = false
	detail.AdjustmentBy = nil
	detail.AdjustmentAt = nil

	// update detail and total workday
	return t.tsDetailsRepo.UpdateTimeSheetDetail(ctx, detail)
}

// service/timesheet_service.go
func (t *timesheetServiceImp) ManualAdjustWorkDay(
	ctx context.Context,
	timesheetDetailID int,
	adjustedWorkDays float64,
	adjustedBy string,
) error {
	// get detail have in database
	detail, err := t.tsDetailsRepo.GetDetailByID(ctx, timesheetDetailID)
	if err != nil {
		return err
	}

	//get timesheet
	if err != nil {
		return errors.New("Lỗi khi lấy timesheetDetail")
	}

	timesheet, err := t.timesheetRepo.GetByID(ctx, detail.TimeSheetID)
	if err != nil {
		return errors.New("Lỗi khi lấy dữ liệu timesheet")
	}

	timesheetList, err := t.tsListRepo.GetByID(ctx, timesheet.TimeSheetListID)
	if err != nil {
		return errors.New("Lỗi khi lấy dữ liệu timesheetList")
	}
	if timesheetList.IsLocked {
		return errors.New("Bảng công được chốt, không thể sửa")
	}

	// save old workday
	if detail.OriginalWorkDays == nil {
		value := detail.WorkDays
		detail.OriginalWorkDays = &value
	}

	// update manual workday and
	detail.WorkDays = adjustedWorkDays
	detail.WorkDaysAdjusted = adjustedWorkDays
	detail.IsManuallyAdjusted = true
	detail.AdjustmentBy = &adjustedBy
	now := time.Now().In(vietnamLoc) // Use Vietnam timezone instead of UTC
	detail.AdjustmentAt = &now

	// update detail
	return t.tsDetailsRepo.UpdateTimeSheetDetail(ctx, detail)
}

func (b *timesheetServiceImp) ExportTimeSheet(ctx context.Context, timesheetListID string) ([]byte, string, error) {
	timesheetList, err := b.tsListRepo.GetForExport(ctx, timesheetListID)
	if err != nil {
		return nil, "", fmt.Errorf("lỗi load dữ liệu: %w", err)
	}
	timesheets := timesheetList.Timesheets
	year := timesheetList.Year
	month := timesheetList.Month
	exporter := utils.NewExcelExporter(fmt.Sprintf("Timesheet_%d_%02d", year, month))

	// get localtime
	timeNow := utils.GetCurrentTimeHCMCity()

	// Các trường cố định
	exporter.RegisterField("employee_id", "Mã MV", "employee_id",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.EmployeeID != "" {
				return ts.EmployeeID
			}
			return ""
		})

	exporter.RegisterField("fullname", "Họ và tên", "fullname",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.Employee != nil && ts.Employee.Fullname != "" {
				return ts.Employee.Fullname
			}
			return ""
		})

	exporter.RegisterField("office", "Văn phòng", "office",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.Office != nil && ts.Office.Name != "" {
				return ts.Office.Name
			}
			return ""
		})

	exporter.RegisterField("department", "Phòng ban", "department",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.Department != nil && ts.Department.Name != "" {
				return ts.Department.Name
			}
			return ""
		})

	exporter.RegisterField("jobTitle", "Chức vụ", "jobTitle",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.Employee != nil && ts.Employee.JobTitle.JobTitle != "" {
				return ts.Employee.JobTitle.JobTitle
			}
			return ""
		})

	exporter.RegisterField("worktype", "Worktype", "worktype",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.Employee != nil && ts.Employee.WorkType != "" {
				return ts.Employee.WorkType
			}
			return ""
		})

	// Tạo các cột cho từng ngày trong tháng
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, timeNow.Location()).Day()
	for day := 1; day <= daysInMonth; day++ {
		// Capture day trong closure
		currentDay := day
		key := fmt.Sprintf("day_%02d", currentDay)
		header := fmt.Sprintf("%02d/%02d/%d", currentDay, month, year)

		exporter.RegisterField(key, header, key, func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.Details == nil {
				return ""
			}

			for _, d := range ts.Details {
				if d.Date.Day() == currentDay &&
					int(d.Date.Month()) == month &&
					d.Date.Year() == year {
					// Kiểm tra nếu có điều chỉnh thủ công
					if d.IsManuallyAdjusted && d.WorkDaysAdjusted > 0 {
						return fmt.Sprintf("%.2f", d.WorkDaysAdjusted)
					}
					// Nếu không có điều chỉnh, dùng WorkDays
					if d.WorkDays > 0 {
						return fmt.Sprintf("%.2f", d.WorkDays)
					}
					// Nếu là ngày nghỉ hoặc holiday
					if d.IsHoliday {
						return "H"
					}
					if d.IsWeekend && !d.IsWorkingDay {
						return "W"
					}
					if d.IsAbsent {
						return "X"
					}
					if d.LeaveType != nil {
						switch *d.LeaveType {
						case variable.LeaveTypeAnnual:
							return "AL"
						case variable.LeaveTypePersonal:
							return "PL"
						case variable.LeaveTypeBusinessTrip:
							return "BT"
						case variable.LeaveTypeRemoteWork:
							return "RW"
						case variable.LeaveTypeUnpaid:
							return "UL"
						default:
							return "L"
						}
					}
					return "0.00"
				}
			}
			return ""
		})
	}

	// Các trường tổng kết
	exporter.RegisterField("total_work_days", "Tổng số công", "total_work_days",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			return fmt.Sprintf("%.2f", ts.TotalWorkDays)
		})

	exporter.RegisterField("late_shifts", "Ca muộn", "late_shifts",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			return ts.LateShifts
		})

	exporter.RegisterField("late_minutes", "Số phút muộn", "late_minutes",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			return ts.TotalLateMinutes
		})

	// Thêm các trường thống kê nghỉ phép
	exporter.RegisterField("annual_leave_days", "Nghỉ phép năm", "annual_leave_days",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.AnnualLeaveDays > 0 {
				return fmt.Sprintf("%.2f", ts.AnnualLeaveDays)
			}
			return ""
		})

	exporter.RegisterField("personal_leave_days", "Nghỉ phép cá nhân", "personal_leave_days",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.PersonalLeaveDays > 0 {
				return fmt.Sprintf("%.2f", ts.PersonalLeaveDays)
			}
			return ""
		})

	exporter.RegisterField("business_trip_days", "Công tác", "business_trip_days",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.BusinessTripDays > 0 {
				return fmt.Sprintf("%.2f", ts.BusinessTripDays)
			}
			return ""
		})

	exporter.RegisterField("remote_work_days", "Làm việc từ xa", "remote_work_days",
		func(item interface{}) any {
			ts := item.(*model.TimeSheet)
			if ts.RemoteWorkDays > 0 {
				return fmt.Sprintf("%.2f", ts.RemoteWorkDays)
			}
			return ""
		})

	// Chuyển đổi sang slice pointer
	ptrs := make([]*model.TimeSheet, len(timesheets))
	for i := range timesheets {
		ptrs[i] = &timesheets[i]
	}

	// Tạo danh sách field theo thứ tự mong muốn
	selected := []string{
		"employee_id", "fullname", "office", "department", "jobTitle", "worktype",
	}

	// Thêm các ngày trong tháng
	for day := 1; day <= daysInMonth; day++ {
		selected = append(selected, fmt.Sprintf("day_%02d", day))
	}

	// Thêm các field tổng kết
	selected = append(selected,
		"total_work_days", "late_shifts", "late_minutes",
		"annual_leave_days", "personal_leave_days",
		"business_trip_days", "remote_work_days")

	data, _, err := exporter.Export(ptrs, selected)
	if err != nil {
		return nil, "", fmt.Errorf("lỗi export excel: %w", err)
	}

	filename := fmt.Sprintf("timesheet_%d_%02d.xlsx", year, month)
	return data, filename, nil
}
