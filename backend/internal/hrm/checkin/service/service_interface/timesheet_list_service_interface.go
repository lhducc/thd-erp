package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimesheetListServiceInterface interface {
	CreateElementOfTimesheetList(ctx context.Context, timesheet *model.TimeSheetList) error
	GetByID(ctx context.Context, id string) (*model.TimeSheetList, error)
	Update(ctx context.Context, timesheet *model.TimeSheetList) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]model.TimeSheetList, int64, error)
	LockedTimesheet(ctx context.Context, timesheet *model.TimeSheetList) error
	ExportCheckinCheckout(ctx context.Context, id string) ([]byte, string, error)
}
