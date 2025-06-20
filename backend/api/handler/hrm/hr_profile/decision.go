package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/repository"
	"strings"

	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DecisionBiz interface {
	CreateDecision(ctx context.Context, data *model.DecisionCreate) (string, error)
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

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("%+v", data), http.StatusBadRequest, data)
			return
		}
		code, err := h.decisionBiz.CreateDecision(c.Request.Context(), &data)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Tạo thành công", http.StatusOK, gin.H{
			"decision_id": code,
		})

	}
}

func (h *DecisionHandler) GetDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		d, err := h.decisionBiz.GetDecision(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		response := model.DecisionResponse{
			DecisionID:       d.DecisionID,
			DecisionName:     d.DecisionName,
			DecisionTypeID:   d.DecisionTypeID,
			DecisionTypeName: "",
			EffectiveDate:    d.EffectiveDate.Format("2006-01-02"),
			SignDate:         d.SignDate.Format("2006-01-02"),
			Condition:        d.Condition,
			Content:          d.Content,
			AttachedFile:     d.AttachedFile,
			CreatedDate:      d.CreatedDate,
		}
		response.Employees = make([]model.EmployeeShort, 0)
		for _, emp := range d.Employees {
			response.Employees = append(response.Employees, model.EmployeeShort{
				EmployeeID: emp.EmployeeID,
				Fullname:   emp.Fullname,
			})
		}
		if d.DecisionType != nil {
			response.DecisionTypeName = d.DecisionType.DecisionType
		}

		utils.ResponseMessage(c, "Thông tin quyết định", http.StatusOK, response)
	}
}

func (h *DecisionHandler) GetAllDecision() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		decisions, err := h.decisionBiz.GetAllDecision(ctx)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		var responses []model.DecisionResponse
		for _, d := range decisions {
			response := model.DecisionResponse{
				DecisionID:       d.DecisionID,
				DecisionName:     d.DecisionName,
				DecisionTypeID:   d.DecisionTypeID,
				DecisionTypeName: "",
				EffectiveDate:    d.EffectiveDate.Format("2006-01-02"),
				SignDate:         d.SignDate.Format("2006-01-02"),
				Condition:        d.Condition,
				Content:          d.Content,
				AttachedFile:     d.AttachedFile,
				CreatedDate:      d.CreatedDate,
			}

			response.Employees = make([]model.EmployeeShort, 0)
			for _, emp := range d.Employees {
				response.Employees = append(response.Employees, model.EmployeeShort{
					EmployeeID: emp.EmployeeID,
					Fullname:   emp.Fullname,
				})
			}
			if d.DecisionType != nil {
				response.DecisionTypeName = d.DecisionType.DecisionType
			}

			responses = append(responses, response)
		}

		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, responses)
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
		fieldStr := c.Query("fields")
		if strings.TrimSpace(fieldStr) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tham số 'fields'"})
			return
		}

		fields := strings.Split(fieldStr, ",")
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
