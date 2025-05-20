package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

type InsuranceRepository interface {
	CreateInsurance(ctx context.Context, insurance *model.Insurance) (*model.Insurance, error)
	UpdateInsurance(ctx context.Context, id string, insurance *model.Insurance) (*model.Insurance, error)
	DeleteInsurance(ctx context.Context, id string) error
	GetInsuranceByID(ctx context.Context, id string) (*model.Insurance, error)
	GetAllInsurance(ctx context.Context) ([]model.Insurance, error)
	GetInsuranceLastByCode(ctx context.Context, insurance *model.Insurance) error
}
type insuranceRepository struct {
	db *gorm.DB
}

func NewInsuranceRepository(db *gorm.DB) InsuranceRepository {
	return &insuranceRepository{db: db}
}

func (i *insuranceRepository) CreateInsurance(ctx context.Context, insurance *model.Insurance) (*model.Insurance, error) {
	if err := i.db.WithContext(ctx).Create(insurance).Error; err != nil {
		return nil, err
	}
	return insurance, nil
}

func (i *insuranceRepository) UpdateInsurance(ctx context.Context, id string, insurance *model.Insurance) (*model.Insurance, error) {
	if err := i.db.WithContext(ctx).Model(&model.Insurance{}).Where("id = ?", id).Updates(insurance).Error; err != nil {
		return nil, err
	}
	return insurance, nil
}

func (i *insuranceRepository) DeleteInsurance(ctx context.Context, id string) error {
	if err := i.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Insurance{}).Error; err != nil {
		return err
	}
	return nil
}

func (i *insuranceRepository) GetInsuranceByID(ctx context.Context, id string) (*model.Insurance, error) {
	var insurance model.Insurance
	if err := i.db.WithContext(ctx).
		Where("id = ?", id).
		First(&insurance).Error; err != nil {
		return nil, err
	}
	return &insurance, nil
}

func (i *insuranceRepository) GetAllInsurance(ctx context.Context) ([]model.Insurance, error) {
	var insurances []model.Insurance
	if err := i.db.WithContext(ctx).
		Find(&insurances).Error; err != nil {
		return nil, err
	}
	return insurances, nil
}

func (i *insuranceRepository) GetInsuranceLastByCode(ctx context.Context, insurance *model.Insurance) error {
	return i.db.WithContext(ctx).
		Order("id DESC").
		First(insurance).Error
}
