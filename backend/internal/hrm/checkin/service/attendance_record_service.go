package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/pkg"
	"erp/backend/pkg/variable"
	"errors"
	"fmt"
	"time"
)

type attendanceRecordService struct {
	repo         repo_interface.AttendanceRecordRepository
	categoryRepo repo_interface.AttendanceCategoryRepository
	officeRepo   usecase.OfficeRepo
}

func NewAttendanceRecordService(repo repo_interface.AttendanceRecordRepository, categoryRepo repo_interface.AttendanceCategoryRepository, officerepo usecase.OfficeRepo) service_interface.AttendanceRecordService {
	return &attendanceRecordService{repo: repo,
		categoryRepo: categoryRepo,
		officeRepo:   officerepo,
	}
}

func (s *attendanceRecordService) CheckCategoryExists(ctx context.Context,
	record *model.AttendanceRecord) (*model.AttendanceCategory, error) {
	category, exists, err := s.categoryRepo.GetIfExists(ctx, *record.AttendanceCategoryID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("category does not exist")
	}
	return category, nil
}

func (s *attendanceRecordService) ValidateAttendanceRecordDistance(ctx context.Context,
	record *model.AttendanceRecord, category *model.AttendanceCategory) error {
	// Parse GPS
	if err := record.ParseGPS(); err != nil {
		return err
	}

	if record.OfficeID == nil || *record.OfficeID == "" {
		return errors.New("office ID is required for auto-approval")
	}

	office, err := s.officeRepo.GetOffice(ctx, record.OfficeID)
	if err != nil {
		return err
	}

	if record.Latitude == nil || record.Longitude == nil || office.Latitude == nil || office.Longitude == nil {
		return errors.New("missing GPS data for location comparison")
	}

	withinScope := utils.IsWithinScope(
		*record.Latitude,
		*record.Longitude,
		*office.Latitude,
		*office.Longitude,
		category.Scope,
	)

	if !withinScope {
		distance := utils.HaversineDistance(
			*record.Latitude,
			*record.Longitude,
			*office.Latitude,
			*office.Longitude,
		)
		return fmt.Errorf("điểm danh cách văn phòng %.2f mét, vượt quá phạm vi cho phép (%d mét)",
			distance, category.Scope)
	}
	record.Status = variable.Approved

	return nil
}

func (s *attendanceRecordService) CreateAttendanceRecord(ctx context.Context, record *model.AttendanceRecord) error {
	return s.repo.Create(ctx, record)
}

func (s *attendanceRecordService) UpdateAttendanceRecord(ctx context.Context, updaterecord *dto.AttendanceRecordUpdate, recordID string) error {
	record, err := s.GetAttendanceRecordByID(ctx, recordID)
	if err != nil {
		return fmt.Errorf("Attendance record not found")
	}
	if record == nil {
		return fmt.Errorf("Attendance record does not exist")
	}
	return s.repo.Update(ctx, updaterecord, recordID)
}

func (s *attendanceRecordService) DeleteAttendanceRecord(ctx context.Context, id string) error {
	// Optional: check if exists first
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errors.New("attendance record not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *attendanceRecordService) GetAttendanceRecordByID(ctx context.Context, id string) (*model.AttendanceRecord, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *attendanceRecordService) GetAttendanceRecordByIDPersonal(ctx context.Context, recordID, employeeID string) (*model.AttendanceRecord, error) {
	return s.repo.GetByIDPersonal(ctx, recordID, employeeID)
}

func (s *attendanceRecordService) ListAttendanceRecordsByEmployee(ctx context.Context, employeeID string) ([]model.AttendanceRecord, error) {
	return s.repo.ListByEmployee(ctx, employeeID)
}

func (s *attendanceRecordService) ListAttendanceRequestsByDateRange(ctx context.Context, employeeID string, from, to string) ([]model.AttendanceRecord, error) {
	fromTime, err1 := ParseDateTime(from)
	toTime, err2 := ParseDateTime(to)

	if err1 != nil || err2 != nil {
		return nil, errors.New("invalid datetime format")
	}

	return s.repo.ListRequestByDateRange(ctx, employeeID, fromTime, toTime)
}

func (s *attendanceRecordService) GetTotalReqOfEmployee(ctx context.Context, employeeID string) (int64, error) {
	return s.repo.CountRequestsByEmployee(ctx, employeeID)
}

func (s *attendanceRecordService) GetHistoryRecordByEmployee(ctx context.Context, employeeID string,
	page int, limit int) ([]model.AttendanceRecord, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	return s.repo.ListHistoryRecordEmployee(ctx, employeeID, limit, offset)
}

func ParseDateTime(value string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime format")
}
