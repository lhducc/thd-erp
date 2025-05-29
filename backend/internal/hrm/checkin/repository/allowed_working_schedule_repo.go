package repository

import (
	"erp/backend/internal/hrm/checkin/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type AllowedWorkingScheduleRepo struct {
	db *gorm.DB
}

func NewAllowedWorkingSchedule(db *gorm.DB) *AllowedWorkingScheduleRepo {
	return &AllowedWorkingScheduleRepo{db: db}
}

func (s *AllowedWorkingScheduleRepo) recoverFromPanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		log.Printf("Panic occurred, transaction rolled back: %v", r)
	}
}

func (h *AllowedWorkingScheduleRepo) CreateAllowedWorkingSchedule(a *model.AllowedWorkingSchedule) error {
	tx := h.db.Begin()

	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	defer h.recoverFromPanic(tx)

	// err := model.ValidateEmployeeReferences(tx, *employee)
	// if err != nil {
	// 	tx.Rollback()
	// 	return fmt.Errorf("Lỗi tạo nhân viên: %w", err)
	// }

	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create holiday: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *AllowedWorkingScheduleRepo) GetAllowedWorkingScheduleId(id string) (model.AllowedWorkingSchedule, error) {
	var holiday model.AllowedWorkingSchedule
	if err := s.db.Where("id = ?", id).First(&holiday).Error; err != nil {
		return model.AllowedWorkingSchedule{}, fmt.Errorf("holiday not found with id %d", id)
	}
	return holiday, nil
}

func (s *AllowedWorkingScheduleRepo) GetAllAllowedWorkingSchedule() ([]model.AllowedWorkingSchedule, error) {
	var Holiday []model.AllowedWorkingSchedule
	if err := s.db.Table("allowed_working_schedule").Where("is_deleted = ?", false).Find(&Holiday).Error; err != nil {
		return nil, err
	}

	return Holiday, nil
}

func (s *AllowedWorkingScheduleRepo) UpdateAllowedWorkingSchedule(id string, holiday model.AllowedWorkingSchedule) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}
	defer s.recoverFromPanic(tx)

	existing, err := s.GetAllowedWorkingScheduleId(id)
	if err != nil {
		tx.Rollback()
		return errors.New("failed to get holiday by id")
	}

	model.UpdateAllowedWorkingScheduleFields(&existing, holiday)

	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update holiday: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *AllowedWorkingScheduleRepo) DeleteAllowedWorkingSchedule(id string) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	defer s.recoverFromPanic(tx)

	var holiday model.AllowedWorkingSchedule
	if err := tx.Where("id = ?", id).First(&holiday).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("holiday with ID %d not found", id)
		}
		return fmt.Errorf("failed to check holiday: %w", err)
	}

	if err := tx.Model(&holiday).Update("is_deleted", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update holiday status to 'Inactive': %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *AllowedWorkingScheduleRepo) GetLastAllowedWorkingScheduleByCode(emp *model.AllowedWorkingSchedule) error {
	return s.db.
		Order("id DESC").
		First(emp).Error
}
