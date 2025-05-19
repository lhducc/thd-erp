package handler

import (
	"context"
	"strings"

	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DecisionBiz interface {
	CreateDecision(context context.Context, data *model.DecisionCreate) error
	GetDecision(ctx context.Context, id string) (*model.Decision, error)
	GetAllDecision(ctx context.Context) ([]model.Decision, error)
	UpdateDecision(ctx context.Context, id string, data *model.DecisionCreate) error
	DeleteDecision(ctx context.Context, id string) error
	ExportDecisionTest(ctx context.Context, selectedFields []string) ([]byte, string, error)
}

type DecisionHandler struct {
	decisionBiz DecisionBiz
	employeeBiz usecase.EmployeeRepo
}

func NewDecisionHandler(db *gorm.DB) *DecisionHandler {
	storeJob := repository.NewDicisionStore(db)
	storeUser := repository.NewUserStore(db)
	biz := usecase.NewDecisionBiz(storeJob, storeUser)

	return &DecisionHandler{
		decisionBiz: biz,
		employeeBiz: storeUser,
	}
}

func (h *DecisionHandler) CreateDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.DecisionCreate

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		if err := h.decisionBiz.CreateDecision(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo thành công", http.StatusOK, nil)
	}
}

func (h *DecisionHandler) GetDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		result, err := h.decisionBiz.GetDecision(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *DecisionHandler) GetAllDecision() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := h.decisionBiz.GetAllDecision(ctx)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}
		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *DecisionHandler) UpdateDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var data model.DecisionCreate
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		if err := h.decisionBiz.UpdateDecision(c.Request.Context(), idParam, &data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Đã cập nhật", http.StatusOK, nil)
	}
}

func (h *DecisionHandler) DeleteDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := h.decisionBiz.DeleteDecision(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Đã xóa", http.StatusOK, nil)
	}
}

func (biz *DecisionHandler) ExportDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := strings.Split(c.Query("fields"), ",")
		data, filename, err := biz.decisionBiz.ExportDecisionTest(c, fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	}
}
