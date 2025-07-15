package service

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
)

type WorkShiftService struct {
	repo repo_interface.WorkShiftRepo
}

func NewWorkShiftService(repo *repository.WorkShiftStore) *WorkShiftService {
	return &WorkShiftService{repo: repo}
}

func (ws *WorkShiftService) CreateWorkShiftService(w *model.WorkShifts) error {
	var (
		code string
		err  error
	)

	// TimeOfDayEnum -> predix
	timeOfDayPrefixes := map[model.TimeOfDayEnum]string{
		model.FullTime:  "CN",
		model.Morning:   "CS",
		model.Afternoon: "CC",
		model.Night:     "CT",
	}
	predix, ok := timeOfDayPrefixes[w.TimeOfDay]
	if !ok {
		return errors.New("thời gian làm việc không hợp lệ: chỉ chấp nhận 'Cả ngày', 'Sáng', 'Chiều', 'Tối'")
	}

	// auto create ID
	code, err = utils.GenerateCode(predix, 3, func() (string, error) {
		var last model.WorkShifts
		err := ws.repo.GetLastWorkShiftByCode(&last, predix)
		if err != nil {
			return "", err
		}
		return last.WorkShiftID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	// check exit start,end time
	wData, err := ws.repo.IsExactTimeRangeExists(w.TimeOfDay, w.StartTime, w.EndTime)
	if err != nil {
		return err
	}
	if wData != nil {
		return fmt.Errorf("khung giờ %s-%s đã tồn tại trong ca buổi '%s' mã ca '%s'",
			w.StartTime, w.EndTime, w.TimeOfDay, wData.WorkShiftID)
	}

	// call repo create new
	w.WorkShiftID = code
	if err := ws.repo.CreateWorkShift(w); err != nil {
		return err
	}

	return nil
}

func (biz *WorkShiftService) GetWorkShiftByIdService(id string) (model.WorkShifts, error) {
	if id == "" {
		return model.WorkShifts{}, errors.New("invalid workshift ID")
	}

	workshift, err := biz.repo.GetWorkShiftById(id)
	if err != nil {
		return model.WorkShifts{}, fmt.Errorf("failed to get workshift: %w", err)
	}

	return workshift, nil
}

func (biz *WorkShiftService) GetAllWorkShiftService() ([]model.WorkShifts, error) {
	workshifts, err := biz.repo.GetAllWorkShift()
	if err != nil {
		return nil, fmt.Errorf("failed to get all workshifts: %w", err)
	}

	return workshifts, nil
}

func (biz *WorkShiftService) UpdateWorkShiftService(id string, workshift *model.WorkShifts) error {
	if id == "" {
		return errors.New("invalid employeeRepo ID")
	}

	if err := biz.repo.UpdateWorkShift(id, workshift); err != nil {
		return fmt.Errorf("failed to update employeeRepo: %w", err)
	}

	return nil
}

func (biz *WorkShiftService) DeleteWorkShiftService(id string) error {
	if id == "" {
		return errors.New("invalid employeeRepo ID")
	}
	if err := biz.repo.DeleteWorkShift(id); err != nil {
		return fmt.Errorf("failed to delete employeeRepo: %w", err)
	}

	return nil
}
