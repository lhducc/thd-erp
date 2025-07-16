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
	Employees []WorkScheduleEmployee `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"employees,omitempty"`
	Weekdays  []WorkScheduleShift    `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"weekdays,omitempty"`

	Office *model.OfficeResponse `gorm:"foreignKey:OfficeID;references:office_id" json:"office"`
}

type WorkScheduleRequest struct {
	WorkScheduleName string                      ` json:"work_schedule_name"`
	OfficeID         string                      `json:"office_id"`
	RepeatType       variable.RepeatTypeEnum     `json:"repeat_type"`
	RepeatCycle      int                         `json:"repeat_cycle"`
	EffectiveDate    time.Time                   `json:"effective_date"`
	ExpirationDate   time.Time                   `json:"expiration_date"`
	Status           variable.StatusWorkSchedule `json:"status"`
	WeekDays         []WorkScheduleShift         `json:"weekdays"`
}

type WorkScheduleEmployee struct {
	EmployeeID     string     `gorm:"column:employee_id;primaryKey;type:varchar" json:"employee_id"`
	WorkScheduleID *int       `gorm:"column:work_schedule_id;primaryKey;type:integer" json:"work_schedule_id"`
	AssignedAt     *time.Time `gorm:"column:assigned_at;default:CURRENT_TIMESTAMP" json:"assigned_at"`
}

type WorkScheduleManager struct {
	EmployeeID     string     `gorm:"column:employee_id;primaryKey;type:varchar" json:"employee_id"`
	WorkScheduleID int        `gorm:"column:work_schedule_id;type:integer" json:"work_schedule_id"`
	AssignedAt     *time.Time `gorm:"column:assigned_at;default:CURRENT_TIMESTAMP" json:"assigned_at"`
}
type AssignEmployeeRequest struct {
	EmployeeIDs []string `json:"employee_ids"`
	ManagerIDs  []string `json:"manager_ids"`

	WorkSchedule WorkSchedule `gorm:"foreignKey:WorkScheduleID;references:WorkScheduleID" json:"work_schedule,omitempty"`
}

type WorkScheduleShift struct {
	WorkScheduleID int                  `gorm:"column:work_schedule_id;type:integer" json:"work_schedule_id"`
	Weekday        variable.WeekdayEnum `gorm:"column:week_day;type:weekday_enum" json:"week_day"`
	WorkShiftID    string               `gorm:"column:workshift_id;type:varchar(20)" json:"workshift_id"`
	Order          int                  `gorm:"column:shift_order;type:integer" json:"order"`

	WorkShift WorkShifts `gorm:"foreignKey:WorkShiftID;references:WorkShiftID" json:"work_shift"`
}

func (WorkScheduleShift) TableName() string {
	return "work_schedule_shift"
}

func (WorkScheduleEmployee) TableName() string {
	return "work_schedule_employee"
}
func (WorkSchedule) TableName() string {
	return "work_schedule"
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
		Status:           req.Status,
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

	if _, ok := struct_support.ValidStatus[req.Status]; !ok {
		return fmt.Errorf("invalid status: %s", req.Status)
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
