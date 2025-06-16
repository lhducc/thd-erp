package repository

import (
	"context"
	model "erp/backend/internal/hrm/hr_profile/model"

	"gorm.io/gorm"
)

type departmentStore struct {
	db *gorm.DB
}

func NewDepartmentStore(db *gorm.DB) *departmentStore {
	return &departmentStore{db: db}
}

func (s *departmentStore) CreateDepartment(context context.Context, data *model.DepartmentCreate) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *departmentStore) GetDepartment(ctx context.Context, id string) (*model.Department, error) {
	var department model.Department
	if err := r.db.WithContext(ctx).Table("department").Preload("Office").
		Where("department_id = ?", id).
		First(&department).Error; err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *departmentStore) GetAllDepartment(ctx context.Context) ([]model.Department, error) {

	var positions []model.Department
	if err := r.db.WithContext(ctx).
		Table("department").
		Preload("Office").
		Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *departmentStore) UpdateDepartment(ctx context.Context, id string, data *model.DepartmentCreate) error {
	return r.db.WithContext(ctx).Table("department").
		Where("department_id = ?", id).
		Updates(data).Error
}

func (r *departmentStore) DeleteDepartment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("department").
		Where("department_id = ?", id).
		Delete(nil).Error
}

func (s *departmentStore) GetLastDepartmentByCode(ctx context.Context, office *model.Department) error {
	return s.db.WithContext(ctx).
		Order("department_id DESC").
		First(office).Error
}
