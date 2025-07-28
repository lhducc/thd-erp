package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type TimeSheetDetailRepoInterface interface {
	Create(ctx context.Context, timesheet *model.TimeSheetDetail) error
}
