package handler

import (
	"context"
	documentTypemodel "erp/backend/internal/hrm/hr_profile/model"
	documentTyperepository "erp/backend/internal/hrm/hr_profile/repository"
	documentTypeusecase "erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type DocumentTypeBiz interface {
	CreateDocumentType(ctx context.Context, data *documentTypemodel.EmployeeDocumentTypeCreate) error
	GetDocumentType(ctx context.Context, id string) (*documentTypemodel.EmployeeDocumentType, error)
	GetAllDocumentType(ctx context.Context) ([]*documentTypemodel.EmployeeDocumentType, error)
	UpdateDocumentType(ctx context.Context, id string, data *documentTypemodel.EmployeeDocumentTypeUpdate) error
	DeleteDocumentType(ctx context.Context, id string) error
	GetEnumDocumentType(ctx context.Context) ([]string, error)
}

type DocumentTypeHandler struct {
	documentTypeBiz DocumentTypeBiz
}

func NewDocumentTypeHandler(db *gorm.DB) *DocumentTypeHandler {
	repo := documentTyperepository.NewDocumentType(db)
	biz := documentTypeusecase.NewDocumentTypeBiz(repo)

	return &DocumentTypeHandler{
		documentTypeBiz: biz,
	}
}

// UpdateDocumentType PUT /document-types/:id
func (h *DocumentTypeHandler) UpdateDocumentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var data documentTypemodel.EmployeeDocumentTypeUpdate

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		if err := h.documentTypeBiz.UpdateDocumentType(c.Request.Context(), id, &data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật loại tài liệu thành công", http.StatusOK, nil)
	}
}

// DeleteDocumentType DELETE /document-types/:id
func (h *DocumentTypeHandler) DeleteDocumentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := h.documentTypeBiz.DeleteDocumentType(c.Request.Context(), id); err != nil {
			utils.ResponseMessage(c, "Xóa loại tài liệu thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa loại tài liệu thành công", http.StatusOK, nil)
	}
}

// GetAllDocumentType GET /document-types
func (h *DocumentTypeHandler) GetAllDocumentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		results, err := h.documentTypeBiz.GetAllDocumentType(c.Request.Context())
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy danh sách loại tài liệu", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy danh sách loại tài liệu thành công", http.StatusOK, results)
	}
}

// GetDocumentTypeById GET /document-types/:id
func (h *DocumentTypeHandler) GetDocumentTypeById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := h.documentTypeBiz.GetDocumentType(c.Request.Context(), id)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy loại tài liệu", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy loại tài liệu thành công", http.StatusOK, result)
	}
}

func (h *DocumentTypeHandler) CreateDocumentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data documentTypemodel.EmployeeDocumentTypeCreate

		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		if err := h.documentTypeBiz.CreateDocumentType(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, "Tạo loại tài liệu thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Tạo loại tài liệu thành công", http.StatusOK, nil)
	}
}

func (h *DocumentTypeHandler) GetEnumDocumentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.documentTypeBiz.GetEnumDocumentType(c.Request.Context())
		if err != nil {
			utils.ResponseMessage(c, "Không lấy được danh sách enum", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Lấy danh sách enum thành công", http.StatusOK, result)
	}
}
