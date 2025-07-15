package model

import (
	"gorm.io/gorm"
	"time"
)

type EmployeeWorkshift struct {
	gorm.Model
	EmployeeID  string     `gorm:"column:employee_id" json:"employee_id"`
	WorkShiftID string     `gorm:"column:work_shift_id" json:"work_shift_id"`
	WorkShift   WorkShifts `gorm:"foreignKey:work_shift_id" json:"work_shift"`
	Date        time.Time  `gorm:"column:date" json:"date"`
}

func (EmployeeWorkshift) TableName() string {
	return "employee_workshift"
}
