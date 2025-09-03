package usecase

import (
	"context"
	checkin_model "erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/model/dto"
	"erp/backend/internal/hrm/hr_profile/repository"
	utils "erp/backend/pkg"
	"erp/backend/pkg/mail"
	"erp/backend/pkg/transaction"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type EmployeeRepo interface {
	CreateEmployee(tx *gorm.DB, employee *model.Employee) error
	UpdateEmployeeWithAccount(tx *gorm.DB, employee *model.Employee, accountID int64) error
	GetUserById(id string) (model.Employee, error)
	GetAllEmployees() ([]model.Employee, error)
	GetAllEmployeesPagination(page, pageSize int, filters map[string]interface{}) ([]model.Employee, int64, error)
	GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error)
	UpdateEmployee(tx *gorm.DB, updatedEmployee model.Employee) error
	DeleteEmployee(id string) error
	GetLastEmployeeByCode(emp *model.Employee) error
	CheckExistEmployees(employeeIDs []string) ([]string, error)
	GetBySchedule(ctx context.Context, scheduleIDs []int, managerID string, filter string) ([]model.Employee, error)
	CheckExists(employeeID string) (bool, error)
	GetScheduleOfEmployee(employeeID string) (*model.Employee, error)
	GetEmployeesByOfficeID(ctx context.Context, officeID string) ([]*model.Employee, error)
	CheckExistEmployeeID(employeeID string) (bool, error)
	GetUserByRoleID(roleID string) ([]model.ManagerResponse, error)
	UpdateStatusEmployee(employeeID, statusChange string) error
}

type AccountRepo interface {
	CreateAccount(tx *gorm.DB, employee *model.Employee, hashedPassword, roleID string) (*model.Account, error)
	GetAccount(ctx context.Context, id int64) (*model.Account, error)
	GetAccountNoCtx(id int64) (*model.Account, error)
	GetAllAccount(ctx context.Context) ([]model.Account, error)
	UpdateAccount(ctx context.Context, id string, data *model.Account) error
	UpdateAccountTrans(tx *gorm.DB, id int64, data *model.Account) error
	DeleteAccount(ctx context.Context, id string) error
	CheckExistEmail(email string) (bool, error)
	CheckFirstLogin(email string) (bool, error)
}

type EmployeeBiz struct {
	db                *gorm.DB
	repo              EmployeeRepo
	account           AccountRepo
	timesheetRepo     repo_interface.TimeSheetRepoInterface
	timesheetListRepo repo_interface.TimesheetListInterface
	departmentRepo    DepartmentRepo
}

