package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"

	"gorm.io/gorm"
)

type PositionStore struct {
	db *gorm.DB
}

func NewPositionStore(db *gorm.DB) *PositionStore {
	return &PositionStore{db: db}
}

func (s *PositionStore) CreatePosition(context context.Context, data *model.Position) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (s *PositionStore) GetPosition(ctx context.Context, id string) (*model.Position, error) {
	var position model.Position
	if err := s.db.WithContext(ctx).Table("position").
		Where("position_id = ?", id).
		First(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func (s *PositionStore) GetAllPositions(ctx context.Context) ([]model.Position, error) {

	var positions []model.Position
	if err := s.db.WithContext(ctx).Table("position").Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (s *PositionStore) UpdatePosition(ctx context.Context, id string, data *model.Position) error {
	return s.db.WithContext(ctx).Table("position").
		Where("position_id = ?", id).
		Updates(data).Error
}

func (s *PositionStore) DeletePosition(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Table("position").
		Where("position_id = ?", id).
		Delete(nil).Error
}

func (s *PositionStore) CheckExistPosition(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.Position{}).Where("position_name = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Position is already")
	}
	return false, nil
}

func (s *PositionStore) GetLastPositionByCode(ctx context.Context, office *model.Position) error {
	return s.db.WithContext(ctx).
		Order("position_id DESC").
		First(office).Error
}

func (s *PositionStore) FindByID(id string) (model.Position, error) {
	var position model.Position
	if err := s.db.Where("position_id = ?", id).First(&position).Error; err != nil {
		return model.Position{}, err
	}
	return position, nil
}
