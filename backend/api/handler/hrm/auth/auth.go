package auth

import (
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

	accessToken, refreshToken, firstLogin, err := r.authUsecase.SignIn(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ResponseMessage(ctx, err.Error(), http.StatusBadRequest, nil)
		return
	}

	ctx.SetCookie("access_token", accessToken, 24*60*60, "/", "", false, true)
	ctx.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", false, true)

	utils.ResponseMessage(ctx, "Đăng nhập thành công", http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"first_login":   firstLogin,
	})
}

func (r *AuthenticationHandler) ChangePassword(ctx *gin.Context) {
	var req types.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseMessage(ctx, err.Error(), http.StatusBadRequest, nil)
		return
	}

	employeeIDValue, exists := ctx.Get("employeeId")
	if !exists {
		utils.ResponseMessage(ctx, "Unauthorized", http.StatusUnauthorized, nil)
		return
	}

	accountID, ok := employeeIDValue.(string)
	if !ok {
		utils.ResponseMessage(ctx, "Invalid account ID", http.StatusUnauthorized, nil)
		return
	}

	// Lấy thông tin employee hiện tại
	employee, err := r.authUsecase.GetEmployeeByID(ctx.Request.Context(), accountID)
	if err != nil {
		utils.ResponseMessage(ctx, "Employee not found", http.StatusNotFound, nil)
		return
	}

	// Kiểm tra mật khẩu cũ
	if err := utils.CheckPassword(req.OldPassword, employee.Password); err != nil {
		utils.ResponseMessage(ctx, "Current password is incorrect", http.StatusBadRequest, nil)
		return
	}

	// Cập nhật mật khẩu mới
	err = r.authUsecase.ChangePassword(ctx.Request.Context(), accountID, req.NewPassword)
	if err != nil {
		utils.ResponseMessage(ctx, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	
	utils.ResponseMessage(ctx, "Change password successfully", http.StatusOK, nil)
}
