package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"

	"gorm.io/gorm"
)

type AllowanceRepository interface {
	Create(ctx context.Context, data *model.Allowance) error
	GetAll(ctx context.Context) ([]model.Allowance, error)
	GetLastCode(ctx context.Context) (string, error)
	GetByID(ctx context.Context, id string) (*model.Allowance, error)
	Update(ctx context.Context, id string, updateData *model.AllowanceCreate) error
	Delete(ctx context.Context, id string) error
}

type allowanceRepo struct {
	db *gorm.DB
}

func NewAllowanceRepo(db *gorm.DB) AllowanceRepository {
	return &allowanceRepo{db}
}

func (r *allowanceRepo) Create(ctx context.Context, data *model.Allowance) error {
	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *allowanceRepo) GetAll(ctx context.Context) ([]model.Allowance, error) {
	var allowances []model.Allowance
	if err := r.db.WithContext(ctx).
		Where("is_deleted = false").
		Find(&allowances).Error; err != nil {
		return nil, err
	}
	return allowances, nil
}

func (r *allowanceRepo) GetLastCode(ctx context.Context) (string, error) {
	var last model.Allowance
	err := r.db.WithContext(ctx).
		Where("id LIKE ?", "PC%").
		Order("CAST(SUBSTRING(id FROM 3) AS INTEGER) DESC").
		//Order("LENGTH(id) DESC, id DESC").
		First(&last).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	return last.ID, nil
}

func (r *allowanceRepo) GetByID(ctx context.Context, id string) (*model.Allowance, error) {
	var allowance model.Allowance
	if err := r.db.WithContext(ctx).
		Where("id = ? AND is_deleted = false", id).
		First(&allowance).Error; err != nil {
		return nil, err
	}
	return &allowance, nil
}

func (r *allowanceRepo) Update(ctx context.Context, id string, updateData *model.AllowanceCreate) error {
	return r.db.WithContext(ctx).
		Model(&model.Allowance{}).
		Where("id = ? AND is_deleted = false", id).
		Updates(map[string]interface{}{
			"allowance_name": updateData.AllowanceName,
			"amount":         updateData.Amount,
			"unit":           updateData.Unit,
			"tax":            updateData.Tax,
		}).Error
}

func (r *allowanceRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&model.Allowance{}).
		Where("id = ? AND is_deleted = false", id).
		Update("is_deleted", true).Error
}
