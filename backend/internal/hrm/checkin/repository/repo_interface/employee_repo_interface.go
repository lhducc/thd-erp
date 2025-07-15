package repo_interface

import (
	"erp/backend/internal/hrm/checkin/model"
	"time"
)

type EmployeeWorkShiftRepo interface {
	Save(assign *model.EmployeeWorkshift) error
	GetAllByEmployeeID(employeeID string) ([]model.EmployeeWorkshift, error)
	Delete(id string) error
	IsExisting(userID string, WorkShiftID string, date time.Time) bool
	GetAll() ([]model.EmployeeWorkshift, error)
	FindByID(id string) (*model.EmployeeWorkshift, error)
}
