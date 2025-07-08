package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type AttendanceRecordService interface {
	CreateAttendanceRecord(ctx context.Context, record *model.AttendanceRecord) error
	UpdateAttendanceRecord(ctx context.Context, record *model.AttendanceRecord) error
	DeleteAttendanceRecord(ctx context.Context, id string) error
	GetAttendanceRecordByID(ctx context.Context, id string) (*model.AttendanceRecord, error)
	ListAttendanceRecordsByEmployee(ctx context.Context, employeeID string) ([]model.AttendanceRecord, error)
	ListAttendanceRecordsByDateRange(ctx context.Context, employeeID string, from, to string) ([]model.AttendanceRecord, error)
	GetTotalReqOfEmployee(ctx context.Context, employeeID string) (int64, error)
	ValidateAttendanceRecordDistance(ctx context.Context, record *model.AttendanceRecord, category *model.AttendanceCategory) error
	CheckCatrgoryExists(ctx context.Context, record *model.AttendanceRecord) (*model.AttendanceCategory, error)
}
