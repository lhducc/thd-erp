package job

import "erp/backend/internal/hrm/checkin/model"

type TimesheetJob struct {
	Timesheet     *model.TimeSheet
	TimesheetList *model.TimeSheetList
	Service       interface { // Khai báo interface chứa hàm tính toán
		CalculateForEmployee(ts *model.TimeSheet, tsl *model.TimeSheetList) error
	}
}
