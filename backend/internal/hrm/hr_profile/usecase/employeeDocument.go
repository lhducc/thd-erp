package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"fmt"
	"strconv"
	"strings"

	"erp/backend/pkg"
)

type EmployeeDocumentRepo interface {
	CreateEmployeeDocument(ctx context.Context, employeeDocument *model.EmployeeDocument) error
	UpdateEmployeeDocument(ctx context.Context, id string, employeeDocument *model.EmployeeDocument) error
	DeleteEmployeeDocument(ctx context.Context, id string) error
	GetEmployeeDocumentById(ctx context.Context, id string) (*model.EmployeeDocument, error)
	GetAllEmployeeDocuments(ctx context.Context) ([]*model.EmployeeDocument, error)
	GetEmployeeDocumentsByEmployeeID(ctx context.Context, employeeID string) ([]*model.EmployeeDocument, error)
	GetEmployeeDocumentsByDocumentTypeID(ctx context.Context, documentTypeID string) ([]*model.EmployeeDocument, error)
	GetLastDocumentByCode(ctx context.Context, documentType *model.EmployeeDocument) error
	GetEmployeeDocumentsPaginated(ctx context.Context, search, status, condition string, offset, limit int) ([]*model.EmployeeDocumentResponse, int64, error)
}

type EmployeeDocumentBiz struct {
	repo EmployeeDocumentRepo
}

func NewEmployeeDocumentBiz(repo EmployeeDocumentRepo) *EmployeeDocumentBiz {
	return &EmployeeDocumentBiz{repo: repo}
}

// CreateEmployeeDocument creates a new employee document with generated DocumentID and current date
func (b *EmployeeDocumentBiz) CreateEmployeeDocument(ctx context.Context, employeeDocument *model.EmployeeDocument) error {
	employeeDocument.CreatedDate = utils.GetCurrentDate()

	code, err := b.GenerateDocumentCode(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate document code: %w", err)
	}
	employeeDocument.DocumentID = code

	if err := employeeDocument.ValidateEmployeeDocument(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if err := b.repo.CreateEmployeeDocument(ctx, employeeDocument); err != nil {
		return fmt.Errorf("failed to create employee document: %w", err)
	}
	return nil
}

// GetEmployeeDocument retrieves an employee document by ID
func (b *EmployeeDocumentBiz) GetEmployeeDocument(ctx context.Context, id string) (*model.EmployeeDocument, error) {
	doc, err := b.repo.GetEmployeeDocumentById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee document: %w", err)
	}
	return doc, nil
}

// GetAllEmployeeDocuments returns all employee documents
func (b *EmployeeDocumentBiz) GetAllEmployeeDocuments(ctx context.Context) ([]*model.EmployeeDocument, error) {
	list, err := b.repo.GetAllEmployeeDocuments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all employee documents: %w", err)
	}
	return list, nil
}

func (b *EmployeeDocumentBiz) GetEmployeeDocumentsPaginated(ctx context.Context, search, status, condition string, offset, limit int) ([]*model.EmployeeDocumentResponse, int64, error) {
	return b.repo.GetEmployeeDocumentsPaginated(ctx, search, status, condition, offset, limit)
}

// GetEmployeeDocumentsByEmployeeID returns documents by employee ID
func (b *EmployeeDocumentBiz) GetEmployeeDocumentsByEmployeeID(ctx context.Context, employeeID string) ([]*model.EmployeeDocument, error) {
	list, err := b.repo.GetEmployeeDocumentsByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents by employee ID: %w", err)
	}
	return list, nil
}

// GetEmployeeDocumentsByDocumentTypeID returns documents by document type ID
func (b *EmployeeDocumentBiz) GetEmployeeDocumentsByDocumentTypeID(ctx context.Context, documentTypeID string) ([]*model.EmployeeDocument, error) {
	list, err := b.repo.GetEmployeeDocumentsByDocumentTypeID(ctx, documentTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents by document type ID: %w", err)
	}
	return list, nil
}

// UpdateEmployeeDocument updates an existing employee document by ID
func (b *EmployeeDocumentBiz) UpdateEmployeeDocument(ctx context.Context, id string, data *model.EmployeeDocument) error {
	if err := data.ValidateEmployeeDocument(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if err := b.repo.UpdateEmployeeDocument(ctx, id, data); err != nil {
		return fmt.Errorf("failed to update employee document: %w", err)
	}
	return nil
}

// DeleteEmployeeDocument deletes an employee document by ID
func (b *EmployeeDocumentBiz) DeleteEmployeeDocument(ctx context.Context, id string) error {
	if err := b.repo.DeleteEmployeeDocument(ctx, id); err != nil {
		return fmt.Errorf("failed to delete employee document: %w", err)
	}
	return nil
}

// GenerateDocumentCode generates next document code with prefix "TL"
func (b *EmployeeDocumentBiz) GenerateDocumentCode(ctx context.Context) (string, error) {
	var lastDocument model.EmployeeDocument
	err := b.repo.GetLastDocumentByCode(ctx, &lastDocument)

	if err != nil || !strings.HasPrefix(lastDocument.DocumentID, "TL") {
		return "TL00001", nil
	}

	numberStr := strings.TrimPrefix(lastDocument.DocumentID, "TL")
	number, err := strconv.Atoi(strings.TrimSpace(numberStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse document type number: %w", err)
	}

	nextNumber := number + 1
	if nextNumber > 99999 {
		return "", fmt.Errorf("maximum document type code reached: TL99999")
	}

	return fmt.Sprintf("TL%05d", nextNumber), nil
}
