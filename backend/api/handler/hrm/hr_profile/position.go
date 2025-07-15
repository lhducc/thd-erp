package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Interface cho biz layer
type PositionBiz interface {
	CreatePosition(ctx context.Context, data *model.Position) error
	GetPosition(ctx context.Context, id string) (*model.Position, error)
	GetAllPositions(ctx context.Context) ([]model.Position, error)
	UpdatePosition(ctx context.Context, id string, data *model.Position) error
	DeletePosition(ctx context.Context, id string) error
}

type PositionHandler struct {
	positionBiz PositionBiz
}

func NewPositionHandler(biz PositionBiz) *PositionHandler {
	return &PositionHandler{
		positionBiz: biz,
	}
}

func (h *PositionHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.Position

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, "Dữ liệu không hợp lệ", http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.positionBiz.CreatePosition(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, "Không thể tạo vị trí", http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, data)
	}
}

func (h *PositionHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		result, err := h.positionBiz.GetPosition(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy vị trí", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *PositionHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.positionBiz.GetAllPositions(c)
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy danh sách vị trí", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy danh sách thành công", http.StatusOK, gin.H{"data": result})
	}
}

func (h *PositionHandler) UpdateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		var data model.Position
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Dữ liệu không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.positionBiz.UpdatePosition(c.Request.Context(), idParam, &data); err != nil {
			utils.ResponseMessage(c, "Cập nhật thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật thành công", http.StatusOK, nil)
	}
}

func (h *PositionHandler) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		if err := h.positionBiz.DeletePosition(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, "Xoá thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xoá thành công", http.StatusOK, nil)
	}
}
