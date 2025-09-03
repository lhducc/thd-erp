package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/pkg/struct_support"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"time"
)

type WorkScheduleRequest struct {
	WorkScheduleName string                    ` json:"work_schedule_name"`
	OfficeID         string                    `json:"office_id"`
	RepeatType       *variable.RepeatTypeEnum  `json:"repeat_type"`
	RepeatCycle      *int                      `json:"repeat_cycle"`
	EffectiveDate    time.Time                 `json:"effective_date"`
	ExpirationDate   *time.Time                `json:"expiration_date"`
	WeekDays         []model.WorkScheduleShift `json:"weekdays"`
}

type AssignManagersRequest struct {
	Managers []ManagerAssignRequest `json:"managers"`
}

type ManagerAssignRequest struct {
	ManagerID string `json:"manager_id"`
	IsReading bool   `json:"is_reading"`
	IsEditing bool   `json:"is_editing"`
}

type AutoRecurringRequest struct {
	IsAutoRecurring bool `json:"is_auto_recurring"`
}

func ConvertToWorkSchedule(req *WorkScheduleRequest) model.WorkSchedule {
	return model.WorkSchedule{
		WorkScheduleName: req.WorkScheduleName,
		OfficeID:         req.OfficeID,
		RepeatType:       req.RepeatType,
		RepeatCycle:      req.RepeatCycle,
		EffectiveDate:    req.EffectiveDate,
		ExpirationDate:   req.ExpirationDate,
		Weekdays:         req.WeekDays,
	}
}

func (req *WorkScheduleRequest) Validate() error {
	if req.WorkScheduleName == "" {
		return errors.New("work_schedule_name is required")
	}

	if req.OfficeID == "" {
		return errors.New("office_id is required")
	}

	if req.RepeatCycle != nil || req.RepeatType != nil {
		if _, ok := struct_support.ValidRepeatTypes[*req.RepeatType]; !ok {
			return fmt.Errorf("invalid repeat_type: %s", req.RepeatType)
		}
	}

	for _, weekday := range req.WeekDays {
		if _, ok := struct_support.ValidWeekdays[weekday.Weekday]; !ok {
			return fmt.Errorf("invalid week_day: %s", weekday.Weekday)
		}
	}

	if req.EffectiveDate.IsZero() {
		return errors.New("effective_date and expiration_date must not be empty")
	}

	if req.ExpirationDate != nil {
		if req.EffectiveDate.After(*req.ExpirationDate) {
			return errors.New("effective_date must be before expiration_date")
		}
	}

	return nil
}
