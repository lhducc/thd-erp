package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/model/dto"
	utils "erp/backend/pkg"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmployeeBiz interface {
	CreateEmployeeWithAccount(ctx context.Context, employee *model.Employee, roleID string) error
	GetUserById(id string) (model.Employee, error)
	UpdateEmployee(ctx context.Context, id string, updatedEmployee model.Employee) error
	DeleteEmployee(id string) error
	GetAllEmployees(page, pageSize int, filters map[string]interface{}) ([]model.Employee, int64, error)
	GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error)
	ExportEmployeeTest(selectedFields []string) ([]byte, string, error)
	GetUserByRoleID(roleID string) ([]model.ManagerResponse, error)
}

type EmployeeHandler struct {
	employeeBiz   EmployeeBiz
	departmentBiz DepartmentBiz
}

func NewEmployeeHandler(biz EmployeeBiz, bizDepartment DepartmentBiz) *EmployeeHandler {
	return &EmployeeHandler{
		employeeBiz:   biz,
		departmentBiz: bizDepartment,
	}
}

func (biz *EmployeeHandler) CreateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataDTO dto.EmployeeDTO

		if err := c.ShouldBind(&dataDTO); err != nil {
			utils.ResponseMessage(c, "Error data", http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		data := dataDTO.ConvertToEmployeeModel()

		if err := model.ValidateEmployee(data); err != nil {
			utils.ResponseMessage(c, "Error data", http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := biz.employeeBiz.CreateEmployeeWithAccount(c.Request.Context(), data, dataDTO.RoleID); err != nil {
			utils.ResponseMessage(c, "Error save db", http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.ResponseMessage(c, "CreateElementOfTimesheetList success", http.StatusCreated, nil)
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
		filterFields := []string{"job_title_id", "department_id", "position_id", "office_id", "work_type"}
		filters := utils.ExtractFilterArrays(ctx, filterFields)

		// get data
		employees, totalRecords, err := h.employeeBiz.GetAllEmployees(page, pageSize, filters)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

		response := gin.H{
			"data":         &employees,
			"totalRecords": totalRecords,
			"page":         page,
			"pageSize":     pageSize,
			"totalPages":   totalPages,
		}

		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, response)
	}
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

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (biz *EmployeeHandler) GetPersonalInfById() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.GetString("employeeId")

		result, err := biz.employeeBiz.GetUserById(idParam)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusNotFound, nil)
			return
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
		ctx := c.Request.Context()

		var data model.Employee
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}
		if err := biz.employeeBiz.UpdateEmployee(ctx, idParam, data); err != nil {
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
func (biz *EmployeeHandler) GetEmployeesByRoleID() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Query("roleID")
		emps, err := biz.employeeBiz.GetUserByRoleID(roleID)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "List Users", http.StatusOK, &emps)
	}
}
