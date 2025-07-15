package service_interface

import (
	"erp/backend/internal/hrm/checkin/model"
)

type WorkshiftRuleService interface {
	Create(assignment model.WorkshiftRule) error
	GetAll() ([]model.WorkshiftRule, error)
	GetByUserID(userID string) ([]model.WorkshiftRule, error)
	Update(id uint, assignment model.WorkshiftRule) error
	Delete(id uint) error
}
