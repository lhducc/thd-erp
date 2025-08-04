package checkin

import (
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WorkScheduleHandler struct {
	service service_interface.WorkScheduleServiceInterface
}

func NewWorkScheduleHandler(service service_interface.WorkScheduleServiceInterface) *WorkScheduleHandler {
	return &WorkScheduleHandler{service: service}
}

func (h *WorkScheduleHandler) CreateWorkSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req dto.WorkScheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		workSchedule := dto.ConvertToWorkSchedule(&req)

		if workSchedule.RepeatCycle != nil && workSchedule.RepeatType != nil {
			workSchedule.IsScheduleAuto = true
		} else {
			workSchedule.IsScheduleAuto = false
		}

		if err := h.service.CreateNewWorkSchedule(ctx, &workSchedule); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo mới lịch làm việc thành công", http.StatusCreated, nil)
	}
}

func (h *WorkScheduleHandler) AddManagerToWorkScheduleAuto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req dto.AssignManagersRequest

		idSchedule := c.Param("work-schedule-id")
		idConvert, err := strconv.Atoi(idSchedule)
		if err != nil {
			utils.ResponseMessage(c, "ID lịch không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.service.AssignEmployeeToWorkScheduleAuto(ctx, &req, idConvert); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Gán quản lý thành công", http.StatusOK, nil)
	}
}
func (h *WorkScheduleHandler) AddManagerToWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req dto.AssignManagersRequest

		idSchedule := c.Param("work-schedule-id")
		idConvert, err := strconv.Atoi(idSchedule)
		if err != nil {
			utils.ResponseMessage(c, "ID lịch không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.service.AssignEmployeeToWorkScheduleRegister(ctx, &req, idConvert); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Gán quản lý thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleHandler) UpdateWorkScheduleAuto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		var req dto.WorkScheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		wSchedule := dto.ConvertToWorkSchedule(&req)

		if err := h.service.UpdateWorkScheduleAuto(ctx, &wSchedule, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật lịch làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleHandler) UpdateWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		var req dto.WorkScheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		wSchedule := dto.ConvertToWorkSchedule(&req)

		if err := h.service.UpdateWorkScheduleRegister(ctx, &wSchedule, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật lịch làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleHandler) DeleteWorkScheduleAuto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.service.DeleteWorkScheduleAuto(ctx, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa lịch làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleHandler) DeleteWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.service.DeleteWorkScheduleRegister(ctx, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa lịch làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleHandler) FetchListWorkScheduleAuto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := h.service.GetAllWorkScheduleAuto(ctx)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Danh sách lịch làm việc", http.StatusOK, &result)
	}
}

func (h *WorkScheduleHandler) FetchListWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := h.service.GetAllWorkScheduleRegister(ctx)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Danh sách lịch làm việc", http.StatusOK, &result)
	}
}

func (h *WorkScheduleHandler) FetchWorkScheduleByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseError(c, "Id không hợp lệ", err, http.StatusBadRequest)
			return
		}
		result, err := h.service.GetWorkScheduleByID(ctx, id)
		if err != nil {
			utils.ResponseError(c, "Lỗi khi lấy dữ liệu", err, http.StatusInternalServerError)
			return
		}
		utils.ResponseMessage(c, "Lấy dữ liệu thành công", http.StatusOK, &result)
	}
}

func (h *WorkScheduleHandler) DeleteManagerFromWorkScheduleAuto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		workScheduleIDStr := c.Param("schedule-id")
		workScheduleID, err := strconv.Atoi(workScheduleIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Mã lịch làm việc không hợp lệ"})
			return
		}

		employeeID := c.Query("manager_id")
		if employeeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Mã quản lý cần xóa khỏi lịch không được trống"})
			return
		}

		err = h.service.DeleteManagerFromWorkScheduleAuto(ctx, employeeID, workScheduleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Quyền quản lý đã xóa khỏi lịch"})
	}
}

func (h *WorkScheduleHandler) DeleteManagerFromWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		workScheduleIDStr := c.Param("schedule-id")
		workScheduleID, err := strconv.Atoi(workScheduleIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Mã lịch làm việc không hợp lệ"})
			return
		}

		employeeID := c.Query("manager_id")
		if employeeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Mã quản lý cần xóa khỏi lịch không được trống"})
			return
		}

		err = h.service.DeleteManagerFromWorkScheduleRegister(ctx, employeeID, workScheduleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Quyền quản lý đã xóa khỏi lịch"})
	}
}

func (h *WorkScheduleHandler) UpdateStatusAutoRecurring() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		scheduleID, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		var req dto.AutoRecurringRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.service.UpdateStatusRecuringSchedule(ctx, scheduleID, req.IsAutoRecurring); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật trạn thái lặp thành công", http.StatusOK, nil)
	}
}
