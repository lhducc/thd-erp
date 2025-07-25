package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"time"
)

type WorkShiftsRequest struct {
	WorkShiftID    string               `json:"workshift_id"`
	WorkShiftName  string               `json:"workshift_name"`
	StartTime      string               `json:"start_time"`
	EndTime        string               `json:"end_time"`
	CheckinFrom    *string              `json:"checkin_from"`
	CheckinTo      *string              `json:"checkin_to"`
	CheckoutFrom   *string              `json:"checkout_from"`
	CheckoutTo     *string              `json:"checkout_to"`
	HasBreak       bool                 `json:"has_break"`
	BreakStart     *string              `json:"break_start"`
	BreakEnd       *string              `json:"break_end"`
	WorkHours      float64              `json:"work_hours"`
	WorkDay        variable.WorkDayEnum `json:"work_day"`
	CoefNormalDay  float64              `json:"coef_normal_day"`
	CoefWeekend    float64              `json:"coef_weekend"`
	CoefHoliday    float64              `json:"coef_holiday"`
	CreatedBy      string               `json:"created_by"`
	EffectiveDate  time.Time            `json:"effective_date"`
	ExpirationDate *time.Time           `json:"expiration_date"`
}

var validWorkDays = map[variable.WorkDayEnum]bool{
	variable.FullDay: true,
	variable.HaftDay: true,
	variable.NoWork:  true,
}

func (w *WorkShiftsRequest) Validate() error {
	if w.WorkShiftID == "" {
		return errors.New("Mã ca làm việc là bắt buộc")
	}

	if w.WorkShiftName == "" {
		return errors.New("Tên ca làm việc là bắt buộc")
	}

	if len(w.WorkShiftName) > 255 {
		return errors.New("Tên ca làm việc không được vượt quá 255 ký tự")
	}

	if w.StartTime == "" {
		return errors.New("Thời gian bắt đầu là bắt buộc")
	}

	if w.EndTime == "" {
		return errors.New("Thời gian kết thúc là bắt buộc")
	}

	if w.WorkHours < 0 {
		return errors.New("Số giờ làm việc phải lớn hơn hoặc bằng 0")
	}

	if _, ok := validWorkDays[w.WorkDay]; !ok {
		return fmt.Errorf("Giá trị work_day không hợp lệ: %s", w.WorkDay)
	}

	if w.CoefNormalDay < 0 {
		return errors.New("Hệ số ngày thường phải lớn hơn hoặc bằng 0")
	}
	if w.CoefWeekend < 0 {
		return errors.New("Hệ số cuối tuần phải lớn hơn hoặc bằng 0")
	}
	if w.CoefHoliday < 0 {
		return errors.New("Hệ số ngày lễ phải lớn hơn hoặc bằng 0")
	}

	if w.EffectiveDate.IsZero() {
		return errors.New("Ngày bắt đầu hiệu lực là bắt buộc")
	}

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

func (ws *WorkShiftsRequest) ConvertToWorkShifts() *model.WorkShifts {
	return &model.WorkShifts{
		WorkShiftID:    ws.WorkShiftID,
		WorkShiftName:  ws.WorkShiftName,
		StartTime:      ws.StartTime,
		EndTime:        ws.EndTime,
		CheckinFrom:    ws.CheckinFrom,
		CheckinTo:      ws.CheckinTo,
		CheckoutFrom:   ws.CheckoutFrom,
		CheckoutTo:     ws.CheckoutTo,
		HasBreak:       ws.HasBreak,
		BreakStart:     ws.BreakStart,
		BreakEnd:       ws.BreakEnd,
		WorkHours:      ws.WorkHours,
		WorkDay:        ws.WorkDay,
		CoefNormalDay:  ws.CoefNormalDay,
		CoefWeekend:    ws.CoefWeekend,
		CoefHoliday:    ws.CoefHoliday,
		CreatedBy:      ws.CreatedBy,
		EffectiveDate:  ws.EffectiveDate,
		ExpirationDate: ws.ExpirationDate,
	}
}
