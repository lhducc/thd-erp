package service

import (
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
}

func NewEmployeeWorkshiftService(
	repo repo_interface.EmployeeWorkShiftRepo,
	employeeRepo usecase.EmployeeRepo,
	workshiftRepo repo_interface.WorkShiftRepo,
) service_interface.EmployeeWorkshiftService {
	return &employeeWorkshiftService{
		repo:          repo,
		employeeRepo:  employeeRepo,
		workshiftRepo: workshiftRepo,
	}
}

func (s *employeeWorkshiftService) GetByUserId(userId string) ([]model.EmployeeWorkshift, error) {
	return s.repo.GetAllByEmployeeID(userId)
}

func (s *employeeWorkshiftService) GetAll() ([]model.EmployeeWorkshift, error) {
	return s.repo.GetAll()
}
func (s *employeeWorkshiftService) Register(employeeID string, workshiftID string, date time.Time) error {
	existing := s.repo.IsExisting(employeeID, workshiftID, date)
	if existing {
		return fmt.Errorf("employee workshift %s already exists", employeeID)
	}
	existingUser, err := s.employeeRepo.GetUserById(employeeID)
	if err != nil {
		return err
	}

	existingWorkshift, err := s.workshiftRepo.GetWorkShiftById(workshiftID)
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

func (s *employeeWorkshiftService) Delete(id uint) error {
	return s.repo.Delete(id)
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
