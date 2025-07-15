package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"fmt"

	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WorkShiftBiz interface {
	CreateWorkShift(data model.WorkShifts) error
	GetWorkShiftById(id string) (model.WorkShifts, error)
	GetAllWorkShift() ([]model.WorkShifts, error)
	UpdateWorkShift(id string, data model.WorkShifts) error
	DeleteWorkShift(id string) error
}

type WorkShiftHandler struct {
	biz service_interface.WorkShiftService
}

func NewWorkShiftHandler(biz service_interface.WorkShiftService) *WorkShiftHandler {
	return &WorkShiftHandler{
		biz: biz,
	}
}

func (h *WorkShiftHandler) CreateWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.WorkShiftsRequest

		// Parse incoming JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		// Convert request DTO to WorkShifts model
		workShifts := model.ConvertToWorkShifts(req)

		// Extract accountId from context
		employeeID, err := utils.ExtractFromContext[string](c, "employeeId")
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		workShifts.CreatedBy = employeeID

		// Validate workshift data before inserting
		if err := workShifts.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		// Create the workshift
		if err := h.biz.CreateWorkShiftService(&workShifts); err != nil {
			utils.ResponseMessage(c, "Failed to create workshift: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Workshift created successfully", http.StatusOK, nil)
	}
}

func (h *WorkShiftHandler) GetWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		workshift, err := h.biz.GetWorkShiftByIdService(id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy ca làm việc", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Chi tiết ca làm việc", http.StatusOK, workshift)
	}
}

func (h *WorkShiftHandler) GetAllWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.biz.GetAllWorkShiftService()
		if err != nil {
			utils.ResponseMessage(c, "Lỗi khi lấy danh sách ca làm việc", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách ca làm việc", http.StatusOK, result)
	}
}

func (h *WorkShiftHandler) UpdateWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req model.WorkShiftsRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu cập nhật không hợp lệ", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}
		// Convert request DTO to WorkShifts model
		workShiftUpdate := model.ConvertToWorkShifts(req)

		// Validate workshift data before inserting
		if err := workShiftUpdate.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.biz.UpdateWorkShiftService(id, &workShiftUpdate); err != nil {
			utils.ResponseMessage(c, "Cập nhật ca làm việc thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật ca làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkShiftHandler) DeleteWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.biz.DeleteWorkShiftService(id); err != nil {
			utils.ResponseMessage(c, "Xóa ca làm việc thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa ca làm việc thành công", http.StatusOK, nil)
	}
}
