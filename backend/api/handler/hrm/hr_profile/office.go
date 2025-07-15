package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OficeBiz interface {
	CreateOffice(context context.Context, data *model.OfficeCreate) error
	GetOffice(ctx context.Context, id *string) (*model.Office, error)
	GetAllOffice(ctx context.Context) ([]model.Office, error)
	UpdateOffice(ctx context.Context, id string, data *model.OfficeCreate) error
	DeleteOffice(ctx context.Context, id string) error
}

type OfficeHandler struct {
	officeBiz OficeBiz
}

func NewOficeHandler(biz OficeBiz) *OfficeHandler {

	return &OfficeHandler{
		officeBiz: biz,
	}
}

func (h *OfficeHandler) CreateOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.OfficeCreate

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, "Xóa bộ phận thất bại ", http.StatusBadRequest, nil)

			return
		}

		if err := h.officeBiz.CreateOffice(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, "Tạo văn phòng thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo văn phòng thành công", http.StatusCreated, nil)
	}
}

// GetOffice xử lý request lấy Office theo ID
func (h *OfficeHandler) GetOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		result, err := h.officeBiz.GetOffice(c.Request.Context(), &idParam)
		if err != nil {
			utils.ResponseMessage(c, "Lấy thông tin thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *OfficeHandler) GetAllOffice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := h.officeBiz.GetAllOffice(ctx)
		if err != nil {
			utils.ResponseMessage(ctx, "Lấy thông tin thất bại", http.StatusNotFound, nil)
		}
		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *OfficeHandler) UpdateOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var data model.OfficeCreate
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Cập nhật văn phòng thất bại", http.StatusBadRequest, nil)
			return
		}

		if err := h.officeBiz.UpdateOffice(c.Request.Context(), idParam, &data); err != nil {
			utils.ResponseMessage(c, "Cập nhật văn phòng thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật văn phòng thành công", http.StatusOK, nil)
	}
}

// DeleteOffice xử lý request xóa Office
func (h *OfficeHandler) DeleteOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := h.officeBiz.DeleteOffice(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, "Xóa văn phòng thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa văn phòng thành công", http.StatusOK, nil)
	}
}
