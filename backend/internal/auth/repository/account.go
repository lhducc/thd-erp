package authRepository

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

// AccountRepository định nghĩa interface cho repo
type AccountRepository interface {
	GetAccountByEmail(ctx context.Context, email string) (hrmmodel.Account, error)
	GetAccountByID(ctx context.Context, id string) (hrmmodel.Account, error)
	UpdateAccount(ctx context.Context, account hrmmodel.Account) error
}

// accountRepo struct chứa db connection và implement interface
type accountRepo struct {
	db *gorm.DB
}

// NewAccountRepository trả về interface AccountRepository
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) GetAccountByEmail(ctx context.Context, email string) (hrmmodel.Account, error) {
	var account hrmmodel.Account
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("login_mail = ?", email).
		First(&account).Error
	return account, err
}

func (r *accountRepo) GetAccountByID(ctx context.Context, id string) (hrmmodel.Account, error) {
	var account hrmmodel.Account
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&account).Error
	return account, err
}

func (r *accountRepo) UpdateAccount(ctx context.Context, account hrmmodel.Account) error {
	return r.db.WithContext(ctx).Save(&account).Error
}
