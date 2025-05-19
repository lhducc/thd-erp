package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"

	"fmt"
)

type PositionRepo interface {
	CreatePosition(context context.Context, data *model.Position) error
	GetPosition(ctx context.Context, id string) (*model.Position, error)
	GetAllPositions(ctx context.Context) ([]model.Position, error)
	UpdatePosition(ctx context.Context, id string, data *model.Position) error
	DeletePosition(ctx context.Context, id string) error
	CheckExistPosition(name string) (bool, error)
	GetLastPositionByCode(ctx context.Context, office *model.Position) error
}

type positionBiz struct {
	repo PositionRepo
}

func NewPositionBiz(repo PositionRepo) *positionBiz {
	return &positionBiz{repo: repo}
}

func (biz *positionBiz) CreatePosition(context context.Context, data *model.Position) error {
	data.CreatedDate = utils.GetCurrentDate()
	check, err := biz.repo.CheckExistPosition(data.Name)
	if err != nil || check {
		return err
	}

	code, err := utils.GenerateCode("VT", 4, func() (string, error) {
		var last model.Position
		err := biz.repo.GetLastPositionByCode(context, &last)
		if err != nil {
			return "", err
		}
		return last.ID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	data.ID = code

	if err := biz.repo.CreatePosition(context, data); err != nil {
		return err
	}

	return nil
}

func (biz *positionBiz) GetPosition(ctx context.Context, id string) (*model.Position, error) {
	office, err := biz.repo.GetPosition(ctx, id)
	if err != nil {
		return nil, err
	}
	return office, nil
}

func (biz *positionBiz) GetAllPositions(ctx context.Context) ([]model.Position, error) {
	positions, err := biz.repo.GetAllPositions(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (biz *positionBiz) UpdatePosition(ctx context.Context, id string, data *model.Position) error {
	if err := biz.repo.UpdatePosition(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *positionBiz) DeletePosition(ctx context.Context, id string) error {
	if err := biz.repo.DeletePosition(ctx, id); err != nil {
		return err
	}
	return nil
}
