package handler

import (
	DecisionTypeModel "erp/backend/internal/hrm/hr_profile/model"
	DecisionTypeUsecase "erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DecisionTypeHandler struct {
	decisionTypeUsecase DecisionTypeUsecase.DecisionTypeUsecase
}

func NewDecisionTypeHandler(decisionTypeUsecase DecisionTypeUsecase.DecisionTypeUsecase) *DecisionTypeHandler {
	return &DecisionTypeHandler{decisionTypeUsecase: decisionTypeUsecase}
}

// CreateDecisionType POST /decision-types
func (h *DecisionTypeHandler) CreateDecisionType() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data DecisionTypeModel.DecisionType

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		createdDecisionType, err := h.decisionTypeUsecase.CreateDecisionType(c.Request.Context(), &data)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo loại quyết định thành công", http.StatusOK, createdDecisionType)
	}
}

// UpdateDecisionType PUT /decision-types/:id
func (h *DecisionTypeHandler) UpdateDecisionType() gin.HandlerFunc {
	return func(c *gin.Context) {
		decisionTypeId := c.Param("id")
		var data DecisionTypeModel.DecisionType

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		updatedDecisionType, err := h.decisionTypeUsecase.UpdateDecisionType(c.Request.Context(), decisionTypeId, &data)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật loại quyết định thành công", http.StatusOK, updatedDecisionType)
	}
}

// DeleteDecisionType DELETE /decision-types/:id
func (h *DecisionTypeHandler) DeleteDecisionType() gin.HandlerFunc {
	return func(c *gin.Context) {
		decisionTypeId := c.Param("id")

		err := h.decisionTypeUsecase.DeleteDecisionType(c.Request.Context(), decisionTypeId)
		if err != nil {
			utils.ResponseMessage(c, "Xóa loại quyết định thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa loại quyết định thành công", http.StatusOK, nil)
	}
}

// GetDecisionTypes GET /decision-types
func (h *DecisionTypeHandler) GetDecisionTypes() gin.HandlerFunc {
	return func(c *gin.Context) {
		decisionTypes, err := h.decisionTypeUsecase.GetAllDecisionTypes(c.Request.Context())
		if err != nil {
			utils.ResponseMessage(c, "Lấy danh sách loại quyết định thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy danh sách loại quyết định thành công", http.StatusOK, decisionTypes)
	}
}

// GetDecisionTypeById GET /decision-types/:id
func (h *DecisionTypeHandler) GetDecisionTypeById() gin.HandlerFunc {
	return func(c *gin.Context) {
		decisionTypeId := c.Param("id")

		decisionType, err := h.decisionTypeUsecase.GetDecisionTypeByID(c.Request.Context(), decisionTypeId)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy loại quyết định", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy loại quyết định thành công", http.StatusOK, decisionType)
	}
}
