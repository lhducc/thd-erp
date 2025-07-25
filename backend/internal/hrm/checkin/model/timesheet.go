package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"time"
)

type StatusEnum int

const ()

type Timesheet struct {
	ID          string    `gorm:"type:varchar;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	OfficeID    string    `gorm:"type:varchar;not null" json:"office_id"`
	Month       int       `gorm:"not null;index;check:month >= 1 AND month <= 12" json:"month"`
	Year        int       `gorm:"not null;index" json:"year"`
	IsFinalized bool      `gorm:"default:false" json:"is_finalized"`
	CreatedBy   string    `gorm:"type:varchar" json:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Office model.Office `gorm:"foreignKey:OfficeID" json:"office"`
}

func (StatusEnum) TableName() string {
	return "timesheets"
}
