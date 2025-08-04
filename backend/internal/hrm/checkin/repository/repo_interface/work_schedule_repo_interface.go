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
	IsExistsByScheduleIDRegister(c context.Context, id int) (bool, error)
	IsExistsByScheduleIDAuto(c context.Context, id int) (bool, error)
	AssignOrUpdateManagers(ctx context.Context, managers []model.WorkScheduleManager) error
	GetAllScheduleAuto(ctx context.Context) ([]model.WorkSchedule, error)
	GetAllScheduleRegister(ctx context.Context) ([]model.WorkSchedule, error)
	GetByID(ctx context.Context, id int) (*model.WorkSchedule, error)
	DeleteManagerFromWorkSchedule(ctx context.Context, managerID string, workScheduleID int) error
	GetListShiftRegister(ctx context.Context, scheduleID *int) ([]model.WorkScheduleShift, error)
	CheckManagerPermission(ctx context.Context, managerID, employeeID string) (*model.WorkScheduleManager, error)
	UpdateStatusRecuringSchedule(ctx context.Context, scheduleID int, isAuto bool) error
}
