package service

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
)

type WorkShiftService struct {
	repo *repository.WorkShiftStore
}

func NewWorkShiftService(repo *repository.WorkShiftStore) *WorkShiftService {
	return &WorkShiftService{repo: repo}
}

func (ws *WorkShiftService) CreateWorkShift(w model.WorkShifts) error {
	code, err := utils.GenerateCode("", 4, func() (string, error) {
		var last model.WorkShifts
		err := ws.repo.GetLastWorkShiftByCode(&last)
		if err != nil {
			return "", err
		}
		return last.WorkShiftId, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	w.WorkShiftId = code
	if err := ws.repo.CreateWorkShift(&w); err != nil {
		return err
	}

	return nil
}
func (biz *WorkShiftService) GetWorkShiftById(id string) (model.WorkShifts, error) {
	if id == "" {
		return model.WorkShifts{}, errors.New("invalid workshift ID")
	}

	workshift, err := biz.repo.GetWorkShiftById(id)
	if err != nil {
		return model.WorkShifts{}, fmt.Errorf("failed to get workshift: %w", err)
	}

	return workshift, nil
}

func (biz *WorkShiftService) GetAllWorkShift() ([]model.WorkShifts, error) {
	workshifts, err := biz.repo.GetAllWorkShift()
	if err != nil {
		return nil, fmt.Errorf("failed to get all workshifts: %w", err)
	}

	return workshifts, nil
}

func (biz *WorkShiftService) UpdateWorkShift(id string, workshift model.WorkShifts) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}

	if err := biz.repo.UpdateWorkShift(id, workshift); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (biz *WorkShiftService) DeleteWorkShift(id string) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}
	if err := biz.repo.DeleteWorkShift(id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}
