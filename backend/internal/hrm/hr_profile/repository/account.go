package repository

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type accountStore struct {
	db *gorm.DB
}

func NewAccountStore(db *gorm.DB) *accountStore {
	return &accountStore{db: db}
}

func (s *accountStore) CreateAccount(employee *hrmmodel.Employee, hashedPassword, roleID string) (*hrmmodel.Account, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// log.Printf("Panic occurred, transaction rolled back: %v", r)
		}
	}()

	var role hrmmodel.Role
	if err := tx.First(&role, "id = ?", roleID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get role %s: %w", roleID, err)
	}

	account := hrmmodel.Account{
		LoginMail:   employee.Email,
		Password:    hashedPassword,
		FirstLogin:  true,
		CreatedDate: time.Now(),
		EmployeeId:  employee.EmployeeID,
		RoleID:      role.ID,
	}

	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &account, nil
}

func (r *accountStore) GetAccount(ctx context.Context, id int64) (*hrmmodel.Account, error) {
	var contract hrmmodel.Account
	if err := r.db.WithContext(ctx).Table("account").
		Where("id = ?", id).
		First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *accountStore) GetAllAccount(ctx context.Context) ([]hrmmodel.Account, error) {

	var positions []hrmmodel.Account
	if err := r.db.WithContext(ctx).Table("account").Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *accountStore) UpdateAccount(ctx context.Context, id string, data *hrmmodel.Account) error {
	return r.db.WithContext(ctx).Table("account").
		Where("id = ?", id).
		Updates(data).Error
}

func (r *accountStore) DeleteAccount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("account").
		Where("id = ?", id).
		Delete(nil).Error
}

func (s *accountStore) CheckExistEmail(email string) (bool, error) {
	var count int64
	if err := s.db.Model(&hrmmodel.Account{}).Where("login_mail = ?", email).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Account is already")
	}
	return false, nil
}

func (s *accountStore) CheckFirstLogin(email string) (bool, error) {
	return false, nil
}
