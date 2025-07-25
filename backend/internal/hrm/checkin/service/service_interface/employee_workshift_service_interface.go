package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"time"
)

type EmployeeWorkshiftService interface {
	GetByUserId(userId string) ([]model.EmployeeWorkshift, error)
	Register(ctx context.Context, employeeID string, workshiftID string, date time.Time) error
	GetAll() ([]model.EmployeeWorkshift, error)
	Delete(id string) error
	Update(employee_workshift_id string, new_workshift_id string) error
	GetListShiftAllowRegister(ctx context.Context, employeeID string) ([]model.WorkScheduleShift, error)
}
