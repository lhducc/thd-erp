package handler

import (
	ContractTypeModel "erp/backend/internal/hrm/hr_profile/model"
	ContractTypeUsecase "erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContractTypeHandler struct {
	contractTypeUsecase *ContractTypeUsecase.ContractTypeUsecase
}

func NewContractTypeHandler(contractTypeUsecase *ContractTypeUsecase.ContractTypeUsecase) *ContractTypeHandler {
	return &ContractTypeHandler{contractTypeUsecase: contractTypeUsecase}
}

// CreateContractType POST /contract-types
func (h *ContractTypeHandler) CreateContractType() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data ContractTypeModel.ContractType

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		created, err := h.contractTypeUsecase.CreateContractType(c.Request.Context(), &data)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Thêm mới loại hợp đồng thành công", http.StatusOK, created)
	}
}

// UpdateContractType PUT /contract-types/:id
func (h *ContractTypeHandler) UpdateContractType() gin.HandlerFunc {
	return func(c *gin.Context) {
		contractTypeId := c.Param("id")
		var data ContractTypeModel.ContractType

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		updated, err := h.contractTypeUsecase.UpdateContractType(c.Request.Context(), contractTypeId, &data)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật loại hợp đồng thành công", http.StatusOK, updated)
	}
}

// DeleteContractType DELETE /contract-types/:id
func (h *ContractTypeHandler) DeleteContractType() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		contractType, err := h.contractTypeUsecase.DeleteContractType(c.Request.Context(), id)
		if err != nil {
			utils.ResponseMessage(c, "Xóa loại hợp đồng thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa loại hợp đồng thành công", http.StatusOK, contractType)
	}
}

// GetContractTypes GET /contract-types
func (h *ContractTypeHandler) GetContractTypes() gin.HandlerFunc {
	return func(c *gin.Context) {
		types, err := h.contractTypeUsecase.GetAllContractTypes(c.Request.Context())
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy danh sách loại hợp đồng", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy danh sách loại hợp đồng thành công", http.StatusOK, types)
	}
}

// GetContractTypeById GET /contract-types/:id
func (h *ContractTypeHandler) GetContractTypeById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := h.contractTypeUsecase.GetContractTypeById(c.Request.Context(), id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy loại hợp đồng", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy loại hợp đồng thành công", http.StatusOK, result)
	}
}
