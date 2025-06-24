package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/repository"
	"strconv"
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
	GetAllDecisionPagination(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Decision, int64, error)
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
		page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		// Filter processing if necessary
		filterFields := []string{"decision_type_id", "condition"}
		filters := make(map[string]interface{})

		for _, field := range filterFields {
			values := ctx.QueryArray(field)
			cleaned := make([]string, 0)
			for _, v := range values {
				for _, part := range strings.Split(v, ",") {
					if trimmed := strings.TrimSpace(part); trimmed != "" {
						cleaned = append(cleaned, trimmed)
					}
				}
			}
			if len(cleaned) > 0 {
				filters[field] = cleaned
			}
		}

		decisions, totalRecords, err := h.decisionBiz.GetAllDecisionPagination(ctx, page, pageSize, filters)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		var decisionResponses []model.DecisionResponse
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
			decisionResponses = append(decisionResponses, response)
		}
		// Calculate the total number of pages
		totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

		utils.ResponseSuccess(ctx, "Danh sách dữ liệu", http.StatusOK, gin.H{
			"data":         decisionResponses,
			"totalRecords": totalRecords,
			"page":         page,
			"pageSize":     pageSize,
			"totalPages":   totalPages,
		})
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
