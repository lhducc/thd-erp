package repository

import (
	"context"
	model "erp/backend/internal/hrm/hr_profile/model"

	"errors"

	"gorm.io/gorm"
)

type OfficeStore struct {
	db *gorm.DB
}

func NewOfficeStore(db *gorm.DB) *OfficeStore {
	return &OfficeStore{db: db}
}

func (s *OfficeStore) CreateOffice(context context.Context, data *model.OfficeCreate) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (s *OfficeStore) GetOffice(ctx context.Context, id *string) (*model.Office, error) {
	var office model.Office
	if err := s.db.WithContext(ctx).Table("office").
		Where("office_id = ?", id).
		First(&office).Error; err != nil {
		return nil, err
	}
	return &office, nil
}

func (s *OfficeStore) GetAllOffice(ctx context.Context) ([]model.Office, error) {

	var positions []model.Office
	if err := s.db.WithContext(ctx).Table("office").Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (s *OfficeStore) UpdateOffice(ctx context.Context, id string, data *model.OfficeCreate) error {
	return s.db.WithContext(ctx).Table("office").
		Where("office_id = ?", id).
		Updates(data).Error
}

func (s *OfficeStore) DeleteOffice(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Table("office").
		Where("office_id = ?", id).
		Delete(nil).Error
}

func (s *OfficeStore) CheckExistName(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.Office{}).Where("office_name = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Office Name is already")
	}
	return false, nil
}

func (s *OfficeStore) GetLastOfficeByCode(ctx context.Context, office *model.Office) error {
	return s.db.WithContext(ctx).
		Order("office_id DESC").
		First(office).Error
}

func (s *OfficeStore) FindByID(id string) (model.Office, error) {
	var office model.Office
	err := s.db.Where("office_id = ?", id).First(&office).Error
	return office, err
}
