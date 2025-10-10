package service

import (
	"context"
	checkinmodel "erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	hrmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"time"
	time2 "time"

	"gorm.io/gorm"
)

type timesheetListService struct {
	timesheetListRepo   repo_interface.TimesheetListInterface
	employeeRepo        *repository.UserStore
	timesheetRepo       repo_interface.TimeSheetRepoInterface
	timesheetDetailRepo repo_interface.TimeSheetDetailRepoInterface
}

func NewTimesheetListSerivce(timesheetListRepo repo_interface.TimesheetListInterface,
	employeeRepo *repository.UserStore,
	timesheetRepo repo_interface.TimeSheetRepoInterface,
	timesheetDetailRepo repo_interface.TimeSheetDetailRepoInterface,
) service_interface.TimesheetListServiceInterface {
	return &timesheetListService{
		timesheetListRepo:   timesheetListRepo,
		employeeRepo:        employeeRepo,
		timesheetRepo:       timesheetRepo,
		timesheetDetailRepo: timesheetDetailRepo,
	}
}

func (s *timesheetListService) CreateElementOfTimesheetList(ctx context.Context, timesheet *checkinmodel.TimeSheetList) error {
	// Check duplicate globally by month/year
	duplicate, err := s.timesheetListRepo.IsDuplicate(ctx, timesheet.Month, timesheet.Year, "")
	if err != nil {
		return fmt.Errorf("lỗi kiểm tra trùng thời gian bảng công: %w", err)
	}
	if duplicate {
		return fmt.Errorf("đã tồn tại bảng công cho văn phòng và thời gian này")
	}

	code, err := utils.GenerateCode("BC", 3, func() (string, error) {
		lastID, err := s.timesheetListRepo.GetLastDecisionByCode(ctx)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
		return lastID, nil
	})
	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}
	timesheet.TimeSheetListID = code
	timesheet.StartDate, timesheet.EndDate = utils.GetStartAndEndDateVNTime(timesheet.Month, timesheet.Year)

	if err = s.timesheetListRepo.Create(ctx, timesheet); err != nil {
		return err
	}

	var success bool
	if success, err = s.createTimeSheetEmployee(ctx, timesheet); err != nil {
		return err
	}
	if !success {
		if err = s.Delete(ctx, timesheet.TimeSheetListID); err != nil {
			return err
		}
	}
	return nil
}

func (s *timesheetListService) GetByID(ctx context.Context, id string) (*checkinmodel.TimeSheetList, error) {
	return s.timesheetListRepo.GetByID(ctx, id)
}

func (s *timesheetListService) Update(ctx context.Context, timesheet *checkinmodel.TimeSheetList) error {
	ts, err := s.GetByID(ctx, timesheet.TimeSheetListID)
	if err != nil {
		return err
	}
	if ts.IsLocked {
		return errors.New("bảng công đã chốt, không thể sửa đổi")
	}
	timeNowVN := utils.GetCurrentTimeHCMCity()
	timesheet.UpdatedAt = &timeNowVN

	return s.timesheetListRepo.Update(ctx, timesheet)
}

func (s *timesheetListService) Delete(ctx context.Context, id string) error {
	ts, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ts.IsLocked {
		return errors.New("bảng công đã chốt, không thể xóa")
	}
	return s.timesheetListRepo.Delete(ctx, id)
}

func (s *timesheetListService) List(ctx context.Context, page, limit int) ([]checkinmodel.TimeSheetList, int64, error) {
	return s.timesheetListRepo.List(ctx, page, limit)
}

