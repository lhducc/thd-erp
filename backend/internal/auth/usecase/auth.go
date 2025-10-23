package authUsecase

import (
	"context"
	"erp/backend/internal/auth/repository"
	"erp/backend/internal/auth/token"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/pkg"
	"fmt"
	"time"
)

type AuthenticationInterface interface {
	SignIn(ctx context.Context, email, password string) (string, string, bool, error)
	ChangePassword(ctx context.Context, accountId, newPassword string) error
	GetEmployeeByID(ctx context.Context, id string) (hrmmodel.Employee, error)
}

type Authentication struct {
	accountRepo authRepository.AccountRepository
	tokenMaker  token.Maker
}

func NewAuthentication(accountRepo authRepository.AccountRepository, tokenMaker token.Maker) AuthenticationInterface {
	return &Authentication{
		accountRepo: accountRepo,
		tokenMaker:  tokenMaker,
	}
}

func (a *Authentication) SignIn(ctx context.Context, email, password string) (string, string, bool, error) {
	account, err := a.accountRepo.GetEmployeeByEmail(ctx, email)
	if err != nil {
		return "", "", false, fmt.Errorf("email or password is incorrect")
	}

	if err := utils.CheckPassword(password, account.Password); err != nil {
		return "", "", false, fmt.Errorf("password is incorrect")
	}

	accessToken, refreshToken, err := a.CreateToken(ctx, account)
	if err != nil {
		return "", "", false, err
	}

	return accessToken, refreshToken, account.FirstLogin, nil
}

func (a *Authentication) ChangePassword(ctx context.Context, employeeID, newPassword string) error {
	account, err := a.accountRepo.GetEmployeeByID(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if !account.FirstLogin {
		return fmt.Errorf("password change is only allowed on first login")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := a.accountRepo.UpdatePasswordAndFirstLogin(ctx, employeeID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

func (a *Authentication) CreateToken(ctx context.Context, emp hrmmodel.Employee) (string, string, error) {
	fullname := emp.Fullname
	roleID := emp.RoleID
	employeeID := emp.EmployeeID

	accessToken, _, err := a.tokenMaker.CreateToken(employeeID, fullname, roleID, 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := a.tokenMaker.CreateToken(employeeID, fullname, roleID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *Authentication) GetEmployeeByID(ctx context.Context, id string) (hrmmodel.Employee, error) {
	return a.accountRepo.GetEmployeeByID(ctx, id)
}
