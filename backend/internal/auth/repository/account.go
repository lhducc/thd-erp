package authRepository

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

type AccountRepository interface {
	GetEmployeeByEmail(ctx context.Context, email string) (hrmmodel.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (hrmmodel.Employee, error)
	UpdatePasswordAndFirstLogin(ctx context.Context, employeeID string, hashedPassword string) error
}

// accountRepo struct chứa db connection và implement interface
type accountRepo struct {
	db *gorm.DB
}

// NewAccountRepository trả về interface AccountRepository
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) GetEmployeeByEmail(ctx context.Context, email string) (hrmmodel.Employee, error) {
	var employee hrmmodel.Employee
	err := r.db.WithContext(ctx).
		Where("email = ? AND status = ?", email, "active").
		First(&employee).Error
	return employee, err
}

func (r *accountRepo) GetEmployeeByID(ctx context.Context, id string) (hrmmodel.Employee, error) {
	var employee hrmmodel.Employee
	err := r.db.WithContext(ctx).
		Where("employee_id = ?", id).
		First(&employee).Error
	return employee, err
}

func (r *accountRepo) UpdatePasswordAndFirstLogin(ctx context.Context, employeeID string, hashedPassword string) error {
	return r.db.WithContext(ctx).
		Model(&hrmmodel.Employee{}).
		Where("employee_id = ?", employeeID).
		Updates(map[string]interface{}{
			"password":    hashedPassword,
			"first_login": false,
		}).Error
}
