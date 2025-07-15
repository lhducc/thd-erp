package repository

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"errors"
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
	result := r.db.Create(assign)
	return result.Error
}

func (r *employeeWorkShiftRepo) Delete(id uint) error {
	result := r.db.Delete(&model.EmployeeWorkshift{}, id)
	return result.Error
}

func (r *employeeWorkShiftRepo) GetAllByEmployeeID(employeeID string) ([]model.EmployeeWorkshift, error) {
	var records []model.EmployeeWorkshift
	result := r.db.Preload("WorkShift").Where("employee_id = ?", employeeID).Find(&records)
	if result.Error != nil {
		return []model.EmployeeWorkshift{}, result.Error
	}
	return records, nil
}

func (r *employeeWorkShiftRepo) IsExisting(userID string, WorkShiftID string, date time.Time) bool {
	var records model.EmployeeWorkshift
	result := r.db.Where("employee_id = ? AND work_shift_id = ? AND date = ?", userID, WorkShiftID, date).First(&records)
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
