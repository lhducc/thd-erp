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
	GetEmployeeByAccountID(ctx context.Context, accountID int64) (hrmmodel.Employee, error)
	GetAccountByEmployeeID(ctx context.Context, employeeID string) (*hrmmodel.Account, error)
}

// accountRepo struct chứa db connection và implement interface
type accountRepo struct {
	db *gorm.DB
}

// NewAccountRepository trả về interface AccountRepository
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{db: db}
}

// tối ưu ở đây
func (r *accountRepo) GetAccountByEmail(ctx context.Context, email string) (hrmmodel.Account, error) {
	var account hrmmodel.Account
	err := r.db.WithContext(ctx).Joins("JOIN employee on employee.employee_id = account.employee_id").
		Where("account.login_mail = ? AND employee.status = ?", email, "active").
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
func (r *accountRepo) GetEmployeeByAccountID(ctx context.Context, accountID int64) (hrmmodel.Employee, error) {
	var employee hrmmodel.Employee
	err := r.db.Raw(`SELECT
    e.employee_id,
    e.full_name,
    e.account_id,
    a.role_id,
    e.created_date
FROM
    employee AS e   
JOIN
    account AS a      
    ON e.account_id = a.id
WHERE
    e.account_id = ?`, accountID).Scan(&employee).Error
	return employee, err
}

func (r *accountRepo) GetAccountByEmployeeID(ctx context.Context, employeeID string) (*hrmmodel.Account, error) {
	var account hrmmodel.Account
	if err := r.db.WithContext(ctx).
		Where("employee_id = ?", employeeID).
		First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
