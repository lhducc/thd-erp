package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"time"
)

type AssignEmployeeWorkshiftRequest struct {
	EmpWorkShift []model.EmployeeWorkshift `json:"emp_work_shift"`
	ScheduleIDs  []int                     `json:"schedule_ids"`
}

type ManagerEmployeeScheduleDTO struct {
	ID            int       `json:"id"`
	EmployeeID    string    `json:"employee_id"`
	FullName      string    `json:"full_name"`
	Manager       string    `json:"manager"`
	WorkshiftID   string    `json:"workshift_id"`
	Date          time.Time `json:"date"`
	WorkshiftName string    `json:"workshift_name"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
}
