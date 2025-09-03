package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"strings"
	"time"
)

type WorkScheduleService struct {
	repo         repo_interface.WorkScheduleRepo
	employeeRepo *repository.UserStore
}

func NewWorkScheduleService(repo repo_interface.WorkScheduleRepo, employeeRepo *repository.UserStore) service_interface.WorkScheduleServiceInterface {
	return &WorkScheduleService{
		repo:         repo,
		employeeRepo: employeeRepo,
	}
}

func (s *WorkScheduleService) CreateNewWorkSchedule(c context.Context, workSchedule *model.WorkSchedule) error {
	var weekdayShift []model.WorkScheduleShift
	wScheduleId := workSchedule.WorkScheduleID
	timeNow := time.Now()

	// Kiểm tra EffectiveDate
	effectiveBeforeNow := workSchedule.EffectiveDate.Before(timeNow)

	// Mặc định expirationAfterNow = true để không ảnh hưởng nếu nil
	expirationAfterNow := true
	if workSchedule.ExpirationDate != nil {
		expirationAfterNow = workSchedule.ExpirationDate.After(timeNow)
	}

	// Check trùng tên lịch làm việc
	exists, err := s.repo.IsExistsByName(c, workSchedule.WorkScheduleName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Tên lịch làm việc đã tồn tại")
	}

	// Tạo danh sách shift cho các weekday
	for _, weekday := range workSchedule.Weekdays {
		weekdayShift = append(weekdayShift, model.WorkScheduleShift{
			Weekday:        weekday.Weekday,
			WorkScheduleID: wScheduleId,
			WorkShiftID:    weekday.WorkShiftID,
		})
	}

	// Xác định trạng thái lịch làm việc
	switch {
	case workSchedule.EffectiveDate.After(timeNow):
		// Ngày hiệu lực trong tương lai → chưa active
		workSchedule.Status = variable.InActive
	case effectiveBeforeNow && expirationAfterNow:
		// Đang trong khoảng hiệu lực
		workSchedule.Status = variable.Active
	case effectiveBeforeNow && !expirationAfterNow:
		// Đã hết hạn
		workSchedule.Status = variable.Expired
	}

	// Lưu xuống DB
	return s.repo.Save(c, workSchedule, weekdayShift)
}

func (s *WorkScheduleService) DeleteWorkScheduleAuto(c context.Context, id int) error {
	exists, err := s.repo.IsExistsByScheduleIDAuto(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc tự động cần xóa không tồn tại")
	}
	return s.repo.Delete(c, id)
}

func (s *WorkScheduleService) DeleteWorkScheduleRegister(c context.Context, id int) error {
	exists, err := s.repo.IsExistsByScheduleIDRegister(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc đăng ký cần xóa không tồn tại")
	}
	return s.repo.Delete(c, id)
}

func (s *WorkScheduleService) UpdateWorkScheduleAuto(c context.Context, workSchedule *model.WorkSchedule, id int) error {
	// Kiểm tra lịch làm việc có tồn tại không
	exists, err := s.repo.IsExistsByScheduleIDAuto(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc cần sửa đổi không tồn tại")
	}

	timeNow := time.Now()
	effectiveBeforeNow := workSchedule.EffectiveDate.Before(timeNow)

	// Mặc định expirationAfterNow = true (nếu nil thì coi như chưa hết hạn)
	expirationAfterNow := true
	if workSchedule.ExpirationDate != nil {
		expirationAfterNow = workSchedule.ExpirationDate.After(timeNow)
	}

	// Xác định trạng thái lịch làm việc
	switch {
	case workSchedule.EffectiveDate.After(timeNow):
		workSchedule.Status = variable.InActive
	case effectiveBeforeNow && expirationAfterNow:
		workSchedule.Status = variable.Active
	case effectiveBeforeNow && !expirationAfterNow:
		workSchedule.Status = variable.Expired
	}

	// Cập nhật xuống DB
	return s.repo.Update(c, workSchedule, id)
}

func (s *WorkScheduleService) UpdateWorkScheduleRegister(c context.Context, workSchedule *model.WorkSchedule, id int) error {
	exists, err := s.repo.IsExistsByScheduleIDRegister(c, id)
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

func createWorkScheduleManagerRecords(managers []dto.ManagerAssignRequest, scheduleID int, now time.Time) []model.WorkScheduleManager {
	var records []model.WorkScheduleManager
	for _, m := range managers {
		records = append(records, model.WorkScheduleManager{
			EmployeeID:     m.ManagerID,
			WorkScheduleID: scheduleID,
			IsReading:      m.IsReading,
			IsEditing:      m.IsEditing,
		})
	}
	return records
}

func (s *WorkScheduleService) AssignEmployeeToWorkScheduleAuto(c context.Context, req *dto.AssignManagersRequest, idSchedule int) error {
	now := time.Now()

	managerRecords := createWorkScheduleManagerRecords(req.Managers, idSchedule, now)

	exists, err := s.repo.IsExistsByScheduleIDAuto(c, idSchedule)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc tự động không tồn tại")
	}

	allEmployeeIDs := make([]string, 0, len(managerRecords))
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

	if err := s.repo.AssignOrUpdateManagers(c, managerRecords); err != nil {
		return err
	}

	return nil
}

func (s *WorkScheduleService) AssignEmployeeToWorkScheduleRegister(c context.Context, req *dto.AssignManagersRequest, idSchedule int) error {
	now := time.Now()

	managerRecords := createWorkScheduleManagerRecords(req.Managers, idSchedule, now)

	exists, err := s.repo.IsExistsByScheduleIDRegister(c, idSchedule)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc đăng ký không tồn tại")
	}

	allEmployeeIDs := make([]string, 0, len(managerRecords))
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

	if err := s.repo.AssignOrUpdateManagers(c, managerRecords); err != nil {
		return err
	}

	return nil
}

func (s *WorkScheduleService) GetAllWorkScheduleRegister(ctx context.Context) ([]model.WorkSchedule, error) {
	return s.repo.GetAllScheduleRegister(ctx)
}

func (s *WorkScheduleService) GetAllWorkScheduleAuto(ctx context.Context) ([]model.WorkSchedule, error) {
	return s.repo.GetAllScheduleAuto(ctx)
}

func (s *WorkScheduleService) GetWorkScheduleByID(ctx context.Context, id int) (*model.WorkSchedule, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WorkScheduleService) validateEmployeeExists(employeeIDs []string) ([]string, error) {
	missing, err := s.employeeRepo.CheckExistEmployees(employeeIDs)
	if err != nil {
		return nil, err
	}
	return missing, nil
}

func (s *WorkScheduleService) DeleteManagerFromWorkScheduleAuto(ctx context.Context, managerID string, workScheduleID int) error {
	// check exist schedule
	exists, err := s.repo.IsExistsByScheduleIDAuto(ctx, workScheduleID)
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

func (s *WorkScheduleService) DeleteManagerFromWorkScheduleRegister(ctx context.Context, managerID string, workScheduleID int) error {
	// check exist schedule
	exists, err := s.repo.IsExistsByScheduleIDRegister(ctx, workScheduleID)
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

func (s *WorkScheduleService) UpdateStatusRecuringSchedule(ctx context.Context, scheduleId int, isAuto bool) error {
	return s.repo.UpdateStatusRecuringSchedule(ctx, scheduleId, isAuto)
}
