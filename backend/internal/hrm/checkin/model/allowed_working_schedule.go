package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

type AllowedWorkingSchedule struct {
	ID        string    `gorm:"type:char(4);primaryKey;column:id" json:"id" validate:"required,len=4"`
	Office    string    `gorm:"type:varchar(50);column:office" json:"office" validate:"required,max=50"`
	Day       time.Time `gorm:"type:date;column:day" json:"day" validate:"required"`
	Shift     string    `gorm:"type:varchar(20);column:shift" json:"shift" validate:"required,max=20"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	IsAllowed bool      `gorm:"type:boolean;default:false;column:is_allowed" json:"is_allowed"`
	IsDeleted bool      `gorm:"type:boolean;column:is_deleted;default:false"`
	HolidayID string    `gorm:"type:char(4);column:holiday_id" json:"holiday_id" validate:"omitempty,len=4"`

	Holiday       *Holiday     `gorm:"foreignKey:HolidayID;references:ID" json:"-"`
	WorkingShifts []WorkShifts `gorm:"foreignKey:AllowedWorkingScheduleID;references:ID" json:"-"`
}

func (AllowedWorkingSchedule) TableName() string { return "allowed_working_schedule" }

func UpdateAllowedWorkingScheduleFields(existing *AllowedWorkingSchedule, updated AllowedWorkingSchedule) {
	existing.Office = updated.Office
	existing.Day = updated.Day
	existing.Shift = updated.Shift
	existing.CreatedAt = updated.CreatedAt
	existing.IsAllowed = updated.IsAllowed
	existing.IsDeleted = updated.IsDeleted
	existing.HolidayID = updated.HolidayID
}

var validate = validator.New()

var fieldVietnamese = map[string]string{
	"ID":        "Mã lịch làm việc",
	"Office":    "Văn phòng",
	"Day":       "Ngày làm việc",
	"Shift":     "Ca làm",
	"HolidayID": "Mã ngày nghỉ",
}

func (a *AllowedWorkingSchedule) Validate() error {
	err := validate.Struct(a)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fieldName := fieldVietnamese[fieldErr.Field()]
			switch fieldErr.Tag() {
			case "required":
				return fmt.Errorf("%s là bắt buộc", fieldName)
			case "len":
				return fmt.Errorf("%s phải có độ dài chính xác là %s ký tự", fieldName, fieldErr.Param())
			case "max":
				return fmt.Errorf("%s không được dài quá %s ký tự", fieldName, fieldErr.Param())
			default:
				return fmt.Errorf("%s không hợp lệ", fieldName)
			}
		}
	}
	return nil
}
