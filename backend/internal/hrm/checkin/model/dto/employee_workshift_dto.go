package dto

import "erp/backend/internal/hrm/checkin/model"

type AssignEmployeeWorkshiftRequest struct {
	EmpWorkShift []model.EmployeeWorkshift `json:"emp_work_shift"`
	ScheduleIDs  []int                     `json:"schedule_ids"`
}
