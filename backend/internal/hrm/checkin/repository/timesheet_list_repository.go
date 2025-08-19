package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type timesheetListRepo struct {
	db *gorm.DB
}

func NewTimesheetListRepo(db *gorm.DB) repo_interface.TimesheetListInterface {
	return &timesheetListRepo{db: db}
}

func (r *timesheetListRepo) Create(ctx context.Context, timesheet *model.TimeSheetList) error {
	return r.db.WithContext(ctx).Create(timesheet).Error
}

func (r *timesheetListRepo) GetByID(ctx context.Context, id string) (*model.TimeSheetList, error) {
	var ts model.TimeSheetList
	err := r.db.WithContext(ctx).
		Preload("Office").Preload("Timesheets").
		Preload("Timesheets.Employee").
		Preload("Timesheets.Employee.JobTitle").
		Preload("Timesheets.Employee.Position").
		Preload("Timesheets.Employee.JobTitle.HierarchyLevel").
		Preload("Timesheets.Details").
		Preload("Timesheets.Department").
		Preload("Timesheets.Details.WorkShift").
		Preload("Timesheets.Details.CheckInRecord").
		Preload("Timesheets.Details.CheckOutRecord").
		First(&ts, "timesheet_list_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	for _, timesheet := range ts.Timesheets {
		if timesheet.Employee != nil {
			if timesheet.Employee.JobTitle == nil {
				log.Println("timesheet.Employee.JobTitle is nil" + timesheet.EmployeeID)
				continue
			} else {
				timesheet.Employee.HierarchyLevel = timesheet.Employee.JobTitle.HierarchyLevel
			}
		}
	}
	return &ts, nil
}

func (r *timesheetListRepo) Update(ctx context.Context, timesheet *model.TimeSheetList) error {
	if err := r.db.WithContext(ctx).
		Model(&model.TimeSheetList{}).
		Where("timesheet_list_id = ?", timesheet.TimeSheetListID).
		Update("timesheet_list_name", timesheet.TimeSheetListName).
		Update("updated_by", timesheet.UpdatedBy).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *timesheetListRepo) UpdateLocked(ctx context.Context, timesheet *model.TimeSheetList) error {
	if err := r.db.WithContext(ctx).
		Model(&model.TimeSheetList{}).
		Where("timesheet_list_id = ?", timesheet.TimeSheetListID).
		Update("is_locked", timesheet.IsLocked).
		Update("locked_at", timesheet.LockedAt).
		Update("locked_by", timesheet.LockedBy).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *timesheetListRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.TimeSheetList{}, "timesheet_list_id = ?", id).Error
}

func (r *timesheetListRepo) List(ctx context.Context, page, limit int) ([]model.TimeSheetList, int64, error) {
	var (
		timesheets []model.TimeSheetList
		total      int64
	)

	if err := r.db.WithContext(ctx).Model(&model.TimeSheetList{}).Count(&total).Error; err != nil {
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

func (s *timesheetListRepo) GetLastDecisionByCode(ctx context.Context) (string, error) {
	var lastID model.TimeSheetList
	err := s.db.WithContext(ctx).
		Order("timesheet_list_id DESC").
		First(&lastID).Error
	if err != nil {
		return "", err
	}
	return lastID.TimeSheetListID, nil
}

func (r *timesheetListRepo) IsDuplicate(ctx context.Context, officeID string, month, year int, timesheetID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&model.TimeSheetList{}).
		Where("office_id = ? AND month = ? AND year = ?", officeID, month, year)
	if strings.TrimSpace(timesheetID) != "" {
		query = query.Where("timesheet_list_id != ? ", timesheetID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *timesheetListRepo) GetTimeSheetByOfficeIDAndTime(officeID string, month, year int) (*model.TimeSheetList, error) {
	var timesheet model.TimeSheetList
	err := r.db.
		Model(&model.TimeSheetList{}).
		Where("office_id = ? AND month = ? AND year = ?", officeID, month, year).
		First(&timesheet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &timesheet, nil
}

func (r *timesheetListRepo) GetForExport(ctx context.Context, id string) (*model.TimeSheetList, error) {
	var ts model.TimeSheetList
	err := r.db.WithContext(ctx).
		Preload("Office").
		Preload("Timesheets", func(db *gorm.DB) *gorm.DB {
			// Sắp xếp timesheet theo EmployeeID để đảm bảo thứ tự
			return db.Order("employee_id ASC")
		}).
		Preload("Timesheets.Employee").
		Preload("Timesheets.Office").
		Preload("Timesheets.Employee.JobTitle").
		Preload("Timesheets.Employee.Position").
		Preload("Timesheets.Employee.JobTitle.HierarchyLevel").
		Preload("Timesheets.Details", func(db *gorm.DB) *gorm.DB {
			// Sắp xếp details theo ngày
			return db.Order("date ASC")
		}).
		Preload("Timesheets.Department").
		Preload("Timesheets.Details.WorkShift").
		Preload("Timesheets.Details.CheckInRecord").
		Preload("Timesheets.Details.CheckOutRecord").
		Preload("Timesheets.Details.AdjustmentUser").
		First(&ts, "timesheet_list_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy timesheet list: %w", err)
	}

	// Set hierarchy level cho employee
	for i := range ts.Timesheets {
		if ts.Timesheets[i].Employee != nil &&
			ts.Timesheets[i].Employee.JobTitle != nil {
			ts.Timesheets[i].Employee.HierarchyLevel = ts.Timesheets[i].Employee.JobTitle.HierarchyLevel
		}
	}

	return &ts, nil
}
