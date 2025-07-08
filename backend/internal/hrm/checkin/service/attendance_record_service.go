package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/pkg"
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

func (s *attendanceRecordService) CheckCatrgoryExists(ctx context.Context,
	record *model.AttendanceRecord) (*model.AttendanceCategory, error) {
	category, exists, err := s.categoryRepo.GetIfExists(ctx, record.AttendanceCategoryID)
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

	// Kiểm tra dữ liệu GPS
	if record.Latitude == nil || record.Longitude == nil || office.Latitude == nil || office.Longitude == nil {
		return errors.New("missing GPS data for location comparison")
	}

	// Tính khoảng cách và kiểm tra phạm vi (scope tính bằng mét)
	withinScope := utils.IsWithinScope(
		*record.Latitude,
		*record.Longitude,
		*office.Latitude,
		*office.Longitude,
		category.Scope,
	)

	if !withinScope {
		// Tính khoảng cách thực tế để hiển thị thông báo lỗi
		distance := utils.HaversineDistance(
			*record.Latitude,
			*record.Longitude,
			*office.Latitude,
			*office.Longitude,
		)
		return fmt.Errorf("điểm danh cách văn phòng %.2f mét, vượt quá phạm vi cho phép (%d mét)",
			distance, category.Scope)
	}
	record.Status = model.Approved

	return nil
}

func (s *attendanceRecordService) CreateAttendanceRecord(ctx context.Context, record *model.AttendanceRecord) error {
	return s.repo.Create(ctx, record)
}

func (s *attendanceRecordService) UpdateAttendanceRecord(ctx context.Context, record *model.AttendanceRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, record)
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

func (s *attendanceRecordService) ListAttendanceRecordsByEmployee(ctx context.Context, employeeID string) ([]model.AttendanceRecord, error) {
	return s.repo.ListByEmployee(ctx, employeeID)
}

func (s *attendanceRecordService) ListAttendanceRecordsByDateRange(ctx context.Context, employeeID string, from, to string) ([]model.AttendanceRecord, error) {
	fromTime, err1 := ParseDateTime(from)
	toTime, err2 := ParseDateTime(to)

	if err1 != nil || err2 != nil {
		return nil, errors.New("invalid datetime format")
	}

	return s.repo.ListByDateRange(ctx, employeeID, fromTime, toTime)
}

func (s *attendanceRecordService) GetTotalReqOfEmployee(ctx context.Context, employeeID string) (int64, error) {
	return s.repo.CountRequestsByEmployee(ctx, employeeID)
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
