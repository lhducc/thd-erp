package repo_interface

import "erp/backend/internal/hrm/checkin/model"

type WorkShiftRepo interface {
	CreateWorkShift(data *model.WorkShifts) error
	GetWorkShiftById(id string) (model.WorkShifts, error)
	GetAllWorkShift() ([]model.WorkShifts, error)
	UpdateWorkShift(id string, data *model.WorkShifts) error
	DeleteWorkShift(id string) error
	GetLastWorkShiftByCode(emp *model.WorkShifts, predix string) error
	IsExactTimeRangeExists(timeOfDay model.TimeOfDayEnum, start, end string) (*model.WorkShifts, error)
}
