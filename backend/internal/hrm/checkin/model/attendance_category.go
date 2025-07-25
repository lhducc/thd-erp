package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg/variable"
	"time"
)

type AttendanceCategory struct {
	AttendanceCategoryID   string                        `gorm:"column:attendance_category_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"attendance_category_id"`
	AttendanceCategoryName string                        `gorm:"column:attendance_category_name;type:varchar(255);not null;unique" json:"attendance_category_name"`
	IsGPS                  bool                          `gorm:"column:is_gps;not null;default:false" json:"is_gps"`
	IsCamera               bool                          `gorm:"column:is_camera;not null;default:false" json:"is_camera"`
	IsCheckLocation        bool                          `gorm:"column:is_check_location;not null;default:false" json:"is_check_location"`
	Scope                  int                           `gorm:"column:scope;type:integer" json:"scope"` // Changed to integer for meters
	Status                 variable.StatusAttendanceEnum `gorm:"column:status;type:varchar(50);not null;default:'active" json:"status"`
	IsDeleted              bool                          `gorm:"column:is_deleted;default:false" json:"is_deleted"`
	OfficeID               string                        `gorm:"column:office_id;type:varchar(50);not null" json:"office_id"`
	CreatedBy              string                        `gorm:"column:created_by;type:varchar(50)" json:"created_by"`
	CreatedAt              time.Time                     `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	AutoApprove            bool                          `gorm:"column:auto_approve;default:false" json:"auto_approve"`

	Office *model.Office `gorm:"foreignKey:OfficeID;references:ID"`
	//Employee *model.Employee `gorm:"foreignKey:CreatedBy;references:EmployeeID"`
}

func (AttendanceCategory) TableName() string {
	return "attendance_categories"
}
