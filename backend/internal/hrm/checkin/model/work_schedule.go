package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg/variable"
	"time"
)

type WorkSchedule struct {
	WorkScheduleID   int                         `gorm:"column:work_schedule_id;primaryKey;autoIncrement" json:"work_schedule_id"`
	WorkScheduleName string                      `gorm:"column:work_schedule_name;type:varchar;index" json:"work_schedule_name"`
	OfficeID         string                      `gorm:"column:office_id;type:varchar" json:"office_id"`
	RepeatType       *variable.RepeatTypeEnum    `gorm:"column:repeat_type;type:repeat_type_enum" json:"repeat_type"`
	RepeatCycle      *int                        `gorm:"column:repeat_cycle;type:integer" json:"repeat_cycle"`
	EffectiveDate    time.Time                   `gorm:"column:effective_date;type:timestamp" json:"effective_date"`
	ExpirationDate   *time.Time                  `gorm:"column:expiration_date;type:timestamp" json:"expiration_date"`
	Status           variable.StatusWorkSchedule `gorm:"column:status;type:status_work_schedule_enum" json:"status"`
	CreatedAt        time.Time                   `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	IsDeleted        bool                        `gorm:"column:is_deleted;default:false" json:"is_deleted"`
	IsAutoRecurring  bool                        `gorm:"column:is_auto_recurring;default:false" json:"is_auto_recurring"`
	IsScheduleAuto   bool                        `gorm:"column:is_schedule_auto;default:false" json:"is_schedule_auto"`

	Managers []WorkScheduleManager `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"managers"`
	Weekdays []WorkScheduleShift   `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"weekdays"`

	Office *model.OfficeResponse `gorm:"foreignKey:OfficeID;references:ID" json:"office"`
}

type WorkScheduleManager struct {
	WorkScheduleManagerID int    `gorm:"column:work_schedule_manager_id;primaryKey;autoIncrement" json:"work_schedule_manager_id"`
	EmployeeID            string `gorm:"column:employee_id;type:varchar;index:idx_emp_sched,unique" json:"employee_id"`
	WorkScheduleID        int    `gorm:"column:work_schedule_id;type:integer;index:idx_emp_sched,unique" json:"work_schedule_id"`
	IsReading             bool   `gorm:"column:is_reading;default:true" json:"is_reading"`
	IsEditing             bool   `gorm:"column:is_editing;default:false" json:"is_editing"`

	Employee *model.ManagerResponse `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"manager"`
}

type WorkScheduleShift struct {
	WorkScheduleShiftID int                  `gorm:"column:work_schedule_shift_id;primaryKey;autoIncrement" json:"work_schedule_shift_id"`
	WorkScheduleID      int                  `gorm:"column:work_schedule_id;type:integer" json:"work_schedule_id"`
	Weekday             variable.WeekdayEnum `gorm:"column:week_day;type:weekday_enum" json:"week_day"`
	WorkShiftID         string               `gorm:"column:workshift_id;type:varchar(20)" json:"workshift_id"`
	Order               int                  `gorm:"column:shift_order;type:integer" json:"order"`

	WorkShift *WorkShifts `gorm:"foreignKey:WorkShiftID;references:WorkShiftID" json:"work_shift"`
}

func (WorkSchedule) TableName() string {
	return "work_schedule"
}

func (WorkScheduleShift) TableName() string {
	return "work_schedule_shift"
}

func (WorkScheduleManager) TableName() string {
	return "work_schedule_manager"
}
