package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"erp/backend/pkg/mail"
	"errors"
	"fmt"
)

type EmployeeRepo interface {
	CreateEmployee(employee *model.Employee) error
	UpdateEmployeeWithAccount(employee *model.Employee, accountID int64) error
	GetUserById(id string) (model.Employee, error)
	GetAllEmployees() ([]model.Employee, error)
	GetAllEmployeesPagination(page, pageSize int, filters map[string]interface{}) ([]model.Employee, int64, error)
	GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error)
	UpdateEmployee(id string, updatedEmployee model.Employee) error
	DeleteEmployee(id string) error
	GetLastEmployeeByCode(emp *model.Employee) error
	CheckExistEmployee(employeeIDs []string) ([]string, error)
}

type AccountRepo interface {
	CreateAccount(employee *model.Employee, hashedPassword, roleID string) (*model.Account, error)
	GetAccount(ctx context.Context, id int64) (*model.Account, error)
	GetAllAccount(ctx context.Context) ([]model.Account, error)
	UpdateAccount(ctx context.Context, id string, data *model.Account) error
	DeleteAccount(ctx context.Context, id string) error
	CheckExistEmail(email string) (bool, error)
	CheckFirstLogin(email string) (bool, error)
}

type EmployeeBiz struct {
	repo    EmployeeRepo
	account AccountRepo
}

func NewEmployeeBiz(store EmployeeRepo, account AccountRepo) *EmployeeBiz {
	return &EmployeeBiz{
		repo:    store,
		account: account,
	}
}

func (s *EmployeeBiz) CreateEmployeeWithAccount(employee model.Employee) error {
	id, err := utils.GenerateCode("THD", 3, func() (string, error) {
		var lastEmployee model.Employee
		err := s.repo.GetLastEmployeeByCode(&lastEmployee)
		if err != nil {
			return "", err
		}
		return lastEmployee.EmployeeID, nil
	})
	if err != nil {
		return fmt.Errorf("không thể tạo mã nhân viên: %w", err)
	}

	employee.EmployeeID = id

	if err := s.repo.CreateEmployee(&employee); err != nil {
		return err
	}

	password, err := utils.GenerateRandomPassword(12)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	roleID := "admin"
	account, err := s.account.CreateAccount(&employee, hashedPassword, roleID)
	if err != nil {
		return err
	}

	if err := mail.SendEmailWithAccountInfo(employee.Email, employee.Fullname, password); err != nil {
		return fmt.Errorf("tạo nhân viên thành công nhưng gửi email thất bại: %w", err)
	}

	if err := s.repo.UpdateEmployeeWithAccount(&employee, account.ID); err != nil {
		return err
	}

	return nil
}

func (biz *EmployeeBiz) GetUserById(id string) (model.Employee, error) {
	if id == "" {
		return model.Employee{}, errors.New("invalid employee ID")
	}

	employee, err := biz.repo.GetUserById(id)
	if err != nil {
		return model.Employee{}, fmt.Errorf("failed to get employee: %w", err)
	}

	return employee, nil
}

func (biz *EmployeeBiz) GetAllEmployees(page, pageSize int, filters map[string]interface{}) ([]model.Employee, int64, error) {
	employees, totalRecords, err := biz.repo.GetAllEmployeesPagination(page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all employees: %w", err)
	}

	return employees, totalRecords, nil
}

func (biz *EmployeeBiz) GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error) {
	employees, err := biz.repo.GetAllEmployeesByStatus(status, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees with status: %w", err)
	}

	return employees, nil
}

func (biz *EmployeeBiz) UpdateEmployee(id string, updatedEmployee model.Employee) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}

	if err := biz.repo.UpdateEmployee(id, updatedEmployee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (biz *EmployeeBiz) DeleteEmployee(id string) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}
	if err := biz.repo.DeleteEmployee(id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}

// get account by id
func (biz *EmployeeBiz) GetAccount(ctx context.Context, accountId int64) (*model.Account, error) {
	account, err := biz.account.GetAccount(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, err
}

func (e *EmployeeBiz) ExportEmployeeTest(selectedFields []string) ([]byte, string, error) {
	exporter := utils.NewExcelExporter("Employees")

	exporter.RegisterField("employee_id", "Mã Nhân Viên", "employee_id", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.EmployeeID
	})

	exporter.RegisterField("fullname", "Họ và tên", "fullname", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.Fullname
	})

	exporter.RegisterField("birthday", "Ngày sinh", "birthday", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.Birthday
	})

	exporter.RegisterField("gender", "Giới tính", "gender", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.Gender
	})

	exporter.RegisterField("phone", "Số điện thoại", "phone", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.PhoneNumber
	})

	exporter.RegisterField("email", "Email", "email", func(item interface{}) any {
		emp := item.(*model.Employee)
		if emp.Account != nil {
			return emp.Account.LoginMail
		}
		return ""
	})

	exporter.RegisterField("manager", "Quản lý", "manager", func(item interface{}) any {
		emp := item.(*model.Employee)
		if emp.Manager != nil {
			return emp.Manager.Fullname
		}
		return ""
	})

	exporter.RegisterField("created_date", "Ngày tạo", "created_date", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.CreatedDate.Format("2006-01-02")
	})
	employees, err := e.repo.GetAllEmployees()
	if err != nil {
		return nil, "", fmt.Errorf("không thể lấy dữ liệu nhân viên: %w", err)
	}

	employeePtrs := make([]*model.Employee, len(employees))
	for i := range employees {
		employeePtrs[i] = &employees[i]
	}

	return exporter.Export(employeePtrs, selectedFields)
}
