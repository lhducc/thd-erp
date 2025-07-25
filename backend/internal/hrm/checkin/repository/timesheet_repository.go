package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"gorm.io/gorm"
	"strings"
)

type timesheetRepo struct {
	db *gorm.DB
}

func NewTimesheetRepo(db *gorm.DB) repo_interface.TimesheetDetailInterface {
	return &timesheetRepo{db: db}
}

func (r *timesheetRepo) Create(ctx context.Context, timesheet *model.Timesheet) error {
	return r.db.WithContext(ctx).Create(timesheet).Error
}

func (r *timesheetRepo) GetByID(ctx context.Context, id string) (*model.Timesheet, error) {
	var ts model.Timesheet
	err := r.db.WithContext(ctx).Preload("Office").First(&ts, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ts, nil
}

func (r *timesheetRepo) Update(ctx context.Context, timesheet *model.Timesheet) error {
	var existing model.Timesheet
	if err := r.db.WithContext(ctx).First(&existing, "id = ?", timesheet.ID).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&existing).
		Updates(timesheet).Error
}

func (r *timesheetRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Timesheet{}, "id = ?", id).Error
}

func (r *timesheetRepo) List(ctx context.Context, page, limit int) ([]model.Timesheet, int64, error) {
	var (
		timesheets []model.Timesheet
		total      int64
	)

	if err := r.db.WithContext(ctx).Model(&model.Timesheet{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Preload("Office").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&timesheets).Error

	return timesheets, total, err
}

func (s *timesheetRepo) GetLastDecisionByCode(ctx context.Context) (string, error) {
	var lastID model.Timesheet
	err := s.db.WithContext(ctx).
		Order("id DESC").
		First(&lastID).Error
	if err != nil {
		return "", err
	}
	return lastID.ID, nil
}

func (r *timesheetRepo) IsDuplicate(ctx context.Context, officeID string, month, year int, timesheetID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&model.Timesheet{}).
		Where("office_id = ? AND month = ? AND year = ?", officeID, month, year)
	if strings.TrimSpace(timesheetID) != "" {
		query = query.Where("id != ? ", timesheetID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
