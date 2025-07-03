package handler

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AllowanceHandler struct {
	biz usecase.AllowanceRepo
}

func NewAllowanceHandler(biz usecase.AllowanceRepo) *AllowanceHandler {
	return &AllowanceHandler{biz: biz}
}

func (h *AllowanceHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.AllowanceCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.biz.Create(c.Request.Context(), &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.ResponseMessage(c, "Tạo phụ cấp thành công", http.StatusOK, nil)
	}
}

func (h *AllowanceHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := h.biz.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		utils.ResponseMessage(c, "Lấy danh sách phụ cấp thành công", http.StatusOK, list)
	}
}

func (h *AllowanceHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		result, err := h.biz.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phụ cấp"})
			return
		}
		utils.ResponseMessage(c, "Lấy chi tiết phụ cấp thành công", http.StatusOK, result)
	}
}

func (h *AllowanceHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req model.AllowanceCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.biz.Update(c.Request.Context(), id, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.ResponseMessage(c, "Cập nhật phụ cấp thành công", http.StatusOK, nil)
	}
}

func (h *AllowanceHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.biz.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.ResponseMessage(c, "Xóa phụ cấp thành công", http.StatusOK, nil)
	}
}
