package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	utils "erp/backend/pkg"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type workScheduleRepoImpl struct {
	db *gorm.DB
}

func NewWorkScheduleRepo(db *gorm.DB) repo_interface.WorkScheduleRepo {
	return &workScheduleRepoImpl{
		db: db,
	}
}

func (r *workScheduleRepoImpl) Save(c context.Context, workSchedule *model.WorkSchedule, weekdayShift []model.WorkScheduleShift) error {
	return utils.WithTransaction(r.db, c, func(ctx context.Context, tx *gorm.DB) error {
		tx = tx.Session(&gorm.Session{
			SkipDefaultTransaction:   true,
			DisableNestedTransaction: true,
		})

		if err := tx.Exec(`
            INSERT INTO work_schedule 
            (work_schedule_name, office_id, repeat_type, repeat_cycle, 
             effective_date, expiration_date, status, is_deleted)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			workSchedule.WorkScheduleName,
			workSchedule.OfficeID,
			workSchedule.RepeatType,
			workSchedule.RepeatCycle,
			workSchedule.EffectiveDate,
			workSchedule.ExpirationDate,
			workSchedule.Status,
			workSchedule.IsDeleted,
		).Error; err != nil {
			return err
		}

		var lastID int
		if err := tx.Raw("SELECT lastval()").Scan(&lastID).Error; err != nil {
			return err
		}

		if len(weekdayShift) > 0 {
			for _, shift := range weekdayShift {
				if err := tx.Exec(`
                    INSERT INTO work_schedule_shift 
                    (work_schedule_id, week_day, workshift_id)
                    VALUES (?, ?, ?)`,
					lastID,
					shift.Weekday,
					shift.WorkShiftID,
				).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
func (r *workScheduleRepoImpl) Delete(c context.Context, id int) error {
	wschedule := &model.WorkSchedule{
		IsDeleted: true,
		Status:    variable.InActive,
	}

	err := r.db.WithContext(c).
		Model(&model.WorkSchedule{}).
		Where("work_schedule_id = ? AND is_deleted = false", id).
		Updates(wschedule).Error
	return err
}

func (r *workScheduleRepoImpl) Update(c context.Context, workSchedule *model.WorkSchedule, id int) error {
	return utils.WithTransaction(r.db, c, func(ctx context.Context, tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Model(&model.WorkSchedule{}).
			Where("work_schedule_id = ? AND is_deleted = ?", id, false).
			Select("WorkScheduleName", "OfficeID", "RepeatType", "RepeatCycle",
				"EffectiveDate", "ExpirationDate", "Status").
			Updates(workSchedule).Error; err != nil {
			return fmt.Errorf("lỗi cập nhật work schedule: %v", err)
		}

		if len(workSchedule.Weekdays) > 0 {
			var existingShifts []model.WorkScheduleShift
			if err := tx.WithContext(ctx).
				Where("work_schedule_id = ?", id).
				Find(&existingShifts).Error; err != nil {
				return err
			}

			existingMap := make(map[string]model.WorkScheduleShift)
			for _, shift := range existingShifts {
				key := fmt.Sprintf("%s-%s", shift.Weekday, shift.WorkShiftID)
				existingMap[key] = shift
			}

			var toAdd []model.WorkScheduleShift
			newMap := make(map[string]bool)

			for i, shift := range workSchedule.Weekdays {
				key := fmt.Sprintf("%s-%s", shift.Weekday, shift.WorkShiftID)
				newMap[key] = true

				if _, exists := existingMap[key]; !exists {
					shift.WorkScheduleID = id
					shift.Order = i
					toAdd = append(toAdd, shift)
				}
			}

			var toDeleteConds []string
			var toDeleteArgs []interface{}

			for _, shift := range existingShifts {
				key := fmt.Sprintf("%s-%s", shift.Weekday, shift.WorkShiftID)
				if !newMap[key] {
					toDeleteConds = append(toDeleteConds, "(week_day = ? AND workshift_id = ?)")
					toDeleteArgs = append(toDeleteArgs, shift.Weekday, shift.WorkShiftID)
				}
			}

			if len(toDeleteConds) > 0 {
				whereClause := fmt.Sprintf("work_schedule_id = ? AND (%s)", strings.Join(toDeleteConds, " OR "))
				toDeleteArgs = append([]interface{}{id}, toDeleteArgs...)

				if err := tx.WithContext(ctx).
					Where(whereClause, toDeleteArgs...).
					Delete(&model.WorkScheduleShift{}).Error; err != nil {
					return err
				}
			}

			if len(toAdd) > 0 {
				if err := tx.WithContext(ctx).Create(&toAdd).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *workScheduleRepoImpl) IsExistsByName(c context.Context, name string) (bool, error) {
	var workSchedule *model.WorkSchedule
	err := r.db.WithContext(c).Model(model.WorkSchedule{}).
		Where("work_schedule_name = ? AND is_deleted = ?", name, false).
		First(&workSchedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

func (r *workScheduleRepoImpl) IsExistsByID(c context.Context, id int) (bool, error) {
	var workSchedule *model.WorkSchedule
	err := r.db.WithContext(c).Model(model.WorkSchedule{}).
		Where("work_schedule_id = ? AND is_deleted = ?", id, false).
		First(&workSchedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

func (r *workScheduleRepoImpl) AssignEmployeesToWorkSchedule(ctx context.Context, employeeRecords []model.WorkScheduleEmployee,
	managerRecord []model.WorkScheduleManager) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context, tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&employeeRecords).Error; err != nil {
			return fmt.Errorf("failed to assign employees: %w", err)
		}
		if err := tx.WithContext(ctx).Create(&managerRecord).Error; err != nil {
			return fmt.Errorf("failed to assign employees: %w", err)
		}
		return nil
	})
}

func (r *workScheduleRepoImpl) CheckExistEmployeeSchedule(ctx context.Context, employeeID string) (bool, error) {
	var req *model.WorkScheduleEmployee
	err := r.db.WithContext(ctx).Model(model.WorkScheduleEmployee{}).
		Where("employee_id = ?", employeeID).First(&req).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *workScheduleRepoImpl) GetAll(ctx context.Context) ([]model.WorkSchedule, error) {
	var schedules []model.WorkSchedule
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false").
		Preload("Office").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *workScheduleRepoImpl) GetByID(ctx context.Context, id int) (*model.WorkSchedule, error) {
	var workSchedule *model.WorkSchedule
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false").
		Preload("Office").
		Preload("Managers").
		Preload("Employees").
		Preload("Weekdays").
		Preload("Weekdays.WorkShift").
		Where("work_schedule_id = ?", id).First(&workSchedule).Error; err != nil {

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return workSchedule, nil
}
