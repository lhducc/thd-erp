package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"time"
)

type TimeSheet struct {
	TimeSheetID     int    `gorm:"column:timesheet_id;primaryKey;autoIncrement" json:"timesheet_id"`
	TimeSheetListID string `gorm:"column:timesheet_list_id;type:varchar;not null;index" json:"timesheet_list_id"`
	EmployeeID      string `gorm:"column:employee_id;type:varchar;not null;index:idx_list_emp" json:"employee_id"` // employee of Office
	Month           int    `gorm:"column:month;type:integer;not null;index:idx_list_emp" json:"month" validate:"required,min=1,max=12"`
	Year            int    `gorm:"column:year;type:integer;not null;index:idx_list_emp" json:"year" validate:"required,min=2020"`
	OfficeID        string `gorm:"column:office_id;type:varchar" json:"office_id"`         // populated from Employee.Department.OfficeID when creating timesheet
	DepartmentID    string `gorm:"column:department_id;type:varchar" json:"department_id"` //get from Employee.DepartmentID

	// Working day summary
	TotalWorkDays float64 `gorm:"column:total_work_days;type:decimal(5,2);default:0" json:"total_work_days"`
	//AdditionalShiftDays float64 `gorm:"column:additional_shift_days;type:decimal(5,2);default:0" json:"additional_shift_days"`
	//TotalDays           float64 `gorm:"column:total_days;type:decimal(5,2);default:0" json:"total_days"`

	// Late shift statistics
	LateShifts       int `gorm:"column:late_shifts;type:integer;default:0" json:"late_shifts"`
	TotalLateMinutes int `gorm:"column:total_late_minutes;type:integer;default:0" json:"total_late_minutes"`

	// Types of leave and remote work (by day) – Not yet processed
	AnnualLeaveDays   float64 `gorm:"column:annual_leave_days;type:decimal(5,2);default:0" json:"annual_leave_days"`
	PersonalLeaveDays float64 `gorm:"column:personal_leave_days;type:decimal(5,2);default:0" json:"personal_leave_days"`
	BusinessTripDays  float64 `gorm:"column:business_trip_days;type:decimal(5,2);default:0" json:"business_trip_days"`
	RemoteWorkDays    float64 `gorm:"column:remote_work_days;type:decimal(5,2);default:0" json:"remote_work_days"`
	UnpaidLeaveDays   float64 `gorm:"column:unpaid_leave_days;type:decimal(5,2);default:0" json:"unpaid_leave_days"`
	OtherLeaveDays    float64 `gorm:"column:other_leave_days;type:decimal(5,2);default:0" json:"other_leave_days"`

	// Types of leave and remote work (by hour) – Not yet processed
	AnnualLeaveHours float64 `gorm:"column:annual_leave_hours;type:decimal(5,2);default:0" json:"annual_leave_hours"`
	RemoteWorkHours  float64 `gorm:"column:remote_work_hours;type:decimal(5,2);default:0" json:"remote_work_hours"`
	UnpaidLeaveHours float64 `gorm:"column:unpaid_leave_hours;type:decimal(5,2);default:0" json:"unpaid_leave_hours"`

	// Status and creation information – Not yet processed
	//Status variable.TimeSheetStatusEnum `gorm:"column:status;type:varchar(20)" json:"status"`

	//Creation information
	CreatedBy string     `gorm:"column:created_by;type:varchar" json:"created_by"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedBy *string    `gorm:"column:updated_by;type:varchar" json:"updated_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	//ApprovedBy *string    `gorm:"column:approved_by;type:varchar" json:"approved_by"`
	//ApprovedAt *time.Time `gorm:"column:approved_at" json:"approved_at"`
	//IsDeleted bool `gorm:"column:is_deleted;type:boolean;default:false" json:"is_deleted"`

	// Relationships
	Employee   *model.EmployeeInforResponse `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
	Office     *model.Office                `gorm:"foreignKey:OfficeID;references:ID" json:"office,omitempty"`
	Department *model.Department            `gorm:"foreignKey:DepartmentID;references:ID" json:"department,omitempty"`
	Details    []TimeSheetDetail            `gorm:"foreignKey:TimeSheetID;references:TimeSheetID;constraint:OnDelete:CASCADE" json:"details,omitempty"`
	Creator    *model.ManagerResponse       `gorm:"foreignKey:CreatedBy;references:EmployeeID" json:"creator,omitempty"`
	Updater    *model.ManagerResponse       `gorm:"foreignKey:UpdatedBy;references:EmployeeID" json:"updater,omitempty"`
	//Approver   *model.ManagerResponse `gorm:"foreignKey:ApprovedBy;references:EmployeeID" json:"approver,omitempty"`
}

func (TimeSheet) TableName() string {
	return "timesheets"
}
