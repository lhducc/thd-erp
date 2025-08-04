package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"time"
)

type EmployeeWorkShiftRepo interface {
	Save(assign *model.EmployeeWorkshift) error
	GetAllByEmployeeID(employeeID string) ([]model.EmployeeWorkshift, error)
	Delete(id string) error
	DeletePersonalShift(id string, employeeID string) error
	IsExisting(userID string, WorkShiftID string, date time.Time) bool
	GetAll() ([]model.EmployeeWorkshift, error)
	FindByID(id string) (*model.EmployeeWorkshift, error)
	GetEmployeeWorkShifts(ctx context.Context, employeeID string) ([]model.EmployeeWorkshift, error)
	GetEmployeeWorkShiftsByMonthYear(ctx context.Context, employeeID string, startDate, endDate time.Time) ([]model.EmployeeWorkshift, error)
	GetByID(id string) (*model.EmployeeWorkshift, error)
	AssignmentShift(ctx context.Context, assigns []model.EmployeeWorkshift, scheduleID []int) error
	CheckShiftConflict(ctx context.Context, employeeID string, workshiftID string, date time.Time) (bool, error)
}
