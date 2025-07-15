package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type EmployeeWorkshiftBiz interface {
	Register(employeeID string, workShift string, date time.Time) error
	Delete(id uint) error
	GetAllByEmployeeID(employeeID string) ([]model.EmployeeWorkshift, error)
}

type EmployeeWorkshiftHandler struct {
	biz service_interface.EmployeeWorkshiftService
}

func NewEmployeeWorkshift(biz service_interface.EmployeeWorkshiftService) *EmployeeWorkshiftHandler {
	return &EmployeeWorkshiftHandler{
		biz: biz,
	}
}

func (h *EmployeeWorkshiftHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.EmployeeWorkshift

		// Parse incoming JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		// Create the workshift
		if err := h.biz.Register(req.EmployeeID, req.WorkShiftID, req.Date); err != nil {
			utils.ResponseMessage(c, "Failed to register workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Register workshift successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.biz.Delete(id); err != nil {
			utils.ResponseMessage(c, "Failed to register workshift: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Register workshift successfully", http.StatusOK, nil)
	}
}

func (h *EmployeeWorkshiftHandler) GetAllByEmployeeID() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employeeID")
		result, err := h.biz.GetByUserId(employeeID)
		// Create the workshift
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
