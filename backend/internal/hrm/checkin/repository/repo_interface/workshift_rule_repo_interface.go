package repo_interface

import (
	"erp/backend/internal/hrm/checkin/model"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

type WorkshiftRuleRepo interface {
	Create(assignment *model.WorkshiftRule) error
	DB() *gorm.DB
	GetAll() ([]model.WorkshiftRule, error)
	GetByUserID(employee hrmmodel.Employee) ([]model.WorkshiftRule, error)
	Update(id uint, assignment model.WorkshiftRule) error
	Delete(id uint) error
	GetByID(id uint) (model.WorkshiftRule, error)
}
