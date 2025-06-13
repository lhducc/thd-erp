package repository

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"errors"

	"gorm.io/gorm"
)

type contractStore struct {
	db *gorm.DB
}

func NewContractStore(db *gorm.DB) *contractStore {
	return &contractStore{db: db}
}

func (s *contractStore) CreateContract(context context.Context, data *hrmmodel.Contract) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *contractStore) GetContract(ctx context.Context, id string) (*hrmmodel.Contract, error) {
	var contract hrmmodel.Contract
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Employee.Department.Office").
		Where("contract_id = ?", id).
		First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *contractStore) GetAllContract(ctx context.Context) ([]hrmmodel.Contract, error) {

	var contracts []hrmmodel.Contract
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Employee.Department.Office").
		Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *contractStore) UpdateContract(ctx context.Context, id string, data *hrmmodel.ContractCreate) error {
	return r.db.WithContext(ctx).Table("contract").
		Where("contract_id = ?", id).
		Updates(data).Error
}

func (r *contractStore) DeleteContract(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("contract").
		Where("contract_id = ?", id).
		Delete(nil).Error
}

func (s *contractStore) CheckExistName(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&hrmmodel.Contract{}).Where("contract_name = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Contract Name is already")
	}
	return false, nil
}

func (s *contractStore) GetLastContractByCode(ctx context.Context, Contract *hrmmodel.Contract) error {
	return s.db.WithContext(ctx).
		Order("contract_id DESC").
		First(Contract).Error
}
