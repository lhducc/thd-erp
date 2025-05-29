package model

import (
	"fmt"

	"github.com/go-playground/validator"
)

type WorkShifts struct {
	WorkShiftId             string `gorm:"type:varchar(4);primaryKey;column:work_shift_id" json:"work_shift_id"`
	WorkShiftName           string `gorm:"type:varchar(100);column:work_shift_name" json:"work_shift_name"`
	StartTime               string `gorm:"column:start_time" example:"13:30" json:"start_time" binding:"required"`
	EndTime                 string `gorm:"column:end_time" example:"17:30" json:"end_time" binding:"required"`
	EquivalentWorkingNumber string `gorm:"type:decimal(4,2);column:equivalent_working_number" json:"equivalent_working_number"`
	Note                    string `gorm:"type:text;column:note" json:"note"`
	IsDeleted               bool   `gorm:"type:boolean;column:is_deleted;default:false"`

	AllowedWorkingScheduleID string                  `gorm:"type:char(4);column:allowed_working_schedule_id" json:"allowed_working_schedule_id"`
	AllowedWorkingSchedule   *AllowedWorkingSchedule `gorm:"foreignKey:AllowedWorkingScheduleID;references:ID"`
}

func (WorkShifts) TableName() string {
	return "work_shifts"
}

func UpdateWorkShiftFields(existing *WorkShifts, updated WorkShifts) {
	existing.WorkShiftName = updated.WorkShiftName
	existing.StartTime = updated.StartTime
	existing.EndTime = updated.EndTime
	existing.EquivalentWorkingNumber = updated.EquivalentWorkingNumber
	existing.Note = updated.Note
}

var validateWS = validator.New()

var fieldWSVietnamese = map[string]string{
	"WorkShiftId":              "Mã ca làm",
	"WorkShiftName":            "Tên ca làm",
	"StartTime":                "Thời gian bắt đầu",
	"EndTime":                  "Thời gian kết thúc",
	"EquivalentWorkingNumber":  "Số giờ làm tương đương",
	"Note":                     "Ghi chú",
	"AllowedWorkingScheduleID": "Mã lịch làm việc cho phép",
}

// Validate kiểm tra dữ liệu và trả về lỗi tiếng Việt
func (w *WorkShifts) Validate() error {
	err := validateWS.Struct(w)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fieldName := fieldWSVietnamese[fieldErr.Field()]
			switch fieldErr.Tag() {
			case "required":
				return fmt.Errorf("%s là bắt buộc", fieldName)
			case "len":
				return fmt.Errorf("%s phải có độ dài chính xác là %s ký tự", fieldName, fieldErr.Param())
			case "max":
				return fmt.Errorf("%s không được vượt quá %s ký tự", fieldName, fieldErr.Param())
			case "datetime":
				return fmt.Errorf("%s phải đúng định dạng HH:mm", fieldName)
			default:
				return fmt.Errorf("%s không hợp lệ", fieldName)
			}
		}
	}
	return nil
}
