package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/repository"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"gorm.io/gorm"
	time2 "time"
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

func (s *timesheetListService) CreateElementOfTimesheetList(ctx context.Context, timesheet *model.TimeSheetList) error {
	duplicate, err := s.timesheetListRepo.IsDuplicate(ctx, timesheet.OfficeID, timesheet.Month, timesheet.Year, "")
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

func (s *timesheetListService) GetByID(ctx context.Context, id string) (*model.TimeSheetList, error) {
	return s.timesheetListRepo.GetByID(ctx, id)
}

func (s *timesheetListService) Update(ctx context.Context, timesheet *model.TimeSheetList) error {
	ts, err := s.GetByID(ctx, timesheet.TimeSheetListID)
	if err != nil {
		return err
	}
	if ts.IsLocked {
		return errors.New("bảng công đã chốt, không thể sửa đổi")
	}
	timeNowVN, err := utils.GetCurrentTimeHCMCity()
	if err != nil {
		return err
	}
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

func (s *timesheetListService) List(ctx context.Context, page, limit int) ([]model.TimeSheetList, int64, error) {
	return s.timesheetListRepo.List(ctx, page, limit)
}

func (s *timesheetListService) createTimeSheetEmployee(ctx context.Context, timeSheetList *model.TimeSheetList) (bool, error) {
	employees, err := s.employeeRepo.GetEmployeesByOfficeID(ctx, timeSheetList.OfficeID)
	if err != nil {
		return false, err
	}

	if len(employees) == 0 {
		return true, nil
	}

	var timesheets []*model.TimeSheet
	for _, employee := range employees {
		timesheet := model.TimeSheet{
			TimeSheetListID: timeSheetList.TimeSheetListID,
			OfficeID:        timeSheetList.OfficeID,
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

func (s *timesheetListService) LockedTimesheet(ctx context.Context, timesheet *model.TimeSheetList) error {
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

func (s *timesheetListService) GetByEmployeeID(ctx context.Context, id string) (*model.TimeSheetList, error) {
	return s.timesheetListRepo.GetByID(ctx, id)
}
