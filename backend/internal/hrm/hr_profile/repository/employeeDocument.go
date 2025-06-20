package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
	"strings"
)

type EmployeeDocument struct {
	db *gorm.DB
}

func NewEmployeeDocument(db *gorm.DB) *EmployeeDocument {
	return &EmployeeDocument{db: db}
}

func (e *EmployeeDocument) CreateEmployeeDocument(ctx context.Context, employeeDocument *model.EmployeeDocument) error {
	return e.db.WithContext(ctx).Create(employeeDocument).Error
}

func (e *EmployeeDocument) UpdateEmployeeDocument(ctx context.Context, id string, employeeDocument *model.EmployeeDocument) error {
	return e.db.WithContext(ctx).
		Model(&model.EmployeeDocument{}).
		Where("document_id = ?", id).
		Updates(employeeDocument).Error
}

func (e *EmployeeDocument) DeleteEmployeeDocument(ctx context.Context, id string) error {
	return e.db.WithContext(ctx).
		Model(&model.EmployeeDocument{}).
		Where("id = ?", id).
		Delete(&model.EmployeeDocument{}).Error
}

func (e *EmployeeDocument) GetEmployeeDocumentById(ctx context.Context, id string) (*model.EmployeeDocument, error) {
	var employeeDocument model.EmployeeDocument
	if err := e.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("DocumentType").
		Preload("Employee").
		First(&employeeDocument).Error; err != nil {
		return nil, err
	}
	return &employeeDocument, nil
}

func (e *EmployeeDocument) GetEmployeeDocumentsPaginated(ctx context.Context, search, status, condition string, offset, limit int) ([]*model.EmployeeDocumentResponse, int64, error) {
	var docs []*model.EmployeeDocumentResponse
	var total int64

	query := e.db.WithContext(ctx).
		Model(&model.EmployeeDocument{}).
		Joins("JOIN employee ON employee.employee_id = employeedocument.employee_id").
		Joins("JOIN employeedocumenttype ON employeedocumenttype.id = employeedocument.document_type_id")
	query = query.Select(
		"employeedocument.document_id",
		"employee.full_name as employee_name",
		"employeedocumenttype.document_type_name as document_type",
		"employeedocument.effective_date",
		"employeedocument.expired_date",
		"employeedocument.condition",
		"employeedocument.employee_id",
	)

	if search != "" {
		query = query.Where("LOWER(employee.full_name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}
	if status != "" {
		query = query.Where("employeedocument.status = ?", status)
	}
	if condition != "" {
		query = query.Where("employeedocument.condition = ?", condition)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).
		Order("employeedocument.created_date DESC").Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

func (e *EmployeeDocument) GetAllEmployeeDocuments(ctx context.Context) ([]*model.EmployeeDocument, error) {
	var list []*model.EmployeeDocument
	if err := e.db.WithContext(ctx).
		Preload("DocumentType").
		Preload("Employee").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (e *EmployeeDocument) GetLastEmployeeDocumentByCode(ctx context.Context, employeeDocument *model.EmployeeDocument) error {
	return e.db.WithContext(ctx).
		Order("id DESC").
		First(employeeDocument).Error
}

func (e *EmployeeDocument) GetEmployeeDocumentsByEmployeeID(ctx context.Context, employeeID string) ([]*model.EmployeeDocument, error) {
	var documents []*model.EmployeeDocument
	if err := e.db.WithContext(ctx).
		Where("employee_id = ?", employeeID).
		Preload("DocumentType").
		Preload("Employee").
		Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (e *EmployeeDocument) GetEmployeeDocumentsByDocumentTypeID(ctx context.Context, documentTypeID string) ([]*model.EmployeeDocument, error) {
	var documents []*model.EmployeeDocument
	if err := e.db.WithContext(ctx).
		Where("document_type_id = ?", documentTypeID).
		Preload("DocumentType").
		Preload("Employee").
		Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (e *EmployeeDocument) GetLastDocumentByCode(ctx context.Context, document *model.EmployeeDocument) error {
	return e.db.WithContext(ctx).
		Order("document_id DESC").
		First(document).Error
}
