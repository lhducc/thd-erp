package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type WorkScheduleRegisterServiceInterface interface {
	CreateNewWorkSchedule(c context.Context, workSchedule *model.WorkScheduleRegister) error
	DeleteWorkSchedule(c context.Context, id int) error
	UpdateWorkSchedule(c context.Context, workSchedule *model.WorkScheduleRegister, id int) error
	AssignEmployeeToWorkSchedule(c context.Context, req *model.AssignEmployeeRequest, id int) error
	GetAllWorkSchedule(ctx context.Context) ([]model.WorkScheduleRegister, error)
	GetWorkScheduleByID(ctx context.Context, id int) (*model.WorkScheduleRegister, error)
	DeleteEmployeeFromWorkScheduleRegister(ctx context.Context, employeeID string, workScheduleID int) error
	DeleteManagerFromWorkScheduleRegister(ctx context.Context, managerID string, workScheduleID int) error
}
