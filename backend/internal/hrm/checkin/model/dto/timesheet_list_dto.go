package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"errors"
)

type AdjustWorkDayReq struct {
	AdjustedWorkDay float64 `json:"adjusted_work_day"`
}

type TimesheetDTO struct {
	Name  string `json:"name"`
	Month int    `json:"month"`
	Year  int    `json:"year"`
}

func (timesheet TimesheetDTO) Validate() error {
	if timesheet.Name == "" {
		return errors.New("Tên bảng công không được để trống")
	}
	if timesheet.Month < 1 || timesheet.Month > 12 {
		return errors.New("Tháng không hợp lệ (phải từ 1 đến 12)")
	}
	if timesheet.Year < 2000 {
		return errors.New("Năm không hợp lệ (phải >= 2000)")
	}
	return nil
}

func (ts *TimesheetDTO) ConvertToBusinessModel() *model.TimeSheetList {
	return &model.TimeSheetList{
		TimeSheetListName: ts.Name,
		Month:             ts.Month,
		Year:              ts.Year,
	}
}
