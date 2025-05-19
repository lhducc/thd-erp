package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"

	"gorm.io/gorm"
)

type positionStore struct {
	db *gorm.DB
}

func NewPositionStore(db *gorm.DB) *positionStore {
	return &positionStore{db: db}
}

func (s *positionStore) CreatePosition(context context.Context, data *model.Position) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *positionStore) GetPosition(ctx context.Context, id string) (*model.Position, error) {
	var position model.Position
	if err := r.db.WithContext(ctx).Table("position").
		Where("position_id = ?", id).
		First(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *positionStore) GetAllPositions(ctx context.Context) ([]model.Position, error) {

	var positions []model.Position
	if err := r.db.WithContext(ctx).Table("position").Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *positionStore) UpdatePosition(ctx context.Context, id string, data *model.Position) error {
	return r.db.WithContext(ctx).Table("position").
		Where("position_id = ?", id).
		Updates(data).Error
}

func (r *positionStore) DeletePosition(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("position").
		Where("position_id = ?", id).
		Delete(nil).Error
}

func (s *positionStore) CheckExistPosition(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.Position{}).Where("position_name = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Position is already")
	}
	return false, nil
}

func (s *positionStore) GetLastPositionByCode(ctx context.Context, office *model.Position) error {
	return s.db.WithContext(ctx).
		Order("position_id DESC").
		First(office).Error
}
