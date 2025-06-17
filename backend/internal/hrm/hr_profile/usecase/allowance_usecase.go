package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	utils "erp/backend/pkg"
	"fmt"
	"time"
)

type AllowanceUsecase interface {
	Create(ctx context.Context, data *model.AllowanceCreate) error
	GetAll(ctx context.Context) ([]model.Allowance, error)
}

type allowanceBiz struct {
	repo repository.AllowanceRepository
}

func NewAllowanceBiz(repo repository.AllowanceRepository) AllowanceUsecase {
	return &allowanceBiz{repo: repo}
}

func (b *allowanceBiz) Create(ctx context.Context, data *model.AllowanceCreate) error {
	code, err := utils.GenerateCodeAllowance("PC", 3, func() (string, error) {
		return b.repo.GetLastCode(ctx)
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	allowance := &model.Allowance{
		ID:            code,
		AllowanceName: data.AllowanceName,
		Tax:           data.Tax,
		Amount:        data.Amount,
		Unit:          data.Unit,
		IsDeleted:     false,
		CreatedDate:   time.Now(),
	}
	return b.repo.Create(ctx, allowance)
}

func (b *allowanceBiz) GetAll(ctx context.Context) ([]model.Allowance, error) {
	return b.repo.GetAll(ctx)
}
