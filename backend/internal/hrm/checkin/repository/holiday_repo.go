package repository

import (
	"erp/backend/internal/hrm/checkin/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type HolidayRepo struct {
	db *gorm.DB
}

func NewHolidayRepo(db *gorm.DB) *HolidayRepo {
	return &HolidayRepo{db: db}
}

func (s *HolidayRepo) recoverFromPanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		log.Printf("Panic occurred, transaction rolled back: %v", r)
	}
}

func (h *HolidayRepo) CreateHoliday(holiday *model.Holiday) error {
	tx := h.db.Begin()

	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	h.recoverFromPanic(tx)

	// err := model.ValidateEmployeeReferences(tx, *employee)
	// if err != nil {
	// 	tx.Rollback()
	// 	return fmt.Errorf("Lỗi tạo nhân viên: %w", err)
	// }

	if err := tx.Create(holiday).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create holiday: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *HolidayRepo) GetHolidayById(id string) (model.Holiday, error) {
	var holiday model.Holiday
	if err := s.db.Where("holiday_id = ?", id).First(&holiday).Error; err != nil {
		return model.Holiday{}, fmt.Errorf("holiday not found with id %d", id)
	}
	return holiday, nil
}

func (s *HolidayRepo) GetAllHoliday() ([]model.Holiday, error) {
	var Holiday []model.Holiday
	if err := s.db.Table("holiday").Where("is_deleted = ?", false).Find(&Holiday).Error; err != nil {
		return nil, err
	}

	return Holiday, nil
}

func (s *HolidayRepo) UpdateHoliday(id string, holiday model.Holiday) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}
	defer s.recoverFromPanic(tx)

	//get holyday i4
	existing, err := s.GetHolidayById(id)
	if err != nil {
		tx.Rollback()
		return errors.New("failed to get holiday by id")
	}

	model.UpdateHolidayFields(&existing, holiday)

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

func (s *HolidayRepo) DeleteHoliday(id string) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	defer s.recoverFromPanic(tx)

	var holiday model.Holiday
	if err := tx.Where("holiday_id = ?", id).First(&holiday).Error; err != nil {
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

func (s *HolidayRepo) GetLastHolidayByCode(emp *model.Holiday) error {
	return s.db.
		Order("holiday_id DESC").
		First(emp).Error
}
