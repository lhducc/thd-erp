package checkin

import (
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ShiftAllocationHandler struct {
	service service_interface.ShiftAllocationServiceInterface
}

func NewShiftAllocationHandler(service service_interface.ShiftAllocationServiceInterface) *ShiftAllocationHandler {
	return &ShiftAllocationHandler{
		service: service,
	}
}

func (biz *ShiftAllocationHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			data  []dto.EmployeeScheduleResponse
			err   error
			param dto.GetShiftAllocationRequest
		)

		if err := c.ShouldBindQuery(&param); err != nil {
			utils.ResponseMessage(c, "Invalid query params", http.StatusBadRequest, nil)
			return
		}

		ctx := c.Request.Context()
		role := c.GetString("role")

		if role != "" && strings.EqualFold(role, "admin") {
			data, err = biz.service.GetListShiftAllocation(ctx, param, "")
			if err != nil {
				utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		} else if strings.EqualFold(role, "manager") {
			managerID := c.GetString("employeeId")
			data, err = biz.service.GetListShiftAllocation(ctx, param, managerID)
			if err != nil {
				utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}
		utils.ResponseMessage(c, "Lấy dữ liệu thành công", http.StatusOK, &data)

	}
}
