package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"strings"
	"time"
)

type WorkScheduleService struct {
	repo repo_interface.WorkScheduleRepo
}

func NewWorkScheduleService(repo repo_interface.WorkScheduleRepo) service_interface.WorkScheduleServiceInterface {
	return &WorkScheduleService{
		repo: repo,
	}
}

func (s *WorkScheduleService) CreateNewWorkSchedule(c context.Context, workSchedule *model.WorkSchedule) error {
	var weekdayShift []model.WorkScheduleShift
	wScheduleId := workSchedule.WorkScheduleID

	exists, err := s.repo.IsExistsByName(c, workSchedule.WorkScheduleName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Tên lịch làm việc đã tồn tại")
	}
	for _, weekday := range workSchedule.Weekdays {
		weekdayShift = append(weekdayShift, model.WorkScheduleShift{
			Weekday:        weekday.Weekday,
			WorkScheduleID: wScheduleId,
			WorkShiftID:    weekday.WorkShiftID,
		})
	}
	return s.repo.Save(c, workSchedule, weekdayShift)
}

func (s *WorkScheduleService) DeleteWorkSchedule(c context.Context, id int) error {
	exists, err := s.repo.IsExistsByID(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc cần xóa không tồn tại")
	}
	return s.repo.Delete(c, id)
}
func (s *WorkScheduleService) UpdateWorkSchedule(c context.Context, workSchedule *model.WorkSchedule, id int) error {
	exists, err := s.repo.IsExistsByID(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc cần sửa đổi không tồn tại")
	}
	return s.repo.Update(c, workSchedule, id)
}

func (s *WorkScheduleService) AssignEmployeeToWorkSchedule(c context.Context, req *model.AssignEmployeeRequest, id int) error {
	var employeeRecords []model.WorkScheduleEmployee
	var managerRecords []model.WorkScheduleManager
	now := time.Now()
	for _, empIDCheck := range req.EmployeeIDs {
		exist, err := s.repo.CheckExistEmployeeSchedule(c, empIDCheck)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("Mã nhân viên `%s` đã tồn tại", empIDCheck)
		}
	}
	for _, empID := range req.EmployeeIDs {
		employeeRecords = append(employeeRecords, model.WorkScheduleEmployee{
			EmployeeID:     empID,
			WorkScheduleID: &id,
			AssignedAt:     &now,
		})
	}
	for _, managerID := range req.ManagerIDs {
		managerRecords = append(managerRecords, model.WorkScheduleManager{
			EmployeeID:     managerID,
			WorkScheduleID: id,
			AssignedAt:     &now,
		})
	}

	exists, err := s.repo.IsExistsByID(c, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Lịch làm việc không tồn tại")
	}

	err = s.repo.AssignEmployeesToWorkSchedule(c, employeeRecords, managerRecords)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkScheduleService) GetAllWorkSchedule(ctx context.Context) ([]model.WorkSchedule, error) {
	return s.repo.GetAll(ctx)
}

func (s *WorkScheduleService) GetWorkScheduleByID(ctx context.Context, id int) (*model.WorkSchedule, error) {
	return s.repo.GetByID(ctx, id)
}

func (e *WorkScheduleService) ExportWorkSchedule(c context.Context, selectedFields []string) ([]byte, string, error) {
	exporter := utils.NewExcelExporter("WorkSchedule")

	exporter.RegisterField("work_schedule_id", "Mã lịch làm việc", "work_schedule_id", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		return ws.WorkScheduleID
	})

	exporter.RegisterField("work_schedule_name", "Tên lịch làm việc", "work_schedule_name", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		return ws.WorkScheduleName
	})

	exporter.RegisterField("office_id", "Mã văn phòng", "office_id", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		return ws.OfficeID
	})

	exporter.RegisterField("office_name", "Tên văn phòng", "office_name", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		if ws.Office != nil {
			return ws.Office.Name
		}
		return ""
	})

	exporter.RegisterField("repeat_type", "Loại lặp lại", "repeat_type", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		switch ws.RepeatType {
		case variable.Weekly:
			return "Hàng tuần"
		case variable.Monthly:
			return "Hàng tháng"
		default:
			return "Không xác định"
		}
	})

	exporter.RegisterField("repeat_cycle", "Chu kỳ lặp", "repeat_cycle", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		return ws.RepeatCycle
	})

	exporter.RegisterField("effective_date", "Ngày hiệu lực", "effective_date", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		return ws.EffectiveDate.Format("2006-01-02")
	})

	exporter.RegisterField("expiration_date", "Ngày hết hạn", "expiration_date", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		return ws.ExpirationDate.Format("2006-01-02")
	})

	exporter.RegisterField("status", "Trạng thái", "status", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		switch ws.Status {
		case variable.Active:
			return "Đang hoạt động"
		case variable.InActive:
			return "Ngừng hoạt động"
		case variable.Expired:
			return "Hết hạn"
		default:
			return "Không xác định"
		}
	})

	exporter.RegisterField("created_at", "Ngày tạo", "created_at", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		return ws.CreatedAt.Format("2006-01-02 15:04:05")
	})

	exporter.RegisterField("weekdays", "Ngày làm việc", "weekdays", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		var days []string
		for _, day := range ws.Weekdays {
			days = append(days, fmt.Sprintf("%s (%s)", day.Weekday, day.WorkShiftID))
		}
		return strings.Join(days, ", ")
	})

	exporter.RegisterField("managers", "Quản lý", "managers", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		var managerNames []string
		for _, manager := range ws.Managers {
			managerNames = append(managerNames, manager.EmployeeID)
		}
		return strings.Join(managerNames, ", ")
	})

	exporter.RegisterField("employees", "Nhân viên", "employees", func(item interface{}) any {
		ws := item.(model.WorkSchedule)
		var employeeNames []string
		for _, emp := range ws.Employees {
			employeeNames = append(employeeNames, emp.EmployeeID)
		}
		return strings.Join(employeeNames, ", ")
	})

	result, err := e.repo.GetAll(c)
	if err != nil {
		return nil, "", fmt.Errorf("không thể lấy dữ liệu lịch làm việc: %w", err)
	}

	items := make([]interface{}, len(result))
	for i, v := range result {
		items[i] = v
	}

	return exporter.Export(items, selectedFields)
}
