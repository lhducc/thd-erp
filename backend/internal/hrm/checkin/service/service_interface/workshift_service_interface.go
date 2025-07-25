package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type WorkShiftService interface {
	CreateWorkShift(ctx context.Context, data *model.WorkShifts) error
	GetWorkShiftById(ctx context.Context, id string) (*model.WorkShifts, error)
	GetAllWorkShift(ctx context.Context) ([]model.WorkShifts, error)
	UpdateWorkShift(ctx context.Context, data *model.WorkShifts) error
	DeleteWorkShift(ctx context.Context, id string) error
}
