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
	// Kiểm tra context đã bị hủy chưa trước khi bắt đầu
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("operation canceled before starting: %w", err)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Tối ưu hóa cập nhật timesheet
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

		// 2. Xử lý batch insert/update cho details
		if len(ts.Details) > 0 {
			// Chuẩn bị dữ liệu batch
			details := make([]model.TimeSheetDetail, 0, len(ts.Details))
			for i := range ts.Details {
				ts.Details[i].TimeSheetID = ts.TimeSheetID
				details = append(details, ts.Details[i])
			}

			// Batch upsert với kích thước phù hợp (100-500 records/batch)
			batchSize := 200
			for i := 0; i < len(details); i += batchSize {
				end := i + batchSize
				if end > len(details) {
					end = len(details)
				}
				batch := details[i:end]

				// Kiểm tra context trước mỗi batch
				if err := ctx.Err(); err != nil {
					return fmt.Errorf("operation canceled during batch processing: %w", err)
				}

				// Thực hiện batch upsert
				err := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "timesheet_id"}, {Name: "date"}},
					DoUpdates: clause.AssignmentColumns([]string{
						"work_shift_id", "is_working_day", "is_weekend",
						"checkin_record_id", "checkout_record_id",
						"work_hours", "work_days", "is_additional_shift",
						"is_late", "late_minutes",
						"leave_type", "leave_hours", "is_absent", "absent_reason",
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
