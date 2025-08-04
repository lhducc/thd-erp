package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"gorm.io/gorm"
)

type timesheetDetailRepo struct {
	db *gorm.DB
}

func NewTimesheetDetailRepo(db *gorm.DB) repo_interface.TimeSheetDetailRepoInterface {
	return &timesheetDetailRepo{db: db}
}

func (t timesheetDetailRepo) Create(ctx context.Context, timesheet *model.TimeSheetDetail) error {
	return t.db.WithContext(ctx).Create(timesheet).Error
}

func (t timesheetDetailRepo) GetDetailsByTimeSheetID(ctx context.Context, timeSheetID int) ([]model.TimeSheetDetail, error) {
	var details []model.TimeSheetDetail
	if err := t.db.WithContext(ctx).Where("timesheet_id = ?", timeSheetID).Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *timesheetDetailRepo) GetDetailByID(ctx context.Context, detailID int) (*model.TimeSheetDetail, error) {
	var detail model.TimeSheetDetail
	err := r.db.WithContext(ctx).Where("timesheet_detail_id = ?", detailID).First(&detail).Error
	return &detail, err
}

func (r *timesheetDetailRepo) UpdateTimeSheetDetail(ctx context.Context, detail *model.TimeSheetDetail) error {
	if err := r.db.WithContext(ctx).Save(detail).Error; err != nil {
		return err
	}
	return nil
}
