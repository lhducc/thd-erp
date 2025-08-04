package usecase

import (
	"context"
	checkin_model "erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	utils "erp/backend/pkg"
	"erp/backend/pkg/mail"
	"errors"
	"fmt"
	"strings"
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
	CheckExistEmployees(employeeIDs []string) ([]string, error)
	GetBySchedule(ctx context.Context, scheduleIDs []int, managerID string, filter string) ([]model.Employee, error)
	CheckExists(employeeID string) (bool, error)
	GetScheduleOfEmployee(employeeID string) (*model.Employee, error)
	GetEmployeesByOfficeID(ctx context.Context, officeID string) ([]*model.Employee, error)
	CheckExistEmployeeID(employeeID string) (bool, error)
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
	repo              EmployeeRepo
	account           AccountRepo
	timesheetRepo     repo_interface.TimeSheetRepoInterface
	timesheetListRepo repo_interface.TimesheetListInterface
	departmentRepo    DepartmentRepo
}

func NewEmployeeBiz(store *repository.UserStore,
	account AccountRepo,
	timesheetRepo repo_interface.TimeSheetRepoInterface,
	timesheetListRepo repo_interface.TimesheetListInterface,
	departmentRepo DepartmentRepo) *EmployeeBiz {
	return &EmployeeBiz{
		repo:              store,
		account:           account,
		timesheetRepo:     timesheetRepo,
		timesheetListRepo: timesheetListRepo,
		departmentRepo:    departmentRepo,
	}
}

func (s *EmployeeBiz) CreateEmployeeWithAccount(ctx context.Context, employee *model.Employee, roleID string) error {

	exists, err := s.repo.CheckExistEmployeeID(employee.EmployeeID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("employee with id " + employee.EmployeeID + " already exists")
	}

	if err := s.repo.CreateEmployee(employee); err != nil {
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

	if strings.TrimSpace(roleID) == "" {
		roleID = "employee"
	}
	account, err := s.account.CreateAccount(employee, hashedPassword, roleID)
	if err != nil {
		return err
	}

	if err := mail.SendEmailWithAccountInfo(employee.Email, employee.Fullname, password); err != nil {
		return fmt.Errorf("tạo nhân viên thành công nhưng gửi email thất bại: %w", err)
	}

	if err := s.repo.UpdateEmployeeWithAccount(employee, account.ID); err != nil {
		return err
	}

	// add employee into timesheet
	timeNow, err := utils.GetCurrentTimeHCMCity()
	if err != nil {
		return err
	}
	month := int(timeNow.Month())
	year := int(timeNow.Year())

	timesheetList, err := s.timesheetListRepo.GetTimeSheetByOfficeIDAndTime(employee.Department.OfficeID, month, year)
	if timesheetList != nil {
		department, err := s.departmentRepo.GetDepartment(ctx, employee.DepartmentID)
		if err != nil {
			return errors.New("Lỗi khi lấy dữ liệu Department")
		}
		timesheet := checkin_model.TimeSheet{
			TimeSheetListID: timesheetList.TimeSheetListID,
			OfficeID:        department.OfficeID,
			Month:           timesheetList.Month,
			Year:            timesheetList.Year,
			DepartmentID:    employee.DepartmentID,
			EmployeeID:      employee.EmployeeID,
			CreatedBy:       timesheetList.CreatedBy,
		}
		if err = s.timesheetRepo.CreateEmployeeTimeSheet(&timesheet); err != nil {
			fmt.Printf(err.Error())
			return errors.New("Lỗi trong quá trình thêm nhân viên vào bảng công")
		}
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
