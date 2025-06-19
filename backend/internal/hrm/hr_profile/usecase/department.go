package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"fmt"
)

type DepartmentRepo interface {
	CreateDepartment(context context.Context, data *model.DepartmentCreate) error
	GetDepartment(ctx context.Context, id string) (*model.Department, error)
	GetAllDepartment(ctx context.Context) ([]model.Department, error)
	UpdateDepartment(ctx context.Context, id string, data *model.DepartmentCreate) error
	DeleteDepartment(ctx context.Context, id string) error
	GetLastDepartmentByCode(ctx context.Context, office *model.Department) error
}

type departmentBiz struct {
	repo DepartmentRepo
	// office OfficeRepo
}

func NewDepartmentBiz(store DepartmentRepo) *departmentBiz {
	return &departmentBiz{repo: store}
}

func (biz *departmentBiz) CreateDepartment(context context.Context, data *model.DepartmentCreate) error {
	data.CreatedDate = utils.GetCurrentDate()
	code, err := utils.GenerateCode("PB", 4, func() (string, error) {
		var last model.Department
		err := biz.repo.GetLastDepartmentByCode(context, &last)
		if err != nil {
			return "", err
		}
		return last.ID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	data.ID = code
	if err := biz.repo.CreateDepartment(context, data); err != nil {
		return err
	}

	return nil
}

func (biz *departmentBiz) GetDepartment(ctx context.Context, id string) (*model.Department, error) {
	office, err := biz.repo.GetDepartment(ctx, id)
	if err != nil {
		return nil, err
	}
	return office, nil
}

func (biz *departmentBiz) GetAllDepartment(ctx context.Context) ([]model.Department, error) {
	positions, err := biz.repo.GetAllDepartment(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (biz *departmentBiz) UpdateDepartment(ctx context.Context, id string, data *model.DepartmentCreate) error {
	if err := biz.repo.UpdateDepartment(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *departmentBiz) DeleteDepartment(ctx context.Context, id string) error {
	if err := biz.repo.DeleteDepartment(ctx, id); err != nil {
		return err
	}
	return nil
}
