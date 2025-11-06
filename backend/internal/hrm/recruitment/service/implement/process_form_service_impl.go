package implement

import (
	"context"
	request "erp/backend/internal/hrm/recruitment/dto/request"
	"erp/backend/internal/hrm/recruitment/model"
	repository "erp/backend/internal/hrm/recruitment/repository/interface"
	service "erp/backend/internal/hrm/recruitment/service/interface"
)

type processFormServiceImpl struct {
	repo repository.ProcessFormRepository
}

func NewProcessFormServiceImpl(repo repository.ProcessFormRepository) service.ProcessFormService {
	return &processFormServiceImpl{repo: repo}
}

func (s *processFormServiceImpl) CreateFormWithStages(ctx context.Context, req request.CreateProcessFormRequest) error {
	return s.repo.CreateFormWithStages(ctx, req)
}

func (s *processFormServiceImpl) UpdateForm(ctx context.Context, req request.UpdateProcessFormRequest) error {
	form, err := s.repo.FindByID(ctx, req.ProcessFormID)
	if err != nil {
		return err
	}
	form.ProcessFormName = req.ProcessFormName
	return s.repo.Update(ctx, form)
}

func (s *processFormServiceImpl) DeleteForm(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *processFormServiceImpl) GetAllForms(ctx context.Context) ([]model.ProcessForm, error) {
	return s.repo.GetAllActive(ctx)
}

func (s *processFormServiceImpl) GetFormById(ctx context.Context, id string) (model.ProcessForm, error) {
	return s.repo.GetById(ctx, id)
}
