package model

import (
	"time"
)

type EmployeeWorkshift struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeID  string    `gorm:"column:employee_id;index:idx_employee_date" json:"employee_id"`
	WorkShiftID string    `gorm:"column:workshift_id;type:varchar" json:"workshift_id"`
	Date        time.Time `gorm:"column:date;index:idx_employee_date;type:timestamp" json:"date"`
	CreatedAt   time.Time `gorm:"column:created_at;" json:"created_at"`

	WorkShift WorkShifts `gorm:"foreignKey:WorkShiftID;references:WorkShiftID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"workshift"`
}

func (EmployeeWorkshift) TableName() string {
	return "employee_workshift"
}
