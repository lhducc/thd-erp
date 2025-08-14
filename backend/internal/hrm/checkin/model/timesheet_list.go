package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"time"
)

type TimeSheetList struct {
	TimeSheetListID   string `gorm:"column:timesheet_list_id;primaryKey;type:varchar" json:"timesheet_list_id"`
	TimeSheetListName string `gorm:"column:timesheet_list_name;type:varchar(255);not null" json:"time_sheet_list_name"`
	OfficeID          string `gorm:"column:office_id;type:varchar;not null;index:idx_office_month_year" json:"office_id"`
	Month             int    `gorm:"column:month;type:integer;not null;index:idx_office_month_year" json:"month" validate:"required,min=1,max=12"`
	Year              int    `gorm:"column:year;type:integer;not null;index:idx_office_month_year" json:"year" validate:"required,min=2020"`

	// General statistics
	//TotalEmployees int `gorm:"column:total_employees;type:integer;default:0" json:"total_employees"`
	//CompletedTimesheets int `gorm:"column:completed_timesheets;type:integer;default:0" json:"completed_timesheets"`
	//PendingTimesheets   int `gorm:"column:pending_timesheets;type:integer;default:0" json:"pending_timesheets"`
	//ApprovedTimesheets  int `gorm:"column:approved_timesheets;type:integer;default:0" json:"approved_timesheets"`

	// dates
	StartDate time.Time `gorm:"column:start_date;type:date;not null" json:"start_date"` // start month
	EndDate   time.Time `gorm:"column:end_date;type:date;not null" json:"end_date"`     // end month
	//DeadlineSubmit  *time.Time `gorm:"column:deadline_submit;type:timestamp" json:"deadline_submit"`
	//DeadlineApprove *time.Time `gorm:"column:deadline_approve;type:timestamp" json:"deadline_approve"`

	//Not yet processed
	//Status variable.TimeSheetListStatusEnum `gorm:"column:status;type:varchar(20);default:'open'" json:"status"`

	// Status and permissions
	IsLocked bool `gorm:"column:is_locked;type:boolean;default:false" json:"is_locked"`

	// Creation and update info
	CreatedBy string     `gorm:"column:created_by;type:varchar;not null" json:"created_by"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedBy *string    `gorm:"column:updated_by;type:varchar" json:"updated_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	LockedBy  *string    `gorm:"column:locked_by;type:varchar" json:"locked_by"`
	LockedAt  *time.Time `gorm:"column:locked_at" json:"locked_at"`
	//IsDeleted bool       `gorm:"column:is_deleted;type:boolean;default:false" json:"is_deleted"`

	// Relationships
	Office     *model.Office   `gorm:"foreignKey:OfficeID;references:ID" json:"office,omitempty"`
	Timesheets []TimeSheet     `gorm:"foreignKey:TimeSheetListID;references:TimeSheetListID;constraint:OnDelete:CASCADE" json:"timesheets,omitempty"`
	Creator    *model.Employee `gorm:"foreignKey:CreatedBy;references:EmployeeID" json:"creator,omitempty"`
	Updater    *model.Employee `gorm:"foreignKey:UpdatedBy;references:EmployeeID" json:"updater,omitempty"`
	LockedUser *model.Employee `gorm:"foreignKey:LockedBy;references:EmployeeID" json:"locked_user,omitempty"`
}

func (TimeSheetList) TableName() string {
	return "timesheet_lists"
}
