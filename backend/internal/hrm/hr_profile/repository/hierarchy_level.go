package repository

import (
	"context"
	model "erp/backend/internal/hrm/hr_profile/model"
	"errors"

	"gorm.io/gorm"
)

type hierarchyLevelStore struct {
	db *gorm.DB
}

func NewhierarchyLevelStore(db *gorm.DB) *hierarchyLevelStore {
	return &hierarchyLevelStore{db: db}
}

func (s *hierarchyLevelStore) CreateHierarchyLevel(context context.Context, data *model.HierarchyLevelCreate) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *hierarchyLevelStore) GetHierarchyLevel(ctx context.Context, id string) (*model.HierarchyLevel, error) {
	var hierarchyLevel model.HierarchyLevel
	if err := r.db.WithContext(ctx).Table("hierarchylevel").
		Where("id = ?", id).
		First(&hierarchyLevel).Error; err != nil {
		return nil, err
	}
	return &hierarchyLevel, nil
}

func (r *hierarchyLevelStore) GetAllHierarchyLevel(ctx context.Context) ([]model.HierarchyLevel, error) {

	var positions []model.HierarchyLevel
	if err := r.db.WithContext(ctx).Table("hierarchylevel").Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *hierarchyLevelStore) UpdateHierarchyLevel(ctx context.Context, id string, data *model.HierarchyLevelCreate) error {
	return r.db.WithContext(ctx).Table("hierarchylevel").
		Where("id = ?", id).
		Updates(data).Error
}

func (r *hierarchyLevelStore) DeleteHierarchyLevel(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("hierarchylevel").
		Where("id = ?", id).
		Delete(nil).Error
}

func (s *hierarchyLevelStore) CheckExistHierarchyLevel(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.HierarchyLevel{}).Where("hierarchy_level = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Hierarchy Level is already")
	}
	return false, nil
}
func (s *hierarchyLevelStore) GetLastHierarchyLevelyCode(ctx context.Context, office *model.HierarchyLevel) error {
	return s.db.WithContext(ctx).
		Order("id DESC").
		First(office).Error
}
