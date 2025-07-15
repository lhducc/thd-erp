package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type WorkScheduleRepo interface {
	Save(c context.Context, workSchedule *model.WorkSchedule, weekdayShift []model.WorkScheduleShift) error
	Delete(c context.Context, id int) error
	Update(c context.Context, workSchedule *model.WorkSchedule, id int) error
	IsExistsByName(c context.Context, name string) (bool, error)
	IsExistsByID(c context.Context, id int) (bool, error)
	AssignEmployeesToWorkSchedule(ctx context.Context, employeeRecords []model.WorkScheduleEmployee, managerRecords []model.WorkScheduleManager) error
	CheckExistEmployeeSchedule(ctx context.Context, employeeID string) (bool, error)
	GetAll(ctx context.Context) ([]model.WorkSchedule, error)
	GetByID(ctx context.Context, id int) (*model.WorkSchedule, error)
}
