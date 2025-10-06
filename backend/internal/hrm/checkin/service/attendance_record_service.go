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
	"github.com/xuri/excelize/v2"
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

func (s *attendanceRecordService) GetHistoryByDate(
	ctx context.Context,
	dateStr string,
) ([]dto.AttendanceRecordHistoryByDate, error) {
	date, err := time.Parse("2006-01-02", dateStr) // format yyyy-mm-dd
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected yyyy-mm-dd")
	}

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	date = date.In(loc)

	return s.repo.ListHistoryByDate(ctx, date)
}

func (s *attendanceRecordService) ExportAttendanceExcel(ctx context.Context, targetDate time.Time) ([]byte, error) {
	records, err := s.repo.ListHistoryByDate(ctx, targetDate)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("Không có dữ liệu chấm công cho ngày %s", targetDate.Format("2006-01-02"))
	}

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	f := excelize.NewFile()
	sheet := "Attendance"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// Header
	headers := []string{"Mã Nhân Viên", "Tên", "Phòng Ban", "Văn Phòng", "Thời gian chấm công", "Đúng giờ", "Trễ giờ"}
	for i, h := range headers {
		col := string(rune('A' + i))
		cell := col + "1"
		f.SetCellValue(sheet, cell, h)
	}

	// Style cho header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDDDDD"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, "A1", "G1", headerStyle)

	// Style chữ đỏ cho giờ đi trễ
	redStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "FF0000"},
	})

	// Biến đếm số lượng đúng giờ và trễ
	onTimeCount := 0
	lateCount := 0

	// Fill data
	for i, r := range records {
		row := i + 2
		f.SetCellValue(sheet, "A"+itoa(row), r.EmployeeID)
		f.SetCellValue(sheet, "B"+itoa(row), r.FullName)
		f.SetCellValue(sheet, "C"+itoa(row), r.DepartmentName)
		f.SetCellValue(sheet, "D"+itoa(row), r.OfficeName)
		timeCell := "E" + itoa(row)
		f.SetCellValue(sheet, timeCell, r.Timestamp.Format("2006-01-02 15:04:05"))

		startTime, err := time.ParseInLocation("15:04:05", r.StartTime, loc) // Format: HH:MM:SS
		if err != nil {
			return nil, fmt.Errorf("invalid start_time format in work shift data")
		}

		// Ghép start_time với ngày targetDate
		shiftStart := time.Date(
			targetDate.Year(),
			targetDate.Month(),
			targetDate.Day(),
			startTime.Hour(),
			startTime.Minute(),
			startTime.Second(),
			0,
			loc,
		)
		// So sánh giờ chấm công với giờ bắt đầu ca (+3 phút grace period)
		shiftStartWithGrace := shiftStart.Add(3 * time.Minute)

		// So sánh giờ chấm công với giờ bắt đầu ca
		onTime := r.Timestamp.Before(shiftStartWithGrace) || r.Timestamp.Equal(shiftStartWithGrace)
		if onTime {
			f.SetCellValue(sheet, "F"+itoa(row), "X")
			onTimeCount++
		} else {
			f.SetCellValue(sheet, "G"+itoa(row), "X")
			f.SetCellStyle(sheet, timeCell, timeCell, redStyle)
			lateCount++
		}
	}

	f.SetColWidth(sheet, "A", "G", 20)

	// Thêm dòng tổng
	totalRow := len(records) + 3
	f.SetCellValue(sheet, "E"+itoa(totalRow), "Tổng")
	f.SetCellValue(sheet, "F"+itoa(totalRow), onTimeCount)
	f.SetCellValue(sheet, "G"+itoa(totalRow), lateCount)

	// Xuất file ra []byte
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
