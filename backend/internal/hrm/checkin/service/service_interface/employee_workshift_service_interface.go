package service_interface

import (
	"erp/backend/internal/hrm/checkin/model"
	"time"
)

type EmployeeWorkshiftService interface {
	GetByUserId(userId string) ([]model.EmployeeWorkshift, error)
	Register(employeeID string, workshiftID string, date time.Time) error
	GetAll() ([]model.EmployeeWorkshift, error)
	Delete(id string) error
	Update(employee_workshift_id string, new_workshift_id string) error
}
