package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Interface
type HierarchyLevelBiz interface {
	CreateHierarchyLevel(ctx context.Context, data *model.HierarchyLevelCreate) error
	GetHierarchyLevel(ctx context.Context, id string) (*model.HierarchyLevel, error)
	UpdateHierarchyLevel(ctx context.Context, id string, data *model.HierarchyLevelCreate) error
	DeleteHierarchyLevel(ctx context.Context, id string) error
	GetAllHierarchyLevel(ctx context.Context) ([]model.HierarchyLevel, error)
}

// Handler
type HierarchyLevelHandler struct {
	hierarchyLevelBiz HierarchyLevelBiz
}

func NewHierarchyLevelHandler(biz HierarchyLevelBiz) *HierarchyLevelHandler {
	return &HierarchyLevelHandler{
		hierarchyLevelBiz: biz,
	}
}

func (h *HierarchyLevelHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.HierarchyLevelCreate

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		if err := h.hierarchyLevelBiz.CreateHierarchyLevel(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo cấp bậc thành công", http.StatusCreated, gin.H{"data": data})
	}
}

func (h *HierarchyLevelHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := h.hierarchyLevelBiz.GetHierarchyLevel(c.Request.Context(), id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy cấp bậc", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *HierarchyLevelHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.hierarchyLevelBiz.GetAllHierarchyLevel(c)
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy danh sách vị trí", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *HierarchyLevelHandler) UpdateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		var data model.HierarchyLevelCreate
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Dữ liệu không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		if err := h.hierarchyLevelBiz.UpdateHierarchyLevel(c.Request.Context(), idParam, &data); err != nil {
			utils.ResponseMessage(c, "Cập nhật cấp bậc thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật thành công", http.StatusOK, nil)
	}
}

func (h *HierarchyLevelHandler) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := h.hierarchyLevelBiz.DeleteHierarchyLevel(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, "Xóa cấp bậc thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa cấp bậc thành công", http.StatusOK, nil)
	}
}
