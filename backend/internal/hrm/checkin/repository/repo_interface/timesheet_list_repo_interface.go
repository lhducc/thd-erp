package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimesheetListInterface interface {
	Create(ctx context.Context, timesheet *model.TimeSheetList) error
	GetByID(ctx context.Context, id string) (*model.TimeSheetList, error)
	Update(ctx context.Context, timesheet *model.TimeSheetList) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]model.TimeSheetList, int64, error)
	GetLastDecisionByCode(ctx context.Context) (string, error)
	IsDuplicate(ctx context.Context, officeID string, month, year int, timesheetID string) (bool, error)
	UpdateLocked(ctx context.Context, timesheet *model.TimeSheetList) error
	GetTimeSheetByOfficeIDAndTime(officeID string, month, year int) (*model.TimeSheetList, error)
}
