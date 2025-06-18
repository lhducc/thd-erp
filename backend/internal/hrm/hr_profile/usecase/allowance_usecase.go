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
	GetByID(ctx context.Context, id string) (*model.Allowance, error)
	Update(ctx context.Context, id string, data *model.AllowanceCreate) error
	Delete(ctx context.Context, id string) error
}
type allowanceBiz struct {
	repo repository.AllowanceRepository
}

func NewAllowanceBiz(repo repository.AllowanceRepository) AllowanceUsecase {
	return &allowanceBiz{repo: repo}
}

func (b *allowanceBiz) Create(ctx context.Context, data *model.AllowanceCreate) error {
	// Validate
	if data.AllowanceName == "" {
		return fmt.Errorf("Tên phụ cấp không được để trống")
	}
	if data.Amount < 0 {
		return fmt.Errorf("Số tiền phụ cấp không được âm")
	}
	if !model.IsValidUnit(model.AllowanceUnit(data.Unit)) {
		return fmt.Errorf("Đơn vị tiền tệ không hợp lệ: %v", data.Unit)
	}

	code, err := utils.GenerateCodeAllowance("PC", 4, func() (string, error) {
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
		Unit:          model.AllowanceUnit(data.Unit),
		IsDeleted:     false,
		CreatedDate:   time.Now(),
	}
	return b.repo.Create(ctx, allowance)
}

func (b *allowanceBiz) GetAll(ctx context.Context) ([]model.Allowance, error) {
	return b.repo.GetAll(ctx)
}

func (b *allowanceBiz) GetByID(ctx context.Context, id string) (*model.Allowance, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *allowanceBiz) Update(ctx context.Context, id string, data *model.AllowanceCreate) error {
	// Validate
	if data.AllowanceName == "" {
		return fmt.Errorf("Tên phụ cấp không được để trống")
	}
	if data.Amount < 0 {
		return fmt.Errorf("Số tiền phụ cấp không được âm")
	}
	if !model.IsValidUnit(model.AllowanceUnit(data.Unit)) {
		return fmt.Errorf("Đơn vị tiền tệ không hợp lệ: %v", data.Unit)
	}

	return b.repo.Update(ctx, id, data)
}

func (b *allowanceBiz) Delete(ctx context.Context, id string) error {
	return b.repo.Delete(ctx, id)
}