func NewEmployeeBiz(db *gorm.DB, store *repository.UserStore,
	account AccountRepo,
	timesheetRepo repo_interface.TimeSheetRepoInterface,
	timesheetListRepo repo_interface.TimesheetListInterface,
	departmentRepo DepartmentRepo) *EmployeeBiz {
	return &EmployeeBiz{
		db:                db,
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
	return transaction.WithTransaction(s.db, ctx, func(ctx context.Context, tx *gorm.DB) error {

		if err := s.repo.CreateEmployee(tx, employee); err != nil {
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
		account, err := s.account.CreateAccount(tx, employee, hashedPassword, roleID)
		if err != nil {
			return err
		}

		if err := s.repo.UpdateEmployeeWithAccount(tx, employee, account.ID); err != nil {
			return err
		}

		// add employee into timesheet
		timeNow := utils.GetCurrentTimeHCMCity()
		month := int(timeNow.Month())
		year := int(timeNow.Year())

		timesheetList, err := s.timesheetListRepo.GetTimeSheetByTime(ctx, month, year)
		if timesheetList != nil {
			departmentID := employee.DepartmentID
			fmt.Printf(departmentID)
			if departmentID == "" {
				return errors.New("giá trị mã phòng ban không hợp lệ")
			}
			department, err := s.departmentRepo.GetDepartment(ctx, departmentID)
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
			if err = s.timesheetRepo.CreateEmployeeTimeSheet(tx, &timesheet); err != nil {
				fmt.Printf(err.Error())
				return errors.New("Lỗi trong quá trình thêm nhân viên vào bảng công")
			}
		}

		if err := mail.SendEmailWithAccountInfo(employee.Email, employee.Fullname, password); err != nil {
			return fmt.Errorf("tạo nhân viên thành công nhưng gửi email thất bại: %w", err)
		}

		return nil
	})
}

func (biz *EmployeeBiz) GetUserById(id string) (*dto.EmployeeResponse, error) {
	if id == "" {
		return nil, errors.New("invalid employee ID")
	}

	employee, err := biz.repo.GetUserById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	var role model.Role
	if employee.AccountID != nil {
		account, err := biz.account.GetAccountNoCtx(*employee.AccountID)
		if err != nil {
			return nil, fmt.Errorf("failed to get account: %w", err)
		}
		if account != nil && account.Role != nil {
			role = *account.Role
		}
	}

	employeeInfo := dto.EmployeeResponse{
		Employee: &employee,
		Role:     &role,
	}
	return &employeeInfo, nil
}

func (biz *EmployeeBiz) GetAllEmployees(page, pageSize int, filters map[string]interface{}) ([]dto.EmployeeResponse, int64, error) {
	employees, totalRecords, err := biz.repo.GetAllEmployeesPagination(page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all employees: %w", err)
	}

	var employeeInfs []dto.EmployeeResponse
	for _, employee := range employees {
		var employeeInf dto.EmployeeResponse
		acc, err := biz.account.GetAccountNoCtx(*employee.AccountID)
		employeeInf.Employee = &employee
		if err != nil {
			return nil, 0, err
		}
		if acc != nil && acc.Role != nil {
			employeeInf.Role = acc.Role
		}
		employeeInfs = append(employeeInfs, employeeInf)
	}

	return employeeInfs, totalRecords, nil
}

func (biz *EmployeeBiz) GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error) {
	employees, err := biz.repo.GetAllEmployeesByStatus(status, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees with status: %w", err)
	}

	return employees, nil
}

func (biz *EmployeeBiz) UpdateEmployee(ctx context.Context, id string, updatedEmployee dto.EmployeeDTO) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}

	emp, err := biz.repo.GetUserById(id)
	if err != nil {
		return err
	}

	return transaction.WithTransaction(biz.db, ctx, func(ctx context.Context, tx *gorm.DB) error {
		if emp.DepartmentID != updatedEmployee.DepartmentID {
			department, err := biz.departmentRepo.GetDepartment(ctx, updatedEmployee.DepartmentID)
			if err != nil {
				return fmt.Errorf("lỗi khi lấy dữ liệu Department: %w", err)
			}

			timeNow := utils.GetCurrentTimeHCMCity()
			month := int(timeNow.Month())
			year := timeNow.Year()

			// get timesheet list
			timesheetList, err := biz.timesheetListRepo.GetTimeSheetByTime(ctx, month, year)
			if err != nil {
				return fmt.Errorf("lỗi khi lấy timesheet list: %w", err)
			}

			if timesheetList != nil {
				timesheet, err := biz.timesheetRepo.CheckExist(ctx, id, timesheetList.Month, timesheetList.Year)
				if err != nil {
					return err
				}

				if timesheet == nil {
					newTs := checkin_model.TimeSheet{
						TimeSheetListID: timesheetList.TimeSheetListID,
						OfficeID:        department.OfficeID,
						Month:           timesheetList.Month,
						Year:            timesheetList.Year,
						DepartmentID:    updatedEmployee.DepartmentID,
						EmployeeID:      updatedEmployee.EmployeeID,
						CreatedBy:       timesheetList.CreatedBy,
					}
					if err = biz.timesheetRepo.CreateEmployeeTimeSheet(tx, &newTs); err != nil {
						return fmt.Errorf("lỗi khi thêm nhân viên vào bảng công: %w", err)
					}
				}
			}
		}

		// update role nếu có
		account, err := biz.GetAccount(ctx, *emp.AccountID)
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}
		if updatedEmployee.RoleID != "" {
			account.RoleID = updatedEmployee.RoleID
			if err := biz.account.UpdateAccountTrans(tx, account.ID, account); err != nil {
				return fmt.Errorf("failed to update account role: %w", err)
			}
		}

		// Convert DTO thành model và update employee
		employee := updatedEmployee.ConvertToEmployeeModel()
		employee.EmployeeID = id
		employee.AccountID = &account.ID
		if err := biz.repo.UpdateEmployee(tx, *employee); err != nil {
			return fmt.Errorf("failed to update employee: %w", err)
		}

		return nil
	})
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

	// exporter.RegisterField("email", "Email", "email", func(item interface{}) any {
	// 	emp := item.(*model.Employee)
	// 	if emp.Account != nil {
	// 		return emp.Account.LoginMail
	// 	}
	// 	return ""
	// })

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

func (e *EmployeeBiz) GetUserByRoleID(roleID string) ([]model.ManagerResponse, error) {
	return e.repo.GetUserByRoleID(roleID)
}

func (e *EmployeeBiz) UpdateStatus(employeeID, statusChange string) error {
	//check exist
	exist, err := e.repo.CheckExists(employeeID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("not_found")
	}
	return e.repo.UpdateStatusEmployee(employeeID, statusChange)
}
