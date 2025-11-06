package handler

import (
	dto "erp/backend/internal/hrm/recruitment/dto/request"
	service "erp/backend/internal/hrm/recruitment/service/interface"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProcessFormHandler struct {
	service service.ProcessFormService
}

func NewProcessFormHandler(service service.ProcessFormService) *ProcessFormHandler {
	return &ProcessFormHandler{service: service}
}

func (h *ProcessFormHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req dto.UpdateProcessFormRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid request: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		req.ProcessFormID = id

		if err := utils.ValidateStruct(req); err != nil {
			utils.ResponseMessage(c, "Validation failed: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx := c.Request.Context()
		if err := h.service.UpdateForm(ctx, req); err != nil {
			utils.ResponseMessage(c, "Lỗi khi cập nhật mẫu quy trình: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật mẫu quy trình thành công", http.StatusOK, nil)
	}
}

func (h *ProcessFormHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()

		if err := h.service.DeleteForm(ctx, id); err != nil {
			utils.ResponseMessage(c, "Lỗi khi xóa mẫu quy trình: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa mẫu quy trình thành công", http.StatusOK, nil)
	}
}

func (h *ProcessFormHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		forms, err := h.service.GetAllForms(ctx)
		if err != nil {
			utils.ResponseMessage(c, "Lấy danh sách mẫu quy trình thất bại: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy danh sách mẫu quy trình thành công", http.StatusOK, &forms)
	}
}

func (h *ProcessFormHandler) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()

		form, err := h.service.GetFormById(ctx, id)
		if err != nil {
			utils.ResponseMessage(c, "Lấy mẫu quy trình thất bại: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy mẫu quy trình thành công", http.StatusOK, &form)
	}
}

func (h *ProcessFormHandler) CreateFormWithStages() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateProcessFormRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Invalid request: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := utils.ValidateStruct(req); err != nil {
			utils.ResponseMessage(c, "Validation failed: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		ctx := c.Request.Context()
		if err := h.service.CreateFormWithStages(ctx, req); err != nil {
			utils.ResponseMessage(c, "Lỗi khi tạo mẫu quy trình: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo mẫu quy trình thành công", http.StatusOK, nil)
	}
}
