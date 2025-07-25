package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"errors"
	"fmt"
	"time"
)

type EmployeeScheduleResponse struct {
	EmployeeID     string              `json:"employee_id"`
	Fullname       string              `json:"full_name"`
	DepartmentName string              `json:"department_name"`
	Schedules      []ScheduleInfo      `json:"schedules"`
	WorkShifts     []WorkShiftTimeline `json:"work_shifts"`
}

func ConvertToScheduleInfor(ws *model.WorkSchedule) *ScheduleInfo {
	return &ScheduleInfo{
		ScheduleID:      ws.WorkScheduleID,
		ScheduleName:    ws.WorkScheduleName,
		IsScheduleAuto:  ws.IsScheduleAuto,
		IsAutoRecurring: ws.IsAutoRecurring,
	}
}
func ConvertToTimeLine(ew *model.EmployeeWorkshift) *WorkShiftTimeline {
	return &WorkShiftTimeline{
		WorkShiftID:   ew.WorkShiftID,
		WorkShiftName: ew.WorkShift.WorkShiftName,
		Date:          ew.Date,
		StartTime:     ew.WorkShift.StartTime,
		EndTime:       ew.WorkShift.EndTime,
	}
}

type ScheduleInfo struct {
	ScheduleID      int    `json:"work_schedule_id"`
	ScheduleName    string `json:"work_schedule_name"`
	IsScheduleAuto  bool   `gorm:"column:is_schedule_auto;default:false" json:"is_schedule_auto"`
	IsAutoRecurring bool   `gorm:"column:is_auto_recurring;default:false" json:"is_auto_recurring"`
}

type WorkShiftTimeline struct {
	Date          time.Time `json:"date"`
	WorkShiftID   string    `json:"work_shift_id"`
	WorkShiftName string    `json:"work_shift_name"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
}

type GetShiftAllocationRequest struct {
	Month       int    `form:"month" validate:"required,min=1,max=12"`
	Year        int    `form:"year" validate:"required,min=2000"`
	Filter      string `form:"filter" enums:"register,auto,no_schedule"`
	ScheduleIDs []int  `form:"schedule_ids"`
}

type ApplyAutoScheduleRequest struct {
	ScheduleID int `json:"schedule_id"`
	Month      int `form:"month" validate:"required,min=1,max=12"`
	Year       int `form:"year" validate:"required,min=2000"`
}

func (v *GetShiftAllocationRequest) ValidateMonthYear() error {
	if v.Month < 1 || v.Month > 12 {
		return errors.New("tháng phải từ 1 đến 12")
	}

	currentYear := time.Now().Year()
	if v.Year < 2000 || v.Year > time.Now().Year()+1 {
		return fmt.Errorf("năm phải từ 2000 đến %d", currentYear+1)
	}
	return nil
}
