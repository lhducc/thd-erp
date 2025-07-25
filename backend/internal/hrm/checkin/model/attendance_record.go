package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg/variable"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AttendanceRecord struct {
	AttendanceRecordID   string              `gorm:"column:attendance_record_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"attendance_record_id"`
	EmployeeID           string              `gorm:"column:employee_id;type:varchar;not null;index" json:"employee_id"`
	Timestamp            time.Time           `gorm:"column:timestamp;type:timestamp;not null" json:"timestamp"`
	GPSLocation          string              `gorm:"-" json:"gps_location"` // GPS coordinates
	Latitude             *float64            `gorm:"column:latitude;" json:"latitude"`
	Longitude            *float64            `gorm:"column:longitude;" json:"longitude"`
	ImageName            string              `gorm:"column:image_name;type:text" json:"image_name"`
	ImageURL             string              `gorm:"-" json:"image_URL"`
	OfficeID             *string             `gorm:"column:office_id;type:varchar;default:null" json:"office_id"`
	AttendanceCategoryID *string             `gorm:"column:category_id;type:uuid" json:"category_id"` // Onsite, WFH, AtOffice,..
	NoteRequest          *string             `gorm:"column:note_request;type:text" json:"note_request"`
	NoteReject           *string             `gorm:"column:note_reject;type:text" json:"note_reject"`
	Status               variable.StatusEnum `gorm:"column:status;default:pending" json:"status"`
	CreatedBy            *string             `gorm:"column:created_by;type:varchar;default:null" json:"created_by"`

	Employee           *model.EmployeeInforResponse `gorm:"foreignKey:EmployeeID;references:EmployeeID"`
	Office             *model.Office                `gorm:"foreignKey:OfficeID;references:ID"`
	CreateByInfo       *model.ManagerResponse       `gorm:"foreignKey:CreatedBy;references:EmployeeID" json:"create_by_info"`
	AttendanceCategory *AttendanceCategory          `gorm:"foreignKey:AttendanceCategoryID;references:AttendanceCategoryID"`
}

func (AttendanceRecord) TableName() string {
	return "attendance_records"
}

func (r *AttendanceRecord) ParseGPS() error {
	if r.GPSLocation == "" {
		return nil // Không có dữ liệu GPS
	}

	parts := strings.Split(r.GPSLocation, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid GPS format. Expected 'latitude,longitude'")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return fmt.Errorf("invalid latitude: %v", err)
	}

	long, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return fmt.Errorf("invalid longitude: %v", err)
	}

	r.Latitude = &lat
	r.Longitude = &long
	return nil
}
