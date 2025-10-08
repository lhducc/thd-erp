package service_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"time"
)

type AttendanceRecordService interface {
	CreateAttendanceRecord(ctx context.Context, record *model.AttendanceRecord) error
	UpdateAttendanceRecord(ctx context.Context, record *dto.AttendanceRecordUpdate, recordID string) error
	DeleteAttendanceRecord(ctx context.Context, id string) error
	GetAttendanceRecordByID(ctx context.Context, id string) (*model.AttendanceRecord, error)
	ListAttendanceRecordsByEmployee(ctx context.Context, employeeID string) ([]model.AttendanceRecord, error)
	ListAttendanceRequestsByDateRange(ctx context.Context, employeeID string, from, to string) ([]model.AttendanceRecord, error)
	GetTotalReqOfEmployee(ctx context.Context, employeeID string) (int64, error)
	ValidateAttendanceRecordDistance(ctx context.Context, record *model.AttendanceRecord, category *model.AttendanceCategory) error
	CheckCategoryExists(ctx context.Context, record *model.AttendanceRecord) (*model.AttendanceCategory, error)
	GetHistoryRecordByEmployee(ctx context.Context, employeeID string, page int, limit int) ([]model.AttendanceRecord, error)
	GetAttendanceRecordByIDPersonal(ctx context.Context, recordID, employeeID string) (*model.AttendanceRecord, error)
	GetHistoryByDate(ctx context.Context, dateStr string) ([]dto.AttendanceRecordHistoryByDate, error)
	ExportAttendanceExcel(ctx context.Context, targetDate time.Time) ([]byte, error)
	GetListHistoryByDateForManager(ctx context.Context, managerID string, dateStr string) ([]dto.AttendanceRecordHistoryByDate, error)
}
