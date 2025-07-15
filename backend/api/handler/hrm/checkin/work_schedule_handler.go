package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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

		var req model.WorkScheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		workSchedule := model.ConvertToWorkSchedule(&req)
		if err := workSchedule.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.service.CreateNewWorkSchedule(ctx, &workSchedule); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo mới lịch làm việc thành công", http.StatusCreated, nil)
	}
}

func (h *WorkScheduleHandler) AddEmployeeToWorkSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req model.AssignEmployeeRequest
		idSchedule := c.Param("work-schedule-id")
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

func (h *WorkScheduleHandler) UpdateWorkSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			utils.ResponseMessage(c, "ID không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		var req model.WorkSchedule
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

func (h *WorkScheduleHandler) DeleteWorkSchedule() gin.HandlerFunc {
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

func (h *WorkScheduleHandler) FetchListWorkSchedule() gin.HandlerFunc {
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
			utils.ResponseError(c, "Lỗi khi lấy dữ liệu", err, http.StatusBadRequest)
			return
		}
		utils.ResponseMessage(c, "Lấy dữ liệu thành công", http.StatusOK, &result)
	}
}

func (h *WorkScheduleHandler) ExportWorkSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := strings.Split(c.Query("fields"), ",")
		data, filename, err := h.service.ExportWorkSchedule(c, fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	}
}
