package checkin

import (
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	utils "erp/backend/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TimesheetHandler struct {
	tsService service_interface.TimeSheetServiceInterface
}

func NewTimesheetHandler(tsService service_interface.TimeSheetServiceInterface) *TimesheetHandler {
	return &TimesheetHandler{tsService: tsService}
}

func (h *TimesheetHandler) GetPersonalTimesheet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		employeeID := c.GetString("employeeId")
		monthParam := c.Query("month")
		yearParam := c.Query("year")
		month, err := strconv.Atoi(monthParam)
		if err != nil {
			fmt.Printf("err: %w", err.Error())
			utils.ResponseMessage(c, "Tháng không hợp lệ", http.StatusBadRequest, nil)
			return
		}
		year, err := strconv.Atoi(yearParam)
		if err != nil {
			fmt.Printf("err: %w", err.Error())
			utils.ResponseMessage(c, "Năm không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		personalTS, err := h.tsService.FindByEmployeeIDAndMonth(ctx, employeeID, month, year)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, &personalTS)
			return
		}
		utils.ResponseMessage(c, "Lấy bảng công cá nhân thành công", http.StatusOK, &personalTS)
	}
}

func (h *TimesheetHandler) CalculationTimeSheet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		timeSheetID := c.Param("id")
		err := h.tsService.CalculatorTimeSheetList(ctx, timeSheetID)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Tính toán thành công", http.StatusOK, nil)
	}
}

func (h *TimesheetHandler) AdjustWorkDayManual() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		employeeID := c.GetString("employeeId")

		detailIDstring := c.Param("id")
		detailID, err := strconv.Atoi(detailIDstring)
		if err != nil {
			fmt.Printf("err: %w", err.Error())
			utils.ResponseMessage(c, "id không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		var req dto.AdjustWorkDayReq
		if err := c.ShouldBind(&req); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}
		if err := h.tsService.ManualAdjustWorkDay(
			ctx,
			detailID,
			req.AdjustedWorkDay,
			employeeID,
		); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Thay đổi thành công", http.StatusOK, nil)
	}
}

func (h *TimesheetHandler) ResetWorkDayAdjustment() gin.HandlerFunc {
	return func(c *gin.Context) {
		detailID := c.Param("id")
		ctx := c.Request.Context()
		id, err := strconv.Atoi(detailID)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		err = h.tsService.ResetWorkDayAdjustment(ctx, id)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Reset thành công", http.StatusOK, nil)
	}
}
