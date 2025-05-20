package handler

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InsuranceHandler struct
type InsuranceHandler struct {
	insuranceUsecase usecase.InsuranceUsecase
}

// NewInsuranceHandler constructor InsuranceHandler
func NewInsuranceHandler(insuranceUsecase usecase.InsuranceUsecase) *InsuranceHandler {
	return &InsuranceHandler{insuranceUsecase: insuranceUsecase}
}

// POST /insurances
func (h *InsuranceHandler) CreateInsurance() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.Insurance

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		created, err := h.insuranceUsecase.CreateInsurance(c.Request.Context(), &data)
		if err != nil {
			utils.ResponseMessage(c, "Không thể tạo bảo hiểm", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Bảo hiểm đã được tạo thành công", http.StatusOK, created)
	}
}

// PUT /insurances/:id
func (h *InsuranceHandler) UpdateInsurance() gin.HandlerFunc {
	return func(c *gin.Context) {
		insuranceId := c.Param("id")
		var data model.Insurance

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		updated, err := h.insuranceUsecase.UpdateInsurance(c.Request.Context(), insuranceId, &data)
		if err != nil {
			utils.ResponseMessage(c, "Không thể cập nhật bảo hiểm", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Bảo hiểm đã được cập nhật thành công", http.StatusOK, updated)
	}
}

// DELETE /insurances/:id
func (h *InsuranceHandler) DeleteInsurance() gin.HandlerFunc {
	return func(c *gin.Context) {
		insuranceId := c.Param("id")

		err := h.insuranceUsecase.DeleteInsurance(c.Request.Context(), insuranceId)
		if err != nil {
			utils.ResponseMessage(c, "Không thể xóa bảo hiểm", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Bảo hiểm đã được xóa thành công", http.StatusOK, nil)
	}
}

// GET /insurances
func (h *InsuranceHandler) GetInsurances() gin.HandlerFunc {
	return func(c *gin.Context) {
		insurances, err := h.insuranceUsecase.GetAllInsurances(c.Request.Context())
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy danh sách bảo hiểm", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách bảo hiểm", http.StatusOK, insurances)
	}
}

// GET /insurances/:id
func (h *InsuranceHandler) GetInsuranceById() gin.HandlerFunc {
	return func(c *gin.Context) {
		insuranceId := c.Param("id")

		insurance, err := h.insuranceUsecase.GetInsuranceByID(c.Request.Context(), insuranceId)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy bảo hiểm", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Thông tin bảo hiểm", http.StatusOK, insurance)
	}
}
