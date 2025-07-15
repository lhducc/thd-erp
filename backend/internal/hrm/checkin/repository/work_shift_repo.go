package repository

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type WorkShiftStore struct {
	db *gorm.DB
}

func NewWorkShiftStore(db *gorm.DB) repo_interface.WorkShiftRepo {
	return &WorkShiftStore{db: db}
}

func (s *WorkShiftStore) recoverFromPanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		log.Printf("Panic occurred, transaction rolled back: %v", r)
	}
}

func (s *WorkShiftStore) FindByID(id string) (model.WorkShifts, error) {
	var workShifts model.WorkShifts
	err := s.db.Where("workshift_id = ?", id).First(&workShifts).Error
	return workShifts, err
}

func (s *WorkShiftStore) CreateWorkShift(data *model.WorkShifts) error {
	tx := s.db.Begin()

	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	defer s.recoverFromPanic(tx)

	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create workshift: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *WorkShiftStore) GetWorkShiftById(id string) (model.WorkShifts, error) {
	var workshift model.WorkShifts
	if err := s.db.Where("workshift_id = ? AND is_deleted = ? ", id, false).First(&workshift).Error; err != nil {
		return model.WorkShifts{}, fmt.Errorf("work shift not found with id %d", id)
	}
	return workshift, nil
}

func (s *WorkShiftStore) GetAllWorkShift() ([]model.WorkShifts, error) {
	var workShifts []model.WorkShifts
	if err := s.db.Where("is_deleted = ?", false).Find(&workShifts).Error; err != nil {
		return nil, err
	}

	return workShifts, nil
}

func (s *WorkShiftStore) UpdateWorkShift(id string, workshift *model.WorkShifts) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("không thể bắt đầu transaction")
	}
	defer s.recoverFromPanic(tx)

	var existing model.WorkShifts
	if err := tx.Where("workshift_id = ?", id).First(&existing).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy ca làm việc với ID: %s", id)
		}
		return fmt.Errorf("lỗi khi truy vấn ca làm việc: %w", err)
	}

	if err := tx.Model(&existing).Updates(workshift).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi cập nhật ca làm việc: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction thất bại: %w", err)
	}

	return nil
}

func (s *WorkShiftStore) DeleteWorkShift(id string) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}
	defer s.recoverFromPanic(tx)

	var workshift model.WorkShifts
	if err := tx.Where("workshift_id = ?", id).First(&workshift).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("workshift with ID %d not found", id)
		}
		return fmt.Errorf("failed to check workshift: %w", err)
	}

	if err := tx.Model(&workshift).Update("is_deleted", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update workshift status to 'Inactive': %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *WorkShiftStore) GetLastWorkShiftByCode(emp *model.WorkShifts, predix string) error {
	return s.db.
		Where("workshift_id LIKE ?", predix+"%").
		Order("workshift_id DESC").
		First(&emp).Error
}

func (s *WorkShiftStore) IsExactTimeRangeExists(timeOfDay model.TimeOfDayEnum, start, end string) (*model.WorkShifts, error) {
	var w model.WorkShifts
	err := s.db.Table(model.WorkShifts{}.TableName()).
		Where("time_of_day = ?", timeOfDay).
		Where("start_time = ? AND end_time = ?", start, end).
		First(&w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}
