package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
)

type WorkScheduleServiceInterface interface {
	CreateNewWorkSchedule(c context.Context, workSchedule *model.WorkSchedule) error
	DeleteWorkScheduleAuto(c context.Context, id int) error
	DeleteWorkScheduleRegister(c context.Context, id int) error
	UpdateWorkScheduleRegister(c context.Context, workSchedule *model.WorkSchedule, id int) error
	UpdateWorkScheduleAuto(c context.Context, workSchedule *model.WorkSchedule, id int) error
	AssignEmployeeToWorkScheduleAuto(c context.Context, req *dto.AssignManagersRequest, id int) error
	AssignEmployeeToWorkScheduleRegister(c context.Context, req *dto.AssignManagersRequest, id int) error
	GetAllWorkScheduleAuto(ctx context.Context) ([]model.WorkSchedule, error)
	GetAllWorkScheduleRegister(ctx context.Context) ([]model.WorkSchedule, error)
	GetWorkScheduleByID(ctx context.Context, id int) (*model.WorkSchedule, error)
	DeleteManagerFromWorkScheduleAuto(ctx context.Context, managerID string, workScheduleID int) error
	DeleteManagerFromWorkScheduleRegister(ctx context.Context, managerID string, workScheduleID int) error
	UpdateStatusRecuringSchedule(ctx context.Context, scheduleId int, isAuto bool) error
}
