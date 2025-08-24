package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg/variable"
	"time"
)

type WorkShifts struct {
	WorkShiftID   string               `gorm:"column:workshift_id;type:varchar(20);primaryKey;" json:"workshift_id"`
	WorkShiftName string               `gorm:"column:workshift_name;type:varchar(255);not null;unique" json:"workshift_name" validate:"required,max=255"`
	StartTime     string               `gorm:"column:start_time;type:timestamp;not null;index:idx_time_range" json:"start_time" validate:"required"`
	EndTime       string               `gorm:"column:end_time;type:timestamp;not null;index:idx_time_range" json:"end_time" validate:"required"`
	CheckinFrom   *string              `gorm:"column:checkin_from;type:timestamp" json:"checkin_from"`
	CheckinTo     *string              `gorm:"column:checkin_to;type:timestamp" json:"checkin_to"`
	CheckoutFrom  *string              `gorm:"column:checkout_from;type:timestamp" json:"checkout_from"`
	CheckoutTo    *string              `gorm:"column:checkout_to;type:timestamp" json:"checkout_to"`
	HasBreak      bool                 `gorm:"column:has_break;type:boolean" json:"has_break"`
	BreakStart    *string              `gorm:"column:break_start;type:timestamp" json:"break_start"`
	BreakEnd      *string              `gorm:"column:break_end;type:timestamp" json:"break_end"`
	WorkHours     float64              `gorm:"column:work_hours;type:decimal(4,2)" json:"work_hours" validate:"gte=0"`
	WorkDay       variable.WorkDayEnum `gorm:"column:work_day;type:int;default:0" json:"work_day" validate:"required"`
	CoefNormalDay float64              `gorm:"column:coef_normal_day;type:decimal(3,2);default:1.00" json:"coef_normal_day" validate:"gte=0"`
	CoefWeekend   float64              `gorm:"column:coef_weekend;type:decimal(3,2);default:1.00" json:"coef_weekend" validate:"gte=0"`
	CoefHoliday   float64              `gorm:"column:coef_holiday;type:decimal(3,2);default:1.00" json:"coef_holiday" validate:"gte=0"`
	CreatedBy     string               `gorm:"column:created_by;not null" json:"created_by"`
	CreatedDate   time.Time            `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	IsDeleted     bool                 `gorm:"column:is_deleted;type:boolean;default:false" json:"is_deleted"`

	EffectiveDate  time.Time  `gorm:"column:effective_date;type:timestamp;not null" json:"effective_date" validate:"required"`
	ExpirationDate *time.Time `gorm:"column:expiration_date;type:timestamp" json:"expiration_date"`

	Creator *model.Employee `gorm:"foreignKey:CreatedBy;references:employee_id" json:"creator,omitempty"`
}

func (WorkShifts) TableName() string {
	return "workshifts"
}
