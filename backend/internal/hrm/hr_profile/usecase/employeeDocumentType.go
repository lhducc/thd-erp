package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"fmt"
	"strconv"
	"strings"

	"erp/backend/pkg"
)

type DocumentTypeRepo interface {
	CreateDocumentType(ctx context.Context, documentType *model.EmployeeDocumentTypeCreate) error
	UpdateDocumentType(ctx context.Context, id string, documentType *model.EmployeeDocumentTypeUpdate) error
	DeleteDocumentType(ctx context.Context, id string) error
	GetDocumentTypeById(ctx context.Context, id string) (*model.EmployeeDocumentType, error)
	GetAllDocumentType(ctx context.Context) ([]*model.EmployeeDocumentType, error)
	GetLastDocumentTypeByCode(ctx context.Context, documentType *model.EmployeeDocumentType) error
}

type DocumentTypeBiz struct {
	repo DocumentTypeRepo
}

func NewDocumentTypeBiz(repo DocumentTypeRepo) *DocumentTypeBiz {
	return &DocumentTypeBiz{repo: repo}
}

// CreateDocumentType creates a new document type with generated ID and current date
func (b *DocumentTypeBiz) CreateDocumentType(ctx context.Context, documentType *model.EmployeeDocumentTypeCreate) error {
	documentType.CreatedDate = utils.GetCurrentDate()

	code, err := b.GenerateDocumentTypeCode(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate document type code: %w", err)
	}
	documentType.ID = code

	if err := b.repo.CreateDocumentType(ctx, documentType); err != nil {
		return fmt.Errorf("failed to create document type: %w", err)
	}
	return nil
}

// GetDocumentType retrieves a document type by its ID
func (b *DocumentTypeBiz) GetDocumentType(ctx context.Context, id string) (*model.EmployeeDocumentType, error) {
	documentType, err := b.repo.GetDocumentTypeById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get document type: %w", err)
	}
	return documentType, nil
}

// GetAllDocumentType returns all document types
func (b *DocumentTypeBiz) GetAllDocumentType(ctx context.Context) ([]*model.EmployeeDocumentType, error) {
	list, err := b.repo.GetAllDocumentType(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all document types: %w", err)
	}
	return list, nil
}

// UpdateDocumentType updates an existing document type by ID
func (b *DocumentTypeBiz) UpdateDocumentType(ctx context.Context, id string, data *model.EmployeeDocumentTypeUpdate) error {
	if err := b.repo.UpdateDocumentType(ctx, id, data); err != nil {
		return fmt.Errorf("failed to update document type: %w", err)
	}
	return nil
}

// DeleteDocumentType deletes a document type by ID
func (b *DocumentTypeBiz) DeleteDocumentType(ctx context.Context, id string) error {
	if err := b.repo.DeleteDocumentType(ctx, id); err != nil {
		return fmt.Errorf("failed to delete document type: %w", err)
	}
	return nil
}

// GenerateDocumentTypeCode generates next available document type code with prefix "TL"
func (b *DocumentTypeBiz) GenerateDocumentTypeCode(ctx context.Context) (string, error) {
	var lastDocumentType model.EmployeeDocumentType
	err := b.repo.GetLastDocumentTypeByCode(ctx, &lastDocumentType)

	// Start from TL00001 if no record or invalid prefix
	if err != nil || !strings.HasPrefix(lastDocumentType.ID, "DT") {
		return "DT00001", nil
	}

	numberStr := strings.TrimPrefix(lastDocumentType.ID, "DT")
	number, err := strconv.Atoi(strings.TrimSpace(numberStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse document type number: %w", err)
	}

	nextNumber := number + 1
	if nextNumber > 99999 {
		return "", fmt.Errorf("maximum document type code reached: DT99999")
	}

	return fmt.Sprintf("DT%05d", nextNumber), nil
}
