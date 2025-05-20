package repository

import (
	"context"
	ContractModel "erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

type ContractTypeRepository interface {
	GetContractType(ctx context.Context) ([]ContractModel.ContractType, error)
	CreateContractType(ctx context.Context, contractType *ContractModel.ContractType) (*ContractModel.ContractType, error)
	UpdateContractType(ctx context.Context, contractTypeId string, contractType *ContractModel.ContractType) error
	DeleteContractType(ctx context.Context, contractTypeId string, contractType *ContractModel.ContractType) error
	GetContractTypeById(ctx context.Context, id string) (*ContractModel.ContractType, error)
	GetLastContractTypeByCode(ctx context.Context, lastContractType *ContractModel.ContractType) error
}

func NewContractTypeRepository(db *gorm.DB) ContractTypeRepository {
	return &contractTypeRepository{db: db}
}

type contractTypeRepository struct {
	db *gorm.DB
}

func (c *contractTypeRepository) GetContractType(ctx context.Context) ([]ContractModel.ContractType, error) {
	var contractTypes []ContractModel.ContractType
	if err := c.db.WithContext(ctx).Find(&contractTypes).Error; err != nil {
		return nil, err
	}
	return contractTypes, nil
}

func (c *contractTypeRepository) CreateContractType(ctx context.Context, contractType *ContractModel.ContractType) (*ContractModel.ContractType, error) {
	if err := c.db.WithContext(ctx).Create(&contractType).Error; err != nil {
		return nil, err
	}
	return contractType, nil
}

func (c *contractTypeRepository) UpdateContractType(ctx context.Context, contractTypeId string, contractType *ContractModel.ContractType) error {
	return c.db.WithContext(ctx).
		Table("contracttype").
		Where("contract_type_id = ?", contractTypeId).
		Updates(contractType).Error
}

func (c *contractTypeRepository) DeleteContractType(ctx context.Context, contractTypeId string, contractType *ContractModel.ContractType) error {
	return c.db.WithContext(ctx).
		Table("contracttype").
		Where("contract_type_id = ?", contractTypeId).
		Updates(contractType).Error
}

func (c *contractTypeRepository) GetContractTypeById(ctx context.Context, id string) (*ContractModel.ContractType, error) {
	var contractType ContractModel.ContractType
	if err := c.db.WithContext(ctx).
		Where("contract_type_id = ?", id).
		First(&contractType).Error; err != nil {
		return nil, err
	}
	return &contractType, nil
}

func (c *contractTypeRepository) GetLastContractTypeByCode(ctx context.Context, lastContractType *ContractModel.ContractType) error {
	return c.db.WithContext(ctx).
		Order("contract_type_id DESC").
		Limit(1).
		Find(lastContractType).Error
}
