package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WorkScheduleRegisterHandler struct {
	service service_interface.WorkScheduleRegisterServiceInterface
}

func NewWorkScheduleRegisterHandler(service service_interface.WorkScheduleRegisterServiceInterface) *WorkScheduleRegisterHandler {
	return &WorkScheduleRegisterHandler{service: service}
}

func (h *WorkScheduleRegisterHandler) CreateWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req model.WorkScheduleRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		workSchedule := model.ConvertToWorkScheduleRegister(&req)
		if err := workSchedule.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.service.CreateNewWorkSchedule(ctx, &workSchedule); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo mới lịch làm việc đăng ký thành công", http.StatusCreated, nil)
	}
}

func (h *WorkScheduleRegisterHandler) AddEmployeeToWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req model.AssignEmployeeRequest
		idSchedule := c.Param("id")
		idConvert, err := strconv.Atoi(idSchedule)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		err = h.service.AssignEmployeeToWorkSchedule(ctx, &req, idConvert)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Gán nhân viên thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleRegisterHandler) UpdateWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		var req model.WorkScheduleRegister
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.service.UpdateWorkSchedule(ctx, &req, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật lịch làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleRegisterHandler) DeleteWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.service.DeleteWorkSchedule(ctx, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa lịch làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkScheduleRegisterHandler) FetchListWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := h.service.GetAllWorkSchedule(ctx)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Danh sách lịch làm việc", http.StatusOK, &result)
	}
}

func (h *WorkScheduleRegisterHandler) FetchWorkScheduleRegisterByID() gin.HandlerFunc {
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
			utils.ResponseError(c, "Lỗi khi lấy dữ liệu", err, http.StatusBadRequest)
			return
		}
		utils.ResponseMessage(c, "Lấy dữ liệu thành công", http.StatusOK, &result)
	}
}

func (h *WorkScheduleRegisterHandler) DeleteEmployeeFromWorkScheduleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		workScheduleIDStr := c.Param("schedule-id")
		workScheduleID, err := strconv.Atoi(workScheduleIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Mã lịch làm việc không hợp lệ"})
			return
		}

		employeeID := c.Query("employee_id")
		if employeeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Mã nhân viên cần xóa khỏi lịch không được trống"})
			return
		}

		err = h.service.DeleteEmployeeFromWorkScheduleRegister(ctx, employeeID, workScheduleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Xóa nhân viên khỏi lịch làm việc thành công"})
	}
}

func (h *WorkScheduleRegisterHandler) DeleteManagerFromWorkScheduleRegister() gin.HandlerFunc {
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
