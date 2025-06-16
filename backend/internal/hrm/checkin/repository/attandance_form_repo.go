package repository

import (
	"erp/backend/internal/hrm/checkin/model"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type AttendanceFormStore struct {
	db *gorm.DB
}

func NewAttendanceFormStore(db *gorm.DB) *AttendanceFormStore {
	return &AttendanceFormStore{db: db}
}

func (r *AttendanceFormStore) recoverFromPanic(tx *gorm.DB) {
	if rec := recover(); rec != nil {
		tx.Rollback()
		log.Printf("Đã xảy ra lỗi, giao dịch bị hủy: %v", rec)
	}
}

func (r *AttendanceFormStore) CreateAttendanceForm(attendanceForm *model.AttendanceForm) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return errors.New("không thể bắt đầu giao dịch")
	}
	defer r.recoverFromPanic(tx)

	if err := tx.Create(attendanceForm).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi tạo bản ghi chấm công: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi commit giao dịch: %w", err)
	}
	return nil
}

func (s *AttendanceFormStore) GetLastAttandanceFormByCode(emp *model.AttendanceForm) error {
	return s.db.
		Order("attendance_form_id DESC").
		First(emp).Error
}

func (r *AttendanceFormStore) GetAttendanceFormByID(id string) (model.AttendanceForm, error) {
	var attendanceForm model.AttendanceForm
	if err := r.db.Where("attendance_form_id = ?", id).First(&attendanceForm).Error; err != nil {
		return model.AttendanceForm{}, fmt.Errorf("không tìm thấy bản ghi chấm công với ID %s", id)
	}
	return attendanceForm, nil
}

// Lấy tất cả bản ghi chấm công
func (r *AttendanceFormStore) GetAllAttendanceForm() ([]model.AttendanceForm, error) {
	var Attendance []model.AttendanceForm
	if err := r.db.Find(&Attendance).Error; err != nil {
		return nil, fmt.Errorf("lỗi khi lấy danh sách chấm công: %w", err)
	}
	return Attendance, nil
}

func (r *AttendanceFormStore) UpdateAttendanceForm(id string, AttendanceRecord model.AttendanceForm) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return errors.New("không thể bắt đầu giao dịch")
	}
	defer r.recoverFromPanic(tx)

	existing, err := r.GetAttendanceFormByID(id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("không tìm thấy bản ghi chấm công với ID %s", id)
	}

	model.UpdateAttendanceFormFields(&existing, AttendanceRecord)

	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi cập nhật bản ghi chấm công: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi commit giao dịch: %w", err)
	}
	return nil
}

func (r *AttendanceFormStore) DisableAttandanceForm(id string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return errors.New("không thể bắt đầu giao dịch")
	}
	defer r.recoverFromPanic(tx)

	var attandanceForm model.AttendanceForm
	if err := tx.Where("attendance_form_id = ?", id).First(&attandanceForm).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("không tìm thấy bản ghi chấm công với ID %s", id)
		}
		return fmt.Errorf("lỗi khi kiểm tra bản ghi chấm công: %w", err)
	}

	if attandanceForm.Status == "Không áp dụng" {
		tx.Rollback()
		return fmt.Errorf("bảng công này đã bị vô hiệu từ trước")
	}

	if err := tx.Model(&attandanceForm).Update("status", "Không áp dụng").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi cập nhật trạng thái bản ghi chấm công: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("lỗi khi commit giao dịch: %w", err)
	}
	return nil
}