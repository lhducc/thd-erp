package handler

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmployeeBiz interface {
	CreateEmployeeWithAccount(employee model.Employee) error
	GetUserById(id string) (model.Employee, error)
	UpdateEmployee(id string, updatedEmployee model.Employee) error
	DeleteEmployee(id string) error
	GetAllEmployees() ([]model.Employee, error)
	GetAllEmployeesByStatus(status string) ([]model.Employee, error)
	ExportEmployeeTest(selectedFields []string) ([]byte, string, error)
}

type EmployeeHandler struct {
	employeeBiz EmployeeBiz
}

func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	store := repository.NewUserStore(db)
	account := repository.NewAccountStore(db)
	biz := usecase.NewEmployeeBiz(store, account)

	return &EmployeeHandler{
		employeeBiz: biz,
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

func (biz *EmployeeHandler) GetAllEmployees() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := biz.employeeBiz.GetAllEmployees()
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (biz *EmployeeHandler) GetAllEmployeeByStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status := ctx.Param("status")

		result, err := biz.employeeBiz.GetAllEmployeesByStatus(status)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, result)
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
