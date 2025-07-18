package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"strings"
	"time"
)

type workScheduleRegisterService struct {
	repo             repo_interface.WorkScheduleRegisterRepo
	autoScheduleRepo repo_interface.WorkScheduleRepo
	employeeRepo     *repository.UserStore
}

func NewWorkScheduleSRegisterervice(repo repo_interface.WorkScheduleRegisterRepo, autoScheduleRepo repo_interface.WorkScheduleRepo, employeeRepo *repository.UserStore) service_interface.WorkScheduleRegisterServiceInterface {
	return &workScheduleRegisterService{
		repo:             repo,
		autoScheduleRepo: autoScheduleRepo,
		employeeRepo:     employeeRepo,
	}
}

func (s *workScheduleRegisterService) CreateNewWorkSchedule(c context.Context, workSchedule *model.WorkScheduleRegister) error {
	var weekdayShift []model.WorkScheduleRegisterShift
	wScheduleId := workSchedule.WorkScheduleRegisterID
	timeNow := time.Now()
	effectiveBeforeNow := workSchedule.EffectiveDate.Before(timeNow)
	expirationAfterNow := workSchedule.ExpirationDate.After(timeNow)

	exists, err := s.repo.IsExistsByName(c, workSchedule.WorkScheduleRegisterName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Tên lịch làm việc đã tồn tại")
	}
	for _, weekday := range workSchedule.Weekdays {
		weekdayShift = append(weekdayShift, model.WorkScheduleRegisterShift{
			Weekday:                weekday.Weekday,
			WorkScheduleRegisterID: wScheduleId,
			WorkShiftID:            weekday.WorkShiftID,
		})
	}

	switch {
	case workSchedule.EffectiveDate.After(timeNow):
		workSchedule.Status = variable.InActive
	case effectiveBeforeNow && expirationAfterNow:
		workSchedule.Status = variable.Active
	case effectiveBeforeNow && !expirationAfterNow:
		workSchedule.Status = variable.Expired
	}

	return s.repo.Save(c, workSchedule, weekdayShift)
}

func (s *workScheduleRegisterService) DeleteWorkSchedule(c context.Context, id int) error {
	exists, err := s.repo.IsExistsByID(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc cần xóa không tồn tại")
	}
	return s.repo.Delete(c, id)
}
func (s *workScheduleRegisterService) UpdateWorkSchedule(c context.Context, workSchedule *model.WorkScheduleRegister, id int) error {
	exists, err := s.repo.IsExistsByID(c, id)
	timeNow := time.Now()
	effectiveBeforeNow := workSchedule.EffectiveDate.Before(timeNow)
	expirationAfterNow := workSchedule.ExpirationDate.After(timeNow)

	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc cần sửa đổi không tồn tại")
	}

	switch {
	case workSchedule.EffectiveDate.After(timeNow):
		workSchedule.Status = variable.InActive
	case effectiveBeforeNow && expirationAfterNow:
		workSchedule.Status = variable.Active
	case effectiveBeforeNow && !expirationAfterNow:
		workSchedule.Status = variable.Expired
	}

	return s.repo.Update(c, workSchedule, id)
}

func (s *workScheduleRegisterService) checkEmployeeExistsInOtherSchedules(c context.Context, employeeIDs []string, idSchedule int) error {
	for _, empID := range employeeIDs {
		wSREmployee, wsr, exist, err := s.repo.CheckExistEmployeeScheduleRegister(c, empID, &idSchedule)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("Nhân viên '%s: %s' đã thuộc lịch đăng ký tên '%s'",
				empID, wSREmployee.Employee.Fullname, wsr.WorkScheduleRegisterName)
		}

		wSEmployee, wSchedule, exist, err := s.autoScheduleRepo.CheckExistEmployeeSchedule(c, empID, nil)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("Nhân viên '%s: %s' đã thuộc lịch tự động tên '%s'",
				empID, wSEmployee.Employee.Fullname,
				wSchedule.WorkScheduleName,
			)
		}
	}
	return nil
}

func createEmployeeRecords(employeeIDs []string, scheduleID int, timeNow time.Time) []model.WorkScheduleRegisterEmployee {
	var records []model.WorkScheduleRegisterEmployee
	for _, empID := range employeeIDs {
		records = append(records, model.WorkScheduleRegisterEmployee{
			EmployeeID:             empID,
			WorkScheduleRegisterID: &scheduleID,
			AssignedAt:             &timeNow,
		})
	}
	return records
}

