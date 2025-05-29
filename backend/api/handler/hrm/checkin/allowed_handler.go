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

type AllowedWorkingScheduleBiz interface {
	CreateAllowedWorkingSchedule(w model.AllowedWorkingSchedule) error
	GetAllowedWorkingScheduleById(id string) (model.AllowedWorkingSchedule, error)
	GetAllAllowedWorkingSchedule() ([]model.AllowedWorkingSchedule, error)
	UpdateAllowedWorkingSchedule(id string, holiday model.AllowedWorkingSchedule) error
	DeleteAllowedWorkingScheduleService(id string) error
}

type AllowedWorkingScheduleHandler struct {
	biz AllowedWorkingScheduleBiz
}

func NewAllowedWorkingScheduleHandler(db *gorm.DB) *AllowedWorkingScheduleHandler {
	repo := repository.NewAllowedWorkingSchedule(db)
	biz := service.NewAllowedWorkingScheduleService(repo)

	return &AllowedWorkingScheduleHandler{
		biz: biz,
	}
}

func (h *AllowedWorkingScheduleHandler) CreateAllowedWorkingSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.AllowedWorkingSchedule
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu" + err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := data.Validate(); err != nil {
			fmt.Println("Dữ liệu đầu vào không hợp lệ: ", err)
		}

		if err := h.biz.CreateAllowedWorkingSchedule(data); err != nil {
			utils.ResponseMessage(c, "Tạo lịch làm việc thất bại" + err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo lịch làm việc thành công", http.StatusOK, nil)
	}
}

func (h *AllowedWorkingScheduleHandler) GetAllowedWorkingSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		result, err := h.biz.GetAllowedWorkingScheduleById(id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy lịch làm việc", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Thông tin lịch làm việc", http.StatusOK, result)
	}
}

func (h *AllowedWorkingScheduleHandler) GetAllAllowedWorkingSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.biz.GetAllAllowedWorkingSchedule()
		if err != nil {
			utils.ResponseMessage(c, "Lấy danh sách thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách lịch làm việc", http.StatusOK, result)
	}
}

func (h *AllowedWorkingScheduleHandler) UpdateAllowedWorkingSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var data model.AllowedWorkingSchedule
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Dữ liệu không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.biz.UpdateAllowedWorkingSchedule(id, data); err != nil {
			utils.ResponseMessage(c, "Cập nhật thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật thành công", http.StatusOK, nil)
	}
}

func (h *AllowedWorkingScheduleHandler) DeleteAllowedWorkingSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.biz.DeleteAllowedWorkingScheduleService(id); err != nil {
			utils.ResponseMessage(c, "Xóa thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa thành công", http.StatusOK, nil)
	}
}
