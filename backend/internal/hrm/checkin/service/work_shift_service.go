package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type WorkShiftService struct {
	repo repo_interface.WorkShiftRepo
}

func NewWorkShiftService(repo repo_interface.WorkShiftRepo) service_interface.WorkShiftService {
	return &WorkShiftService{repo: repo}
}

func (sv *WorkShiftService) CreateWorkShift(ctx context.Context, w *model.WorkShifts) error {
	_, err := sv.repo.GetWorkShiftById(ctx, w.WorkShiftID)
	if err == nil {
		return errors.New("Mã ca làm việc đã tồn tại")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	exist, err := sv.repo.CheckExistByName(ctx, w.WorkShiftName, "")
	if err != nil {
		return errors.New("Lỗi hệ thống, Không thể kiểm tra tồn tại của tên ca")
	}
	if exist {
		return errors.New("Tên ca đã tồn tại trong hệ thống dữ liệu")
	}

	isDup, err := sv.repo.IsDuplicateTimeRange(ctx, w.StartTime, w.EndTime, "")
	if err != nil {
		fmt.Errorf("Lỗi: " + err.Error())
		return errors.New("Lỗi hệ thống khi kiểm tra trùng lặp khung thời gian bắt đầu, kết thúc")
	}
	if isDup {
		return fmt.Errorf("Ca làm việc với StartTime %s và EndTime %s đã tồn tại", w.StartTime, w.EndTime)
	}

	if err := sv.repo.SaveWorkShift(ctx, w); err != nil {
		fmt.Printf(err.Error())
		return errors.New("Lỗi hệ thống")
	}
	return nil
}

func (biz *WorkShiftService) GetWorkShiftById(ctx context.Context, id string) (*model.WorkShifts, error) {
	if id == "" {
		return nil, errors.New("invalid workshift ID")
	}

	workshift, err := biz.repo.GetWorkShiftById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workshift: %w", err)
	}

	return workshift, nil
}

func (biz *WorkShiftService) GetAllWorkShift(ctx context.Context) ([]model.WorkShifts, error) {
	workshifts, err := biz.repo.GetAllWorkShift(ctx)
	if err != nil {
		fmt.Printf("failed to get all workshifts: %w", err)
		return nil, errors.New("Lỗi hệ thống")
	}

	return workshifts, nil
}

func (sv *WorkShiftService) UpdateWorkShift(ctx context.Context, workshift *model.WorkShifts) error {
	_, err := sv.repo.GetWorkShiftById(ctx, workshift.WorkShiftID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("Mã ca làm việc cần cập nhật không tồn tại")
	}
	if err != nil {
		fmt.Printf("failed to get workshift: %w", err)
		return fmt.Errorf("Lỗi khi tìm ca làm việc")
	}

	exist, err := sv.repo.CheckExistByName(ctx, workshift.WorkShiftName, workshift.WorkShiftID)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("Tên ca đã tồn tại trong hệ thống dữ liệu")
	}

	isDup, err := sv.repo.IsDuplicateTimeRange(ctx, workshift.StartTime, workshift.EndTime, workshift.WorkShiftID)
	if err != nil {
		fmt.Errorf("Lỗi: " + err.Error())
		return errors.New("Lỗi hệ thống khi kiểm tra tùng lặp khung thời gian bắt đầu, kết thúc")
	}
	if isDup {
		return fmt.Errorf("ca làm việc với StartTime %s và EndTime %s đã tồn tại", workshift.StartTime, workshift.EndTime)
	}

	if err := sv.repo.SaveWorkShift(ctx, workshift); err != nil {
		fmt.Printf("failed to save workshift: %w", err)
		return errors.New("Lỗi hệ thống")
	}

	return nil
}

func (biz *WorkShiftService) DeleteWorkShift(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("ID không hợp lệ")
	}
	if err := biz.repo.DeleteWorkShift(ctx, id); err != nil {
		fmt.Printf("failed to delete employeeRepo: %w", err)
		return errors.New("Lỗi hệ thống")
	}

	return nil
}