func createManagerRecords(managers []model.ManagerAssignRequest, scheduleID int, now time.Time) []model.WorkScheduleRegisterManager {
	var records []model.WorkScheduleRegisterManager
	for _, m := range managers {
		records = append(records, model.WorkScheduleRegisterManager{
			EmployeeID:             m.ManagerID,
			WorkScheduleRegisterID: scheduleID,
			IsReading:              m.IsReading,
			IsEditing:              m.IsEditing,
		})
	}
	return records
}

func (s *workScheduleRegisterService) AssignEmployeeToWorkSchedule(c context.Context, req *model.AssignEmployeeRequest, idSchedule int) error {
	timeNow := time.Now()

	if err := s.checkEmployeeExistsInOtherSchedules(c, req.EmployeeIDs, idSchedule); err != nil {
		return err
	}

	exists, err := s.repo.IsExistsByID(c, idSchedule)
	if err != nil {
		return fmt.Errorf("Lỗi khi tìm lịch làm việc: %w", err)
	}
	if !exists {
		return errors.New("Không tìm thấy lịch làm việc")
	}

	employeeRecords := createEmployeeRecords(req.EmployeeIDs, idSchedule, timeNow)
	managerRecords := createManagerRecords(req.Managers, idSchedule, timeNow)

	allEmployeeIDs := make([]string, 0, len(employeeRecords)+len(managerRecords))
	for _, emp := range employeeRecords {
		allEmployeeIDs = append(allEmployeeIDs, emp.EmployeeID)
	}
	for _, mgr := range managerRecords {
		allEmployeeIDs = append(allEmployeeIDs, mgr.EmployeeID)
	}

	missing, err := s.validateEmployeeExists(allEmployeeIDs)
	if err != nil {
		return err
	}
	if missing != nil {
		return fmt.Errorf("không tìm thấy các nhân viên có mã: %s", strings.Join(missing, ", "))
	}

	if err := s.repo.AssignOrUpdateEmployeesAndManagers(c, employeeRecords, managerRecords); err != nil {
		return err
	}

	return nil
}

func (s *workScheduleRegisterService) GetAllWorkSchedule(ctx context.Context) ([]model.WorkScheduleRegister, error) {
	return s.repo.GetAll(ctx)
}

func (s *workScheduleRegisterService) GetWorkScheduleByID(ctx context.Context, id int) (*model.WorkScheduleRegister, error) {
	return s.repo.GetByID(ctx, id)
}

//func (s *workScheduleRegisterService) RemoveEmployeesFromWorkSchedule(
//	ctx context.Context,
//	employeeIDs []string,
//	managerIDs []string,
//	scheduleID int,
//) error {
//	exists, err := s.repo.IsExistsByID(ctx, scheduleID)
//	if err != nil {
//		return err
//	}
//	if !exists {
//		return errors.New("Không tìm thấy lịch làm việc")
//	}
//
//	return s.repo.DeleteEmployeesFromWorkSchedule(ctx, employeeIDs, managerIDs, scheduleID)
//}

func (s *workScheduleRegisterService) validateEmployeeExists(employeeIDs []string) ([]string, error) {
	missing, err := s.employeeRepo.CheckExistEmployee(employeeIDs)
	if err != nil {
		return nil, err
	}
	return missing, nil
}

func (s *workScheduleRegisterService) DeleteEmployeeFromWorkScheduleRegister(ctx context.Context, employeeID string, workScheduleID int) error {
	// check exist schedule
	exists, err := s.repo.IsExistsByID(ctx, workScheduleID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc cần sửa đổi không tồn tại")
	}

	err = s.repo.DeleteEmployeeFromWorkSchedule(ctx, employeeID, workScheduleID)
	if err != nil {
		return err
	}
	return nil
}

func (s *workScheduleRegisterService) DeleteManagerFromWorkScheduleRegister(ctx context.Context, managerID string, workScheduleID int) error {
	// check exist schedule
	exists, err := s.repo.IsExistsByID(ctx, workScheduleID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc cần sửa đổi không tồn tại")
	}

	err = s.repo.DeleteManagerFromWorkSchedule(ctx, managerID, workScheduleID)
	if err != nil {
		return err
	}
	return nil
}
