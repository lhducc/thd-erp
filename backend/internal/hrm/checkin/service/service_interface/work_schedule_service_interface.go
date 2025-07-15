package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type WorkScheduleServiceInterface interface {
	CreateNewWorkSchedule(c context.Context, workSchedule *model.WorkSchedule) error
	DeleteWorkSchedule(c context.Context, id int) error
	UpdateWorkSchedule(c context.Context, workSchedule *model.WorkSchedule, id int) error
	AssignEmployeeToWorkSchedule(c context.Context, req *model.AssignEmployeeRequest, id int) error
	GetAllWorkSchedule(ctx context.Context) ([]model.WorkSchedule, error)
	GetWorkScheduleByID(ctx context.Context, id int) (*model.WorkSchedule, error)
	ExportWorkSchedule(c context.Context, selectedFields []string) ([]byte, string, error)
}
