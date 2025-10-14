package checkin

import (
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"fmt"

	utils "erp/backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		var req dto.WorkShiftsRequest
		ctx := c.Request.Context()

		// Parse incoming JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		// Convert request DTO to WorkShifts model
		workShifts := req.ConvertToWorkShifts()

		// Extract accountId from context
		employeeID, err := utils.ExtractFromContext[string](c, "employeeId")
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		workShifts.CreatedBy = employeeID

		if err := h.biz.CreateWorkShift(ctx, workShifts); err != nil {
			utils.ResponseMessage(c, "Lỗi khi tạo mới ca: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Ca đã được tạo thành công", http.StatusOK, nil)
	}
}

func (h *WorkShiftHandler) GetWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()
		workshift, err := h.biz.GetWorkShiftById(ctx, id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy ca làm việc", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Chi tiết ca làm việc", http.StatusOK, workshift)
	}
}

func (h *WorkShiftHandler) GetAllWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := h.biz.GetAllWorkShift(ctx)
		if err != nil {
			utils.ResponseMessage(c, "Lỗi khi lấy danh sách ca làm việc: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách ca làm việc", http.StatusOK, result)
	}
}

func (h *WorkShiftHandler) UpdateWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()
		var req dto.WorkShiftsRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu cập nhật không hợp lệ", http.StatusBadRequest, nil)
			fmt.Printf("Lỗi: %s", err.Error())
			return
		}
		req.WorkShiftID = id

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}
		// Convert request DTO to WorkShifts model
		workShiftUpdate := req.ConvertToWorkShifts()

		if err := h.biz.UpdateWorkShift(ctx, workShiftUpdate); err != nil {
			utils.ResponseMessage(c, "Lỗi khi update ca: "+err.Error(), http.StatusInternalServerError, nil)
			fmt.Printf(err.Error())
			return
		}

		utils.ResponseMessage(c, "Cập nhật ca làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkShiftHandler) DeleteWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		id := c.Param("id")

		if err := h.biz.DeleteWorkShift(ctx, id); err != nil {
			utils.ResponseMessage(c, "Xóa ca làm việc thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa ca làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkShiftHandler) GetWorkshiftForRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		employeeID := c.GetString("employeeId")
		result, err := h.biz.GetListShiftForRegister(ctx, employeeID)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		fmt.Println(result)
		utils.ResponseMessage(c, "Danh sách ca làm việc được phép đăng ký của nhân viên", http.StatusOK, &result)
	}
}

func (h *WorkShiftHandler) GetWorkshiftInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		employeeID := c.GetString("employeeId")
		result, err := h.biz.GetWorkshiftInfo(ctx, employeeID)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		fmt.Println(result)
		utils.ResponseMessage(c, "Thông tin lịch làm việc của nhân viên", http.StatusOK, &result)
	}
}
