package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model/dto"
)

type ShiftAllocationServiceInterface interface {
	GetListShiftAllocation(ctx context.Context, req dto.GetShiftAllocationRequest, managerID string) ([]dto.EmployeeScheduleResponse, error)
}
