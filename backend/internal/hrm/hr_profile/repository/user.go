package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
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

func (s *UserStore) CheckExistEmployees(employeeIDs []string) ([]string, error) {
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

func (s *UserStore) GetBySchedule(ctx context.Context, scheduleIDs []int, managerID string, filter string) ([]model.Employee, error) {
	var employees []model.Employee
	query := s.db.WithContext(ctx).Model(model.Employee{}).
		Select("employee.employee_id", "employee.full_name", "employee.schedule_id", "employee.department_id").
		Preload("Department")

	if strings.TrimSpace(managerID) != "" {
		query = query.Joins("JOIN work_schedule_manager ON employee.schedule_id = work_schedule_manager.work_schedule_id").
			Where("work_schedule_manager.employee_id = ?", managerID)
	}

	if filter != "" {
		switch filter {
		case "register":
			query = query.Joins("JOIN work_schedule ON employee.schedule_id = work_schedule.work_schedule_id").
				Where("work_schedule.is_schedule_auto = ?", false)
		case "auto":
			query = query.Joins("JOIN work_schedule ON employee.schedule_id = work_schedule.work_schedule_id").
				Where("work_schedule.is_schedule_auto = ?", true)
		case "no_schedule":
			query = query.Where("employee.schedule_id IS NULL")
		}
	}

	// Xử lý lọc theo scheduleIDs nếu có
	if len(scheduleIDs) > 0 {
		query = query.Where("employee.schedule_id IN (?)", scheduleIDs)
	}

	if err := query.Find(&employees).Error; err != nil {
		fmt.Printf("Lỗi DB: %s", err.Error())
		return nil, errors.New("Lỗi hệ thống, truy vấn người dùng theo lịch làm việc thất bại")
	}

	return employees, nil
}

func (s *UserStore) CheckExists(employeeID string) (bool, error) {
	var employees model.Employee
	err := s.db.
		Where("employee_id = ?", employeeID).
		First(&employees).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get employees: %w", err)
	}
	return true, nil
}

func (s *UserStore) GetScheduleOfEmployee(employeeID string) (*model.Employee, error) {
	var employee model.Employee
	err := s.db.
		Where("employee_id = ?", employeeID).
		Select("schedule_id").
		First(&employee).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	return &employee, nil
}
