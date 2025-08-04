package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.WithContext(ctx).Model(&model.Role{}).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
