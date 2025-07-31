package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/pkg/pointer"
	"errors"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

var vietnamLoc = func() *time.Location {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return loc
}()

type timesheetServiceImp struct {
	timesheetRepo        repo_interface.TimeSheetRepoInterface
	tsListRepo           repo_interface.TimesheetListInterface
	attendanceRecordRepo repo_interface.AttendanceRecordRepository
	empWSRepo            repo_interface.EmployeeWorkShiftRepo
}

func NewTimesheetService(
	timesheetRepo repo_interface.TimeSheetRepoInterface,
	tsListRepo repo_interface.TimesheetListInterface,
	attendanceRecordRepo repo_interface.AttendanceRecordRepository,
	empWSRepo repo_interface.EmployeeWorkShiftRepo,
) service_interface.TimeSheetServiceInterface {
	return &timesheetServiceImp{
		timesheetRepo:        timesheetRepo,
		tsListRepo:           tsListRepo,
		attendanceRecordRepo: attendanceRecordRepo,
		empWSRepo:            empWSRepo,
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

	// Tạo context có thể hủy
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Kênh cho các job và kết quả
	jobs := make(chan model.TimeSheet, len(tsl.Timesheets))
	results := make(chan error, len(tsl.Timesheets))

	// Số lượng worker (tối ưu dựa trên số lượng job)
	numWorkers := 5
	if numWorkers > len(tsl.Timesheets) {
		numWorkers = len(tsl.Timesheets)
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Khởi động worker
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for ts := range jobs {
				// Kiểm tra context trước khi xử lý
				if ctx.Err() != nil {
					results <- nil
					return
				}

				// Tính toán cho từng nhân viên
				err := t.calculateForEmployee(ctx, ts, tsl)
				results <- err
			}
		}()
	}

	// Gửi job vào kênh
	for _, ts := range tsl.Timesheets {
		jobs <- ts
	}
	close(jobs)

	// Đợi worker hoàn thành và đóng kênh kết quả
	go func() {
		wg.Wait()
		close(results)
	}()

	// Thu thập kết quả
	var firstErr error
	for i := 0; i < len(tsl.Timesheets); i++ {
		err := <-results
		if err != nil && firstErr == nil {
			firstErr = err
			cancel() // Hủy context khi có lỗi đầu tiên
		}
	}

	return firstErr
}

// Hàm tính toán cho từng nhân viên (được gọi song song)
func (t *timesheetServiceImp) calculateForEmployee(
	ctx context.Context,
	ts model.TimeSheet,
	tsl *model.TimeSheetList,
) error {
	// Kiểm tra context trước các thao tác tốn kém
	if ctx.Err() != nil {
		return ctx.Err()
	}

	records, err := t.attendanceRecordRepo.ListHistoryRecordApproveByEmpID(ctx, ts.EmployeeID, tsl.StartDate, tsl.EndDate)
	if err != nil {
		return err
	}

	// Kiểm tra context sau mỗi thao tác I/O
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
		dateKey := NormalizeToDay(shift.Date, vietnamLoc)
		shiftMap[dateKey] = shift
	}

	details := make([]model.TimeSheetDetail, 0)
	totalWorkDays := 0.0
	totalLateMinutes := 0
	lateShifts := 0

	currentDate := tsl.StartDate.In(vietnamLoc)
	endDate := tsl.EndDate.In(vietnamLoc)

	currentDate = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, vietnamLoc)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, vietnamLoc)
	for currentDate.Before(tsl.EndDate) || currentDate.Equal(endDate) {
		// Kiểm tra context trước mỗi lần lặp
		if ctx.Err() != nil {
			return ctx.Err()
		}

		dateKey := NormalizeToDay(currentDate, vietnamLoc)

		// Chỉ xử lý nếu có ca làm việc
		if shift, exists := shiftMap[dateKey]; exists {
			detail := model.TimeSheetDetail{
				TimeSheetID:  ts.TimeSheetID,
				Date:         dateKey,
				DayOfWeek:    int(dateKey.Weekday()),
				IsWorkingDay: true,
				IsWeekend:    isWeekend(dateKey),
				WorkShiftID:  &shift.WorkShiftID,
			}

			// Lấy khung thời gian check-in/out từ ca làm việc
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

			// Xử lý nếu có bản ghi check-out hợp lệ
			if len(checkOutRecords) > 0 {
				detail.CheckOutRecordID = &checkOutRecords[0].AttendanceRecordID
			}

			// Tính toán nếu có cả check-in và check-out
			if len(checkInRecords) > 0 && len(checkOutRecords) > 0 {
				firstCheckIn := checkInRecords[0].Timestamp.In(vietnamLoc)
				lastCheckOut := checkOutRecords[0].Timestamp.In(vietnamLoc)

				// caculator real time working
				workHours := lastCheckOut.Sub(firstCheckIn).Hours()
				detail.WorkHours = math.Round(workHours*100) / 100

				// convert string time to time.Time
				startTime, err := time.ParseInLocation("15:04:05", shift.WorkShift.StartTime, dateKey.Location())
				if err != nil {
					currentDate = currentDate.AddDate(0, 0, 1)
					continue
				}

				shiftStart := time.Date(
					dateKey.Year(),
					dateKey.Month(),
					dateKey.Day(),
					startTime.Hour(),
					startTime.Minute(),
					0, 0, dateKey.Location(),
				)

				// check
				hasValidCheckIn := checkinFrom != nil && checkinTo != nil &&
					firstCheckIn.After(checkinFrom.Add(-1*time.Minute)) &&
					firstCheckIn.Before(checkinTo.Add(1*time.Minute))

				hasValidCheckOut := checkoutFrom != nil && checkoutTo != nil &&
					lastCheckOut.After(checkoutFrom.Add(-1*time.Minute)) &&
					lastCheckOut.Before(checkoutTo.Add(1*time.Minute))

				if hasValidCheckIn && hasValidCheckOut {
					detail.WorkDays = float64(shift.WorkShift.WorkDay)
					totalWorkDays += float64(shift.WorkShift.WorkDay)

					if firstCheckIn.After(shiftStart) {
						lateMinutes := int(firstCheckIn.Sub(shiftStart).Minutes())
						detail.LateMinutes = lateMinutes
						detail.IsLate = true
						totalLateMinutes += lateMinutes
						lateShifts++
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

	// Kiểm tra context trước khi lưu
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

	parseTime := func(timeStr *string) (*time.Time, error) {
		if timeStr == nil || *timeStr == "" {
			return nil, nil
		}

		parsed, err := time.ParseInLocation("15:04:05", *timeStr, location)
		if err != nil {
			return nil, err
		}

		result := time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			parsed.Hour(),
			parsed.Minute(),
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

func NormalizeToDay(t time.Time, loc *time.Location) time.Time {
	localTime := t.In(loc)
	dayStart := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, loc)
	return dayStart
}
