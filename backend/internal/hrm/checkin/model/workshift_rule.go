package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
)

type WorkshiftRule struct {
	WorkshiftRuleID uint `gorm:"primaryKey;column:workshift_rule_id;autoIncrement" json:"workshift_rule_id"`

	Name string `gorm:"type:varchar(255);not null" json:"name"`

	Office     []model.Office     `gorm:"many2many:workshift_rule_offices;joinForeignKey:WorkshiftRuleID;joinReferences:OfficeID" json:"offices"`
	Department []model.Department `gorm:"many2many:workshift_rule_departments;joinForeignKey:WorkshiftRuleID;joinReferences:DepartmentID" json:"departments"`
	Position   []model.Position   `gorm:"many2many:workshift_rule_positions;joinForeignKey:WorkshiftRuleID;joinReferences:PositionID" json:"positions"`
	JobTitle   []model.JobTitle   `gorm:"many2many:workshift_rule_job_titles;joinForeignKey:WorkshiftRuleID;joinReferences:JobTitleID" json:"job_title"`

	StartDate string `gorm:"not null" json:"start_date"`
	EndDate   string `gorm:"not null" json:"end_date"`
	//Description string `gorm:"type:text;not null" json:"description"`

	WorkShifts []WorkShifts `gorm:"many2many:workshift_rule_workshifts;" json:"work_shifts"`
}

func (WorkshiftRule) TableName() string { return "workshift_rule" }
