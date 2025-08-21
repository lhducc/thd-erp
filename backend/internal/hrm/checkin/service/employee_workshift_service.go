package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"fmt"
	"time"
)

type employeeWorkshiftService struct {
	repo          repo_interface.EmployeeWorkShiftRepo
	employeeRepo  usecase.EmployeeRepo
	workshiftRepo repo_interface.WorkShiftRepo
	scheduleRepo  repo_interface.WorkScheduleRepo
}

func NewEmployeeWorkshiftService(
	repo repo_interface.EmployeeWorkShiftRepo,
	employeeRepo usecase.EmployeeRepo,
	workshiftRepo repo_interface.WorkShiftRepo,
	scheduleRepo repo_interface.WorkScheduleRepo,
) service_interface.EmployeeWorkshiftService {
	return &employeeWorkshiftService{
		repo:          repo,
		employeeRepo:  employeeRepo,
		workshiftRepo: workshiftRepo,
		scheduleRepo:  scheduleRepo,
	}
}

func (s *employeeWorkshiftService) GetByUserIdAndMonthYear(ctx context.Context, userId string, month, year int) ([]model.EmployeeWorkshift, error) {
	// Ngày bắt đầu là ngày 1 của tháng
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Ngày kết thúc là ngày cuối tháng
	endDate := startDate.AddDate(0, 1, -1)

	return s.repo.GetEmployeeWorkShiftsByMonthYear(ctx, userId, startDate, endDate)
}

func (s *employeeWorkshiftService) GetAll() ([]model.EmployeeWorkshift, error) {
	return s.repo.GetAll()
}

func (s *employeeWorkshiftService) Register(ctx context.Context, employeeID string, workshiftID string, date time.Time) error {
	existing := s.repo.IsExisting(employeeID, workshiftID, date)
	if existing {
		return fmt.Errorf("employeeRepo workshift %s already exists", employeeID)
	}
	existingUser, err := s.employeeRepo.GetUserById(employeeID)
	if err != nil {
		return err
	}

	existingWorkshift, err := s.workshiftRepo.GetWorkShiftById(ctx, workshiftID)
	if err != nil {
		return err
	}

	if s.repo.Save(&model.EmployeeWorkshift{
		EmployeeID:  existingUser.EmployeeID,
		WorkShiftID: existingWorkshift.WorkShiftID,
		Date:        date,
	}) != nil {
		return err
	}
	return nil
}

func (s *employeeWorkshiftService) RegisterMany(assigns []*model.EmployeeWorkshift) error {
	if len(assigns) == 0 {
		return fmt.Errorf("no workshift assignments provided")
	}

	var newAssigns []*model.EmployeeWorkshift

	for _, assign := range assigns {
		existing := s.repo.IsExisting(assign.EmployeeID, assign.WorkShiftID, assign.Date)
		if !existing {
			newAssigns = append(newAssigns, assign)
		}
	}

	if len(newAssigns) == 0 {
		return nil
	}

	if err := s.repo.SaveMany(newAssigns); err != nil {
		return fmt.Errorf("failed to save workshift assignments: %w", err)
	}

	return nil
}

func (s *employeeWorkshiftService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *employeeWorkshiftService) CheckManagerPermission(ctx context.Context, managerID, employeeID string) (*model.WorkScheduleManager, error) {
	record, err := s.scheduleRepo.CheckManagerPermission(ctx, managerID, employeeID)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *employeeWorkshiftService) DeleteByManager(idEmpShift string) error {
	return s.repo.Delete(idEmpShift)
}

func (s *employeeWorkshiftService) DeletePersonalShift(id string, employeeID string) error {
	return s.repo.DeletePersonalShift(id, employeeID)
}

func (s *employeeWorkshiftService) Update(employeeWorkshiftId string, newWorkshiftId string) error {
	result, err := s.repo.FindByID(employeeWorkshiftId)
	if err != nil {
		return err
	}
	result.WorkShiftID = newWorkshiftId
	err = s.repo.Save(result)
	if err != nil {
		return err
	}
	return nil
}

func (s *employeeWorkshiftService) GetListShiftAllowRegister(ctx context.Context, employeeID string) ([]model.WorkScheduleShift, error) {
	employee, err := s.employeeRepo.GetScheduleOfEmployee(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	if employee.ScheduleID == nil {
		return nil, nil
	}

	exists, err := s.scheduleRepo.IsExistsByScheduleIDRegister(ctx, *employee.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to check schedule existence: %w", err)
	}
	if !exists {
		return nil, nil
	}

	shifts, err := s.scheduleRepo.GetListShiftRegister(ctx, employee.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}

	return shifts, nil
}

func (s *employeeWorkshiftService) Assign(ctx context.Context, empWorkshifts []model.EmployeeWorkshift, scheduleIDs []int) error {
	var employeeWorkshifts []model.EmployeeWorkshift
	for _, empWS := range empWorkshifts {
		existingUser, err := s.employeeRepo.GetUserById(empWS.EmployeeID)
		if err != nil {
			return err
		}

		existingWorkshift, err := s.workshiftRepo.GetWorkShiftById(ctx, empWS.WorkShiftID)
		if err != nil {
			return err
		}

		exist, err := s.repo.CheckShiftConflict(ctx, empWS.EmployeeID, empWS.WorkShiftID, empWS.Date)
		if err != nil {
			return err
		}
		// if exist, continue with another record
		if exist {
			continue
		}
		employeeWorkshifts = append(employeeWorkshifts, model.EmployeeWorkshift{
			EmployeeID:  existingUser.EmployeeID,
			WorkShiftID: existingWorkshift.WorkShiftID,
			Date:        empWS.Date,
		})
	}
	if len(employeeWorkshifts) > 0 {
		if err := s.repo.AssignmentShift(ctx, employeeWorkshifts, scheduleIDs); err != nil {
			return err
		}
	}

	return nil
}

func (s *employeeWorkshiftService) GetAllEmployeeWorkshiftsByMonthYear(ctx context.Context, month, year int) ([]model.EmployeeWorkshift, error) {
	// Ngày bắt đầu là ngày 1 của tháng
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Ngày kết thúc là ngày cuối tháng
	endDate := startDate.AddDate(0, 1, -1)

	return s.repo.GetAllEmployeeWorkShiftsByMonthYear(ctx, startDate, endDate)
}
