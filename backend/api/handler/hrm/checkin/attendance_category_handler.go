package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	checkinrepo "erp/backend/internal/hrm/checkin/repository"
	"erp/backend/internal/hrm/checkin/service"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type AttendanceCategoryHandler struct {
	biz service_interface.AttendanceCategoryService
}

func NewAttendanceCategoryHandler(db *gorm.DB) *AttendanceCategoryHandler {
	repo := checkinrepo.NewAttendanceCategoryRepository(db)
	biz := service.NewAttendanceCategoryService(repo)

	return &AttendanceCategoryHandler{
		biz: biz,
	}
}

func (h *AttendanceCategoryHandler) CreateAttendanceCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.AttendanceCategoryRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		attendanceCategory := model.ConvertToAttendanceCategory(req)

		employeeID, err := utils.ExtractFromContext[string](c, "employeeId")
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		attendanceCategory.CreatedBy = employeeID

		if err := attendanceCategory.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.biz.CreateAttendanceCategory(c, &attendanceCategory); err != nil {
			utils.ResponseMessage(c, "Failed to create attendance category: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance category created successfully", http.StatusOK, nil)
	}
}

func (h *AttendanceCategoryHandler) GetAttendanceCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		category, err := h.biz.GetAttendanceCategoryByID(c, id)
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

		categories, err := h.biz.ListAttendanceCategoriesByOffice(c, id)
		if err != nil {
			utils.ResponseMessage(c, "Attendance category not found", http.StatusNotFound, nil)
			return
		}
		utils.ResponseMessage(c, "Attendance category of Office", http.StatusOK, categories)
	}
}

func (h *AttendanceCategoryHandler) GetAllAttendanceCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.biz.ListAllAttendanceCategories(c)
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
		var req model.AttendanceCategoryRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid update data", http.StatusBadRequest, nil)
			fmt.Printf("error: %s", err.Error())
			return
		}

		categoryUpdate := model.ConvertToAttendanceCategory(req)
		categoryUpdate.AttendanceCategoryID = id

		if err := categoryUpdate.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.biz.UpdateAttendanceCategory(c, &categoryUpdate); err != nil {
			utils.ResponseMessage(c, "Failed to update attendance category", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance category updated successfully", http.StatusOK, nil)
	}
}

func (h *AttendanceCategoryHandler) DeleteAttendanceCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.biz.DeleteAttendanceCategory(c, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance category deleted successfully", http.StatusOK, nil)
	}
}
