package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/pkg/timeonly"
)

type WorkShiftRepo interface {
	GetWorkShiftById(ctx context.Context, id string) (*model.WorkShifts, error)
	GetAllWorkShift(ctx context.Context) ([]model.WorkShifts, error)
	DeleteWorkShift(ctx context.Context, id string) error
	GetLastWorkShiftByCode(ctx context.Context, emp *model.WorkShifts, predix string) error
	SaveWorkShift(ctx context.Context, workshift *model.WorkShifts) error
	CheckExistByName(ctx context.Context, nameWS string, excludeID string) (bool, error)
	IsDuplicateTimeRange(ctx context.Context, startTime timeonly.TimeOnly, endTime timeonly.TimeOnly, excludeID string) (bool, error)
}
