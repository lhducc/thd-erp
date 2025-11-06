package service

import (
	"context"
	request "erp/backend/internal/hrm/recruitment/dto/request"
	"erp/backend/internal/hrm/recruitment/model"
)

type ProcessFormService interface {
	CreateFormWithStages(ctx context.Context, req request.CreateProcessFormRequest) error
	UpdateForm(ctx context.Context, req request.UpdateProcessFormRequest) error
	DeleteForm(ctx context.Context, id string) error
	GetAllForms(ctx context.Context) ([]model.ProcessForm, error)
	GetFormById(ctx context.Context, id string) (model.ProcessForm, error)
}
