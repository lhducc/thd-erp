package repository

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/store"
	"errors"
	"gorm.io/gorm"
)

type ContractStore struct {
	db *gorm.DB
}

var _ store.ContractRepo = (*ContractStore)(nil)

func NewContractStore(db *gorm.DB) *ContractStore {
	return &ContractStore{db: db}
}

func (r *ContractStore) WithTransaction(ctx context.Context, fn func(txRepo store.ContractRepo) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &ContractStore{db: tx}
		return fn(txRepo)
	})
}

func (r *ContractStore) CreateContract(ctx context.Context, contract *hrmmodel.Contract) error {
	return r.db.WithContext(ctx).Create(contract).Error
}

func (r *ContractStore) GetContract(ctx context.Context, id string) (*hrmmodel.Contract, error) {
	var contract hrmmodel.Contract
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Employee.Department.Office").
		Preload("Allowances").
		Where("contract_id = ?", id).
		First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}
func (r *ContractStore) GetContractByEmployeeID(ctx context.Context, employeeId string) ([]hrmmodel.Contract, error) {
	var contracts []hrmmodel.Contract
	if err := r.db.WithContext(ctx).
		Preload("ContractType").
		Where("employee_id = ?", employeeId).
		Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *ContractStore) GetAllContract(ctx context.Context) ([]hrmmodel.Contract, error) {

	var contracts []hrmmodel.Contract
	if err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Employee.Department").
		Preload("Employee.Department.Office").
		Preload("Allowances").
		Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *ContractStore) UpdateContract(ctx context.Context, id string, data *hrmmodel.ContractCreate) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table("contract").
		Where("contract_id = ?", id).
		Updates(map[string]interface{}{
			"effective_date":   data.EffectiveDate,
			"expired_date":     data.ExpiredDate,
			"sign_date":        data.SignDate,
			"note":             data.Note,
			"attached_file":    data.AttachedFile,
			"condition":        data.Condition,
			"created_date":     data.CreatedDate,
			"contract_type_id": data.ContractTypeId,
			"approve_status":   data.ApproveStatus,
			"employee_id":      data.Manager,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var contract hrmmodel.Contract
	if err := tx.First(&contract, "contract_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(data.AllowanceIDs) > 0 {
		var allowances []*hrmmodel.Allowance
		if err := tx.Where("id IN ?", data.AllowanceIDs).Find(&allowances).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&contract).Association("Allowances").Replace(allowances); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *ContractStore) DeleteContract(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("contract").
		Where("contract_id = ?", id).
		Delete(nil).Error
}

func (s *ContractStore) CheckExistName(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&hrmmodel.Contract{}).Where("contract_name = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Contract Name is already")
	}
	return false, nil
}

func (s *ContractStore) GetLastContractByCode(ctx context.Context, Contract *hrmmodel.Contract) error {
	return s.db.WithContext(ctx).
		Order("contract_id DESC").
		First(Contract).Error
}

func (s *ContractStore) CreateContractAllowance(ctx context.Context, contract *hrmmodel.ContractAllowance) error {
	return s.db.WithContext(ctx).Create(&contract).Error
}

func (r *ContractStore) UpdateApproveStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).
		Table("contract").
		Where("contract_id = ?", id).
		Update("approve_status", status).Error
}
