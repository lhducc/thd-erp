package checkin

import (
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AttendanceCategoryHandler struct {
	biz service_interface.AttendanceCategoryService
}

func NewAttendanceCategoryHandler(biz service_interface.AttendanceCategoryService) *AttendanceCategoryHandler {
	return &AttendanceCategoryHandler{
		biz: biz,
	}
}

func (h *AttendanceCategoryHandler) CreateAttendanceCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req dto.AttendanceCategoryRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		attendanceCategory := dto.ConvertToAttendanceCategory(req)

		employeeID, err := utils.ExtractFromContext[string](c, "employeeId")
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		attendanceCategory.CreatedBy = employeeID

		if err := h.biz.CreateAttendanceCategory(ctx, &attendanceCategory); err != nil {
			utils.ResponseMessage(c, "Failed to create attendance category: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance category created successfully", http.StatusOK, nil)
	}
}

func (h *AttendanceCategoryHandler) GetAttendanceCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()

		category, err := h.biz.GetAttendanceCategoryByID(ctx, id)
		if err != nil {
			utils.ResponseMessage(c, "Attendance category not found", http.StatusNotFound, nil)
			return
		}
		utils.ResponseMessage(c, "Attendance category details", http.StatusOK, category)
	}
}

func (h *AttendanceCategoryHandler) GetAttendanceCategoryByOfficeID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("officeId")
		ctx := c.Request.Context()

		categories, err := h.biz.ListAttendanceCategoriesByOffice(ctx, id)
		if err != nil {
			utils.ResponseMessage(c, "Attendance category not found", http.StatusNotFound, nil)
			return
		}
		utils.ResponseMessage(c, "Attendance category of Office", http.StatusOK, categories)
	}
}

func (h *AttendanceCategoryHandler) GetAllAttendanceCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := h.biz.ListAllAttendanceCategories(ctx)
		if err != nil {
			utils.ResponseMessage(c, "Error getting attendance categories", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "List of attendance categories", http.StatusOK, result)
	}
}

func (h *AttendanceCategoryHandler) UpdateAttendanceCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()
		var req dto.AttendanceCategoryRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid update data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		categoryUpdate := dto.ConvertToAttendanceCategory(req)
		categoryUpdate.AttendanceCategoryID = id

		if err := h.biz.UpdateAttendanceCategory(ctx, &categoryUpdate); err != nil {
			utils.ResponseMessage(c, "Failed to update attendance category", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance category updated successfully", http.StatusOK, nil)
	}
}

func (h *AttendanceCategoryHandler) DeleteAttendanceCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()

		if err := h.biz.DeleteAttendanceCategory(ctx, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance category deleted successfully", http.StatusOK, nil)
	}
}
