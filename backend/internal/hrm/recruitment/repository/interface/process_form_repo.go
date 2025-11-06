package repository

import (
	"context"
	dto "erp/backend/internal/hrm/recruitment/dto/request"
	"erp/backend/internal/hrm/recruitment/model"
)

type ProcessFormRepository interface {
	CreateFormWithStages(ctx context.Context, req dto.CreateProcessFormRequest) error
	Update(ctx context.Context, form *model.ProcessForm) error
	FindByID(ctx context.Context, id string) (*model.ProcessForm, error)
	GetAllActive(ctx context.Context) ([]model.ProcessForm, error)
	GetById(ctx context.Context, id string) (model.ProcessForm, error)
	Delete(ctx context.Context, id string) error
}
