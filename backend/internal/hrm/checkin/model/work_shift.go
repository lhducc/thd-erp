package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

type WorkDayEnum string
type TimeOfDayEnum string

const (
	FullDay WorkDayEnum = "1"
	HaftDay WorkDayEnum = "0.5"
	NoWork  WorkDayEnum = "0"
)
const (
	FullTime  TimeOfDayEnum = "Cả ngày"
	Morning   TimeOfDayEnum = "Sáng"
	Night     TimeOfDayEnum = "Tối"
	Afternoon TimeOfDayEnum = "Chiều"
)

type WorkShifts struct {
	WorkShiftID   string        `gorm:"column:workshift_id;type:varchar(5);primaryKey;" json:"workshift_id"`
	WorkShiftName string        `gorm:"column:workshift_name;type:varchar(255);not null" json:"workshift_name" validate:"required,max=255"`
	StartTime     string        `gorm:"column:start_time;type:time;not null" json:"start_time" validate:"required"`
	EndTime       string        `gorm:"column:end_time;type:time;not null" json:"end_time" validate:"required"`
	CheckinFrom   *string       `gorm:"column:checkin_from;type:time" json:"checkin_from"`
	CheckinTo     *string       `gorm:"column:checkin_to;type:time" json:"checkin_to"`
	CheckoutFrom  *string       `gorm:"column:checkout_from;type:time" json:"checkout_from"`
	CheckoutTo    *string       `gorm:"column:checkout_to;type:time" json:"checkout_to"`
	HasBreak      bool          `gorm:"column:has_break;type:boolean" json:"has_break"`
	BreakStart    *string       `gorm:"column:break_start;type:time" json:"break_start"`
	BreakEnd      *string       `gorm:"column:break_end;type:time" json:"break_end"`
	WorkHours     float64       `gorm:"column:work_hours;type:decimal(4,2)" json:"work_hours" validate:"gte=0"`
	WorkDay       WorkDayEnum   `gorm:"column:work_day;type:work_day_enum;default:0" json:"work_day" validate:"required"`
	CoefNormalDay float64       `gorm:"column:coef_normal_day;type:decimal(3,2);default:1.00" json:"coef_normal_day" validate:"gte=0"`
	CoefWeekend   float64       `gorm:"column:coef_weekend;type:decimal(3,2);default:1.00" json:"coef_weekend" validate:"gte=0"`
	CoefHoliday   float64       `gorm:"column:coef_holiday;type:decimal(3,2);default:1.00" json:"coef_holiday" validate:"gte=0"`
	CreatedBy     string        `gorm:"column:created_by;not null" json:"created_by"`
	CreatedDate   time.Time     `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	IsDeleted     bool          `gorm:"column:is_deleted;type:boolean;default:false" json:"is_deleted"`
	TimeOfDay     TimeOfDayEnum `gorm:"column:time_of_day;type:varchar(50);" json:"time_of_day"`

	EffectiveDate  time.Time  `gorm:"column:effective_date;type:timestamp;not null" json:"effective_date" validate:"required"`
	ExpirationDate *time.Time `gorm:"column:expiration_date;type:timestamp" json:"expiration_date"`

	Creator *model.Employee `gorm:"foreignKey:CreatedBy;references:employee_id" json:"creator,omitempty"`
}
type WorkShiftsRequest struct {
	WorkShiftName  string        `json:"workshift_name"`
	StartTime      string        `json:"start_time"`
	EndTime        string        `json:"end_time"`
	CheckinFrom    *string       `json:"checkin_from"`
	CheckinTo      *string       `json:"checkin_to"`
	CheckoutFrom   *string       `json:"checkout_from"`
	CheckoutTo     *string       `json:"checkout_to"`
	HasBreak       bool          `json:"has_break"`
	BreakStart     *string       `json:"break_start"`
	BreakEnd       *string       `json:"break_end"`
	WorkHours      float64       `json:"work_hours"`
	WorkDay        WorkDayEnum   `json:"work_day"`
	CoefNormalDay  float64       `json:"coef_normal_day"`
	CoefWeekend    float64       `json:"coef_weekend"`
	CoefHoliday    float64       `json:"coef_holiday"`
	CreatedBy      string        `json:"created_by"`
	TimeOfDay      TimeOfDayEnum `json:"time_of_day"`
	EffectiveDate  time.Time     `json:"effective_date"`
	ExpirationDate *time.Time    `json:"expiration_date"`
}

func (WorkShifts) TableName() string {
	return "workshifts"
}

var validateWS = validator.New()

var fieldWSVietnamese = map[string]string{
	"WorkShiftID":    "Mã ca làm",
	"WorkShiftName":  "Tên ca làm",
	"StartTime":      "Thời gian bắt đầu",
	"EndTime":        "Thời gian kết thúc",
	"CheckinFrom":    "Giờ cho phép check-in từ",
	"CheckinTo":      "Giờ cho phép check-in đến",
	"CheckoutFrom":   "Giờ cho phép check-out từ",
	"CheckoutTo":     "Giờ cho phép check-out đến",
	"HasBreak":       "Có giờ nghỉ giữa ca",
	"BreakStart":     "Giờ bắt đầu nghỉ",
	"BreakEnd":       "Giờ kết thúc nghỉ",
	"WorkHours":      "Số giờ làm việc",
	"WorkDay":        "Số ngày công tương đương",
	"CoefNormalDay":  "Hệ số ngày thường",
	"CoefWeekend":    "Hệ số cuối tuần",
	"CoefHoliday":    "Hệ số ngày lễ",
	"CreatedBy":      "Người tạo",
	"CreatedDate":    "Ngày tạo",
	"TimeOfDay":      "Buổi",
	"EffectiveDate":  "Ngày hiệu lực",
	"ExpirationDate": "Ngày hết hiệu lực",
}

// Validate kiểm tra dữ liệu và trả về lỗi tiếng Việt
func (w *WorkShifts) Validate() error {
	err := validateWS.Struct(w)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fieldName := fieldWSVietnamese[fieldErr.Field()]
			if fieldName == "" {
				fieldName = fieldErr.Field() // fallback nếu chưa định nghĩa trong map
			}
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
				return fmt.Errorf("%s không hợp lệ (%s)", fieldName, fieldErr.Tag())
			}
		}
	}
	//check effective & expiration dates
	if w.ExpirationDate != nil {
		if w.EffectiveDate.After(*w.ExpirationDate) {
			return errors.New("Ngày bắt đầu hiệu lực không được sau ngày hết hiệu lực")
		}
		if w.EffectiveDate.Equal(*w.ExpirationDate) {
			return errors.New("Ngày bắt đầu và ngày hết hiệu lực không được trùng nhau")
		}
	}
	return nil
}

func ConvertToWorkShifts(workShiftsCreate WorkShiftsRequest) WorkShifts {
	shifts := WorkShifts{
		WorkShiftName:  workShiftsCreate.WorkShiftName,
		StartTime:      workShiftsCreate.StartTime,
		EndTime:        workShiftsCreate.EndTime,
		CheckinFrom:    workShiftsCreate.CheckinFrom,
		CheckinTo:      workShiftsCreate.CheckinTo,
		CheckoutFrom:   workShiftsCreate.CheckoutFrom,
		CheckoutTo:     workShiftsCreate.CheckoutTo,
		HasBreak:       workShiftsCreate.HasBreak,
		BreakStart:     workShiftsCreate.BreakStart,
		BreakEnd:       workShiftsCreate.BreakEnd,
		WorkHours:      workShiftsCreate.WorkHours,
		WorkDay:        workShiftsCreate.WorkDay,
		CoefNormalDay:  workShiftsCreate.CoefNormalDay,
		CoefWeekend:    workShiftsCreate.CoefWeekend,
		CoefHoliday:    workShiftsCreate.CoefHoliday,
		CreatedBy:      workShiftsCreate.CreatedBy,
		TimeOfDay:      workShiftsCreate.TimeOfDay,
		EffectiveDate:  workShiftsCreate.EffectiveDate,
		ExpirationDate: workShiftsCreate.ExpirationDate,
	}
	return shifts
}
