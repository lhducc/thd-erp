package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"errors"
	"gorm.io/gorm"
)

type timesheetRepo struct {
	db *gorm.DB
}

func NewTimesheetRepo(db *gorm.DB) repo_interface.TimeSheetRepoInterface {
	return &timesheetRepo{db: db}
}

func (t timesheetRepo) Create(ctx context.Context, timesheets []*model.TimeSheet) error {
	return t.db.WithContext(ctx).Create(timesheets).Error
}

func (r *timesheetRepo) Update(timesheet *model.TimeSheet) error {
	return r.db.Save(timesheet).Error
}

func (r *timesheetRepo) FindByID(id int) (*model.TimeSheet, error) {
	var ts model.TimeSheet
	if err := r.db.Preload("Details").First(&ts, id).Error; err != nil {
		return nil, err
	}
	return &ts, nil
}

func (r *timesheetRepo) FindByEmployeeAndMonth(employeeID string, month, year int) (*model.TimeSheet, error) {
	var ts model.TimeSheet
	err := r.db.Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).
		Preload("Details").
		First(&ts).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ts, err
}
