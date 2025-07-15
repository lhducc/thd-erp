package repository

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

type workshiftRuleRepository struct {
	db *gorm.DB
}

func NewWorkshiftRuleRepository(db *gorm.DB) repo_interface.WorkshiftRuleRepo {
	return &workshiftRuleRepository{
		db: db,
	}
}

func (w *workshiftRuleRepository) DB() *gorm.DB {
	return w.db
}

func (w *workshiftRuleRepository) Create(assignment *model.WorkshiftRule) error {
	return w.db.Create(assignment).Error
}

func (w *workshiftRuleRepository) GetByID(id uint) (model.WorkshiftRule, error) {
	var workshift model.WorkshiftRule
	err := w.db.Where("id = ?", id).First(&workshift).Error
	return workshift, err
}

func (w *workshiftRuleRepository) GetAll() ([]model.WorkshiftRule, error) {
	var rules []model.WorkshiftRule
	err := w.db.Preload("Office").Preload("Department").Preload("Position").Preload("JobTitle").Find(&rules).Error
	return rules, err
}

func (w *workshiftRuleRepository) GetByUserID(employee hrmmodel.Employee) ([]model.WorkshiftRule, error) {
	var rules []model.WorkshiftRule

	err := w.db.
		Preload("Office").
		Preload("Department").
		Preload("Position").
		Preload("JobTitle").
		Preload("WorkShifts").
		Joins("LEFT JOIN workshift_rule_positions ON workshift_rule.workshift_rule_id = workshift_rule_positions.workshift_rule_id").
		Joins("LEFT JOIN workshift_rule_departments ON workshift_rule.workshift_rule_id = workshift_rule_departments.workshift_rule_id").
		Joins("LEFT JOIN workshift_rule_job_titles ON workshift_rule.workshift_rule_id = workshift_rule_job_titles.workshift_rule_id").
		Where("workshift_rule_positions.position_id = ? OR workshift_rule_departments.department_id = ? OR workshift_rule_job_titles.job_title_id = ?",
			employee.PositionID, employee.DepartmentID, employee.JobTitleID).
		Find(&rules).Error

	if err != nil {
		return nil, err
	}

	return rules, nil
}

func (w *workshiftRuleRepository) Update(id uint, assignment model.WorkshiftRule) error {
	var existing model.WorkshiftRule
	err := w.db.Where("id = ?", id).First(&existing).Error
	if err != nil {
		return err
	}
	err = w.db.Model(&existing).Updates(assignment).Error
	if err != nil {
		return err
	}
	return nil
}

func (w *workshiftRuleRepository) Delete(id uint) error {
	var existing model.WorkshiftRule
	err := w.db.Where("id = ?", id).First(&existing).Error
	if err != nil {
		return err
	}
	return w.db.Where("id = ?", id).Delete(&existing).Error
}
