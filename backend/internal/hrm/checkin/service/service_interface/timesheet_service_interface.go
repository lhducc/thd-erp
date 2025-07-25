package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimesheetServiceInterface interface {
	Create(ctx context.Context, timesheet *model.Timesheet) error
	GetByID(ctx context.Context, id string) (*model.Timesheet, error)
	Update(ctx context.Context, timesheet *model.Timesheet) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]model.Timesheet, int64, error)
}
