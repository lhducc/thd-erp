package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"fmt"
)

type HierarchyLevelRepo interface {
	CreateHierarchyLevel(context context.Context, data *model.HierarchyLevelCreate) error
	GetHierarchyLevel(ctx context.Context, id string) (*model.HierarchyLevel, error)
	GetAllHierarchyLevel(ctx context.Context) ([]model.HierarchyLevel, error)
	UpdateHierarchyLevel(ctx context.Context, id string, data *model.HierarchyLevelCreate) error
	DeleteHierarchyLevel(ctx context.Context, id string) error
	CheckExistHierarchyLevel(name string) (bool, error)
	GetLastHierarchyLevelyCode(ctx context.Context, office *model.HierarchyLevel) error
}

type hierarchyLevelBiz struct {
	repo HierarchyLevelRepo
}

func NewHierarchyLevelBiz(repo HierarchyLevelRepo) *hierarchyLevelBiz {
	return &hierarchyLevelBiz{repo: repo}
}

func (biz *hierarchyLevelBiz) CreateHierarchyLevel(context context.Context, data *model.HierarchyLevelCreate) error {
	data.CreatedDate = utils.GetCurrentDate()
	check, err := biz.repo.CheckExistHierarchyLevel(data.HierarchyLevel)

	if err != nil || check {
		return err
	}

	code, err := utils.GenerateCode("CB", 4, func() (string, error) {
		var last model.HierarchyLevel
		err := biz.repo.GetLastHierarchyLevelyCode(context, &last)
		if err != nil {
			return "", err
		}
		return last.ID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	data.ID = code
	if err := biz.repo.CreateHierarchyLevel(context, data); err != nil {
		return err
	}

	return nil
}

func (biz *hierarchyLevelBiz) GetHierarchyLevel(ctx context.Context, id string) (*model.HierarchyLevel, error) {
	office, err := biz.repo.GetHierarchyLevel(ctx, id)
	if err != nil {
		return nil, err
	}
	return office, nil
}

func (biz *hierarchyLevelBiz) GetAllHierarchyLevel(ctx context.Context) ([]model.HierarchyLevel, error) {
	positions, err := biz.repo.GetAllHierarchyLevel(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (biz *hierarchyLevelBiz) UpdateHierarchyLevel(ctx context.Context, id string, data *model.HierarchyLevelCreate) error {
	if err := biz.repo.UpdateHierarchyLevel(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *hierarchyLevelBiz) DeleteHierarchyLevel(ctx context.Context, id string) error {
	if err := biz.repo.DeleteHierarchyLevel(ctx, id); err != nil {
		return err
	}
	return nil
}
