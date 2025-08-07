package service

import (
	"context"
	checkin_model "erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"errors"
	"strings"
	"time"
)

type shiftAllocationService struct {
	scheduleRepo repo_interface.WorkScheduleRepo
	employeeRepo *repository.UserStore
	empShiftRepo repo_interface.EmployeeWorkShiftRepo
}

func NewShiftAllocationService(workSchRepo repo_interface.WorkScheduleRepo,
	employeeRepo *repository.UserStore, empShiftRepo repo_interface.EmployeeWorkShiftRepo) service_interface.ShiftAllocationServiceInterface {
	return &shiftAllocationService{
		scheduleRepo: workSchRepo,
		employeeRepo: employeeRepo,
		empShiftRepo: empShiftRepo,
	}
}

func (s *shiftAllocationService) GetListShiftAllocation(ctx context.Context, req dto.GetShiftAllocationRequest, managerID string) ([]dto.EmployeeScheduleResponse, error) {
	var employeeSchedules []dto.EmployeeScheduleResponse
	var (
		employees          []model.Employee
		err                error
		startDate, endDate time.Time
	)

	now := time.Now().Local()
	year := req.Year
	month := req.Month

	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}

	startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate = startDate.AddDate(0, 1, -1)

	if strings.TrimSpace(managerID) != "" {
		exist, err := s.employeeRepo.CheckExists(managerID)
		if err != nil {
			return nil, errors.New("Lỗi truy vấn quản lý lịch làm việc")
		}
		if exist {
			employees, err = s.employeeRepo.GetBySchedule(ctx, req.ScheduleIDs, managerID, req.Filter)
		}
	} else {
		employees, err = s.employeeRepo.GetBySchedule(ctx, req.ScheduleIDs, "", req.Filter)
	}
	if err != nil {
		return nil, err
	}
	if employees == nil {
		return nil, nil
	}

	for _, employee := range employees {
		var scheduleInfors []checkin_model.WorkSchedule
		var workShiftTimelines []dto.WorkShiftTimeline

		if employee.ScheduleID != nil {
			scheduleInfor, err := s.scheduleRepo.GetByID(ctx, *employee.ScheduleID)
			if err != nil {
				return nil, errors.New("Lỗi hệ thống, lỗi khi truy vấn lịch làm việc thuộc nhân viên")
			}
			if scheduleInfor != nil {
				scheduleInfors = append(scheduleInfors, *scheduleInfor)
			}
		}

		employeeWorkshift, err := s.empShiftRepo.GetEmployeeWorkShiftsByMonthYear(ctx, employee.EmployeeID, startDate, endDate)
		if err != nil {
			return nil, errors.New("Lỗi hệ thống, lỗi khi truy vấn lịch làm việc thuộc nhân viên")
		}
		for _, workShift := range employeeWorkshift {
			wsTimeline := dto.ConvertToTimeLine(&workShift)
			workShiftTimelines = append(workShiftTimelines, *wsTimeline)
		}
		if employee.Department != nil {
			empSchedule := dto.EmployeeScheduleResponse{
				EmployeeID:     employee.EmployeeID,
				Fullname:       employee.Fullname,
				WorkShifts:     workShiftTimelines,
				Schedules:      scheduleInfors,
				DepartmentName: employee.Department.Name,
			}
			employeeSchedules = append(employeeSchedules, empSchedule)
		}
	}

	return employeeSchedules, nil
}
