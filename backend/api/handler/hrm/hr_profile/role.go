package handler

import (
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleHandler struct {
	roleUsecase *usecase.RoleUsecase
}

func NewRoleHandler(usecase *usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		roleUsecase: usecase,
	}
}

func (handler *RoleHandler) GetAllRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		results, err := handler.roleUsecase.FetchAllRole(ctx)
		if err != nil {
			utils.ResponseMessage(c, "Lỗi khi lấy danh sách phân quyền", http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "OK", http.StatusOK, results)
	}

}
