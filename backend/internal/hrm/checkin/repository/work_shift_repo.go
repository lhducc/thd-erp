package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type WorkShiftStore struct {
	db *gorm.DB
}

func NewWorkShiftStore(db *gorm.DB) repo_interface.WorkShiftRepo {
	return &WorkShiftStore{db: db}
}

func (s *WorkShiftStore) GetWorkShiftById(ctx context.Context, id string) (*model.WorkShifts, error) {
	var workshift model.WorkShifts
	if err := s.db.WithContext(ctx).Where("workshift_id = ? AND is_deleted = ? ", id, false).First(&workshift).Error; err != nil {
		return nil, err
	}
	return &workshift, nil
}

func (s *WorkShiftStore) GetAllWorkShift(ctx context.Context) ([]model.WorkShifts, error) {
	var workShifts []model.WorkShifts
	if err := s.db.WithContext(ctx).Where("is_deleted = ?", false).Find(&workShifts).Error; err != nil {
		return nil, err
	}

	return workShifts, nil
}

func (s *WorkShiftStore) SaveWorkShift(ctx context.Context, workshift *model.WorkShifts) error {
	var existing model.WorkShifts
	err := s.db.WithContext(ctx).Where("workshift_id = ?", workshift.WorkShiftID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&workshift).Error; err != nil {
				return err
			}
		} else {
			return fmt.Errorf("lỗi khi kiểm tra ca làm việc: %w", err)
		}
	} else {
		if err := s.db.Model(&existing).Updates(workshift).Error; err != nil {
			return fmt.Errorf("lỗi khi cập nhật ca làm việc: %w", err)
		}
	}
	return nil
}

func (s *WorkShiftStore) DeleteWorkShift(ctx context.Context, id string) error {
	var workshift model.WorkShifts
	if err := s.db.WithContext(ctx).Where("workshift_id = ?", id).First(&workshift).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("workshift with ID %d not found", id)
		}
		return fmt.Errorf("failed to check workshift: %w", err)
	}

	if err := s.db.Model(&workshift).Update("is_deleted", true).Error; err != nil {
		return fmt.Errorf("failed to update workshift status to 'Inactive': %w", err)
	}
	return nil
}

func (s *WorkShiftStore) GetLastWorkShiftByCode(ctx context.Context, emp *model.WorkShifts, predix string) error {
	return s.db.WithContext(ctx).
		Where("workshift_id LIKE ?", predix+"%").
		Order("workshift_id DESC").
		First(&emp).Error
}

func (s *WorkShiftStore) CheckExistByName(ctx context.Context, nameWS string, excludeID string) (bool, error) {
	var workShifts model.WorkShifts
	query := s.db.WithContext(ctx).Where("workshift_name = ?", nameWS)

	if excludeID != "" {
		query = query.Where("workshift_id != ?", excludeID)
	}

	err := query.First(&workShifts).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *WorkShiftStore) IsDuplicateTimeRange(ctx context.Context, startTime, endTime, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.WorkShifts{}).
		Where("start_time = ? AND end_time = ? AND is_deleted = false", startTime, endTime)

	if excludeID != "" {
		query = query.Where("workshift_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check duplicate time range: %w", err)
	}

	return count > 0, nil
}
