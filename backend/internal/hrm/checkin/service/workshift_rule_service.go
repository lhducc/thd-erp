package service

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/repository/repo_interface"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"fmt"
	"gorm.io/gorm"
)

type WorkshiftRuleService struct {
	repo           repo_interface.WorkshiftRuleRepo
	employeeRepo   *repository.UserStore
	workshiftRepo  repo_interface.WorkShiftRepo
	departmentRepo *repository.DepartmentStore
	jobTitleRepo   *repository.JobTitleStore
	positionRepo   *repository.PositionStore
	officeRepo     *repository.OfficeStore
}

func NewWorkShiftRuleService(
	repo repo_interface.WorkshiftRuleRepo,
	employeeRepo *repository.UserStore,
	workshiftRepo repo_interface.WorkShiftRepo,
	departmentRepo *repository.DepartmentStore,
	jobTitleRepo *repository.JobTitleStore,
	positionRepo *repository.PositionStore,
	officeRepo *repository.OfficeStore,
) *WorkshiftRuleService {
	return &WorkshiftRuleService{
		repo:           repo,
		employeeRepo:   employeeRepo,
		workshiftRepo:  workshiftRepo,
		departmentRepo: departmentRepo,
		jobTitleRepo:   jobTitleRepo,
		positionRepo:   positionRepo,
		officeRepo:     officeRepo,
	}
}

func (w *WorkshiftRuleService) Create(input model.WorkshiftRule) error {
	// Validate Offices
	if len(input.Office) == 0 {
		return fmt.Errorf("at least one office is required")
	}
	var offices []hrmmodel.Office
	for _, office := range input.Office {
		dbOffice, err := w.officeRepo.FindByID(office.ID)
		if err != nil {
			return fmt.Errorf("office ID %s not found: %w", office.ID, err)
		}
		offices = append(offices, dbOffice)
	}

	// Validate Departments
	var departments []hrmmodel.Department
	for _, dept := range input.Department {
		dbDept, err := w.departmentRepo.FindByID(dept.ID)
		if err != nil {
			return fmt.Errorf("department ID %s not found: %w", dept.ID, err)
		}
		departments = append(departments, dbDept)
	}

	// Validate Positions
	var positions []hrmmodel.Position
	for _, pos := range input.Position {
		dbPos, err := w.positionRepo.FindByID(pos.ID)
		if err != nil {
			return fmt.Errorf("position ID %s not found: %w", pos.ID, err)
		}
		positions = append(positions, dbPos)
	}

	// Validate JobTitles
	var jobTitles []hrmmodel.JobTitle
	for _, jt := range input.JobTitle {
		dbJT, err := w.jobTitleRepo.FindByID(jt.JobTitleID)
		if err != nil {
			return fmt.Errorf("job title ID %s not found: %w", jt.JobTitleID, err)
		}
		jobTitles = append(jobTitles, dbJT)
	}

	// Validate WorkShifts
	var workshifts []model.WorkShifts
	for _, ws := range input.WorkShifts {
		dbWS, err := w.workshiftRepo.FindByID(ws.WorkShiftID)
		if err != nil {
			return fmt.Errorf("workshift ID %d not found: %w", ws.WorkShiftID, err)
		}
		workshifts = append(workshifts, dbWS)
	}

	db := w.repo.DB()

	newRule := model.WorkshiftRule{
		Name:      input.Name,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		//Description: input.Description,

		Office:     offices,
		Department: departments,
		Position:   positions,
		JobTitle:   jobTitles,
		WorkShifts: workshifts,
	}

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&newRule).Error
	if err != nil {
		return fmt.Errorf("failed to create workshift rule with associations: %w", err)
	}

	return nil
}

func (w *WorkshiftRuleService) GetAll() ([]model.WorkshiftRule, error) {
	return w.repo.GetAll()
}

func (w *WorkshiftRuleService) GetByUserID(userID string) ([]model.WorkshiftRule, error) {
	employee, err := w.employeeRepo.GetUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	rules, err := w.repo.GetByUserID(employee)
	if err != nil {
		return nil, fmt.Errorf("failed to get workshift rules: %w", err)
	}

	return rules, nil
}

func (w *WorkshiftRuleService) Update(id uint, assignment model.WorkshiftRule) error {
	return w.repo.Update(id, assignment)
}

func (w *WorkshiftRuleService) Delete(id uint) error {
	return w.repo.Delete(id)
}
