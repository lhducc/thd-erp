package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"time"
)

type EmployeeWorkshiftService interface {
	GetByUserIdAndMonthYear(ctx context.Context, userId string, month, year int) ([]model.EmployeeWorkshift, error)
	Register(ctx context.Context, employeeID string, workshiftID string, date time.Time) error
	GetAll() ([]model.EmployeeWorkshift, error)
	Delete(id string) error
	DeletePersonalShift(id string, employeeID string) error
	Update(employee_workshift_id string, new_workshift_id string) error
	GetListShiftAllowRegister(ctx context.Context, employeeID string) ([]model.WorkScheduleShift, error)
	DeleteByManager(idEmpShift string) error
	CheckManagerPermission(ctx context.Context, managerID, employeeID string) (*model.WorkScheduleManager, error)
}
