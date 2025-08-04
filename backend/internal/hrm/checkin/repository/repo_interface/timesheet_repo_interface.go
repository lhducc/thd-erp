package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimeSheetRepoInterface interface {
	CreateTimeSheets(ctx context.Context, timesheets []*model.TimeSheet) error
	FindByEmployeeAndMonth(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error)
	UpdateTimesheetAndCreateDetail(ctx context.Context, timesheets *model.TimeSheet) error
	GetByID(ctx context.Context, timesheetID int) (*model.TimeSheet, error)
	Update(ctx context.Context, timesheet *model.TimeSheet) error
	CreateEmployeeTimeSheet(timesheet *model.TimeSheet) error
}
