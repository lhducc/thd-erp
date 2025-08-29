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
	// IsDuplicate checks whether a timesheet list exists for the given month/year.
	// When timesheets are global (not per office) we only check month/year.
	IsDuplicate(ctx context.Context, month, year int, timesheetID string) (bool, error)
	UpdateLocked(ctx context.Context, timesheet *model.TimeSheetList) error
	// GetTimeSheetByTime returns the timesheet list for the given month/year (global, not per office).
	GetTimeSheetByTime(ctx context.Context, month, year int) (*model.TimeSheetList, error)
	GetForExport(ctx context.Context, id string) (*model.TimeSheetList, error)
}
