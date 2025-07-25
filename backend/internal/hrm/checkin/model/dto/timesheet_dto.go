package dto

import (
	"erp/backend/internal/hrm/checkin/model"
	"errors"
)

type TimesheetDTO struct {
	Name        string `json:"name"`
	OfficeID    string `json:"office_id"`
	Month       int    `json:"month"`
	Year        int    `json:"year"`
	IsFinalized bool   `json:"is_finalized"`
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

func (ts *TimesheetDTO) ConvertToBusinessModel() *model.Timesheet {
	return &model.Timesheet{
		Name:        ts.Name,
		OfficeID:    ts.OfficeID,
		Month:       ts.Month,
		Year:        ts.Year,
		IsFinalized: ts.IsFinalized,
	}
}
