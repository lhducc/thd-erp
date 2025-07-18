package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type WorkScheduleRegisterRepo interface {
	Save(c context.Context, workSchedule *model.WorkScheduleRegister, weekdayShift []model.WorkScheduleRegisterShift) error
	Delete(c context.Context, id int) error
	Update(c context.Context, wschedule *model.WorkScheduleRegister, id int) error
	IsExistsByName(c context.Context, name string) (bool, error)
	IsExistsByID(c context.Context, id int) (bool, error)
	AssignOrUpdateEmployeesAndManagers(ctx context.Context, employees []model.WorkScheduleRegisterEmployee, managers []model.WorkScheduleRegisterManager) error
	CheckExistEmployeeScheduleRegister(ctx context.Context, employeeID string, idSchedule *int) (*model.WorkScheduleRegisterEmployee, *model.WorkScheduleRegister, bool, error)
	GetAll(ctx context.Context) ([]model.WorkScheduleRegister, error)
	GetByID(ctx context.Context, id int) (*model.WorkScheduleRegister, error)
	DeleteManagerFromWorkSchedule(ctx context.Context, managerID string, workScheduleID int) error
	DeleteEmployeeFromWorkSchedule(ctx context.Context, employeeID string, workScheduleID int) error
}
