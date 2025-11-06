package handler

import (
	"net/http"
	"strconv"

	reqdto "erp/backend/internal/hrm/recruitment/dto/request"
	dto "erp/backend/internal/hrm/recruitment/dto/response"
	svc "erp/backend/internal/hrm/recruitment/service/interface"
	utils "erp/backend/pkg"

	"github.com/gin-gonic/gin"
)

type ProcessStageHandler struct {
	biz svc.ProcessStageService
}

func NewProcessStageHandler(biz svc.ProcessStageService) *ProcessStageHandler {
	return &ProcessStageHandler{biz: biz}
}

// Request DTOs are defined in dto/request package and reused by handlers

func (h *ProcessStageHandler) CreateStage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reqdto.CreateStageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		ctx := c.Request.Context()

		result, err := h.biz.CreateStage(ctx, req.ProcessFormID, req.StageName)
		if err != nil {
			// Return the service error directly; do not try to classify specialized sentinel errors.
			utils.ResponseMessage(c, "Lỗi khi tạo vòng tuyển dụng mới: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		resp := dto.ProcessStageResponse{
			ProcessStageID: result.ProcessStageID,
			StageName:      result.StageName,
			StageOrder:     result.StageOrder,
		}

		utils.ResponseSuccess(c, "Tạo vòng tuyển dụng mới thành công", http.StatusCreated, resp)
	}
}

func (h *ProcessStageHandler) UpdateStage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req reqdto.UpdateStageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ResponseMessage(c, "Dữ liệu đầu vào không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		processStageIDStr := c.Param("id")
		processStageID, errConv := strconv.Atoi(processStageIDStr)
		if errConv != nil {
			utils.ResponseMessage(c, "process_stage_id không hợp lệ", http.StatusBadRequest, nil)
			return
		}
		ctx := c.Request.Context()

		result, err := h.biz.UpdateStage(ctx, processStageID, req.StageName)
		if err != nil {
			utils.ResponseMessage(c, "Lỗi khi cập nhật vòng tuyển dụng: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		resp := dto.ProcessStageResponse{
			ProcessStageID: result.ProcessStageID,
			StageName:      result.StageName,
			StageOrder:     result.StageOrder,
		}

		utils.ResponseSuccess(c, "Cập nhật vòng tuyển dụng thành công", http.StatusOK, resp)
	}
}

func (h *ProcessStageHandler) DeleteStage() gin.HandlerFunc {
	return func(c *gin.Context) {
		processStageIDStr := c.Param("id")
		processStageID, errConv := strconv.Atoi(processStageIDStr)
		if errConv != nil {
			utils.ResponseMessage(c, "process_stage_id không hợp lệ", http.StatusBadRequest, nil)
			return
		}
		ctx := c.Request.Context()

		if err := h.biz.DeleteStage(ctx, processStageID); err != nil {
			utils.ResponseMessage(c, "Lỗi khi xóa vòng tuyển dụng: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseSuccess(c, "Xóa vòng tuyển dụng thành công", http.StatusOK, nil)
	}
}
