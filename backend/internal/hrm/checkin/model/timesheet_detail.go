package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg/variable"
	"time"
)

type TimeSheetDetail struct {
	TimeSheetDetailID int       `gorm:"column:timesheet_detail_id;primaryKey;autoIncrement" json:"timesheet_detail_id"`
	TimeSheetID       int       `gorm:"column:timesheet_id;not null;uniqueIndex:idx_timesheet_date" json:"timesheet_id"`
	Date              time.Time `gorm:"column:date;type:date;not null;uniqueIndex:idx_timesheet_date" json:"date"`
	DayOfWeek         int       `gorm:"column:day_of_week;type:integer" json:"day_of_week"` // 1=monday, 7=sunday

	//  Work shift information
	WorkShiftID  *string `gorm:"column:work_shift_id;type:varchar(20)" json:"work_shift_id"`
	IsWorkingDay bool    `gorm:"column:is_working_day;type:boolean;default:true" json:"is_working_day"`
	IsHoliday    bool    `gorm:"column:is_holiday;type:boolean;default:false" json:"is_holiday"`
	IsWeekend    bool    `gorm:"column:is_weekend;type:boolean;default:false" json:"is_weekend"`

	// Attendance information
	CheckInRecordID  *string `gorm:"column:checkin_record_id;type:uuid" json:"checkin_record_id"`
	CheckOutRecordID *string `gorm:"column:checkout_record_id;type:uuid" json:"checkout_record_id"`

	// Workday calculation
	WorkHours float64 `gorm:"column:work_hours;type:decimal(4,2);default:0" json:"work_hours"`
	WorkDays  float64 `gorm:"column:work_days;type:decimal(3,2);default:0" json:"work_days"`
	//IsAdditionalShift bool    `gorm:"column:is_additional_shift;type:boolean;default:false" json:"is_additional_shift"`

	// Late information
	IsLate      bool `gorm:"column:is_late;type:boolean;default:false" json:"is_late"`
	LateMinutes int  `gorm:"column:late_minutes;type:integer;default:0" json:"late_minutes"`

	// Leave and status info – Not yet processed
	LeaveType    *variable.LeaveTypeEnum `gorm:"column:leave_type;type:varchar(50)" json:"leave_type"`
	LeaveHours   float64                 `gorm:"column:leave_hours;type:decimal(4,2);default:0" json:"leave_hours"`
	IsAbsent     bool                    `gorm:"column:is_absent;type:boolean;default:false" json:"is_absent"`
	AbsentReason *string                 `gorm:"column:absent_reason;type:text" json:"absent_reason"`

	//Notes and manual adjustments
	WorkDaysAdjusted   float64    `gorm:"column:work_day_adjusted;type:decimal(3,2);default:0" json:"work_day_adjusted"`
	OriginalWorkDays   *float64   `gorm:"column:original_work_days;type:decimal(3,2)" json:"original_work_days"`
	IsManuallyAdjusted bool       `gorm:"column:is_manually_adjusted;type:boolean;default:false" json:"is_manually_adjusted"`
	AdjustmentBy       *string    `gorm:"column:adjustment_by;type:varchar" json:"adjustment_by"`
	AdjustmentAt       *time.Time `gorm:"column:adjustment_at" json:"adjustment_at"`

	// Relationships
	WorkShift      *WorkShifts            `gorm:"foreignKey:WorkShiftID;references:WorkShiftID" json:"work_shift,omitempty"`
	CheckInRecord  *AttendanceRecord      `gorm:"foreignKey:CheckInRecordID;references:AttendanceRecordID" json:"checkin_record,omitempty"`
	CheckOutRecord *AttendanceRecord      `gorm:"foreignKey:CheckOutRecordID;references:AttendanceRecordID" json:"checkout_record,omitempty"`
	AdjustmentUser *model.ManagerResponse `gorm:"foreignKey:AdjustmentBy;references:EmployeeID" json:"adjustment_user,omitempty"`
}

func (TimeSheetDetail) TableName() string {
	return "timesheet_details"
}
