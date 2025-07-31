package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type employeeWorkShiftRepo struct {
	db *gorm.DB
}

func NewEmployeeWorkShiftRepo(db *gorm.DB) repo_interface.EmployeeWorkShiftRepo {
	return &employeeWorkShiftRepo{
		db: db,
	}
}

func (r *employeeWorkShiftRepo) GetAll() ([]model.EmployeeWorkshift, error) {
	var employeeWorkShifts []model.EmployeeWorkshift
	result := r.db.Preload("WorkShift").Find(&employeeWorkShifts)
	if result.Error != nil {
		return []model.EmployeeWorkshift{}, result.Error
	}
	return employeeWorkShifts, nil
}

func (r *employeeWorkShiftRepo) Save(assign *model.EmployeeWorkshift) error {
	result := r.db.Save(assign)
	return result.Error
}

func (r *employeeWorkShiftRepo) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&model.EmployeeWorkshift{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee_workshift with id %s not found", id)
	}
	return nil
}

func (r *employeeWorkShiftRepo) GetAllByEmployeeID(employeeID string) ([]model.EmployeeWorkshift, error) {
	var records []model.EmployeeWorkshift
	result := r.db.Where("employee_id = ?", employeeID).
		Model(&model.EmployeeWorkshift{}).
		Preload("WorkShift").Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (r *employeeWorkShiftRepo) IsExisting(userID string, WorkShiftID string, date time.Time) bool {
	var records model.EmployeeWorkshift
	result := r.db.Where("employee_id = ? AND workshift_id = ? AND date = ?", userID, WorkShiftID, date).First(&records)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		return true
	}

	return true
}

func (r *employeeWorkShiftRepo) FindByID(id string) (*model.EmployeeWorkshift, error) {
	var records model.EmployeeWorkshift
	result := r.db.Preload("WorkShift").Where("id = ?", id).First(&records)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		}
	}
	return &records, nil
}

func (r *employeeWorkShiftRepo) GetEmployeeWorkShifts(ctx context.Context, employeeID string) ([]model.EmployeeWorkshift, error) {
	var empShift []model.EmployeeWorkshift

	err := r.db.WithContext(ctx).
		Model(&model.EmployeeWorkshift{}).
		Preload("WorkShift").
		Where("employee_id = ?", employeeID).
		Order("date").
		Find(&empShift).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get work shifts: %w", err)
	}

	return empShift, nil
}

func (r *employeeWorkShiftRepo) GetEmployeeWorkShiftsByMonthYear(ctx context.Context, employeeID string, startDate, endDate time.Time) ([]model.EmployeeWorkshift, error) {
	var workShifts []model.EmployeeWorkshift

	err := r.db.WithContext(ctx).
		Model(&model.EmployeeWorkshift{}).
		Where("employee_id = ?", employeeID).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Preload("WorkShift").
		Order("date").
		Find(&workShifts).Error

	if err != nil {
		return nil, err
	}

	return workShifts, nil
}
