package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type attendanceRecordRepository struct {
	db *gorm.DB
}

func NewAttendanceRecordRepository(db *gorm.DB) repo_interface.AttendanceRecordRepository {
	return &attendanceRecordRepository{db: db}
}

func (r *attendanceRecordRepository) Create(ctx context.Context, record *model.AttendanceRecord) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context) error {
		if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
			return fmt.Errorf("failed to create attendance record: %w", err)
		}
		return nil
	})
}

func (r *attendanceRecordRepository) Update(ctx context.Context, record *model.AttendanceRecord) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context) error {
		result := r.db.WithContext(ctx).
			Model(&model.AttendanceRecord{}).
			Where("attendance_record_id = ?", record.AttendanceRecordID).
			Updates(record)

		if result.Error != nil {
			return fmt.Errorf("failed to update attendance record: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return errors.New("no records updated, record may not exist")
		}
		return nil
	})
}

func (r *attendanceRecordRepository) Delete(ctx context.Context, id string) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context) error {
		result := r.db.WithContext(ctx).
			Where("attendance_record_id = ?", id).
			Delete(&model.AttendanceRecord{})

		if result.Error != nil {
			return fmt.Errorf("failed to delete attendance record: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return errors.New("no record deleted, maybe not found")
		}
		return nil
	})
}

func (r *attendanceRecordRepository) GetByID(ctx context.Context, id string) (*model.AttendanceRecord, error) {
	var record model.AttendanceRecord
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Office").
		Preload("AttendanceCategory").
		Joins("JOIN attendance_categories ON attendance_records.category_id = attendance_categories.attendance_category_id").
		Where("attendance_records.attendance_record_id = ? AND attendance_categories.auto_approve = ?",
			id, false).
		First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("attendance record not found or is auto-approved: %w", err)
		}
		return nil, fmt.Errorf("failed to get attendance record: %w", err)
	}
	return &record, nil
}

func (r *attendanceRecordRepository) ListByEmployee(ctx context.Context, employeeID string) ([]model.AttendanceRecord, error) {
	var records []model.AttendanceRecord
	err := r.db.WithContext(ctx).
		Joins("JOIN attendance_categories ON attendance_records.category_id = attendance_categories.attendance_category_id").
		Where("attendance_records.employee_id = ? AND attendance_categories.auto_approve = ?",
			employeeID, false).
		Order("attendance_records.timestamp DESC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list attendance records by employee: %w", err)
	}
	return records, nil
}

func (r *attendanceRecordRepository) CountRequestsByEmployee(ctx context.Context, employeeID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.AttendanceRecord{}).
		Joins("JOIN attendance_categories ON attendance_records.category_id = attendance_categories.attendance_category_id").
		Where("attendance_records.employee_id = ? AND attendance_categories.auto_approve = ?", employeeID, false).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count attendance requests: %w", err)
	}
	return count, nil
}

func (r *attendanceRecordRepository) ListByDateRange(ctx context.Context, employeeID string, from, to time.Time) ([]model.AttendanceRecord, error) {
	var records []model.AttendanceRecord
	err := r.db.WithContext(ctx).
		Joins("JOIN attendance_categories ON attendance_records.category_id = attendance_categories.attendance_category_id").
		Where("attendance_records.employee_id = ? AND attendance_records.timestamp >= ? AND attendance_records.timestamp <= ? AND attendance_categories.auto_approve = ?",
			employeeID, from, to, false).
		Order("attendance_records.timestamp ASC").
		Find(&records).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list attendance records by date range: %w", err)
	}
	return records, nil
}
