package service

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
)

type allowedWorkingScheduleService struct {
	repo *repository.AllowedWorkingScheduleRepo
}

func NewAllowedWorkingScheduleService(repo *repository.AllowedWorkingScheduleRepo) *allowedWorkingScheduleService {
	return &allowedWorkingScheduleService{repo: repo}
}

func (ws *allowedWorkingScheduleService) CreateAllowedWorkingSchedule(w model.AllowedWorkingSchedule) error {
	code, err := utils.GenerateCode("", 4, func() (string, error) {
		var last model.AllowedWorkingSchedule
		err := ws.repo.GetLastAllowedWorkingScheduleByCode(&last)
		if err != nil {
			return "", err
		}
		return last.ID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	w.ID = code
	if err := ws.repo.CreateAllowedWorkingSchedule(&w); err != nil {
		return err
	}

	return nil
}
func (biz *allowedWorkingScheduleService) GetAllowedWorkingScheduleById(id string) (model.AllowedWorkingSchedule, error) {
	if id == "" {
		return model.AllowedWorkingSchedule{}, errors.New("invalid holiday ID")
	}

	holiday, err := biz.repo.GetAllowedWorkingScheduleId(id)
	if err != nil {
		return model.AllowedWorkingSchedule{}, fmt.Errorf("failed to get holiday: %w", err)
	}

	return holiday, nil
}

func (biz *allowedWorkingScheduleService) GetAllAllowedWorkingSchedule() ([]model.AllowedWorkingSchedule, error) {
	holidays, err := biz.repo.GetAllAllowedWorkingSchedule()
	if err != nil {
		return nil, fmt.Errorf("failed to get all holidays: %w", err)
	}

	return holidays, nil
}

func (biz *allowedWorkingScheduleService) UpdateAllowedWorkingSchedule(id string, holiday model.AllowedWorkingSchedule) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}

	if err := biz.repo.UpdateAllowedWorkingSchedule(id, holiday); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (biz *allowedWorkingScheduleService) DeleteAllowedWorkingScheduleService(id string) error {
	if id == "" {
		return errors.New("invalid holiday ID")
	}
	if err := biz.repo.DeleteAllowedWorkingSchedule(id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}
