package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
)

// Validate enum value
func IsValid(s variable.StatusAttendanceEnum) bool {
	switch s {
	case variable.StatusActive, variable.StatusInactive:
		return true
	default:
		return false
	}
}

func (a *AttendanceCategoryRequest) Validate() error {
	if a.AttendanceCategoryName == "" {
		return errors.New("attendance_category_name is required")
	}

	if !IsValid(a.Status) {
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
	AttendanceCategoryName string                        `json:"attendance_category_name" validate:"required"`
	IsGPS                  bool                          `json:"is_gps"`
	IsCamera               bool                          `json:"is_camera"`
	IsCheckLocation        bool                          `json:"is_check_location"`
	Scope                  int                           `json:"scope"` // Changed to int
	Status                 variable.StatusAttendanceEnum `json:"status" validate:"required"`
	OfficeID               string                        `json:"office_id" validate:"required"`
}

func ConvertToAttendanceCategory(req AttendanceCategoryRequest) model.AttendanceCategory {
	return model.AttendanceCategory{
		AttendanceCategoryName: req.AttendanceCategoryName,
		IsGPS:                  req.IsGPS,
		IsCamera:               req.IsCamera,
		IsCheckLocation:        req.IsCheckLocation,
		Scope:                  req.Scope,
		Status:                 req.Status,
		OfficeID:               req.OfficeID,
	}
}
