package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"time"
)

type AttendanceRecordRepository interface {
	Create(ctx context.Context, record *model.AttendanceRecord) error
	Update(ctx context.Context, record *model.AttendanceRecordUpdate, id string) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.AttendanceRecord, error)
	ListByEmployee(ctx context.Context, employeeID string) ([]model.AttendanceRecord, error)
	ListByDateRange(ctx context.Context, employeeID string, from, to time.Time) ([]model.AttendanceRecord, error)
	CountRequestsByEmployee(ctx context.Context, employeeID string) (int64, error)
}
