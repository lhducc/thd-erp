package handler

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmployeeBiz interface {
	CreateEmployeeWithAccount(employee model.Employee) error
	GetUserById(id string) (model.Employee, error)
	UpdateEmployee(id string, updatedEmployee model.Employee) error
	DeleteEmployee(id string) error
	GetAllEmployees(page, pageSize int, filters map[string]interface{}) ([]model.Employee, int64, error)
	GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error)
	ExportEmployeeTest(selectedFields []string) ([]byte, string, error)
}

type EmployeeHandler struct {
	employeeBiz   EmployeeBiz
	departmentBiz DepartmentBiz
}

func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	store := repository.NewUserStore(db)
	account := repository.NewAccountStore(db)
	biz := usecase.NewEmployeeBiz(store, account)

	// Khởi tạo Department Usecase
	departmentStore := repository.NewDepartmentStore(db)
	departmentBiz := usecase.NewDepartmentBiz(departmentStore)

	return &EmployeeHandler{
		employeeBiz:   biz,
		departmentBiz: departmentBiz,
	}
}

func (biz *EmployeeHandler) CreateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.Employee

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, "Error data", http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := model.ValidateEmployee(data); err != nil {
			utils.ResponseMessage(c, "Error data", http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := biz.employeeBiz.CreateEmployeeWithAccount(data); err != nil {
			utils.ResponseMessage(c, "Error save db", http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.ResponseMessage(c, "Create success", http.StatusCreated, nil)
	}
}

func (h *EmployeeHandler) GetAllEmployees() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		// filter
		filterFields := []string{"job_title_id", "department_id", "position_id", "office_id"}
		filters := make(map[string]interface{})
		utils.ExtractFilterArrays(ctx, filterFields)

		// get data
		employees, totalRecords, err := h.employeeBiz.GetAllEmployees(page, pageSize, filters)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		//Add department information
		newEmployees := h.getAllDepartmentByListEmployee(ctx, employees)

		// caculator the total of page number
		totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

		response := gin.H{
			"data":         &newEmployees,
			"totalRecords": totalRecords,
			"page":         page,
			"pageSize":     pageSize,
			"totalPages":   totalPages,
		}

		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, response)
	}
}

// get department by list employeeAdd commentMore actions
func (h *EmployeeHandler) getAllDepartmentByListEmployee(ctx *gin.Context, employees []model.Employee) *[]model.Employee {
	for i, emp := range employees {
		deparmentID := emp.DepartmentID
		if deparmentID == "" {
			continue
		}
		dept, err := h.departmentBiz.GetDepartment(ctx.Request.Context(), emp.DepartmentID)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi khi lấy department: %s", err.Error()), http.StatusNotFound, nil)
			return nil
		}
		employees[i].Department = dept
	}
	return &employees
}

func (biz *EmployeeHandler) GetAllEmployeeByStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status := ctx.Param("status")

		page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		employees, err := biz.employeeBiz.GetAllEmployeesByStatus(status, page, pageSize)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusNotFound, nil)
			return
		}

		//Add department information
		for i, emp := range employees {
			if emp.DepartmentID == "" {
				// Bỏ qua, không cần load department
				continue
			}
			dept, err := biz.departmentBiz.GetDepartment(ctx.Request.Context(), emp.DepartmentID)
			if err != nil {
				utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi khi lấy department: %s", err.Error()), http.StatusNotFound, nil)
				return
			}
			employees[i].Department = dept
		}

		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, employees)
	}
}

func (biz *EmployeeHandler) GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		result, err := biz.employeeBiz.GetUserById(idParam)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusNotFound, nil)
			return
		}

		if result.DepartmentID != "" {
			var Department *model.Department
			Department, err = biz.departmentBiz.GetDepartment(c.Request.Context(), result.DepartmentID)
			if err != nil {
				utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusNotFound, nil)
				return
			}
			result.Department = Department
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (biz *EmployeeHandler) ExportEmployees() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := strings.Split(c.Query("fields"), ",")
		data, filename, err := biz.employeeBiz.ExportEmployeeTest(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	}
}

func (biz *EmployeeHandler) UpdateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var data model.Employee
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}
		if err := biz.employeeBiz.UpdateEmployee(idParam, data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Dữ liệu đã được cập nhật", http.StatusOK, nil)
	}
}
func (biz *EmployeeHandler) DeleteEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		if err := biz.employeeBiz.DeleteEmployee(idParam); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Deleted", http.StatusOK, nil)
	}
}
