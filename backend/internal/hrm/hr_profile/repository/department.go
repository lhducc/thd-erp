package repository

import (
	"context"
	model "erp/backend/internal/hrm/hr_profile/model"

	"gorm.io/gorm"
)

type DepartmentStore struct {
	db *gorm.DB
}

func NewDepartmentStore(db *gorm.DB) *DepartmentStore {
	return &DepartmentStore{db: db}
}

func (s *DepartmentStore) CreateDepartment(context context.Context, data *model.DepartmentCreate) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (s *DepartmentStore) GetDepartment(ctx context.Context, id string) (*model.Department, error) {
	var department model.Department
	if err := s.db.WithContext(ctx).Table("department").Preload("Office").
		Where("department_id = ?", id).
		First(&department).Error; err != nil {
		return nil, err
	}
	return &department, nil
}

func (s *DepartmentStore) GetAllDepartment(ctx context.Context) ([]model.Department, error) {

	var positions []model.Department
	if err := s.db.WithContext(ctx).
		Table("department").
		Preload("Office").
		Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (s *DepartmentStore) UpdateDepartment(ctx context.Context, id string, data *model.DepartmentCreate) error {
	return s.db.WithContext(ctx).Table("department").
		Where("department_id = ?", id).
		Updates(data).Error
}

func (s *DepartmentStore) DeleteDepartment(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Table("department").
		Where("department_id = ?", id).
		Delete(nil).Error
}

func (s *DepartmentStore) GetLastDepartmentByCode(ctx context.Context, office *model.Department) error {
	return s.db.WithContext(ctx).
		Order("department_id DESC").
		First(office).Error
}

func (s *DepartmentStore) GetDepartmentByOfficeID(officeID string) (model.Department, error) {
	var department model.Department
	err := s.db.Where("office_id = ?", officeID).Find(&department).Error
	return department, err
}
