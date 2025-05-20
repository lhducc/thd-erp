package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	hrmrepository "erp/backend/internal/hrm/hr_profile/repository"
	hrmbiz "erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type EmployeeDocumentBiz interface {
	CreateEmployeeDocument(ctx context.Context, data *model.EmployeeDocument) error
	GetEmployeeDocument(ctx context.Context, id string) (*model.EmployeeDocument, error)
	GetAllEmployeeDocuments(ctx context.Context) ([]*model.EmployeeDocument, error)
	UpdateEmployeeDocument(ctx context.Context, id string, data *model.EmployeeDocument) error
	DeleteEmployeeDocument(ctx context.Context, id string) error
}

type EmployeeDocumentHandler struct {
	employeeDocumentBiz EmployeeDocumentBiz
}

func NewEmployeeDocumentHandler(db *gorm.DB) *EmployeeDocumentHandler {
	repo := hrmrepository.NewEmployeeDocument(db)
	biz := hrmbiz.NewEmployeeDocumentBiz(repo)

	return &EmployeeDocumentHandler{
		employeeDocumentBiz: biz,
	}
}

// CreateEmployeeDocument POST /employee-documents
func (h *EmployeeDocumentHandler) CreateEmployeeDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.EmployeeDocument

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		if err := h.employeeDocumentBiz.CreateEmployeeDocument(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo tài liệu nhân viên thành công", http.StatusOK, nil)
	}
}

// UpdateEmployeeDocument PUT /employee-documents/:id
func (h *EmployeeDocumentHandler) UpdateEmployeeDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var data model.EmployeeDocument

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		if err := h.employeeDocumentBiz.UpdateEmployeeDocument(c.Request.Context(), id, &data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật tài liệu nhân viên thành công", http.StatusOK, nil)
	}
}

// DeleteEmployeeDocument DELETE /employee-documents/:id
func (h *EmployeeDocumentHandler) DeleteEmployeeDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.employeeDocumentBiz.DeleteEmployeeDocument(c.Request.Context(), id); err != nil {
			utils.ResponseMessage(c, "Xóa tài liệu nhân viên thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa tài liệu nhân viên thành công", http.StatusOK, nil)
	}
}

// GetAllEmployeeDocuments GET /employee-documents
func (h *EmployeeDocumentHandler) GetAllEmployeeDocuments() gin.HandlerFunc {
	return func(c *gin.Context) {
		results, err := h.employeeDocumentBiz.GetAllEmployeeDocuments(c.Request.Context())
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy danh sách tài liệu nhân viên", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy danh sách tài liệu nhân viên thành công", http.StatusOK, results)
	}
}

// GetEmployeeDocumentById GET /employee-documents/:id
func (h *EmployeeDocumentHandler) GetEmployeeDocumentById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := h.employeeDocumentBiz.GetEmployeeDocument(c.Request.Context(), id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy tài liệu nhân viên", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy tài liệu nhân viên thành công", http.StatusOK, result)
	}
}
