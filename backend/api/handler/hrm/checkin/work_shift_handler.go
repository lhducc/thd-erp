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

type WorkShiftBiz interface {
	CreateWorkShift(data model.WorkShifts) error
	GetWorkShiftById(id string) (model.WorkShifts, error)
	GetAllWorkShift() ([]model.WorkShifts, error)
	UpdateWorkShift(id string, data model.WorkShifts) error
	DeleteWorkShift(id string) error
}

type WorkShiftHandler struct {
	biz WorkShiftBiz
}

func NewWorkShiftHandler(db *gorm.DB) *WorkShiftHandler {
	repo := repository.NewWorkShiftStore(db)
	biz := service.NewWorkShiftService(repo)

	return &WorkShiftHandler{biz: biz}
}

func (h *WorkShiftHandler) CreateWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.WorkShifts
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := req.Validate(); err != nil {
			fmt.Println("Dữ liệu đầu vào không hợp lệ: ", err)
		}

		if err := h.biz.CreateWorkShift(req); err != nil {
			utils.ResponseMessage(c, "Tạo ca làm việc thất bại"+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo ca làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkShiftHandler) GetWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		workshift, err := h.biz.GetWorkShiftById(id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy ca làm việc", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Chi tiết ca làm việc", http.StatusOK, workshift)
	}
}

func (h *WorkShiftHandler) GetAllWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.biz.GetAllWorkShift()
		if err != nil {
			utils.ResponseMessage(c, "Lỗi khi lấy danh sách ca làm việc", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách ca làm việc", http.StatusOK, result)
	}
}

func (h *WorkShiftHandler) UpdateWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req model.WorkShifts

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu cập nhật không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.biz.UpdateWorkShift(id, req); err != nil {
			utils.ResponseMessage(c, "Cập nhật ca làm việc thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật ca làm việc thành công", http.StatusOK, nil)
	}
}

func (h *WorkShiftHandler) DeleteWorkShift() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.biz.DeleteWorkShift(id); err != nil {
			utils.ResponseMessage(c, "Xóa ca làm việc thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa ca làm việc thành công", http.StatusOK, nil)
	}
}
