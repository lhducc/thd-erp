package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimeSheetDetailRepoInterface interface {
	Create(ctx context.Context, timesheet *model.TimeSheetDetail) error
	GetDetailsByTimeSheetID(ctx context.Context, timeSheetID int) ([]model.TimeSheetDetail, error)
	GetDetailByID(ctx context.Context, detailID int) (*model.TimeSheetDetail, error)
	UpdateTimeSheetDetail(ctx context.Context, detail *model.TimeSheetDetail) error
}
