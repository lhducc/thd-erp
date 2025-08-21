package model

import (
	"errors"
	"time"
)

// DecisionType struct represents the type of decision
type DecisionType struct {
	DecisionTypeID string    `gorm:"type:char(6);primaryKey" json:"decision_type_id"` // ID of decision type, e.g., "LQ0001"
	DecisionType   string    `gorm:"type:varchar(100);not null" json:"decision_type"` // Name of decision type
	Description    string    `gorm:"type:text" json:"description"`                    // Description of decision type
	CreatedDate    time.Time `gorm:"not null" json:"created_date"`                    // Date when the record was created
}

func (DecisionType) TableName() string {
	return "decisiontype"
}

// ValidateDecisionType checks required fields and validates DecisionGroup enum and CreatedDate
func (dt DecisionType) ValidateDecisionType() error {
	if dt.DecisionTypeID == "" {
		return errors.New("Mã loại quyết định không được để trống")
	}
	if dt.DecisionType == "" {
		return errors.New("Tên loại quyết định không được để trống")
	}
	if dt.CreatedDate.IsZero() {
		return errors.New("Ngày tạo không được để trống")
	}

	return nil
}
