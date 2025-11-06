package implement

import (
	"context"
	recmodel "erp/backend/internal/hrm/recruitment/model"
	repository "erp/backend/internal/hrm/recruitment/repository/interface"
	"errors"

	"gorm.io/gorm"
)

type processStageRepoImpl struct {
	db *gorm.DB
}

func NewProcessStageRepositoryImpl(db *gorm.DB) repository.ProcessStageRepository {
	return &processStageRepoImpl{db: db}
}

func (r *processStageRepoImpl) Save(ctx context.Context, stage *recmodel.ProcessStage) error {
	if stage == nil {
		return errors.New("nil stage")
	}
	return r.db.WithContext(ctx).Create(stage).Error
}

func (r *processStageRepoImpl) GetLastOrder(ctx context.Context, processFormID string) (int, error) {
	var last recmodel.ProcessStage
	err := r.db.WithContext(ctx).
		Where("process_form_id = ?", processFormID).
		Order("stage_order DESC").
		Limit(1).
		Find(&last).Error
	if err != nil {
		return 0, err
	}
	if last.ProcessStageID == 0 {
		return 0, nil
	}
	return last.StageOrder, nil
}

func (r *processStageRepoImpl) ExistsProcessForm(ctx context.Context, processFormID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&recmodel.ProcessForm{}).
		Where("process_form_id = ?", processFormID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *processStageRepoImpl) FindByID(ctx context.Context, processStageID int) (*recmodel.ProcessStage, error) {
	var ps recmodel.ProcessStage
	err := r.db.WithContext(ctx).Where("process_stage_id = ?", processStageID).First(&ps).Error
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

func (r *processStageRepoImpl) UpdateStageName(ctx context.Context, processStageID int, newName string) error {
	return r.db.WithContext(ctx).Model(&recmodel.ProcessStage{}).
		Where("process_stage_id = ?", processStageID).
		Update("stage_name", newName).Error
}

func (r *processStageRepoImpl) DeleteByID(ctx context.Context, processStageID int) error {
	return r.db.WithContext(ctx).Where("process_stage_id = ?", processStageID).Delete(&recmodel.ProcessStage{}).Error
}

func (r *processStageRepoImpl) DecrementOrdersAfter(ctx context.Context, processFormID string, fromOrder int) error {
	return r.db.WithContext(ctx).
		Model(&recmodel.ProcessStage{}).
		Where("process_form_id = ? AND stage_order > ?", processFormID, fromOrder).
		Update("stage_order", gorm.Expr("stage_order - ?", 1)).Error
}

func (r *processStageRepoImpl) ExistsStageName(ctx context.Context, processFormID string, stageName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&recmodel.ProcessStage{}).
		Where("process_form_id = ? AND stage_name = ?", processFormID, stageName).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
