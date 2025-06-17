package repository

import (
	"context"
	documentTypemodel "erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

type DocumentType struct {
	db *gorm.DB
}

func NewDocumentType(db *gorm.DB) *DocumentType {
	return &DocumentType{db: db}
}

func (d *DocumentType) CreateDocumentType(ctx context.Context, documentType *documentTypemodel.EmployeeDocumentTypeCreate) error {
	return d.db.Model(&documentTypemodel.EmployeeDocumentType{}).WithContext(ctx).Create(documentType).Error
}

func (d *DocumentType) UpdateDocumentType(ctx context.Context, id string, documentType *documentTypemodel.EmployeeDocumentTypeUpdate) error {
	return d.db.WithContext(ctx).
		Model(&documentTypemodel.EmployeeDocumentType{}).
		Where("id = ?", id).
		Updates(documentType).Error
}

func (d *DocumentType) DeleteDocumentType(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).
		Model(&documentTypemodel.EmployeeDocumentType{}).
		Where("id = ?", id).
		Delete(&documentTypemodel.EmployeeDocumentType{}).Error
}

func (d *DocumentType) GetDocumentTypeById(ctx context.Context, id string) (*documentTypemodel.EmployeeDocumentType, error) {
	var documentType documentTypemodel.EmployeeDocumentType
	if err := d.db.WithContext(ctx).
		Where("id = ?", id).
		First(&documentType).Error; err != nil {
		return nil, err
	}
	return &documentType, nil
}

func (d *DocumentType) GetAllDocumentType(ctx context.Context) ([]*documentTypemodel.EmployeeDocumentType, error) {
	var list []*documentTypemodel.EmployeeDocumentType
	if err := d.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DocumentType) GetLastDocumentTypeByCode(ctx context.Context, documentType *documentTypemodel.EmployeeDocumentType) error {
	return d.db.WithContext(ctx).
		Order("id DESC").
		First(documentType).Error
}
