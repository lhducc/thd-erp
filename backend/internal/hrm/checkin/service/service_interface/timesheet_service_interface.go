package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimeSheetServiceInterface interface {
	FindByEmployeeIDAndMonth(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error)
	CalculatorTimeSheetList(ctx context.Context, timesheetListID string) error
	ResetWorkDayAdjustment(ctx context.Context, timesheetDetailID int) error
	ManualAdjustWorkDay(ctx context.Context, timesheetDetailID int, adjustedWorkDays float64, adjustedBy string) error
	ExportTimeSheet(ctx context.Context, timesheetListID string) ([]byte, string, error)
}
