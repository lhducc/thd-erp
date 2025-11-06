package repository

import (
	"context"
	recmodel "erp/backend/internal/hrm/recruitment/model"
)

type ProcessStageRepository interface {
    Save(ctx context.Context, stage *recmodel.ProcessStage) error
    GetLastOrder(ctx context.Context, processFormID string) (int, error)
	ExistsProcessForm(ctx context.Context, processFormID string) (bool, error)
	ExistsStageName(ctx context.Context, processFormID string, stageName string) (bool, error)
	FindByID(ctx context.Context, processStageID int) (*recmodel.ProcessStage, error)
	UpdateStageName(ctx context.Context, processStageID int, newName string) error
	DeleteByID(ctx context.Context, processStageID int) error	
	DecrementOrdersAfter(ctx context.Context, processFormID string, fromOrder int) error
}


