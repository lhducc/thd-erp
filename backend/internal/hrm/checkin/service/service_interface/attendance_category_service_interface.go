package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

// AttendanceCategoryService defines the interface for attendance category operations
type AttendanceCategoryService interface {
	CreateAttendanceCategory(ctx context.Context, category *model.AttendanceCategory) error
	UpdateAttendanceCategory(ctx context.Context, category *model.AttendanceCategory) error
	DeleteAttendanceCategory(ctx context.Context, id string) error
	GetAttendanceCategoryByID(ctx context.Context, id string) (*model.AttendanceCategory, error)
	ListAttendanceCategoriesByOffice(ctx context.Context, employeeID string) ([]model.AttendanceCategory, error)
	ListAllAttendanceCategories(ctx context.Context) ([]model.AttendanceCategory, error)
}
