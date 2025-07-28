package service

import (
	"context"
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	"erp/backend/internal/hrm/checkin/service/service_interface"
)

type timesheetServiceImp struct {
	timesheetRepo repo_interface.TimeSheetRepoInterface
}

func NewTimesheetService(timesheetRepo repo_interface.TimeSheetRepoInterface) service_interface.TimeSheetServiceInterface {
	return &timesheetServiceImp{
		timesheetRepo: timesheetRepo,
	}
}

func (t timesheetServiceImp) FindByEmployeeIDAndMonth(ctx context.Context, employeeID string, month, year int) (*model.TimeSheet, error) {
	return t.timesheetRepo.FindByEmployeeAndMonth(ctx, employeeID, month, year)
}
