package repository

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) recoverFromPanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		log.Printf("Panic occurred, transaction rolled back: %v", r)
	}
}

func (s *UserStore) CreateEmployee(employee *model.Employee) error {
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

func (s *UserStore) GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error) {
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
		Preload("Department").
		Preload("Department.Office").
		Find(&employees)

	if result.Error != nil {
		return []model.Employee{}, fmt.Errorf("failed to get employees with status %s: %v", status, result.Error)
	}
	return employees, nil
}

func (s *UserStore) GetUserById(id string) (model.Employee, error) {
	var employee model.Employee
	if err := s.db.Where("employee_id = ?", id).
		Preload("Account").
		Preload("Manager").
		Preload("JobTitle").
		Preload("Position").
		Preload("Department").
		Preload("Department.Office").
		First(&employee).Error; err != nil {
		return model.Employee{}, fmt.Errorf("employee not found with id %s", id)
	}
	return employee, nil
}

func (s *UserStore) GetAllEmployees() ([]model.Employee, error) {
	var employees []model.Employee
	if err := s.db.
		Order("employee_id ASC").
		Preload("Account").
		Preload("Manager").
		Preload("JobTitle").
		Preload("Position").
		Preload("Department").
		Preload("Department.Office").
		Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	return employees, nil
}

func (s *UserStore) GetAllEmployeesPagination(page, pageSize int, filters map[string]interface{}) ([]model.Employee, int64, error) {
	var employees []model.Employee
	var totalRecords int64

	db := s.db.Model(&model.Employee{})

	for key, value := range filters {
		switch key {
		case "office_id":
			if arr, ok := value.([]string); ok && len(arr) > 0 {
				db = db.Joins("JOIN department d ON employee.department_id = d.department_id").
					Where("d.office_id IN (?)", arr)
			}
		case "department_id":
			if arr, ok := value.([]string); ok && len(arr) > 0 {
				db = db.Where("employee.department_id IN (?)", arr)
			}
		default:
			if arr, ok := value.([]string); ok && len(arr) > 0 {
				db = db.Where(fmt.Sprintf("employee.%s IN (?)", key), arr)
			} else if value != nil && value != "" {
				db = db.Where(fmt.Sprintf("employee.%s = ?", key), value)
			}
		}
	}

	// total records
	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := db.Order("employee.employee_id ASC").
		Limit(pageSize).
		Offset(offset).
		Preload("Account").
		Preload("Manager").
		Preload("JobTitle").
		Preload("Position").
		Preload("Department").
		Preload("Department.Office").
		Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get employees: %w", err)
	}
	return employees, totalRecords, nil
}

func (s *UserStore) UpdateEmployee(id string, updatedEmployee model.Employee) error {
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

func (s *UserStore) DeleteEmployee(id string) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var employee model.Employee
	if err := tx.Where("employee_id = ?", id).First(&employee).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("employee with ID %s not found", id)
		}
		return fmt.Errorf("failed to check employee: %w", err)
	}

	if employee.Status == "inactive" {
		tx.Rollback()
		return fmt.Errorf("employee status is deleted")
	}

	if err := tx.Model(&employee).Update("status", "inactive").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update employee status to 'inactive': %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *UserStore) GetLastEmployeeByCode(emp *model.Employee) error {
	return s.db.
		Order("employee_id DESC").
		First(emp).Error
}

func (s *UserStore) UpdateEmployeeWithAccount(employee *model.Employee, accountID int64) error {
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

func (s *UserStore) CheckExistEmployee(employeeIDs []string) ([]string, error) {
	var employees []model.Employee

	if err := s.db.
		Where("employee_id IN ?", employeeIDs).
		Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	// Map để đối chiếu nhanh
	found := make(map[string]bool)
	for _, emp := range employees {
		found[emp.EmployeeID] = true
	}

	// Tìm những ID không tồn tại
	var missing []string
	for _, id := range employeeIDs {
		if !found[id] {
			missing = append(missing, id)
		}
	}
	return missing, nil
}
