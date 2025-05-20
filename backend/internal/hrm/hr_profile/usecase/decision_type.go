package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"erp/backend/pkg"
)

type DecisionTypeUsecase interface {
	CreateDecisionType(ctx context.Context, dt *model.DecisionType) (*model.DecisionType, error)
	UpdateDecisionType(ctx context.Context, id string, dt *model.DecisionType) (*model.DecisionType, error)
	DeleteDecisionType(ctx context.Context, id string) error
	GetAllDecisionTypes(ctx context.Context) ([]model.DecisionType, error)
	GetDecisionTypeByID(ctx context.Context, id string) (*model.DecisionType, error)
}

type decisionTypeUsecase struct {
	repo repository.DecisionTypeRepository
}

func NewDecisionTypeUsecase(repo repository.DecisionTypeRepository) *decisionTypeUsecase {
	return &decisionTypeUsecase{repo: repo}
}

func (u *decisionTypeUsecase) CreateDecisionType(ctx context.Context, dt *model.DecisionType) (*model.DecisionType, error) {
	code, err := u.generateDecisionTypeCode(ctx)
	if err != nil {
		return nil, err
	}
	dt.DecisionTypeID = code
	dt.CreatedDate = utils.GetCurrentDate()

	if err := dt.ValidateDecisionType(); err != nil {
		return nil, err
	}

	return u.repo.CreateDecisionType(ctx, dt)
}

func (u *decisionTypeUsecase) UpdateDecisionType(ctx context.Context, id string, dt *model.DecisionType) (*model.DecisionType, error) {
	existing, err := u.repo.GetDecisionTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("decision type not found")
	}
	return u.repo.UpdateDecisionType(ctx, id, dt)
}

func (u *decisionTypeUsecase) DeleteDecisionType(ctx context.Context, id string) error {
	existing, err := u.repo.GetDecisionTypeByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("decision type not found")
	}
	return u.repo.DeleteDecisionType(ctx, id)
}

func (u *decisionTypeUsecase) GetAllDecisionTypes(ctx context.Context) ([]model.DecisionType, error) {
	return u.repo.GetAllDecisionTypes(ctx)
}

func (u *decisionTypeUsecase) GetDecisionTypeByID(ctx context.Context, id string) (*model.DecisionType, error) {
	return u.repo.GetDecisionTypeByID(ctx, id)
}

func (u *decisionTypeUsecase) generateDecisionTypeCode(ctx context.Context) (string, error) {
	var last model.DecisionType
	err := u.repo.GetLastDecisionTypeByCode(ctx, &last)
	if err != nil || !strings.HasPrefix(last.DecisionTypeID, "TL") {
		return "TL00001", nil
	}

	numStr := strings.TrimSpace(strings.TrimPrefix(last.DecisionTypeID, "TL"))
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse decision code number: %w", err)
	}

	if next := num + 1; next <= 99999 {
		return fmt.Sprintf("TL%05d", next), nil
	}

	return "", fmt.Errorf("maximum office code reached: TL99999")
}
