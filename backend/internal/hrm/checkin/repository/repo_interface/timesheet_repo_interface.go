package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimeSheetRepoInterface interface {
	Create(ctx context.Context, timesheet []*model.TimeSheet) error
	FindByEmployeeAndMonth(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error)
}
