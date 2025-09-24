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
	SignIn(ctx context.Context, email, password string) (string, string, error)
	ChangePassword(ctx context.Context, accountId, newPassword string) error
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

func (a *Authentication) SignIn(ctx context.Context, email, password string) (string, string, error) {
	account, err := a.accountRepo.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("email or password is incorrect")
	}

	if err := utils.CheckPassword(password, account.Password); err != nil {
		return "", "", fmt.Errorf("password is incorrect")
	}

	accessToken, refreshToken, err := a.CreateToken(ctx, account)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *Authentication) ChangePassword(ctx context.Context, accountId, newPassword string) error {
	account, err := a.accountRepo.GetAccountByID(ctx, accountId)
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

	account.Password = hashedPassword
	account.FirstLogin = false

	if err := a.accountRepo.UpdateAccount(ctx, account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

func (a *Authentication) CreateToken(ctx context.Context, account hrmmodel.Account) (string, string, error) {
	var fullname string

	emp, err := a.accountRepo.GetEmployeeByAccountID(ctx, account.ID) // Lấy thông tin employee
	if err != nil {
		return "", "", fmt.Errorf("lỗi khi lấy thông tin nhân viên theo account id: %w", err)
	}
	fullname = emp.Fullname

	accessToken, _, err := a.tokenMaker.CreateToken(emp.EmployeeID, fullname, account.RoleID, 24*time.Hour) //1day
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := a.tokenMaker.CreateToken(emp.EmployeeID, fullname, account.RoleID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
