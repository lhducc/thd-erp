package auth

import (
	"context"
	"erp/backend/internal/auth/types"
	authusecase "erp/backend/internal/auth/usecase"
	"erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthenticationHandler struct {
	authUsecase authusecase.AuthenticationInterface
}

func NewAuthenticationHandler(authUsecase authusecase.AuthenticationInterface) *AuthenticationHandler {
	return &AuthenticationHandler{authUsecase: authUsecase}
}

func (r *AuthenticationHandler) SignIn(ctx *gin.Context) {
	var req types.SignInRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseMessage(ctx, err.Error(), http.StatusBadRequest, nil)
		return
	}

	accessToken, refreshToken, err := r.authUsecase.SignIn(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ResponseMessage(ctx, err.Error(), http.StatusBadRequest, nil)
		return
	}

	ctx.SetCookie("access_token", accessToken, 24*60*60, "/", "", false, true)
	ctx.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", false, true)

	utils.ResponseMessage(ctx, "Đăng nhập thành công", http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (r *AuthenticationHandler) ChangePassword(ctx *gin.Context) {
	var req types.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseMessage(ctx, err.Error(), http.StatusBadRequest, nil)
		return
	}

	accountIDValue, exists := ctx.Get("accountId")
	if !exists {
		utils.ResponseMessage(ctx, "Unauthorized", http.StatusUnauthorized, nil)
		return
	}

	accountID, ok := accountIDValue.(string)
	if !ok {
		utils.ResponseMessage(ctx, "Invalid account ID", http.StatusUnauthorized, nil)
		return
	}

	err := r.authUsecase.ChangePassword(context.Background(), accountID, req.NewPassword)
	if err != nil {
		utils.ResponseMessage(ctx, err.Error(), http.StatusBadRequest, nil)
		return
	}

	utils.ResponseMessage(ctx, "Change password successfully", http.StatusOK, nil)
}
