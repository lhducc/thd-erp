package officeusecase

import (
	"context"
	officemodel "erp/backend/internal/hrm/office/model"
	utils "erp/backend/pkg"
	"fmt"
	"strconv"
	"strings"
)

type OfficeRepo interface {
	CreateOffice(context context.Context, data *officemodel.OfficeCreate) error
	GetOffice(ctx context.Context, id string) (*officemodel.Office, error)
	GetAllOffice(ctx context.Context) ([]officemodel.Office, error)
	UpdateOffice(ctx context.Context, id string, data *officemodel.OfficeCreate) error
	DeleteOffice(ctx context.Context, id string) error
	CheckExistName(name string) (bool, error)
	GetLastOfficeByCode(ctx context.Context, office *officemodel.Office) error
}

type officeBiz struct {
	repo OfficeRepo
}

func NewOfficeBiz(store OfficeRepo) *officeBiz {
	return &officeBiz{repo: store}
}

func (biz *officeBiz) CreateOffice(context context.Context, data *officemodel.OfficeCreate) error {
	data.CreatedDate = utils.GetCurrentDate()
	check, err := biz.repo.CheckExistName(data.Name)
	if err != nil || check {
		return err
	}

	code, err := biz.GenerateOfficeCode(context)
	if err != nil {
		return fmt.Errorf("failed to generate office code: %w", err)
	}

	data.ID = code

	if err := biz.repo.CreateOffice(context, data); err != nil {
		return err
	}

	return nil
}

func (biz *officeBiz) GetOffice(ctx context.Context, id string) (*officemodel.Office, error) {
	office, err := biz.repo.GetOffice(ctx, id)
	if err != nil {
		return nil, err
	}
	return office, nil
}
func (biz *officeBiz) UpdateOffice(ctx context.Context, id string, data *officemodel.OfficeCreate) error {
	if err := biz.repo.UpdateOffice(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *officeBiz) DeleteOffice(ctx context.Context, id string) error {
	if err := biz.repo.DeleteOffice(ctx, id); err != nil {
		return err
	}
	return nil
}

func (biz *officeBiz) GetAllOffice(ctx context.Context) ([]officemodel.Office, error) {
	positions, err := biz.repo.GetAllOffice(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (b *officeBiz) GenerateOfficeCode(ctx context.Context) (string, error) {
	var lastOffice officemodel.Office
	err := b.repo.GetLastOfficeByCode(ctx, &lastOffice)
	fmt.Print(lastOffice.ID)
	if err != nil {
		if err.Error() == "record not found" {
			return "VP0001", nil
		}
		return "", fmt.Errorf("failed to get last office code: %w", err)
	}

	lastCode := lastOffice.ID
	fmt.Print(lastCode)
	if !strings.HasPrefix(lastCode, "VP") {
		return "", fmt.Errorf("invalid office code format: %s", lastCode)
	}

	numberStr := strings.TrimPrefix(lastCode, "VP")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse office code number: %w", err)
	}

	nextNumber := number + 1
	if nextNumber > 9999 {
		return "", fmt.Errorf("maximum office code reached: VP99999")
	}

	return fmt.Sprintf("VP%04d", nextNumber), nil
}
