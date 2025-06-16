package repository

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *userStore {
	return &userStore{db: db}
}

func (s *userStore) recoverFromPanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		log.Printf("Panic occurred, transaction rolled back: %v", r)
	}
}

func (s *userStore) CreateEmployee(employee *model.Employee) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	s.recoverFromPanic(tx)

	err := model.ValidateEmployeeReferences(tx, *employee)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Lỗi tạo nhân viên: %w", err)
	}

	if err := tx.Create(employee).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create employee: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *userStore) GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error) {
	var employees []model.Employee
	offset := (page - 1) * pageSize

	result := s.db.
		Where("status = ?", status).
		Order("employee_id ASC").
		Limit(pageSize).
		Offset(offset).
		Preload("Account").
		Preload("Manager").
		Preload("JobTitle").
		Preload("Position").
		Find(&employees)

	if result.Error != nil {
		return []model.Employee{}, fmt.Errorf("failed to get employees with status %s: %v", status, result.Error)
	}
	return employees, nil
}

func (s *userStore) GetUserById(id string) (model.Employee, error) {
	var employee model.Employee
	if err := s.db.Where("employee_id = ?", id).
		Preload("Account").
		Preload("Manager").
		Preload("JobTitle").
		Preload("Position").
		First(&employee).Error; err != nil {
		return model.Employee{}, fmt.Errorf("employee not found with id %d", id)
	}
	return employee, nil
}

func (s *userStore) GetAllEmployees() ([]model.Employee, error) {
	var employees []model.Employee
	if err := s.db.
		Order("employee_id ASC").
		Preload("Account").
		Preload("Manager").
		Preload("JobTitle").
		Preload("Position").
		Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	return employees, nil
}

func (s *userStore) GetAllEmployeesPagination(page, pageSize int) ([]model.Employee, error) {
	var employees []model.Employee
	offset := (page - 1) * pageSize

	if err := s.db.
		Order("employee_id ASC").
		Limit(pageSize).
		Offset(offset).
		Preload("Account").
		Preload("Manager").
		Preload("JobTitle").
		Preload("Position").
		Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	return employees, nil
}

func (s *userStore) UpdateEmployee(id string, updatedEmployee model.Employee) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	defer s.recoverFromPanic(tx)

	existingEmployee, err := s.GetUserById(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	model.UpdateEmployeeFields(&existingEmployee, updatedEmployee)

	if err := tx.Save(&existingEmployee).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update employee: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *userStore) DeleteEmployee(id string) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	s.recoverFromPanic(tx)

	var employee model.Employee
	if err := s.db.Where("employee_id = ?", id).First(&employee).Error; err != nil {
		s.db.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("employee with ID %d not found", id)
		}
		return fmt.Errorf("failed to check employee: %w", err)
	}

	if employee.Status == "Inactive" {
		s.db.Rollback()
		return fmt.Errorf("employee status is deleted")
	}

	if err := s.db.Model(&employee).Update("status", "Inactive").Error; err != nil {
		s.db.Rollback()
		return fmt.Errorf("failed to update employee status to 'Inactive': %w", err)
	}

	if err := s.db.Commit().Error; err != nil {
		s.db.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *userStore) GetLastEmployeeByCode(emp *model.Employee) error {
	return s.db.
		Order("employee_id DESC").
		First(emp).Error
}

func (s *userStore) UpdateEmployeeWithAccount(employee *model.Employee, accountID int) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	s.recoverFromPanic(tx)

	employee.AccountID = &accountID
	if err := tx.Save(employee).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update employee with account ID: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
