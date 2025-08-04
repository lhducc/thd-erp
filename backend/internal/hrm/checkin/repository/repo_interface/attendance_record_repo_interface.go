package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"time"
)

type AttendanceRecordRepository interface {
	Create(ctx context.Context, record *model.AttendanceRecord) error
	Update(ctx context.Context, record *dto.AttendanceRecordUpdate, id string) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.AttendanceRecord, error)
	ListByEmployee(ctx context.Context, employeeID string) ([]model.AttendanceRecord, error)
	ListRequestByDateRange(ctx context.Context, employeeID string, from, to time.Time) ([]model.AttendanceRecord, error)
	CountRequestsByEmployee(ctx context.Context, employeeID string) (int64, error)
	ListHistoryRecordEmployee(ctx context.Context, employeeID string, limit int, offset int) ([]model.AttendanceRecord, error)
	ListHistoryRecordApproveByEmpID(ctx context.Context, employeeID string, startTime time.Time, endTime time.Time) ([]model.AttendanceRecord, error)
	GetByIDPersonal(ctx context.Context, recordId, employeeId string) (*model.AttendanceRecord, error)
}
