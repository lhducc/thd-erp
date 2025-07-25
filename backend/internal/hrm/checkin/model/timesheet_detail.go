package model

import (
	"github.com/google/uuid"
	"time"
)

type TimesheetDetail struct {
	ID             int        `gorm:"type:serial;primaryKey" json:"id"`
	TimesheetID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"timesheet_id"`
	EmployeeID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"employee_id"`
	WorkDate       time.Time  `gorm:"type:date;not null;index" json:"work_date"`
	ExpectedShift  string     `gorm:"type:varchar(50)" json:"expected_shift"` // ca dự kiến (VD: "08:00-17:00")
	ActualCheckin  *time.Time `json:"actual_checkin"`
	ActualCheckout *time.Time `json:"actual_checkout"`
	WorkingHours   float64    `gorm:"type:numeric(5,2);default:0" json:"working_hours"`
	WorkdayPoint   float64    `gorm:"type:numeric(3,2);default:0" json:"workday_point"` // VD: 1.0 công, 0.5 công,...
	Note           string     `gorm:"type:text" json:"note"`
}

func (TimesheetDetail) TableName() string {
	return "timesheet_details"
}
