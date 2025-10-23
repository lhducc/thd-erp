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
	GetEmployeesByManager(ctx context.Context, managerID string) ([]model.Employee, error)
	GetAllEmployeesActivePagination(page, pageSize int, filters map[string]interface{}) ([]model.Employee, int64, error)
}

type EmployeeBiz struct {
	db                *gorm.DB
	repo              EmployeeRepo
	timesheetRepo     repo_interface.TimeSheetRepoInterface
	timesheetListRepo repo_interface.TimesheetListInterface
	departmentRepo    DepartmentRepo
}

func NewEmployeeBiz(db *gorm.DB,
	store *repository.UserStore,
	timesheetRepo repo_interface.TimeSheetRepoInterface,
	timesheetListRepo repo_interface.TimesheetListInterface,
	departmentRepo DepartmentRepo) *EmployeeBiz {
	return &EmployeeBiz{
		db:                db,
		repo:              store,
		timesheetRepo:     timesheetRepo,
		timesheetListRepo: timesheetListRepo,
		departmentRepo:    departmentRepo,
	}
}

func (s *EmployeeBiz) CreateEmployee(ctx context.Context, employee *model.Employee, roleID string) error {
	exists, err := s.repo.CheckExistEmployeeID(employee.EmployeeID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("employee with id %s already exists", employee.EmployeeID)
	}

	return transaction.WithTransaction(s.db, ctx, func(ctx context.Context, tx *gorm.DB) error {
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

		employee.Password = hashedPassword
		employee.RoleID = roleID
		employee.Status = "active"

		if err := s.repo.CreateEmployee(tx, employee); err != nil {
			return err
		}

		// Thêm vào bảng công
		timeNow := utils.GetCurrentTimeHCMCity()
		month := int(timeNow.Month())
		year := int(timeNow.Year())

		timesheetList, err := s.timesheetListRepo.GetTimeSheetByTime(ctx, month, year)
		if err == nil && timesheetList != nil {
			department, err := s.departmentRepo.GetDepartment(ctx, employee.DepartmentID)
			if err != nil {
				return fmt.Errorf("lỗi khi lấy dữ liệu department: %w", err)
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
				return fmt.Errorf("lỗi khi thêm nhân viên vào bảng công: %w", err)
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

	role := model.Role{}
	if employee.Role != nil {
		role = *employee.Role
	}

	return &dto.EmployeeResponse{
		Employee: &employee,
		Role:     &role,
	}, nil
}

func (biz *EmployeeBiz) GetAllEmployees(page, pageSize int, filters map[string]interface{}) ([]dto.EmployeeResponse, int64, error) {
	employees, total, err := biz.repo.GetAllEmployeesPagination(page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all employees: %w", err)
	}

	responses := make([]dto.EmployeeResponse, 0, len(employees))
	for _, e := range employees {
		responses = append(responses, dto.EmployeeResponse{
			Employee: &e,
			Role:     e.Role,
		})
	}
	return responses, total, nil
}

func (biz *EmployeeBiz) GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error) {
	employees, err := biz.repo.GetAllEmployeesByStatus(status, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees with status: %w", err)
	}

	return employees, nil
}

func (biz *EmployeeBiz) UpdateEmployee(ctx context.Context, id string, updated dto.EmployeeDTO) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}

	existing, err := biz.repo.GetUserById(id)
	if err != nil {
		return err
	}

	return transaction.WithTransaction(biz.db, ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Cập nhật phòng ban nếu thay đổi
		if existing.DepartmentID != updated.DepartmentID {
			department, err := biz.departmentRepo.GetDepartment(ctx, updated.DepartmentID)
			if err != nil {
				return fmt.Errorf("lỗi khi lấy dữ liệu Department: %w", err)
			}
			timeNow := utils.GetCurrentTimeHCMCity()
			month := int(timeNow.Month())
			year := timeNow.Year()
			timesheetList, err := biz.timesheetListRepo.GetTimeSheetByTime(ctx, month, year)
			if err == nil && timesheetList != nil {
				ts, _ := biz.timesheetRepo.CheckExist(ctx, id, timesheetList.Month, timesheetList.Year)
				if ts == nil {
					newTS := checkin_model.TimeSheet{
						TimeSheetListID: timesheetList.TimeSheetListID,
						OfficeID:        department.OfficeID,
						Month:           timesheetList.Month,
						Year:            timesheetList.Year,
						DepartmentID:    updated.DepartmentID,
						EmployeeID:      updated.EmployeeID,
						CreatedBy:       timesheetList.CreatedBy,
					}
					if err := biz.timesheetRepo.CreateEmployeeTimeSheet(tx, &newTS); err != nil {
						return fmt.Errorf("lỗi khi thêm nhân viên vào bảng công: %w", err)
					}
				}
			}
		}

		// Cập nhật role nếu có
		if updated.RoleID != "" {
			existing.RoleID = updated.RoleID
		}

		// Convert DTO -> Model
		emp := updated.ConvertToEmployeeModel()
		emp.EmployeeID = id
		emp.RoleID = existing.RoleID

		if err := biz.repo.UpdateEmployee(tx, *emp); err != nil {
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

func (biz *EmployeeBiz) GetEmployeesByManager(ctx context.Context, managerID string) ([]dto.ManagerEmployeeDTO, error) {
	if strings.TrimSpace(managerID) == "" {
		return nil, errors.New("invalid manager ID")
	}

	employees, err := biz.repo.GetEmployeesByManager(ctx, managerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by manager: %w", err)
	}

	// Map sang DTO
	var results []dto.ManagerEmployeeDTO
	for _, e := range employees {
		dtoItem := dto.ManagerEmployeeDTO{
			EmployeeID:     e.EmployeeID,
			FullName:       e.Fullname,
			PhoneNumber:    e.PhoneNumber,
			Email:          e.Email,
			PositionName:   e.Position.Name,
			JobTitle:       e.JobTitle.JobTitle,
			DepartmentName: e.Department.Name,
			OfficeName:     e.Department.Office.Name,
		}
		results = append(results, dtoItem)
	}

	return results, nil
}

func (biz *EmployeeBiz) GetAllEmployeesActive(page, pageSize int, filters map[string]interface{}) ([]dto.EmployeeResponse, int64, error) {
	employees, total, err := biz.repo.GetAllEmployeesActivePagination(page, pageSize, filters)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.EmployeeResponse
	for _, e := range employees {
		res = append(res, dto.EmployeeResponse{
			Employee: &e,
			Role:     e.Role,
		})
	}
	return res, total, nil
}
