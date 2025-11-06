package implement

import (
	"context"
	dto "erp/backend/internal/hrm/recruitment/dto/request"
	"erp/backend/internal/hrm/recruitment/model"
	repository "erp/backend/internal/hrm/recruitment/repository/interface"
	"gorm.io/gorm"
)

type processFormRepositoryImpl struct {
	db *gorm.DB
}

func NewProcessFormRepositoryImpl(db *gorm.DB) repository.ProcessFormRepository {
	return &processFormRepositoryImpl{db: db}
}

func (r *processFormRepositoryImpl) Update(ctx context.Context, form *model.ProcessForm) error {
	return r.db.WithContext(ctx).Save(form).Error
}

func (r *processFormRepositoryImpl) FindByID(ctx context.Context, id string) (*model.ProcessForm, error) {
	var form model.ProcessForm
	if err := r.db.WithContext(ctx).First(&form, "process_form_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *processFormRepositoryImpl) GetAllActive(ctx context.Context) ([]model.ProcessForm, error) {
	var forms []model.ProcessForm
	err := r.db.WithContext(ctx).Preload("Stages").Where("is_deleted = ?", false).Find(&forms).Error
	return forms, err
}

func (r *processFormRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.ProcessForm{}).
		Where("process_form_id = ?", id).
		Update("is_deleted", true).Error
}

func (r *processFormRepositoryImpl) GetById(ctx context.Context, id string) (model.ProcessForm, error) {
	var form model.ProcessForm
	err := r.db.WithContext(ctx).
		Preload("Stages").
		Where("process_form_id = ? AND is_deleted = ?", id, false).
		First(&form).Error
	return form, err
}

func (r *processFormRepositoryImpl) CreateFormWithStages(ctx context.Context, req dto.CreateProcessFormRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		form := model.ProcessForm{
			ProcessFormName: req.ProcessFormName,
		}

		if err := tx.Create(&form).Error; err != nil {
			return err
		}

		for i, s := range req.Stages {
			stage := model.ProcessStage{
				StageName:     s.StageName,
				StageOrder:    i + 1,
				ProcessFormID: form.ProcessFormID,
			}
			if err := tx.Create(&stage).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
