package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
)

type WorkShiftService interface {
	CreateWorkShift(ctx context.Context, data *model.WorkShifts) error
	GetWorkShiftById(ctx context.Context, id string) (*model.WorkShifts, error)
	GetAllWorkShift(ctx context.Context) ([]model.WorkShifts, error)
	UpdateWorkShift(ctx context.Context, data *model.WorkShifts) error
	DeleteWorkShift(ctx context.Context, id string) error
	GetListShiftForRegister(ctx context.Context, employeeID string) ([]model.WorkScheduleShift, error)
	GetWorkshiftInfo(ctx context.Context, employeeID string) (*dto.WorkScheduleInfo, error)
}
