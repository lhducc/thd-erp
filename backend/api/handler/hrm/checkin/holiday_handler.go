package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository"
	"erp/backend/internal/hrm/checkin/service"
	utils "erp/backend/pkg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HolidayBiz interface {
	CreateHoliday(w model.Holiday) error
	GetHolidayById(id string) (model.Holiday, error)
	GetAllHoliday() ([]model.Holiday, error)
	UpdateHoliday(id string, holiday model.Holiday) error
	DeleteHoliday(id string) error
}

type HolidayHandler struct {
	biz HolidayBiz
}

func NewHolidayHandler(db *gorm.DB) *HolidayHandler {
	repo := repository.NewHolidayRepo(db)
	biz := service.NewHolidayService(repo)

	return &HolidayHandler{
		biz: biz,
	}
}

func (h *HolidayHandler) CreateHoliday() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.Holiday
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu ngày nghỉ: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := data.Validate(); err != nil {
			fmt.Println("Dữ liệu đầu vào không hợp lệ: ", err)
		}

		if err := h.biz.CreateHoliday(data); err != nil {
			utils.ResponseMessage(c, "Tạo ngày nghỉ thất bại"+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo ngày nghỉ thành công", http.StatusOK, nil)
	}
}

// GetHoliday handles GET /holidays/:id
func (h *HolidayHandler) GetHoliday() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		result, err := h.biz.GetHolidayById(id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy ngày nghỉ", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Thông tin ngày nghỉ", http.StatusOK, result)
	}
}

// GetAllHoliday handles GET /holidays
func (h *HolidayHandler) GetAllHoliday() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.biz.GetAllHoliday()
		if err != nil {
			utils.ResponseMessage(c, "Lấy danh sách ngày nghỉ thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách ngày nghỉ", http.StatusOK, result)
	}
}

// UpdateHoliday handles PUT /holidays/:id
func (h *HolidayHandler) UpdateHoliday() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var data model.Holiday

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Dữ liệu không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.biz.UpdateHoliday(id, data); err != nil {
			utils.ResponseMessage(c, "Cập nhật ngày nghỉ thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật ngày nghỉ thành công", http.StatusOK, nil)
	}
}

// DeleteHoliday handles DELETE /holidays/:id
func (h *HolidayHandler) DeleteHoliday() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.biz.DeleteHoliday(id); err != nil {
			utils.ResponseMessage(c, "Xóa ngày nghỉ thất bại"+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa ngày nghỉ thành công", http.StatusOK, nil)
	}
}
