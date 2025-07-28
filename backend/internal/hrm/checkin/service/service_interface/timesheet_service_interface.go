package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimeSheetServiceInterface interface {
	FindByEmployeeIDAndMonth(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error)
}
