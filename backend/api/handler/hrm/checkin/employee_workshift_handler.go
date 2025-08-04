package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type EmployeeWorkshiftHandler struct {
	biz service_interface.EmployeeWorkshiftService
}

func NewEmployeeWorkshift(biz service_interface.EmployeeWorkshiftService) *EmployeeWorkshiftHandler {
	return &EmployeeWorkshiftHandler{
		biz: biz,
	}
}

func (h *EmployeeWorkshiftHandler) RegisterPersonal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		employeeID := c.GetString("employeeId")
		var req model.EmployeeWorkshift

		// Parse incoming JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		// CreateElementOfTimesheetList the workshift
		if err := h.biz.Register(ctx, employeeID, req.WorkShiftID, req.Date); err != nil {
			utils.ResponseMessage(c, "Failed to register workshift: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "RegisterPersonal workshift successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req model.EmployeeWorkshift
		// Parse incoming JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		role := c.GetString("role")
		managerID := c.GetString("employeeId")
		if role == "manager" {
			record, err := h.biz.CheckManagerPermission(ctx, managerID, req.EmployeeID)
			if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("error: %s", err.Error())
				utils.ResponseMessage(c, "Lỗi xác thực quyền thao tác của quản lý", http.StatusBadRequest, nil)
				return
			}
			if !record.IsEditing {
				utils.ResponseMessage(c, "Không có quyền thao tác chỉnh sửa dữ liệu", http.StatusForbidden, nil)
			}
		}

		// CreateElementOfTimesheetList the workshift
		if err := h.biz.Register(ctx, req.EmployeeID, req.WorkShiftID, req.Date); err != nil {
			utils.ResponseMessage(c, "Failed to register workshift: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "RegisterPersonal workshift successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.biz.Delete(id); err != nil {
			utils.ResponseMessage(c, "Failed to delete employee_workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Delete workshift successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) DeleteByManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		empWSId := c.Param("id")
		role := c.GetString("role")
		managerID := c.GetString("employeeId")
		ctx := c.Request.Context()
		if role == "manager" {
			record, err := h.biz.CheckManagerPermission(ctx, managerID, empWSId)
			if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("error: %s", err.Error())
				utils.ResponseMessage(c, "Lỗi xác thực quyền thao tác của quản lý", http.StatusBadRequest, nil)
				return
			}
			if !record.IsEditing {
				utils.ResponseMessage(c, "Manager không có quyền chỉnh sửa", http.StatusForbidden, nil)
				return
			}
		}
		if err := h.biz.DeleteByManager(empWSId); err != nil {
			utils.ResponseMessage(c, "Failed to delete employee_workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Delete emp_workshift by manager successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) DeletePersonalShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		employeeID := c.GetString("employeeId")

		if err := h.biz.DeletePersonalShift(id, employeeID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.ResponseMessage(c, "Không tìm thấy dữ liệu cần xóa", http.StatusNotFound, nil)
				return
			}
			utils.ResponseMessage(c, "Failed to delete employee_workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Delete workshift personal successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) GetAllByEmployeeID() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employeeID")
		monthstr := c.Query("month")
		yearstr := c.Query("year")

		month, err := strconv.Atoi(monthstr)
		if err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
		}
		year, err := strconv.Atoi(yearstr)
		if err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
		}

		result, err := h.biz.GetByUserIdAndMonthYear(c.Request.Context(), employeeID, month, year)
		if err != nil {
			utils.ResponseMessage(c, "Failed to get workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Get workshift successfully", http.StatusOK, result)
	}
}

func (h *EmployeeWorkshiftHandler) GetAllPersonal() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.GetString("employeeId")
		monthstr := c.Query("month")
		yearstr := c.Query("year")

		month, err := strconv.Atoi(monthstr)
		if err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
		}
		year, err := strconv.Atoi(yearstr)
		if err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
		}

		result, err := h.biz.GetByUserIdAndMonthYear(c.Request.Context(), employeeID, month, year)
		if err != nil {
			utils.ResponseMessage(c, "Failed to get workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Get workshift successfully", http.StatusOK, result)
	}
}

func (h *EmployeeWorkshiftHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.biz.GetAll()
		if err != nil {
			utils.ResponseMessage(c, "Failed to get workshift: "+err.Error(), http.StatusBadRequest, nil)
		}
		utils.ResponseMessage(c, "Get all workshift successfully", http.StatusOK, result)
	}
}

func (h *EmployeeWorkshiftHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ID from URL parameter
		id := c.Param("id")

		// Bind request body
		var req struct {
			WorkShiftID string `json:"work_shift_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			return
		}

		// Call service
		if err := h.biz.Update(id, req.WorkShiftID); err != nil {
			utils.ResponseMessage(c, "Failed to update workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Update workshift successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) GetListShiftAllowRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		employeeID := c.GetString("employeeId")

		// Call service
		var (
			shifts []model.WorkScheduleShift
			err    error
		)
		if shifts, err = h.biz.GetListShiftAllowRegister(ctx, employeeID); err != nil {
			utils.ResponseMessage(c, "Failed: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "successfully", http.StatusOK, shifts)
	}
}

func (h *EmployeeWorkshiftHandler) AssignEmpShifts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req dto.AssignEmployeeWorkshiftRequest
		// Parse incoming JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		// validate
		if len(req.EmpWorkShift) == 0 {
			utils.ResponseMessage(c, "Empty shift assignment list", http.StatusBadRequest, nil)
			return
		}

		// CreateElementOfTimesheetList the workshift
		if err := h.biz.Assign(ctx, req.EmpWorkShift, req.ScheduleIDs); err != nil {
			utils.ResponseMessage(c, "Failed to register workshift: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Assign workshift successfully", http.StatusOK, nil)
	}
}
