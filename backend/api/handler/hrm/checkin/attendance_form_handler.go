package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository"
	"erp/backend/internal/hrm/checkin/service"
	utils "erp/backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttandanceFormBiz interface {
	CreateAttandanceForm(w model.AttendanceForm) error
	GetAttandanceFormById(id string) (model.AttendanceForm, error)
	GetAllAttandanceForm() ([]model.AttendanceForm, error)
	UpdateAttandanceForm(id string, holiday model.AttendanceForm) error
	DisableAttandanceForm(id string) error
}

type AttendanceFormHandler struct {
	biz AttandanceFormBiz
}

func NewAttendanceFormHandler(db *gorm.DB) *AttendanceFormHandler {
	repo := repository.NewAttendanceFormStore(db)
	biz := service.NewAttandanceFormService(repo)

	return &AttendanceFormHandler{
		biz: biz,
	}
}

func (h *AttendanceFormHandler) CreateAttendanceForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.AttendanceForm
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu bảng công: "+err.Error(), http.StatusBadRequest, nil)
			return
		}
		if err := data.Validate(); err != nil {
			utils.ResponseMessage(c, "Lỗi dữ liệu", http.StatusBadRequest, gin.H{
				"Chi tiết lỗi": err.Error(),
			})
			return
		}

		if err := h.biz.CreateAttandanceForm(data); err != nil {
			utils.ResponseMessage(c, "Tạo bảng công thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo bảng công thành công", http.StatusOK, nil)
	}
}

// GetAttendanceForm handles GET /attandance-form/:id
func (h *AttendanceFormHandler) GetAttendanceForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		result, err := h.biz.GetAttandanceFormById(id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy bảng công", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Thông tin bảng công", http.StatusOK, result)
	}
}

// GetAllAttendanceForm handles GET /attandance-form
func (h *AttendanceFormHandler) GetAllAttendanceForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.biz.GetAllAttandanceForm()
		if err != nil {
			utils.ResponseMessage(c, "Lấy danh sách bảng công thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách bảng công", http.StatusOK, result)
	}
}

// UpdateAttendanceForm handles PUT /attandance-form/:id
func (h *AttendanceFormHandler) UpdateAttendanceForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var data model.AttendanceForm

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Dữ liệu không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.biz.UpdateAttandanceForm(id, data); err != nil {
			utils.ResponseMessage(c, "Cập nhật bảng công thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật bảng công thành công", http.StatusOK, nil)
	}
}

// DisableAttendanceForm handles DELETE /attandance-form/:id
func (h *AttendanceFormHandler) DisableAttendanceForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.biz.DisableAttandanceForm(id); err != nil {
			utils.ResponseMessage(c, "Vô hiệu hóa bảng công thất bại "+ err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Vô hiệu hóa bảng công thành công", http.StatusOK, nil)
	}
}