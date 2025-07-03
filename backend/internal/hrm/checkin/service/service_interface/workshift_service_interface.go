package service_interface

import "erp/backend/internal/hrm/checkin/model"

type WorkShiftService interface {
	CreateWorkShiftService(data *model.WorkShifts) error
	GetWorkShiftByIdService(id string) (model.WorkShifts, error)
	GetAllWorkShiftService() ([]model.WorkShifts, error)
	UpdateWorkShiftService(id string, data *model.WorkShifts) error
	DeleteWorkShiftService(id string) error
}
