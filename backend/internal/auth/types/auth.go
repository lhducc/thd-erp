package types

type SignInRequest struct {
	Email    string `json:"email" binding:"required" validate:"email" example:"admin123@gmail.com"`
	Password string `json:"password" binding:"required" example:"toiyeuTHD123@"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"current_password" binding:"required" example:"old_password123"`
	NewPassword string `json:"new_password" binding:"required" example:"new_password123"`
}
