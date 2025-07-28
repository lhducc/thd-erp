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
