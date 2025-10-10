package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	utils "erp/backend/pkg/transaction"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
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

	result := r.db.
		Joins("JOIN employee e ON e.employee_id = employee_workshift.employee_id AND e.status = 'active'").
		Preload("WorkShift").
		Find(&employeeWorkShifts)

	if result.Error != nil {
		return []model.EmployeeWorkshift{}, result.Error
	}

	return employeeWorkShifts, nil
}

func (r *employeeWorkShiftRepo) Save(assign *model.EmployeeWorkshift) error {
	result := r.db.Save(assign)
	return result.Error
}

func (r *employeeWorkShiftRepo) SaveMany(assigns []*model.EmployeeWorkshift) error {
	result := r.db.Save(&assigns)
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

func (r *employeeWorkShiftRepo) DeletePersonalShift(id string, employeeID string) error {
	result := r.db.Where("id = ? AND employee_id = ?", id, employeeID).Delete(&model.EmployeeWorkshift{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
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
		Table("employee_workshift AS ew").
		Joins("JOIN employee AS e ON e.employee_id = ew.employee_id AND e.status = 'active'").
		Where("ew.employee_id = ?", employeeID).
		Where("ew.date BETWEEN ? AND ?", startDate, endDate).
		Preload("WorkShift").
		Order("ew.date").
		Find(&workShifts).Error

	if err != nil {
		return nil, err
	}

	return workShifts, nil
}

func (r *employeeWorkShiftRepo) GetByID(id string) (*model.EmployeeWorkshift, error) {
	var record model.EmployeeWorkshift
	result := r.db.Where("id = ?", id).First(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return &record, nil
}

func (r *employeeWorkShiftRepo) AssignmentShift(ctx context.Context, assigns []model.EmployeeWorkshift, scheduleId []int) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context, tx *gorm.DB) error {
		// 1. Lưu các assignment shift
		if err := tx.WithContext(ctx).Model(&model.EmployeeWorkshift{}).Save(&assigns).Error; err != nil {
			return fmt.Errorf("failed to save employee workshifts: %w", err)
		}

		// 2. Cập nhật trạng thái auto recurring
		query := tx.WithContext(ctx).Model(&model.WorkSchedule{}).
			Where("is_schedule_auto = ?", true)

		// Tối ưu cho cả trường hợp 1 hoặc nhiều scheduleId
		if len(scheduleId) > 0 {
			if len(scheduleId) == 1 {
				query = query.Where("work_schedule_id = ?", scheduleId[0])
			} else {
				query = query.Where("work_schedule_id IN (?)", scheduleId)
			}

		}
		if err := query.Update("is_auto_recurring", true).Error; err != nil {
			return fmt.Errorf("failed to update work schedules: %w", err)
		}

		return nil
	})
}

func (r *employeeWorkShiftRepo) CheckShiftConflict(ctx context.Context, employeeID string, workshiftID string, date time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.EmployeeWorkshift{}).
		Where("employee_id = ? AND workshift_id = ? AND date = ?",
			employeeID,
			workshiftID,
			date.Format("2006-01-02")).
		Count(&count).Error

	return count > 0, err
}

func (r *employeeWorkShiftRepo) GetAllEmployeeWorkShiftsByMonthYear(ctx context.Context, startDate, endDate time.Time) ([]model.EmployeeWorkshift, error) {
	var workshifts []model.EmployeeWorkshift

	err := r.db.WithContext(ctx).
		Model(&model.EmployeeWorkshift{}).
		Joins("JOIN employee AS e ON e.employee_id = employee_workshift.employee_id AND e.status = 'active'").
		Where("employee_workshift.date BETWEEN ? AND ?", startDate, endDate).
		Preload("WorkShift").
		Order("employee_workshift.employee_id, employee_workshift.date").
		Find(&workshifts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get employee workshifts: %w", err)
	}

	return workshifts, nil
}

func (r *employeeWorkShiftRepo) GetEmployeeWorkShiftsByManager(ctx context.Context, managerID string, targetDate string) ([]dto.ManagerEmployeeScheduleDTO, error) {
	var results []dto.ManagerEmployeeScheduleDTO

	t, _ := time.Parse("2006-01-02", targetDate)

	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	query := `
				SELECT
					ew.id,
					ew.employee_id,
					e.full_name,
					e.manager,
					ew.workshift_id,
					ew.date,
					ws.workshift_name,
					ws.start_time,
					ws.end_time
				FROM employee_workshift AS ew
				JOIN employee AS e ON ew.employee_id = e.employee_id
				JOIN workshifts AS ws ON ew.workshift_id = ws.workshift_id
				WHERE e.manager = ?
				  AND DATE(ew.date) BETWEEN ? AND ?
				ORDER BY ew.date ASC
`
	if err := r.db.WithContext(ctx).Raw(query, managerID, firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02")).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get employee workshifts by manager: %w", err)
	}

	return results, nil
}
