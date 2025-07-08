package repo_interface

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
)

type AttendanceCategoryRepository interface {
	Create(ctx context.Context, ac *model.AttendanceCategory) error
	Update(ctx context.Context, ac *model.AttendanceCategory) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.AttendanceCategory, error)
	ListByOffice(ctx context.Context, officeID string) ([]model.AttendanceCategory, error)
	GetList(ctx context.Context) ([]model.AttendanceCategory, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByNameAndNotID(ctx context.Context, name string, excludeID string) (bool, error)
	GetIfExists(ctx context.Context, id string) (*model.AttendanceCategory, bool, error)
}
