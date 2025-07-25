package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type timesheetService struct {
	timesheetRepo repo_interface.TimesheetDetailInterface
}

func NewTimesheetSerivce(timesheetRepo repo_interface.TimesheetDetailInterface) service_interface.TimesheetServiceInterface {
	return &timesheetService{
		timesheetRepo: timesheetRepo,
	}
}

func (s *timesheetService) Create(ctx context.Context, timesheet *model.Timesheet) error {
	duplicate, err := s.timesheetRepo.IsDuplicate(ctx, timesheet.OfficeID, timesheet.Month, timesheet.Year, "")
	if err != nil {
		return fmt.Errorf("lỗi kiểm tra trùng thời gian bảng công: %w", err)
	}
	if duplicate {
		return fmt.Errorf("đã tồn tại bảng công cho văn phòng và thời gian này")
	}

	code, err := utils.GenerateCode("BC", 3, func() (string, error) {
		lastID, err := s.timesheetRepo.GetLastDecisionByCode(ctx)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
		return lastID, nil
	})
	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}
	timesheet.ID = code

	return s.timesheetRepo.Create(ctx, timesheet)
}

func (s *timesheetService) GetByID(ctx context.Context, id string) (*model.Timesheet, error) {
	return s.timesheetRepo.GetByID(ctx, id)
}

func (s *timesheetService) Update(ctx context.Context, timesheet *model.Timesheet) error {
	duplicate, err := s.timesheetRepo.IsDuplicate(ctx, timesheet.OfficeID, timesheet.Month, timesheet.Year, timesheet.ID)
	if err != nil {
		return fmt.Errorf("lỗi kiểm tra trùng thời gian bảng công: %w", err)
	}
	if duplicate {
		return fmt.Errorf("đã tồn tại bảng công cho văn phòng và thời gian này")
	}
	ts, err := s.GetByID(ctx, timesheet.ID)
	if err != nil {
		return err
	}
	if ts.IsFinalized {
		return errors.New("bảng công đã chốt, không thể sửa đổi")
	}

	return s.timesheetRepo.Update(ctx, timesheet)
}

func (s *timesheetService) Delete(ctx context.Context, id string) error {
	ts, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ts.IsFinalized {
		return errors.New("bảng công đã chốt, không thể xóa")
	}
	return s.timesheetRepo.Delete(ctx, id)
}

func (s *timesheetService) List(ctx context.Context, page, limit int) ([]model.Timesheet, int64, error) {
	return s.timesheetRepo.List(ctx, page, limit)
}
