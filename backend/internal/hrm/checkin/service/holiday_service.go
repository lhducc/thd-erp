package service

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
)

type holidayService struct {
	repo *repository.HolidayRepo
}

func NewHolidayService(repo *repository.HolidayRepo) *holidayService {
	return &holidayService{repo: repo}
}

func (ws *holidayService) CreateHoliday(w model.Holiday) error {
	code, err := utils.GenerateCode("", 4, func() (string, error) {
		var last model.Holiday
		err := ws.repo.GetLastHolidayByCode(&last)
		if err != nil {
			return "", err
		}
		return last.HolidayID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	w.HolidayID = code
	if err := ws.repo.CreateHoliday(&w); err != nil {
		return err
	}

	return nil
}
func (biz *holidayService) GetHolidayById(id string) (model.Holiday, error) {
	if id == "" {
		return model.Holiday{}, errors.New("invalid holiday ID")
	}

	holiday, err := biz.repo.GetHolidayById(id)
	if err != nil {
		return model.Holiday{}, fmt.Errorf("failed to get holiday: %w", err)
	}

	return holiday, nil
}

func (biz *holidayService) GetAllHoliday() ([]model.Holiday, error) {
	holidays, err := biz.repo.GetAllHoliday()
	if err != nil {
		return nil, fmt.Errorf("failed to get all holidays: %w", err)
	}

	return holidays, nil
}

func (biz *holidayService) UpdateHoliday(id string, holiday model.Holiday) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}

	if err := biz.repo.UpdateHoliday(id, holiday); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (biz *holidayService) DeleteHoliday(id string) error {
	if id == "" {
		return errors.New("invalid holiday ID")
	}
	if err := biz.repo.DeleteHoliday(id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}
