package model

import (
	"errors"
	"time"
)

// DecisionGroup ENUM for decision group type
type DecisionGroup string

const (
	GroupAward       DecisionGroup = "Hình thức khen thưởng" // Group: Award
	GroupDiscipline  DecisionGroup = "Hình thức kỷ luật"     // Group: Discipline
	GroupTransfer    DecisionGroup = "Lý do điều chuyển"     // Group: Transfer Reason
	GroupReception   DecisionGroup = "Lý do tiếp nhận"       // Group: Reception Reason
	GroupAppointment DecisionGroup = "Lý do bổ nhiệm"        // Group: Appointment Reason
	GroupDismissal   DecisionGroup = "Lý do miễn nhiệm"      // Group: Dismissal Reason
	GroupContractEnd DecisionGroup = "Lý do chấm dứt HĐLĐ"   // Group: End of Labor Contract Reason
)

// DecisionType struct represents the type of decision
type DecisionType struct {
	DecisionTypeID string        `gorm:"type:char(6);primaryKey" json:"decision_type_id"`         // ID of decision type, e.g., "LQ0001"
	DecisionType   string        `gorm:"type:varchar(100);not null" json:"decision_type"`         // Name of decision type
	DecisionGroup  DecisionGroup `gorm:"type:decision_group_enum;not null" json:"decision_group"` // Decision group (e.g., Award, Discipline, etc.)
	Description    string        `gorm:"type:text" json:"description"`                            // Description of decision type
	CreatedDate    time.Time     `gorm:"not null" json:"created_date"`                            // Date when the record was created

	// One-to-many relationship: One DecisionType can have many Decisions
	Decisions []Decision `gorm:"foreignKey:DecisionTypeID" json:"decisions,omitempty"`
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

	validGroups := map[DecisionGroup]bool{
		GroupAward: true, GroupDiscipline: true, GroupTransfer: true,
		GroupReception: true, GroupAppointment: true,
		GroupDismissal: true, GroupContractEnd: true,
	}

	if !validGroups[dt.DecisionGroup] {
		return errors.New("Nhóm quyết định không hợp lệ")
	}

	if dt.CreatedDate.IsZero() {
		return errors.New("Ngày tạo không được để trống")
	}

	return nil
}
