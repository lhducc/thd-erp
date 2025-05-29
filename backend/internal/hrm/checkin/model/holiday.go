package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

type Holiday struct {
	HolidayID               string    `gorm:"primaryKey;type:varchar(4);column:holiday_id" json:"holiday_id" validate:"required,len=4"`
	HolidayName             string    `gorm:"type:varchar(100);column:holiday_name" json:"holiday_name" validate:"required,max=100"`
	HolidayType             string    `gorm:"type:varchar(50);column:holiday_type" json:"holiday_type" validate:"required,max=50"`
	EquivalentWorkingNumber float64   `gorm:"type:decimal(4,2);column:equivalent_working_number" json:"equivalent_working_number" validate:"required,gte=0"`
	Date                    time.Time `gorm:"type:date;column:date" json:"date" validate:"required"`
	Shift                   string    `gorm:"type:varchar(50);column:shift" json:"shift" validate:"required,max=50"`
	Note                    string    `gorm:"type:text;column:note" json:"note" validate:"omitempty"`
	CreatedAt               time.Time `gorm:"column:created_at" json:"created_at"`
	IsDeleted               bool      `gorm:"type:boolean;column:is_deleted;default:false" json:"is_deleted"`

	AllowedWorkingSchedules []AllowedWorkingSchedule `gorm:"foreignKey:HolidayID;references:HolidayID" json:"-"`
}

func (Holiday) TableName() string { return "holiday" }

func UpdateHolidayFields(existing *Holiday, holiday Holiday) {
	// Update the fields from the input holiday
	existing.HolidayName = holiday.HolidayName
	existing.HolidayType = holiday.HolidayType
	existing.EquivalentWorkingNumber = holiday.EquivalentWorkingNumber
	existing.Date = holiday.Date
	existing.Shift = holiday.Shift
	existing.Note = holiday.Note
	existing.IsDeleted = holiday.IsDeleted
}

// Validator instance
var validateHoliday = validator.New()

// Bản đồ ánh xạ field -> tên tiếng Việt
var fieldHolidayVietnamese = map[string]string{
	"HolidayID":               "Mã ngày nghỉ",
	"HolidayName":             "Tên ngày nghỉ",
	"HolidayType":             "Loại ngày nghỉ",
	"EquivalentWorkingNumber": "Số giờ làm tương đương",
	"Date":                    "Ngày nghỉ",
	"Shift":                   "Ca nghỉ",
	"Note":                    "Ghi chú",
}

func (h *Holiday) Validate() error {
	err := validateHoliday.Struct(h)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fieldName := fieldHolidayVietnamese[fieldErr.Field()]
			switch fieldErr.Tag() {
			case "required":
				return fmt.Errorf("%s là bắt buộc", fieldName)
			case "len":
				return fmt.Errorf("%s phải có độ dài chính xác là %s ký tự", fieldName, fieldErr.Param())
			case "max":
				return fmt.Errorf("%s không được dài quá %s ký tự", fieldName, fieldErr.Param())
			case "gte":
				return fmt.Errorf("%s phải lớn hơn hoặc bằng %s", fieldName, fieldErr.Param())
			default:
				return fmt.Errorf("%s không hợp lệ", fieldName)
			}
		}
	}
	return nil
}
