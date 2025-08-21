package repository

import (
	"context"
	model "erp/backend/internal/hrm/hr_profile/model"
	"errors"

	"gorm.io/gorm"
)

type JobTitleStore struct {
	db *gorm.DB
}

func NewJobTitleStore(db *gorm.DB) *JobTitleStore {
	return &JobTitleStore{db: db}
}

func (s *JobTitleStore) CreateJobTitle(context context.Context, data *model.JobTitleCreate) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (s *JobTitleStore) GetJobTitleById(ctx context.Context, id string) (*model.JobTitle, error) {
	var jobTitle model.JobTitle
	if err := s.db.WithContext(ctx).Table("jobtitle").
		Preload("HierarchyLevel").
		Where("job_title_id = ?", id).
		First(&jobTitle).Error; err != nil {
		return nil, err
	}
	return &jobTitle, nil
}

func (s *JobTitleStore) GetAllJobTitle(ctx context.Context) ([]model.JobTitle, error) {

	var positions []model.JobTitle
	if err := s.db.WithContext(ctx).Table("jobtitle").
		Preload("HierarchyLevel").
		Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (s *JobTitleStore) UpdateJobTitleById(ctx context.Context, id string, data *model.JobTitleCreate) error {
	return s.db.WithContext(ctx).Table("jobtitle").
		Where("job_title_id = ?", id).
		Updates(data).Error
}

func (s *JobTitleStore) DeleteJobTitleById(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Table("jobtitle").
		Where("job_title_id = ?", id).
		Delete(nil).Error
}

func (s *JobTitleStore) CheckExistJobTitle(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.JobTitle{}).Where("job_title = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Job title is already")
	}
	return false, nil
}

func (s *JobTitleStore) GetLastJobTitleByCode(ctx context.Context, office *model.JobTitle) error {
	return s.db.WithContext(ctx).
		Order("job_title_id DESC").
		First(office).Error
}

func (s *JobTitleStore) FindByID(id string) (model.JobTitle, error) {
	var jobTitle model.JobTitle
	if err := s.db.Where("job_title_id = ?", id).First(&jobTitle).Error; err != nil {
		return model.JobTitle{}, err
	}
	return jobTitle, nil
}
