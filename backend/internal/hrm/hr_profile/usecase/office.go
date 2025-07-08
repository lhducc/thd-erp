package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"fmt"
)

type OfficeRepo interface {
	CreateOffice(context context.Context, data *model.OfficeCreate) error
	GetOffice(ctx context.Context, id *string) (*model.Office, error)
	GetAllOffice(ctx context.Context) ([]model.Office, error)
	UpdateOffice(ctx context.Context, id string, data *model.OfficeCreate) error
	DeleteOffice(ctx context.Context, id string) error
	CheckExistName(name string) (bool, error)
	GetLastOfficeByCode(ctx context.Context, office *model.Office) error
}

type officeBiz struct {
	repo OfficeRepo
}

func NewOfficeBiz(store OfficeRepo) *officeBiz {
	return &officeBiz{repo: store}
}

func (biz *officeBiz) CreateOffice(context context.Context, data *model.OfficeCreate) error {
	data.CreatedDate = utils.GetCurrentDate()
	check, err := biz.repo.CheckExistName(data.Name)
	if err != nil || check {
		return err
	}

	code, err := utils.GenerateCode("VP", 3, func() (string, error) {
		var last model.Office
		err := biz.repo.GetLastOfficeByCode(context, &last)
		if err != nil {
			return "", err
		}
		return last.ID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	data.ID = code

	if err := biz.repo.CreateOffice(context, data); err != nil {
		return err
	}

	return nil
}

func (biz *officeBiz) GetOffice(ctx context.Context, id *string) (*model.Office, error) {
	office, err := biz.repo.GetOffice(ctx, id)
	if err != nil {
		return nil, err
	}
	return office, nil
}
func (biz *officeBiz) UpdateOffice(ctx context.Context, id string, data *model.OfficeCreate) error {
	if err := biz.repo.UpdateOffice(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *officeBiz) DeleteOffice(ctx context.Context, id string) error {
	if err := biz.repo.DeleteOffice(ctx, id); err != nil {
		return err
	}
	return nil
}

func (biz *officeBiz) GetAllOffice(ctx context.Context) ([]model.Office, error) {
	positions, err := biz.repo.GetAllOffice(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}
