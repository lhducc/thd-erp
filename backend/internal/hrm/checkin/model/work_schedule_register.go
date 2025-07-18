package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg/struct_support"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"time"
)

type WorkScheduleRegister struct {
	WorkScheduleRegisterID   int                         `gorm:"column:work_schedule_register_id;primaryKey;autoIncrement" json:"work_schedule_register_id"`
	WorkScheduleRegisterName string                      `gorm:"column:work_schedule_register_name;type:varchar;index" json:"work_schedule_register_name"`
	OfficeID                 string                      `gorm:"column:office_id;type:varchar" json:"office_id"`
	EffectiveDate            time.Time                   `gorm:"column:effective_date;type:timestamp" json:"effective_date"`
	ExpirationDate           time.Time                   `gorm:"column:expiration_date;type:timestamp" json:"expiration_date"`
	Status                   variable.StatusWorkSchedule `gorm:"column:status;type:status_work_schedule_enum" json:"status"`
	CreatedAt                time.Time                   `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	IsDeleted                bool                        `gorm:"column:is_deleted;default:false" json:"is_deleted"`

	Managers  []WorkScheduleRegisterManager  `gorm:"foreignKey:WorkScheduleRegisterID;references:WorkScheduleRegisterID" json:"managers"`
	Employees []WorkScheduleRegisterEmployee `gorm:"foreignKey:WorkScheduleRegisterID;references:WorkScheduleRegisterID" json:"employees"`
	Weekdays  []WorkScheduleRegisterShift    `gorm:"foreignKey:WorkScheduleRegisterID;references:WorkScheduleRegisterID" json:"weekdays"`

	Office *model.OfficeResponse `gorm:"foreignKey:OfficeID;references:ID" json:"office"`
}

type WorkScheduleRegisterRequest struct {
	WorkScheduleRegisterName string                      `json:"work_schedule_register_name"`
	OfficeID                 string                      `json:"office_id"`
	EffectiveDate            time.Time                   `json:"effective_date"`
	ExpirationDate           time.Time                   `json:"expiration_date"`
	WeekDays                 []WorkScheduleRegisterShift `json:"weekdays"`
}

type WorkScheduleRegisterEmployee struct {
	WorkScheduleRegisterEmployeeID int        `gorm:"column:work_schedule_register_employee_id;primaryKey;autoIncrement" json:"work_schedule_register_employee_id"`
	EmployeeID                     string     `gorm:"column:employee_id;type:varchar;not null;unique" json:"employee_id"` // unique là đúng trong case của bạn
	WorkScheduleRegisterID         *int       `gorm:"column:work_schedule_register_id;type:integer;not null" json:"work_schedule_register_id"`
	AssignedAt                     *time.Time `gorm:"column:assigned_at;default:CURRENT_TIMESTAMP" json:"assigned_at"`

	Employee *model.ManagerResponse `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee"`
}

type WorkScheduleRegisterManager struct {
	WorkScheduleRegisterManagerID int    `gorm:"column:work_schedule_register_manager_id;primaryKey;autoIncrement" json:"work_schedule_register_manager_id"`
	EmployeeID                    string `gorm:"column:employee_id;type:varchar;index:idx_emp_sched_reg,unique" json:"employee_id"`
	WorkScheduleRegisterID        int    `gorm:"column:work_schedule_register_id;type:integer;index:idx_emp_sched_reg,unique" json:"work_schedule_register_id"`
	IsReading                     bool   `gorm:"column:is_reading;default:true" json:"is_reading"`
	IsEditing                     bool   `gorm:"column:is_editing;default:false" json:"is_editing"`

	Employee *model.ManagerResponse `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"manager"`
}

type WorkScheduleRegisterShift struct {
	WorkScheduleRegisterShiftID int                  `gorm:"column:work_schedule_register_shift_id;primaryKey;autoIncrement" json:"work_schedule_register_shift_id"`
	WorkScheduleRegisterID      int                  `gorm:"column:work_schedule_register_id;type:integer" json:"work_schedule_register_id"`
	Weekday                     variable.WeekdayEnum `gorm:"column:week_day;type:weekday_enum" json:"week_day"`
	WorkShiftID                 string               `gorm:"column:workshift_id;type:varchar(20)" json:"workshift_id"`
	Order                       int                  `gorm:"column:shift_order;type:integer" json:"order"`

	WorkShift *WorkShifts `gorm:"foreignKey:WorkShiftID;references:WorkShiftID" json:"work_shift"`
}

func (WorkScheduleRegister) TableName() string {
	return "work_schedule_register"
}

func (WorkScheduleRegisterShift) TableName() string {
	return "work_schedule_register_shift"
}

func (WorkScheduleRegisterEmployee) TableName() string {
	return "work_schedule_register_employee"
}

func (WorkScheduleRegisterManager) TableName() string {
	return "work_schedule_register_manager"
}

func ConvertToWorkScheduleRegister(req *WorkScheduleRegisterRequest) WorkScheduleRegister {
	return WorkScheduleRegister{
		WorkScheduleRegisterName: req.WorkScheduleRegisterName,
		OfficeID:                 req.OfficeID,
		EffectiveDate:            req.EffectiveDate,
		ExpirationDate:           req.ExpirationDate,
		Weekdays:                 req.WeekDays,
	}
}

func (req *WorkScheduleRegister) Validate() error {
	if req.WorkScheduleRegisterName == "" {
		return errors.New("work_schedule_register_name is required")
	}

	if req.OfficeID == "" {
		return errors.New("office_id is required")
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
