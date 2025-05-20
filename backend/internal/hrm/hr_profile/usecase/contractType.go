package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	hrmrepository "erp/backend/internal/hrm/hr_profile/repository"
	"errors"
	"fmt"
	"strconv"

	"erp/backend/pkg"
	"gorm.io/gorm"
)

type ContractType interface {
	CreateContractType(ctx context.Context, contractType *model.ContractType) (*model.ContractType, error)
	UpdateContractType(ctx context.Context, id string, contractType *model.ContractType) (*model.ContractType, error)
	DeleteContractType(ctx context.Context, id string) (*model.ContractType, error)
	GetAllContractTypes(ctx context.Context) ([]model.ContractType, error)
	GetContractTypeById(ctx context.Context, id string) (*model.ContractType, error)
}

type ContractTypeUsecase struct {
	contractTypeRepo hrmrepository.ContractTypeRepository
}

func NewContractTypeUsecase(repo hrmrepository.ContractTypeRepository) *ContractTypeUsecase {
	return &ContractTypeUsecase{contractTypeRepo: repo}
}

func (u *ContractTypeUsecase) CreateContractType(ctx context.Context, contractType *model.ContractType) (*model.ContractType, error) {
	code, err := u.GenerateContractTypeCode(ctx)
	if err != nil {
		return nil, err
	}

	contractType.ContractTypeID = code
	contractType.CreatedDate = utils.GetCurrentDate()
	contractType.IsDelete = false

	if err := contractType.ValidateContractType(); err != nil {
		return nil, err
	}

	return u.contractTypeRepo.CreateContractType(ctx, contractType)
}

func (u *ContractTypeUsecase) UpdateContractType(ctx context.Context, id string, contractType *model.ContractType) (*model.ContractType, error) {
	if err := u.contractTypeRepo.UpdateContractType(ctx, id, contractType); err != nil {
		return nil, err
	}
	return contractType, nil
}

func (u *ContractTypeUsecase) DeleteContractType(ctx context.Context, id string) (*model.ContractType, error) {
	existing, err := u.contractTypeRepo.GetContractTypeById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contract type with id %s not found", id)
		}
		return nil, err
	}

	existing.IsDelete = true
	if err := u.contractTypeRepo.UpdateContractType(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *ContractTypeUsecase) GetAllContractTypes(ctx context.Context) ([]model.ContractType, error) {
	return u.contractTypeRepo.GetContractType(ctx)
}

func (u *ContractTypeUsecase) GetContractTypeById(ctx context.Context, id string) (*model.ContractType, error) {
	result, err := u.contractTypeRepo.GetContractTypeById(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("contract type not found")
	}
	return result, nil
}

func (u *ContractTypeUsecase) GenerateContractTypeCode(ctx context.Context) (string, error) {
	var lastContractType model.ContractType
	if err := u.contractTypeRepo.GetLastContractTypeByCode(ctx, &lastContractType); err != nil || lastContractType.ContractTypeID == "" {
		return "LH000001", nil
	}

	numStr := lastContractType.ContractTypeID[2:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse contract type code number: %w", err)
	}

	next := num + 1
	if next > 999999 {
		return "", fmt.Errorf("maximum contract type code reached: LH999999")
	}

	return fmt.Sprintf("LH%06d", next), nil
}
