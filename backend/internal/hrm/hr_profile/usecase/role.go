package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
)

type RoleUsecase struct {
	repo *repository.RoleRepo
}

func NewRoleUsecase(repo *repository.RoleRepo) *RoleUsecase {
	return &RoleUsecase{
		repo: repo,
	}
}

func (s *RoleUsecase) FetchAllRole(ctx context.Context) ([]model.Role, error) {
	return s.repo.GetAll(ctx)
}
