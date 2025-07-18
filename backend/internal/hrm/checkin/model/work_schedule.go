package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg/struct_support"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"time"
)

type WorkSchedule struct {
	WorkScheduleID   int                         `gorm:"column:work_schedule_id;primaryKey;autoIncrement" json:"work_schedule_id"`
	WorkScheduleName string                      `gorm:"column:work_schedule_name;type:varchar;index" json:"work_schedule_name"`
	OfficeID         string                      `gorm:"column:office_id;type:varchar" json:"office_id"`
	RepeatType       variable.RepeatTypeEnum     `gorm:"column:repeat_type;type:repeat_type_enum" json:"repeat_type"`
	RepeatCycle      int                         `gorm:"column:repeat_cycle;type:integer" json:"repeat_cycle"`
	EffectiveDate    time.Time                   `gorm:"column:effective_date;type:timestamp" json:"effective_date"`
	ExpirationDate   time.Time                   `gorm:"column:expiration_date;type:timestamp" json:"expiration_date"`
	Status           variable.StatusWorkSchedule `gorm:"column:status;type:status_work_schedule_enum" json:"status"`
	CreatedAt        time.Time                   `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	IsDeleted        bool                        `gorm:"column:is_deleted;default:false" json:"is_deleted"`

	Managers  []WorkScheduleManager  `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"managers"`
	Employees []WorkScheduleEmployee `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"employees"`
	Weekdays  []WorkScheduleShift    `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"weekdays"`

	Office *model.OfficeResponse `gorm:"foreignKey:OfficeID;references:ID" json:"office"`
}

type WorkScheduleRequest struct {
	WorkScheduleName string                  ` json:"work_schedule_name"`
	OfficeID         string                  `json:"office_id"`
	RepeatType       variable.RepeatTypeEnum `json:"repeat_type"`
	RepeatCycle      int                     `json:"repeat_cycle"`
	EffectiveDate    time.Time               `json:"effective_date"`
	ExpirationDate   time.Time               `json:"expiration_date"`
	WeekDays         []WorkScheduleShift     `json:"weekdays"`
}

type WorkScheduleEmployee struct {
	WorkScheduleEmployeeID int        `gorm:"column:work_schedule_employee_id;primaryKey;autoIncrement" json:"work_schedule_employee_id"`
	EmployeeID             string     `gorm:"column:employee_id;type:varchar;unique;index" json:"employee_id"`
	WorkScheduleID         *int       `gorm:"column:work_schedule_id;type:integer;index" json:"work_schedule_id"`
	AssignedAt             *time.Time `gorm:"column:assigned_at;default:CURRENT_TIMESTAMP" json:"assigned_at"`

	Employee *model.ManagerResponse `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee"`
}

type WorkScheduleManager struct {
	WorkScheduleManagerID int    `gorm:"column:work_schedule_manager_id;primaryKey;autoIncrement" json:"work_schedule_manager_id"`
	EmployeeID            string `gorm:"column:employee_id;type:varchar;index:idx_emp_sched,unique" json:"employee_id"`
	WorkScheduleID        int    `gorm:"column:work_schedule_id;type:integer;index:idx_emp_sched,unique" json:"work_schedule_id"`
	IsReading             bool   `gorm:"column:is_reading;default:true" json:"is_reading"`
	IsEditing             bool   `gorm:"column:is_editing;default:false" json:"is_editing"`

	Employee *model.ManagerResponse `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"manager"`
}
type AssignEmployeeRequest struct {
	EmployeeIDs []string               `json:"employee_ids"`
	Managers    []ManagerAssignRequest `json:"managers"`
}
type ManagerAssignRequest struct {
	ManagerID string `json:"manager_id"`
	IsReading bool   `json:"is_reading"`
	IsEditing bool   `json:"is_editing"`
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

func (WorkScheduleEmployee) TableName() string {
	return "work_schedule_employee"
}

func (WorkScheduleManager) TableName() string {
	return "work_schedule_manager"
}

func ConvertToWorkSchedule(req *WorkScheduleRequest) WorkSchedule {
	return WorkSchedule{
		WorkScheduleName: req.WorkScheduleName,
		OfficeID:         req.OfficeID,
		RepeatType:       req.RepeatType,
		RepeatCycle:      req.RepeatCycle,
		EffectiveDate:    req.EffectiveDate,
		ExpirationDate:   req.ExpirationDate,
		Weekdays:         req.WeekDays,
	}
}

func (req *WorkSchedule) Validate() error {
	if req.WorkScheduleName == "" {
		return errors.New("work_schedule_name is required")
	}

	if req.OfficeID == "" {
		return errors.New("office_id is required")
	}

	if _, ok := struct_support.ValidRepeatTypes[req.RepeatType]; !ok {
		return fmt.Errorf("invalid repeat_type: %s", req.RepeatType)
	}

	for _, weekday := range req.Weekdays {
		if _, ok := struct_support.ValidWeekdays[weekday.Weekday]; !ok {
			return fmt.Errorf("invalid week_day: %s", weekday.Weekday)
		}
	}

	if req.EffectiveDate.IsZero() || req.ExpirationDate.IsZero() {
		return errors.New("effective_date and expiration_date must not be empty")
	}

	if req.EffectiveDate.After(req.ExpirationDate) {
		return errors.New("effective_date must be before expiration_date")
	}

	return nil
}
