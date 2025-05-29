package repository

import (
	"erp/backend/internal/hrm/checkin/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type WorkShiftStore struct {
	db *gorm.DB
}

func NewWorkShiftStore(db *gorm.DB) *WorkShiftStore {
	return &WorkShiftStore{db: db}
}

func (s *WorkShiftStore) recoverFromPanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		log.Printf("Panic occurred, transaction rolled back: %v", r)
	}
}

func (s *WorkShiftStore) CreateWorkShift(workshift *model.WorkShifts) error {
	tx := s.db.Begin()

	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	defer s.recoverFromPanic(tx)

	// err := model.ValidateEmployeeReferences(tx, *employee)
	// if err != nil {
	// 	tx.Rollback()
	// 	return fmt.Errorf("Lỗi tạo nhân viên: %w", err)
	// }

	if err := tx.Create(workshift).Error; err != nil {
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
	if err := s.db.Where("work_shift_id = ?", id).First(&workshift).Error; err != nil {
		return model.WorkShifts{}, fmt.Errorf("work shift not found with id %d", id)
	}
	return workshift, nil
}

func (s *WorkShiftStore) GetAllWorkShift() ([]model.WorkShifts, error) {
	var workShifts []model.WorkShifts
	if err := s.db.Table("work_shifts").Where("is_deleted = ?", false).Find(&workShifts).Error; err != nil {
		return nil, err
	}

	return workShifts, nil
}

func (s *WorkShiftStore) UpdateWorkShift(id string, workshift model.WorkShifts) error {
	tx := s.db.Begin()

	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	s.recoverFromPanic(tx)
	w, err := s.GetWorkShiftById(id)
	if err != nil {
		return errors.New("failed to get work shift by id")
	}

	model.UpdateWorkShiftFields(&w, workshift)

	if err := tx.Save(&w).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update employee: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
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
	if err := tx.Where("work_shift_id = ?", id).First(&workshift).Error; err != nil {
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

func (s *WorkShiftStore) GetLastWorkShiftByCode(emp *model.WorkShifts) error {
	return s.db.
		Order("work_shift_id DESC").
		First(emp).Error
}
