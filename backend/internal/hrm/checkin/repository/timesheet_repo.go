package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type timesheetRepo struct {
	db *gorm.DB
}

func NewTimesheetRepo(db *gorm.DB) repo_interface.TimeSheetRepoInterface {
	return &timesheetRepo{db: db}
}

func (t timesheetRepo) CreateTimeSheets(ctx context.Context, timesheets []*model.TimeSheet) error {
	return t.db.WithContext(ctx).Create(timesheets).Error
}

func (t timesheetRepo) CreateEmployeeTimeSheet(tx *gorm.DB, timesheet *model.TimeSheet) error {
	return tx.Create(timesheet).Error
}

func (r *timesheetRepo) FindByID(id int) (*model.TimeSheet, error) {
	var ts model.TimeSheet
	if err := r.db.Preload("Details").First(&ts, id).Error; err != nil {
		return nil, err
	}
	return &ts, nil
}

func (r *timesheetRepo) FindByEmployeeAndMonth(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error) {
	var ts model.TimeSheet
	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).
		Preload("Details").
		First(&ts).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ts, err
}

func (r *timesheetRepo) UpdateTimesheetAndCreateDetail(ctx context.Context, ts *model.TimeSheet) error {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("operation canceled before starting: %w", err)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update main timesheet fields
		updateFields := map[string]interface{}{
			"total_work_days":    ts.TotalWorkDays,
			"late_shifts":        ts.LateShifts,
			"total_late_minutes": ts.TotalLateMinutes,
			"updated_at":         time.Now(),
		}

		if err := tx.Model(&model.TimeSheet{}).
			Where("timesheet_id = ?", ts.TimeSheetID).
			Updates(updateFields).Error; err != nil {
			return fmt.Errorf("failed to update timesheet: %w", err)
		}

		// 2. Handle batch upsert for timesheet details
		if len(ts.Details) > 0 {
			details := make([]model.TimeSheetDetail, 0, len(ts.Details))
			for i := range ts.Details {
				ts.Details[i].TimeSheetID = ts.TimeSheetID
				details = append(details, ts.Details[i])
			}

			batchSize := 200
			for i := 0; i < len(details); i += batchSize {
				end := i + batchSize
				if end > len(details) {
					end = len(details)
				}
				batch := details[i:end]

				// Check context before each batch
				if err := ctx.Err(); err != nil {
					return fmt.Errorf("operation canceled during batch processing: %w", err)
				}

				// Perform batch upsert with manual adjustment fields
				err := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "timesheet_id"}, {Name: "date"}},
					DoUpdates: clause.AssignmentColumns([]string{
						"work_shift_id", "is_working_day", "is_weekend",
						"checkin_record_id", "checkout_record_id",
						"work_hours", "work_days", "is_additional_shift",
						"is_late", "late_minutes",
						"leave_type", "leave_hours", "is_absent", "absent_reason",
						// Preserve manual adjustment fields
						"is_manually_adjusted", "original_work_days",
						"work_day_adjusted", "adjustment_by", "adjustment_at",
					}),
				}).Create(&batch).Error

				if err != nil {
					return fmt.Errorf("batch upsert failed at position %d-%d: %w", i, end-1, err)
				}
			}
		}

		return nil
	})
}

func (r *timesheetRepo) GetByID(ctx context.Context, timesheetID int) (*model.TimeSheet, error) {
	var ts model.TimeSheet
	if err := r.db.WithContext(ctx).Where("timesheet_id = ?", timesheetID).First(&ts).Error; err != nil {
		return nil, err
	}
	return &ts, nil
}

func (r *timesheetRepo) Update(ctx context.Context, timesheet *model.TimeSheet) error {
	return r.db.WithContext(ctx).Model(&model.TimeSheet{}).Updates(timesheet).Error
}

func (r *timesheetRepo) FindByEmployeeForExport(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error) {
	var ts model.TimeSheet
	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).
		Preload("Details", func(db *gorm.DB) *gorm.DB {
			return db.Order("date ASC") // Sắp xếp details theo ngày
		}).
		Preload("Employee").
		Preload("Employee.JobTitle.HierarchyLevel").
		Preload("Office").
		Preload("Department").
		First(&ts).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ts, err
}
