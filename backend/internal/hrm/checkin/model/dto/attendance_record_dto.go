package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"strings"
	"time"
)

type AttendanceRecordCreate struct {
	Timestamp            time.Time `form:"timestamp" gorm:"column:timestamp;type:timestamp;not null"`
	GPSLocation          string    `form:"gps_location" gorm:"column:location;type:text"` // GPS coordinates
	OfficeID             *string   `form:"office_id" gorm:"column:office_id;type:varchar"`
	AttendanceCategoryID string    `form:"category_id" gorm:"column:category_id;type:uuid;not null"` // Onsite, WFH, AtOffice,..
	NoteRequest          *string   `form:"note_request" gorm:"column:note_request;type:text"`
}

type AttendanceRecordCreateByAdmin struct {
	EmployeeID string    `json:"employee_id"`
	Timestamp  time.Time `json:"timestamp"`
}

type AttendanceRecordUpdate struct {
	Status     variable.StatusEnum `gorm:"column:status;default:Pending" json:"status"`
	NoteReject *string             `gorm:"column:note_reject;type:text" json:"note_reject"`
}

type AttendanceRecordHistoryByDate struct {
	EmployeeID     string     `gorm:"employee_id" json:"employee_id"`
	FullName       string     `gorm:"full_name" json:"full_name"`
	OfficeName     string     `gorm:"office_name" json:"office_name"`
	DepartmentName string     `gorm:"department_name" json:"department_name"`
	Timestamp      *time.Time `gorm:"timestamp" json:"timestamp"`
	WorkShiftID    string     `gorm:"workshift_id" json:"workshift_id"`
	WorkShiftName  string     `gorm:"workshift_name" json:"workshift_name"`
	StartTime      string     `gorm:"start_time" json:"start_time"`
	CheckinTo      string     `gorm:"checkin_to" json:"checkin_to"`
}

func (r *AttendanceRecordCreate) Validate() error {
	if r.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}

	if r.AttendanceCategoryID == "" {
		return errors.New("category ID is required")
	}

	return nil
}

func (s *AttendanceRecordUpdate) Validate() error {
	switch s.Status {
	case variable.Pending, variable.Rejected, variable.Approved:
		return nil
	default:
		return fmt.Errorf("invalid status: %s", s.Status)
	}
}

func ConvertToAttendanceRecordStruct(record *AttendanceRecordCreate) model.AttendanceRecord {
	return model.AttendanceRecord{
		Timestamp:            record.Timestamp,
		GPSLocation:          record.GPSLocation,
		OfficeID:             record.OfficeID,
		AttendanceCategoryID: &record.AttendanceCategoryID,
		NoteRequest:          record.NoteRequest,
	}
}

func (r *AttendanceRecordCreateByAdmin) TableName() string {
	return "attendance_records"
}

func (r *AttendanceRecordCreateByAdmin) Validate() error {
	if strings.TrimSpace(r.EmployeeID) == "" {
		return errors.New("Mã nhân viên là bắt buộc")
	}
	if r.Timestamp.IsZero() {
		return errors.New("Thời gian là bắt buộc")
	}
	return nil
}
