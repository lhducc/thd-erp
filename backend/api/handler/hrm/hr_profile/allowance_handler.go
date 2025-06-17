package handler

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AllowanceHandler struct {
	biz usecase.AllowanceUsecase
}

func NewAllowanceHandler(biz usecase.AllowanceUsecase) *AllowanceHandler {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
