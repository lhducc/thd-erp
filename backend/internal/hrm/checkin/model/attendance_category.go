package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"time"
)

type StatusAttendanceEnum string

const (
	StatusActive   StatusAttendanceEnum = "active"
	StatusInactive StatusAttendanceEnum = "inactive"
)

// Validate enum value
func (s StatusAttendanceEnum) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive:
		return true
	default:
		return false
	}
}

type AttendanceCategory struct {
	AttendanceCategoryID   string               `gorm:"column:attendance_category_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"attendance_category_id"`
	AttendanceCategoryName string               `gorm:"column:attendance_category_name;type:varchar(255);not null;unique" json:"attendance_category_name"`
	IsGPS                  bool                 `gorm:"column:is_gps;not null;default:false" json:"is_gps"`
	IsCamera               bool                 `gorm:"column:is_camera;not null;default:false" json:"is_camera"`
	IsCheckLocation        bool                 `gorm:"column:is_check_location;not null;default:false" json:"is_check_location"`
	Scope                  int                  `gorm:"column:scope;type:integer" json:"scope"` // Changed to integer for meters
	Status                 StatusAttendanceEnum `gorm:"column:status;type:varchar(50);not null;default:'active" json:"status"`
	IsDeleted              bool                 `gorm:"column:is_deleted;default:false" json:"is_deleted"`
	OfficeID               string               `gorm:"column:office_id;type:varchar(50);not null" json:"office_id"`
	CreatedBy              string               `gorm:"column:created_by;type:varchar(50)" json:"created_by"`
	CreatedAt              time.Time            `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	AutoApprove            bool                 `gorm:"column:auto_approve;default:false" json:"auto_approve"`

	Office   *model.Office   `gorm:"foreignKey:OfficeID;references:ID"`
	Employee *model.Employee `gorm:"foreignKey:CreatedBy;references:EmployeeID"`
}

func (AttendanceCategory) TableName() string {
	return "attendance_categories"
}

// Validate struct AttendanceCategory
func (a *AttendanceCategory) Validate() error {
	if a.AttendanceCategoryName == "" {
		return errors.New("attendance_category_name is required")
	}

	if !a.Status.IsValid() {
		return fmt.Errorf("invalid status: %s", a.Status)
	}

	if a.OfficeID == "" {
		return errors.New("office_id is required")
	}

	if a.IsCheckLocation && a.Scope <= 0 {
		return errors.New("scope must be positive number when is_check_location is true")
	}
	if a.IsCheckLocation && a.IsGPS == false {
		return errors.New("is_gps must be true when is_check_location is true")
	}

	return nil
}

type AttendanceCategoryRequest struct {
	AttendanceCategoryName string               `json:"attendance_category_name" validate:"required"`
	IsGPS                  bool                 `json:"is_gps"`
	IsCamera               bool                 `json:"is_camera"`
	IsCheckLocation        bool                 `json:"is_check_location"`
	Scope                  int                  `json:"scope"` // Changed to int
	Status                 StatusAttendanceEnum `json:"status" validate:"required"`
	OfficeID               string               `json:"office_id" validate:"required"`
}

func ConvertToAttendanceCategory(req AttendanceCategoryRequest) AttendanceCategory {
	return AttendanceCategory{
		AttendanceCategoryName: req.AttendanceCategoryName,
		IsGPS:                  req.IsGPS,
		IsCamera:               req.IsCamera,
		IsCheckLocation:        req.IsCheckLocation,
		Scope:                  req.Scope,
		Status:                 req.Status,
		OfficeID:               req.OfficeID,
	}
}