func (s *timesheetListService) createTimeSheetEmployee(ctx context.Context, timeSheetList *checkinmodel.TimeSheetList) (bool, error) {
	// fetch all employees (global timesheet)
	employeesAll, err := s.employeeRepo.GetAllEmployees()
	if err != nil {
		return false, err
	}

	// convert to pointers to match previous behavior
	var employees []*hrmodel.Employee
	for i := range employeesAll {
		employees = append(employees, &employeesAll[i])
	}
	if len(employees) == 0 {
		return true, nil
	}

	var timesheets []*checkinmodel.TimeSheet
	for _, employee := range employees {
		// office for each timesheet is taken from employee's Department.OfficeID when available
		officeID := ""
		if employee.Department != nil && employee.Department.OfficeID != "" {
			officeID = employee.Department.OfficeID
		}
		timesheet := checkinmodel.TimeSheet{
			TimeSheetListID: timeSheetList.TimeSheetListID,
			OfficeID:        officeID,
			Month:           timeSheetList.Month,
			Year:            timeSheetList.Year,
			DepartmentID:    employee.DepartmentID,
			EmployeeID:      employee.EmployeeID,
			CreatedBy:       timeSheetList.CreatedBy,
		}
		timesheets = append(timesheets, &timesheet)
	}

	if len(timesheets) > 0 {
		err = s.timesheetRepo.CreateTimeSheets(ctx, timesheets)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *timesheetListService) LockedTimesheet(ctx context.Context, timesheet *checkinmodel.TimeSheetList) error {
	ts, err := s.GetByID(ctx, timesheet.TimeSheetListID)
	if err != nil {
		return err
	}
	if ts.IsLocked {
		return errors.New("bảng công đã chốt, không thể sửa đổi")
	}
	timeNow := time2.Now()
	timesheet.LockedAt = &timeNow
	timesheet.IsLocked = true

	return s.timesheetListRepo.UpdateLocked(ctx, timesheet)
}

func (s *timesheetListService) GetByEmployeeID(ctx context.Context, id string) (*checkinmodel.TimeSheetList, error) {
	return s.timesheetListRepo.GetByID(ctx, id)
}

func (s *timesheetListService) ExportCheckinCheckout(ctx context.Context, id string) ([]byte, string, error) {
	ts, err := s.timesheetListRepo.GetForExport(ctx, id)
	if err != nil {
		return nil, "", err
	}

	if ts == nil || len(ts.Timesheets) == 0 {
		return nil, "", fmt.Errorf("không có dữ liệu để xuất")
	}

	activeTimesheets := make([]checkinmodel.TimeSheet, 0)
	for _, t := range ts.Timesheets {
		if t.Employee != nil && t.Employee.Status == "active" {
			activeTimesheets = append(activeTimesheets, t)
		}
	}
	ts.Timesheets = activeTimesheets

	if len(ts.Timesheets) == 0 {
		return nil, "", fmt.Errorf("không có nhân viên active để xuất")
	}

	f := excelize.NewFile()
	sheet := "Checkin-Checkout"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	start := ts.StartDate
	end := ts.EndDate

	// Header thêm cột Phòng ban, Văn phòng
	header := []string{"Tên nhân viên", "Phòng ban", "Văn phòng"}
	dates := []time.Time{}
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		header = append(header, d.Format("02/01"))
		dates = append(dates, d)
	}

	// Style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D9E1F2"}, Pattern: 1},
	})
	weekendStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E2EFDA"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	lateStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFC7CE"}, Pattern: 1},
		Font:      &excelize.Font{Color: "#9C0006"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	normalStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// Header
	for i, title := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, title)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Dữ liệu
	for rowIdx, t := range ts.Timesheets {
		var (
			name       = "-"
			department = "-"
			office     = "-"
		)

		if t.Employee != nil {
			name = t.Employee.Fullname
		}
		if t.Department != nil {
			department = t.Department.Name
			// nếu Department có OfficeName thì lấy
			if t.Department.Office != nil {
				office = t.Department.Office.Name
			}
		} else if t.Office != nil {
			office = t.Office.Name
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx+2), name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx+2), department)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx+2), office)

		for colIdx, day := range dates {
			cell, _ := excelize.CoordinatesToCellName(colIdx+4, rowIdx+2)
			val := "-"
			isLate := false

			for _, d := range t.Details {
				if d.Date.Equal(day) {
					checkin, checkout := "", ""
					if d.CheckInRecord != nil {
						checkin = d.CheckInRecord.Timestamp.Format("15:04")
					}
					if d.CheckOutRecord != nil {
						checkout = d.CheckOutRecord.Timestamp.Format("15:04")
					}
					if checkin != "" || checkout != "" {
						val = fmt.Sprintf("%s - %s", checkin, checkout)
					}
					isLate = d.IsLate
					break
				}
			}

			f.SetCellValue(sheet, cell, val)

			switch day.Weekday() {
			case time.Saturday, time.Sunday:
				f.SetCellStyle(sheet, cell, cell, weekendStyle)
			default:
				if isLate {
					f.SetCellStyle(sheet, cell, cell, lateStyle)
				} else {
					f.SetCellStyle(sheet, cell, cell, normalStyle)
				}
			}
		}
	}

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
	for colIdx := range header {
		colLetter, _ := excelize.ColumnNumberToName(colIdx + 1)
		_ = f.SetColWidth(sheet, colLetter, colLetter, 18)
	}

	fileName := fmt.Sprintf("timesheet_checkin_checkout_%d_%d.xlsx", ts.Month, ts.Year)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", fmt.Errorf("lỗi tạo file excel: %w", err)
	}
	return buf.Bytes(), fileName, nil
}
