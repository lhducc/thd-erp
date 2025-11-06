package service

import (
	"context"
	recmodel "erp/backend/internal/hrm/recruitment/model"
)

type ProcessStageService interface {
	CreateStage(ctx context.Context, processFormID string, stageName string) (*recmodel.ProcessStage, error)
	UpdateStage(ctx context.Context, processStageID int, newName string) (*recmodel.ProcessStage, error)
	DeleteStage(ctx context.Context, processStageID int) error
}
