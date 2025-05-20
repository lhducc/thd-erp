package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	InsuranceRepository "erp/backend/internal/hrm/hr_profile/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"erp/backend/pkg"
)

// InsuranceUsecase defines methods for insurance business logic
type InsuranceUsecase interface {
	CreateInsurance(ctx context.Context, insurance *model.Insurance) (*model.Insurance, error)
	UpdateInsurance(ctx context.Context, insuranceID string, insurance *model.Insurance) (*model.Insurance, error)
	DeleteInsurance(ctx context.Context, insuranceID string) error
	GetAllInsurances(ctx context.Context) ([]model.Insurance, error)
	GetInsuranceByID(ctx context.Context, insuranceID string) (*model.Insurance, error)
}

// insuranceUsecase implements InsuranceUsecase interface
type insuranceUsecase struct {
	insuranceRepo InsuranceRepository.InsuranceRepository
}

// NewInsuranceUsecase returns a new instance of insuranceUsecase
func NewInsuranceUsecase(repo InsuranceRepository.InsuranceRepository) *insuranceUsecase {
	return &insuranceUsecase{insuranceRepo: repo}
}

// CreateInsurance validates and creates a new insurance record
func (u *insuranceUsecase) CreateInsurance(ctx context.Context, insurance *model.Insurance) (*model.Insurance, error) {
	if err := insurance.ValidateInsurance(); err != nil {
		return nil, err
	}

	code, err := u.GenerateInsuranceCode(ctx)
	if err != nil {
		return nil, err
	}
	insurance.ID = code
	insurance.CreatedDate = utils.GetCurrentDate()

	return u.insuranceRepo.CreateInsurance(ctx, insurance)
}

// UpdateInsurance validates and updates an existing insurance record
func (u *insuranceUsecase) UpdateInsurance(ctx context.Context, insuranceID string, insurance *model.Insurance) (*model.Insurance, error) {
	existing, err := u.insuranceRepo.GetInsuranceByID(ctx, insuranceID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("insurance not found")
	}

	if err := insurance.ValidateInsurance(); err != nil {
		return nil, err
	}

	return u.insuranceRepo.UpdateInsurance(ctx, insuranceID, insurance)
}

// DeleteInsurance deletes an insurance record if it exists
func (u *insuranceUsecase) DeleteInsurance(ctx context.Context, insuranceID string) error {
	existing, err := u.insuranceRepo.GetInsuranceByID(ctx, insuranceID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("insurance not found")
	}
	return u.insuranceRepo.DeleteInsurance(ctx, insuranceID)
}

// GetAllInsurances returns all insurance records
func (u *insuranceUsecase) GetAllInsurances(ctx context.Context) ([]model.Insurance, error) {
	return u.insuranceRepo.GetAllInsurance(ctx)
}

// GetInsuranceByID returns an insurance record by ID
func (u *insuranceUsecase) GetInsuranceByID(ctx context.Context, insuranceID string) (*model.Insurance, error) {
	return u.insuranceRepo.GetInsuranceByID(ctx, insuranceID)
}

// GenerateInsuranceCode generates the next insurance code with prefix "BH"
func (u *insuranceUsecase) GenerateInsuranceCode(ctx context.Context) (string, error) {
	var lastInsurance model.Insurance
	err := u.insuranceRepo.GetInsuranceLastByCode(ctx, &lastInsurance)

	// Start from "BH00001" if no record or invalid prefix
	if err != nil || !strings.HasPrefix(lastInsurance.ID, "BH") {
		return "BH00001", nil
	}

	numberStr := strings.TrimPrefix(lastInsurance.ID, "BH")
	number, err := strconv.Atoi(strings.TrimSpace(numberStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse insurance code number: %w", err)
	}

	nextNumber := number + 1
	if nextNumber > 99999 {
		return "", errors.New("maximum insurance code reached: BH99999")
	}

	return fmt.Sprintf("BH%05d", nextNumber), nil
}
