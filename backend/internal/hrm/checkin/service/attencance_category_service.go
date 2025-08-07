package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/internal/hrm/hr_profile/repository"
	"errors"
)

type attendanceCategoryService struct {
	repo         repo_interface.AttendanceCategoryRepository
	employeeRepo *repository.UserStore
}

func NewAttendanceCategoryService(repo repo_interface.AttendanceCategoryRepository, employeeRepo *repository.UserStore,
) service_interface.AttendanceCategoryService {
	return &attendanceCategoryService{repo: repo,
		employeeRepo: employeeRepo}
}

func (s *attendanceCategoryService) CreateAttendanceCategory(ctx context.Context, category *model.AttendanceCategory) error {
	// Check duplicate name
	if exists, err := s.repo.ExistsByName(ctx, category.AttendanceCategoryName); err != nil {
		return err
	} else if exists {
		return errors.New("tên loại chấm công đã tồn tại")
	}
	category.AutoApprove = category.IsCheckLocation && category.IsGPS
	return s.repo.Create(ctx, category)
}

func (s *attendanceCategoryService) UpdateAttendanceCategory(ctx context.Context, category *model.AttendanceCategory) error {
	// Check record exists
	if _, err := s.repo.GetByID(ctx, category.AttendanceCategoryID); err != nil {
		return errors.New("không tìm thấy loại chấm công")
	}

	// Check duplicate name (exclude current record)
	if exists, err := s.repo.ExistsByNameAndNotID(ctx, category.AttendanceCategoryName, category.AttendanceCategoryID); err != nil {
		return err
	} else if exists {
		return errors.New("tên loại chấm công đã tồn tại")
	}

	return s.repo.Update(ctx, category)
}

func (s *attendanceCategoryService) DeleteAttendanceCategory(ctx context.Context, id string) error {
	// Check record exists
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errors.New("Not found attendance-category need to delete")
	}

	// if inUse, err := s.repo.IsInUse(ctx, id); err != nil {
	//     return err
	// } else if inUse {
	//     return errors.New("không thể xóa loại chấm công đang được sử dụng")
	// }

	return s.repo.Delete(ctx, id)
}

func (s *attendanceCategoryService) GetAttendanceCategoryByID(ctx context.Context, id string) (*model.AttendanceCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *attendanceCategoryService) ListAttendanceCategoriesByOffice(ctx context.Context, employeeID string) ([]model.AttendanceCategory, error) {
	employee, err := s.employeeRepo.GetUserById(employeeID)
	if err != nil {
		return nil, err
	}
	if employee.Department.OfficeID == "" {
		return nil, nil
	}
	return s.repo.ListByOffice(ctx, employee.Department.OfficeID)
}

func (s *attendanceCategoryService) ListAllAttendanceCategories(ctx context.Context) ([]model.AttendanceCategory, error) {
	return s.repo.GetList(ctx)
}
