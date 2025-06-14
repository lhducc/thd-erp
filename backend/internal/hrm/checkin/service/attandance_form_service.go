package service

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
)

type attandanceFormService struct {
	repo *repository.AttendanceFormStore
}

func NewAttandanceFormService(repo *repository.AttendanceFormStore) *attandanceFormService {
	return &attandanceFormService{repo: repo}
}

func (ws *attandanceFormService) CreateAttandanceForm(w model.AttendanceForm) error {
	code, err := utils.GenerateCode("", 4, func() (string, error) {
		var last model.AttendanceForm
		err := ws.repo.GetLastAttandanceFormByCode(&last)
		if err != nil {
			return "", err
		}
		return last.AttendanceFormID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	w.AttendanceFormID = code
	if err := ws.repo.CreateAttendanceForm(&w); err != nil {
		return err
	}

	return nil
}

func (biz *attandanceFormService) GetAttandanceFormById(id string) (model.AttendanceForm, error) {
	if id == "" {
		return model.AttendanceForm{}, errors.New("invalid form ID")
	}

	holiday, err := biz.repo.GetAttendanceFormByID(id)
	if err != nil {
		return model.AttendanceForm{}, fmt.Errorf("failed to get form: %w", err)
	}

	return holiday, nil
}

func (biz *attandanceFormService) GetAllAttandanceForm() ([]model.AttendanceForm, error) {
	holidays, err := biz.repo.GetAllAttendanceForm()
	if err != nil {
		return nil, fmt.Errorf("failed to get all form: %w", err)
	}

	return holidays, nil
}

func (biz *attandanceFormService) UpdateAttandanceForm(id string, holiday model.AttendanceForm) error {
	if id == "" {
		return errors.New("invalid form ID")
	}

	if err := biz.repo.UpdateAttendanceForm(id, holiday); err != nil {
		return fmt.Errorf("failed to update form: %w", err)
	}

	return nil
}

func (biz *attandanceFormService) DisableAttandanceForm(id string) error {
	if id == "" {
		return errors.New("invalid form ID")
	}
	if err := biz.repo.DisableAttandanceForm(id); err != nil {
		return fmt.Errorf("failed to disable form: %w", err)
	}

	return nil
}