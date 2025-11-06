package implement

import (
	"context"
	"fmt"

	recmodel "erp/backend/internal/hrm/recruitment/model"
	repository "erp/backend/internal/hrm/recruitment/repository/interface"
	service "erp/backend/internal/hrm/recruitment/service/interface"
)

type processStageService struct {
	repo repository.ProcessStageRepository
}

func NewProcessStageServiceImpl(repo repository.ProcessStageRepository) service.ProcessStageService {
	return &processStageService{repo: repo}
}

func (s *processStageService) CreateStage(ctx context.Context, processFormID string, stageName string) (*recmodel.ProcessStage, error) {
	if stageName == "" {
		return nil, fmt.Errorf("stage name required")
	}

	exists, err := s.repo.ExistsProcessForm(ctx, processFormID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("process form not found")
	}
	lastOrder, err := s.repo.GetLastOrder(ctx, processFormID)
	if err != nil {
		return nil, err
	}

	nameExists, err := s.repo.ExistsStageName(ctx, processFormID, stageName)
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, fmt.Errorf("stage name already exists")
	}

	newOrder := lastOrder + 1

	ps := &recmodel.ProcessStage{
		ProcessFormID: processFormID,
		StageName:     stageName,
		StageOrder:    newOrder,
	}

	if err := s.repo.Save(ctx, ps); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *processStageService) UpdateStage(ctx context.Context, processStageID int, newName string) (*recmodel.ProcessStage, error) {
	if newName == "" {
		return nil, fmt.Errorf("stage name required")
	}

	ps, err := s.repo.FindByID(ctx, processStageID)
	if err != nil {
		return nil, err
	}

	if ps.StageName == newName {
		return ps, nil
	}

	exists, err := s.repo.ExistsStageName(ctx, ps.ProcessFormID, newName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("stage name already exists")
	}

	if err := s.repo.UpdateStageName(ctx, processStageID, newName); err != nil {
		return nil, err
	}

	ps.StageName = newName
	return ps, nil
}

func (s *processStageService) DeleteStage(ctx context.Context, processStageID int) error {
	ps, err := s.repo.FindByID(ctx, processStageID)
	if err != nil {
		return err
	}

	deletedOrder := ps.StageOrder

	if err := s.repo.DeleteByID(ctx, processStageID); err != nil {
		return err
	}

	if err := s.repo.DecrementOrdersAfter(ctx, ps.ProcessFormID, deletedOrder); err != nil {
		return err
	}

	return nil
}
