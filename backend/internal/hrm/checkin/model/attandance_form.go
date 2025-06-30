package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type AttendanceForm struct {
	AttendanceFormID    string    `gorm:"type:char(4);primaryKey;column:attendance_form_id" validate:"required,len=4" json:"attendance_form_id"`
	AttendanceTableName string    `gorm:"type:varchar(100);not null;column:attendance_table_name" validate:"required,max=100" json:"attendance_table_name"`
	AppliedObject       string    `gorm:"type:varchar(50);not null;column:applied_object" validate:"required,max=50" json:"applied_object"`
	WorkingDays         int16     `gorm:"type:smallint;not null;column:working_days" validate:"required" json:"working_days"`
	Office              string    `gorm:"type:varchar(50);not null;column:office" validate:"required,max=50" json:"office"`
	Department          string    `gorm:"type:varchar(50);not null;column:department" validate:"required,max=50" json:"department"`
	TimeWorkType        string    `gorm:"not null;column:time_work_type" validate:"required,oneof='Theo ca' 'Theo ngày'" json:"time_work_type"`
	FlexWorkType        string    `gorm:"not null;column:flex_work_type" validate:"required,oneof='Theo tháng' 'Linh hoạt'" json:"flex_work_type"`
	Status              string    `gorm:"not null;column:status" validate:"required,oneof='Đang áp dụng' 'Không áp dụng' 'Chưa áp dụng'" json:"status"`
	EffectiveDate       time.Time `gorm:"not null;column:effective_date" validate:"required" json:"effective_date"`
	CreatedAt           time.Time `gorm:"not null;column:created_at" validate:"required" json:"created_at"`
}

func (AttendanceForm) TableName() string { return "attendance_forms" }

var FieldVietnamese = map[string]string{
	"AttendanceFormID":    "Mã mẫu chấm công",
	"AttendanceTableName": "Tên bảng chấm công",
	"AppliedObject":       "Đối tượng áp dụng",
	"WorkingDays":         "Số ngày làm việc",
	"Office":              "Văn phòng",
	"Department":          "Phòng ban",
	"TimeWorkType":        "Loại thời gian làm việc",
	"FlexWorkType":        "Loại làm việc linh hoạt",
	"Status":              "Trạng thái",
	"EffectiveDate":       "Ngày hiệu lực",
	"CreatedAt":           "Ngày tạo",
}
var validateform = validator.New()

func (a *AttendanceForm) Validate() error {
	err := validateform.Struct(a)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fieldName := FieldVietnamese[fieldErr.Field()]
			switch fieldErr.Tag() {
			case "required":
				return fmt.Errorf("%s là bắt buộc", fieldName)
			case "len":
				return fmt.Errorf("%s phải có độ dài chính xác là %s ký tự", fieldName, fieldErr.Param())
			case "max":
				return fmt.Errorf("%s không được dài quá %s ký tự", fieldName, fieldErr.Param())
			case "oneof":
				return fmt.Errorf("%s không hợp lệ", fieldName)
			default:
				return fmt.Errorf("%s không hợp lệ", fieldName)
			}
		}
	}
	return nil
}

func UpdateAttendanceFormFields(existing *AttendanceForm, updated AttendanceForm) {
	existing.AppliedObject = updated.AppliedObject
	existing.AttendanceTableName = updated.AttendanceTableName
	existing.EffectiveDate = updated.EffectiveDate
	existing.FlexWorkType = updated.FlexWorkType
	existing.Department = updated.Department
	existing.Office = updated.Office
	existing.TimeWorkType = updated.TimeWorkType
	existing.WorkingDays = updated.WorkingDays
}

