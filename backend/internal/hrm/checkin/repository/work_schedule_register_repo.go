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
	"gorm.io/gorm/clause"
	"strings"
)

type workScheduleRegisterRepoImpl struct {
	db *gorm.DB
}

func NewWorkScheduleRegisterRepo(db *gorm.DB) repo_interface.WorkScheduleRegisterRepo {
	return &workScheduleRegisterRepoImpl{
		db: db,
	}
}

func (r *workScheduleRegisterRepoImpl) Save(c context.Context, workSchedule *model.WorkScheduleRegister, weekdayShift []model.WorkScheduleRegisterShift) error {
	return utils.WithTransaction(r.db, c, func(ctx context.Context, tx *gorm.DB) error {
		tx = tx.Session(&gorm.Session{
			SkipDefaultTransaction:   true,
			DisableNestedTransaction: true,
		})

		if err := tx.Exec(`
           INSERT INTO work_schedule_register
           (work_schedule_register_name, office_id, effective_date,
             expiration_date, status, is_deleted)
           VALUES (?, ?, ?, ?, ?, ?)`,
			workSchedule.WorkScheduleRegisterName,
			workSchedule.OfficeID,
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
                   INSERT INTO work_schedule_register_shift
                   (work_schedule_register_id, week_day, workshift_id)
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

func (r *workScheduleRegisterRepoImpl) Delete(c context.Context, id int) error {
	wschedule := &model.WorkScheduleRegister{
		IsDeleted: true,
		Status:    variable.InActive,
	}
	return utils.WithTransaction(r.db, c, func(ctx context.Context, tx *gorm.DB) error {
		if err := tx.WithContext(c).Model(&model.WorkScheduleRegisterEmployee{}).Where("work_schedule_register_id = ?", id).
			Delete(&model.WorkScheduleRegisterEmployee{}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(c).Model(&model.WorkScheduleRegisterManager{}).Where("work_schedule_register_id = ?", id).
			Delete(&model.WorkScheduleRegisterManager{}).Error; err != nil {
			return err
		}
		if err := tx.WithContext(c).
			Model(&model.WorkScheduleRegister{}).
			Where("work_schedule_register_id = ? AND is_deleted = false", id).
			Updates(wschedule).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *workScheduleRegisterRepoImpl) Update(c context.Context, wschedule *model.WorkScheduleRegister, id int) error {
	return utils.WithTransaction(r.db, c, func(ctx context.Context, tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Model(&model.WorkScheduleRegister{}).
			Where("work_schedule_register_id = ? AND is_deleted = ?", id, false).
			Select("WorkScheduleRegisterName", "OfficeID",
				"EffectiveDate", "ExpirationDate", "Status").
			Updates(wschedule).Error; err != nil {
			return fmt.Errorf("lỗi cập nhật work schedule: %v", err)
		}

		if len(wschedule.Weekdays) > 0 {
			var existingShifts []model.WorkScheduleRegisterShift
			if err := tx.WithContext(ctx).
				Where("work_schedule_register_id = ?", id).
				Find(&existingShifts).Error; err != nil {
				return err
			}

			existingMap := make(map[string]model.WorkScheduleRegisterShift)
			for _, shift := range existingShifts {
				key := fmt.Sprintf("%s-%s", shift.Weekday, shift.WorkShiftID)
				existingMap[key] = shift
			}

			var toAdd []model.WorkScheduleRegisterShift
			newMap := make(map[string]bool)

			for i, shift := range wschedule.Weekdays {
				key := fmt.Sprintf("%s-%s", shift.Weekday, shift.WorkShiftID)
				newMap[key] = true

				if _, exists := existingMap[key]; !exists {
					shift.WorkScheduleRegisterID = id
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
				whereClause := fmt.Sprintf("work_schedule_register_id = ? AND (%s)", strings.Join(toDeleteConds, " OR "))
				toDeleteArgs = append([]interface{}{id}, toDeleteArgs...)

				if err := tx.WithContext(ctx).
					Where(whereClause, toDeleteArgs...).
					Delete(&model.WorkScheduleRegisterShift{}).Error; err != nil {
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

func (r *workScheduleRegisterRepoImpl) IsExistsByName(c context.Context, name string) (bool, error) {
	var workSchedule *model.WorkScheduleRegister
	err := r.db.WithContext(c).Model(model.WorkScheduleRegister{}).
		Where("work_schedule_register_name = ? AND is_deleted = ?", name, false).
		First(&workSchedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

func (r *workScheduleRegisterRepoImpl) IsExistsByID(c context.Context, id int) (bool, error) {
	var workSchedule *model.WorkScheduleRegister
	err := r.db.WithContext(c).Model(model.WorkScheduleRegister{}).
		Where("work_schedule_register_id = ? AND is_deleted = ?", id, false).
		First(&workSchedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

func (r *workScheduleRegisterRepoImpl) AssignOrUpdateEmployeesAndManagers(
	ctx context.Context,
	employees []model.WorkScheduleRegisterEmployee,
	managers []model.WorkScheduleRegisterManager,
) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Upsert cho Employees
		if len(employees) > 0 {
			if err := tx.WithContext(ctx).
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "employee_id"}},
					DoUpdates: clause.AssignmentColumns([]string{"work_schedule_register_id", "assigned_at"}),
				}).
				Create(&employees).Error; err != nil {
				return fmt.Errorf("failed to upsert employees: %w", err)
			}
		}

		// Upsert cho Managers
		if len(managers) > 0 {
			if err := tx.WithContext(ctx).
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "employee_id"}, {Name: "work_schedule_register_id"}},
					DoUpdates: clause.AssignmentColumns([]string{"is_reading", "is_editing"}),
				}).
				Create(&managers).Error; err != nil {
				return fmt.Errorf("failed to upsert managers: %w", err)
			}
		}

		return nil
	})
}

func (r *workScheduleRegisterRepoImpl) CheckExistEmployeeScheduleRegister(ctx context.Context, employeeID string, idSchedule *int) (*model.WorkScheduleRegisterEmployee, *model.WorkScheduleRegister, bool, error) {
	var req model.WorkScheduleRegisterEmployee
	var wsr model.WorkScheduleRegister
	var err error
	if idSchedule != nil {
		err = r.db.WithContext(ctx).Model(model.WorkScheduleRegisterEmployee{}).
			Preload("Employee").
			Where("employee_id = ? AND work_schedule_register_id != ?", employeeID, idSchedule).First(&req).Error
	} else {
		err = r.db.WithContext(ctx).Model(model.WorkScheduleRegisterEmployee{}).
			Preload("Employee").
			Where("employee_id = ?", employeeID).First(&req).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, nil
	}
	if err != nil {
		return nil, nil, false, err
	}
	err = r.db.WithContext(ctx).Model(model.WorkScheduleRegister{}).
		Where("work_schedule_register_id = ?", req.WorkScheduleRegisterID).First(&wsr).Error
	if err != nil {
		return nil, nil, false, err
	}
	return &req, &wsr, true, nil
}

func (r *workScheduleRegisterRepoImpl) GetAll(ctx context.Context) ([]model.WorkScheduleRegister, error) {
	var schedules []model.WorkScheduleRegister
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false").
		Preload("Office").
		Preload("Managers").
		Preload("Managers.Employee").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *workScheduleRegisterRepoImpl) GetByID(ctx context.Context, id int) (*model.WorkScheduleRegister, error) {
	var workSchedule *model.WorkScheduleRegister
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false").
		Preload("Office").
		Preload("Managers").
		Preload("Employees").
		Preload("Weekdays").
		Preload("Managers.Employee").
		Preload("Employees.Employee").
		Preload("Weekdays.WorkShift").
		Where("work_schedule_register_id = ?", id).First(&workSchedule).Error; err != nil {

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return workSchedule, nil
}

func (r *workScheduleRegisterRepoImpl) DeleteManagerFromWorkSchedule(
	ctx context.Context,
	managerID string,
	workScheduleID int) error {
	var wsEmp model.WorkScheduleRegisterManager

	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND work_schedule_register_id = ?", managerID, workScheduleID).
		First(&wsEmp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("quản lý %s không tồn tại ở lịch đăng ký'", managerID)
	}

	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Where("employee_id = ? AND work_schedule_register_id = ?", managerID, workScheduleID).
		Delete(&model.WorkScheduleRegisterManager{}).Error; err != nil {
		return fmt.Errorf("Lỗi khi xóa quản lý khỏi lịch: %w", err)
	}

	return nil
}

func (r *workScheduleRegisterRepoImpl) DeleteEmployeeFromWorkSchedule(
	ctx context.Context,
	employeeID string,
	workScheduleID int,
) error {
	var wsEmp model.WorkScheduleRegisterEmployee

	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND work_schedule_register_id = ?", employeeID, workScheduleID).
		First(&wsEmp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("Nhân viên %s không tồn tại ở lịch đăng ký", employeeID)
	}

	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Where("employee_id = ? AND work_schedule_register_id = ?", employeeID, workScheduleID).
		Delete(&model.WorkScheduleRegisterEmployee{}).Error; err != nil {
		return fmt.Errorf("Lỗi khi xóa nhân viên khỏi lịch: %w", err)
	}

	return nil
}
