package repository

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type attendanceCategoryRepository struct {
	db *gorm.DB
}

func NewAttendanceCategoryRepository(db *gorm.DB) repo_interface.AttendanceCategoryRepository {
	return &attendanceCategoryRepository{db: db}
}

func (r *attendanceCategoryRepository) Create(ctx context.Context, ac *model.AttendanceCategory) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context) error {
		if err := r.db.WithContext(ctx).Create(ac).Error; err != nil {
			return fmt.Errorf("failed to create attendance category: %w", err)
		}
		return nil
	})
}

func (r *attendanceCategoryRepository) Update(ctx context.Context, ac *model.AttendanceCategory) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context) error {
		result := r.db.WithContext(ctx).
			Model(&model.AttendanceCategory{}).
			Where("attendance_category_id = ?", ac.AttendanceCategoryID).
			Updates(ac)

		if result.Error != nil {
			return fmt.Errorf("failed to update attendance category: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.New("no records updated, record may not exist")
		}

		return nil
	})
}

func (r *attendanceCategoryRepository) Delete(ctx context.Context, id string) error {
	return utils.WithTransaction(r.db, ctx, func(ctx context.Context) error {
		result := r.db.WithContext(ctx).
			Model(&model.AttendanceCategory{}).
			Where("attendance_category_id = ?", id).
			Update("is_deleted", true)

		if result.Error != nil {
			return fmt.Errorf("failed to soft delete attendance category: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.New("no records deleted, record may not exist")
		}

		return nil
	})
}

func (r *attendanceCategoryRepository) GetByID(ctx context.Context, id string) (*model.AttendanceCategory, error) {
	var ac model.AttendanceCategory
	err := r.db.WithContext(ctx).
		Preload("Office").
		Where("attendance_category_id = ? AND is_deleted = false", id).
		First(&ac).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("attendance category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get attendance category: %w", err)
	}
	return &ac, nil
}

func (r *attendanceCategoryRepository) ListByOffice(ctx context.Context, officeID string) ([]model.AttendanceCategory, error) {
	var list []model.AttendanceCategory
	err := r.db.WithContext(ctx).
		Where("office_id = ? AND is_deleted = false AND status = ?", officeID, model.StatusActive).
		Find(&list).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list attendance categories by office: %w", err)
	}
	return list, nil
}

func (r *attendanceCategoryRepository) GetList(ctx context.Context) ([]model.AttendanceCategory, error) {
	var list []model.AttendanceCategory
	err := r.db.WithContext(ctx).
		Where("is_deleted = false").
		Find(&list).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get attendance categories list: %w", err)
	}
	return list, nil
}

func (r *attendanceCategoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.AttendanceCategory{}).
		Where("attendance_category_name = ? AND is_deleted = false", name).
		Count(&count).Error
	return count > 0, err
}

func (r *attendanceCategoryRepository) ExistsByNameAndNotID(ctx context.Context, name string, excludeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.AttendanceCategory{}).
		Where("attendance_category_name = ? AND attendance_category_id != ? AND is_deleted = false", name, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *attendanceCategoryRepository) GetIfExists(ctx context.Context, id string) (*model.AttendanceCategory, bool, error) {
	var category model.AttendanceCategory
	err := r.db.WithContext(ctx).
		Model(&model.AttendanceCategory{}).
		Where("attendance_category_id = ? AND is_deleted = false", id).
		First(&category).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &category, true, nil
}
