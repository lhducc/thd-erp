package export

import "erp/backend/internal/hrm/checkin/model"

type TimeSheetExport struct {
	*model.TimeSheet
	DayDetails          map[string]float64 // Key: "day_01", Value: work_days
	AdditionalShiftDays float64
	Deduction           float64
	// Thêm các trường từ Employee
	FullName   string
	Position   string
	WorkType   string
	Department string
	Office     string
}
