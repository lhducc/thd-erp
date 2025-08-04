package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	hrm_model "erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg/transaction"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
             effective_date, expiration_date, status, is_deleted, is_schedule_auto)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			workSchedule.WorkScheduleName,
			workSchedule.OfficeID,
			workSchedule.RepeatType,
			workSchedule.RepeatCycle,
			workSchedule.EffectiveDate,
			workSchedule.ExpirationDate,
			workSchedule.Status,
			workSchedule.IsDeleted,
			workSchedule.IsScheduleAuto,
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
	return utils.WithTransaction(r.db, c, func(ctx context.Context, tx *gorm.DB) error {
		if err := tx.WithContext(c).Model(&hrm_model.Employee{}).Where("schedule_id = ?", id).
			Update("schedule_id", nil).Error; err != nil {
			return err
		}
		if err := tx.WithContext(c).Model(&model.WorkScheduleManager{}).Where("work_schedule_id = ?", id).
			Delete(&model.WorkScheduleManager{}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(c).
			Model(&model.WorkSchedule{}).
			Where("work_schedule_id = ? AND is_deleted = false", id).
			Updates(wschedule).Error; err != nil {
			return err
		}
		return nil
	})
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

func (r *workScheduleRepoImpl) IsExistsByScheduleIDAuto(c context.Context, id int) (bool, error) {
	var workSchedule *model.WorkSchedule
	err := r.db.WithContext(c).Model(model.WorkSchedule{}).
		Where("work_schedule_id = ? AND is_deleted = ? AND is_schedule_auto = true", id, false).
		First(&workSchedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

func (r *workScheduleRepoImpl) IsExistsByScheduleIDRegister(c context.Context, id int) (bool, error) {
	var workSchedule *model.WorkSchedule
	err := r.db.WithContext(c).Model(model.WorkSchedule{}).
		Where("work_schedule_id = ? AND is_deleted = ? AND is_schedule_auto = false", id, false).
		First(&workSchedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

func (r *workScheduleRepoImpl) AssignOrUpdateManagers(
	ctx context.Context,
	managers []model.WorkScheduleManager,
) error {
	// Upsert cho Managers
	if len(managers) > 0 {
		if err := r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "employee_id"}, {Name: "work_schedule_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"is_reading", "is_editing"}),
			}).
			Create(&managers).Error; err != nil {
			return fmt.Errorf("failed to upsert managers: %w", err)
		}
	}

	return nil
}

func (r *workScheduleRepoImpl) GetAllScheduleAuto(ctx context.Context) ([]model.WorkSchedule, error) {
	var schedules []model.WorkSchedule
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false AND is_schedule_auto = true").
		Preload("Office").
		Preload("Managers").
		Preload("Managers.Employee").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *workScheduleRepoImpl) GetAllScheduleRegister(ctx context.Context) ([]model.WorkSchedule, error) {
	var schedules []model.WorkSchedule
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false AND is_schedule_auto = false").
		Preload("Office").
		Preload("Managers").
		Preload("Managers.Employee").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *workScheduleRepoImpl) GetByID(ctx context.Context, id int) (*model.WorkSchedule, error) {
	var workSchedule model.WorkSchedule
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false").
		Preload("Office").
		Preload("Managers").
		Preload("Weekdays").
		Preload("Weekdays.WorkShift").
		Preload("Managers.Employee").
		Where("work_schedule_id = ?", id).First(&workSchedule).Error; err != nil {

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &workSchedule, nil
}

func (r *workScheduleRepoImpl) DeleteManagerFromWorkSchedule(
	ctx context.Context,
	managerID string,
	workScheduleID int) error {
	var wsEmp model.WorkScheduleManager

	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND work_schedule_id = ?", managerID, workScheduleID).
		First(&wsEmp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("Quản lý %s không tồn tại ở lịch tự động", managerID)
	}

	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Where("employee_id = ? AND work_schedule_id = ?", managerID, workScheduleID).
		Delete(&model.WorkScheduleManager{}).Error; err != nil {
		return fmt.Errorf("Lỗi khi xóa quản lý khỏi lịch: %w", err)
	}

	return nil
}

func (r *workScheduleRepoImpl) GetListShiftRegister(ctx context.Context, scheduleID *int) ([]model.WorkScheduleShift, error) {
	var scheduleShifts []model.WorkScheduleShift

	err := r.db.WithContext(ctx).
		Joins("JOIN work_schedule ws ON ws.work_schedule_id = work_schedule_shift.work_schedule_id").
		Where("ws.is_schedule_auto = false").
		Where("work_schedule_shift.work_schedule_id = ?", scheduleID).
		Preload("WorkShift").
		Find(&scheduleShifts).Error

	if err != nil {
		return nil, err
	}
	return scheduleShifts, nil
}

func (r *workScheduleRepoImpl) CheckManagerPermission(ctx context.Context, managerID, employeeID string) (*model.WorkScheduleManager, error) {
	var result model.WorkScheduleManager

	err := r.db.WithContext(ctx).
		Table("work_schedule_manager AS m1").
		Joins("JOIN work_schedule ws ON ws.work_schedule_id = m1.work_schedule_id").
		Joins("JOIN work_schedule_manager AS m2 ON m2.work_schedule_id = ws.work_schedule_id").
		Where("m1.employee_id = ?", managerID).
		Where("m2.employee_id = ?", employeeID).
		Where("ws.is_deleted = false").
		Select("m1.*").
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
