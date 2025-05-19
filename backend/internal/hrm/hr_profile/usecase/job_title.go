package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"fmt"
)

type JobTitleRepo interface {
	CreateJobTitle(context context.Context, data *model.JobTitleCreate) error
	GetJobTitleById(ctx context.Context, id string) (*model.JobTitle, error)
	GetAllJobTitle(ctx context.Context) ([]model.JobTitle, error)
	UpdateJobTitleById(ctx context.Context, id string, data *model.JobTitleCreate) error
	DeleteJobTitleById(ctx context.Context, id string) error
	CheckExistJobTitle(name string) (bool, error)
	GetLastJobTitleByCode(ctx context.Context, office *model.JobTitle) error
}

type jobTitleBiz struct {
	repo    JobTitleRepo
	hieRepo HierarchyLevelRepo
}

func NewJobTitleBiz(repo JobTitleRepo, hieRepo HierarchyLevelRepo) *jobTitleBiz {
	return &jobTitleBiz{
		repo:    repo,
		hieRepo: hieRepo,
	}
}

func (biz *jobTitleBiz) CreateJobTitle(context context.Context, data *model.JobTitleCreate) error {
	data.CreatedDate = utils.GetCurrentDate()
	check, err := biz.repo.CheckExistJobTitle(data.JobTitle)
	if err != nil || check {
		return err
	}

	code, err := utils.GenerateCode("CV", 4, func() (string, error) {
		var last model.JobTitle
		err := biz.repo.GetLastJobTitleByCode(context, &last)
		if err != nil {
			return "", err
		}
		return last.JobTitleID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	data.JobTitleID = code

	_, err = biz.hieRepo.GetHierarchyLevel(context, data.HierarchyLevelID)
	if err != nil {
		return err
	}

	if err := biz.repo.CreateJobTitle(context, data); err != nil {
		return err
	}

	return nil
}

func (biz *jobTitleBiz) GetJobTitleById(ctx context.Context, id string) (*model.JobTitle, error) {
	office, err := biz.repo.GetJobTitleById(ctx, id)
	if err != nil {
		return nil, err
	}
	return office, nil
}

func (biz *jobTitleBiz) GetAllJobTitle(ctx context.Context) ([]model.JobTitle, error) {
	positions, err := biz.repo.GetAllJobTitle(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (biz *jobTitleBiz) UpdateJobTitleById(ctx context.Context, id string, data *model.JobTitleCreate) error {
	if err := biz.repo.UpdateJobTitleById(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *jobTitleBiz) DeleteJobTitleById(ctx context.Context, id string) error {
	if err := biz.repo.DeleteJobTitleById(ctx, id); err != nil {
		return err
	}
	return nil
}
