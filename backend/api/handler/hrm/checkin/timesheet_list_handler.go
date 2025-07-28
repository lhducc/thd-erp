package checkin

import (
	model "erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/pkg"
	"erp/backend/pkg/struct_support"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TimesheetListHandler struct {
	service service_interface.TimesheetListServiceInterface
}

func NewTimesheetListHandler(svc service_interface.TimesheetListServiceInterface) *TimesheetListHandler {
	return &TimesheetListHandler{service: svc}
}

func (h *TimesheetListHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TimesheetDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input", http.StatusBadRequest, nil)
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		model := req.ConvertToBusinessModel()
		model.CreatedBy = c.GetString("employeeId")

		if err := h.service.CreateElementOfTimesheetList(c.Request.Context(), model); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo timesheet_list thành công", http.StatusCreated, model.TimeSheetListID)
	}
}

func (h *TimesheetListHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TimesheetDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input", http.StatusBadRequest, nil)
			return
		}

		model := req.ConvertToBusinessModel()
		model.TimeSheetListID = c.Param("id")
		updateBy := c.GetString("employeeId")
		model.UpdatedBy = &updateBy

		if err := h.service.Update(c.Request.Context(), model); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật timesheet_list thành công", http.StatusOK, model.TimeSheetListID)
	}
}

func (h *TimesheetListHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		data, err := h.service.GetByID(c.Request.Context(), id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy timesheet", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy dữ liệu thành công", http.StatusOK, data)
	}
}

func (h *TimesheetListHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}

		data, total, err := h.service.List(c.Request.Context(), page, limit)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		totalPages := int((total + int64(limit) - 1) / int64(limit))
		res := struct_support.PaginatedResponse{
			Data:       data,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			Total:      total,
		}

		utils.ResponseMessage(c, "Lấy danh sách thành công", http.StatusOK, &res)
	}
}

func (h *TimesheetListHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.service.Delete(c.Request.Context(), id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xoá timesheet_list thành công", http.StatusNoContent, nil)
	}
}

func (h *TimesheetListHandler) LockedTimeSheet() gin.HandlerFunc {
	return func(c *gin.Context) {

		lockeModel := &model.TimeSheetList{}
		lockeModel.TimeSheetListID = c.Param("id")
		employeeID := c.GetString("employeeId")
		lockeModel.LockedBy = &employeeID

		if err := h.service.LockedTimesheet(c.Request.Context(), lockeModel); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Chốt công thành công", http.StatusOK, lockeModel.TimeSheetListID)
	}
}
