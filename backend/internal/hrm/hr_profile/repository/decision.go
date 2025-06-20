package repository

import (
	"context"
	model "erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type decisionStore struct {
	db *gorm.DB
}

func NewDicisionStore(db *gorm.DB) *decisionStore {
	return &decisionStore{db: db}
}

func (s *decisionStore) CreateDecision(ctx context.Context, data *model.Decision, employeeIDs []string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		for _, empID := range employeeIDs {
			link := model.DecisionEmployee{
				DecisionID: data.DecisionID,
				EmployeeID: empID,
			}
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *decisionStore) GetDecision(ctx context.Context, id string) (*model.Decision, error) {
	var decision model.Decision
	if err := r.db.WithContext(ctx).
		Preload("Employees").
		Preload("DecisionType").
		Where("decision_id = ?", id).
		First(&decision).Error; err != nil {
		return nil, err
	}
	return &decision, nil
}

func (r *decisionStore) GetAllDecision(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Decision, int64, error) {
	var (
		decisions    []model.Decision
		totalRecords int64
	)

	db := r.db.WithContext(ctx).Model(&model.Decision{})

	for key, value := range filters {
		if arr, ok := value.([]string); ok && len(arr) > 0 {
			db = db.Where(fmt.Sprintf("decision.%s IN (?)", key), arr)
		} else if value != nil && value != "" {
			db = db.Where(fmt.Sprintf("decision.%s = ?", key), value)
		}
	}

	// Count total records after filter
	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := db.
		Preload("Employees").
		Preload("DecisionType").
		Offset(offset).
		Limit(pageSize).
		Find(&decisions).Error; err != nil {
		return nil, 0, err
	}

	return decisions, totalRecords, nil
}

func (r *decisionStore) UpdateDecision(ctx context.Context, id string, data *model.DecisionCreate) error {
	return r.db.WithContext(ctx).Table("decision").
		Where("decision_id = ?", id).
		Updates(data).Error
}

func (r *decisionStore) DeleteDecision(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("decision").
		Where("decision_id = ?", id).
		Delete(nil).Error
}

func (s *decisionStore) CheckExistName(name string) (bool, error) {
	var count int64
	if err := s.db.Model(&model.Decision{}).Where("decision_name = ?", name).Count(&count).Error; err != nil {
		return true, err
	}
	if count > 0 {
		return true, errors.New("Decision Name is already")
	}
	return false, nil
}

func (s *decisionStore) GetLastDecisionByCode(ctx context.Context, decision *model.Decision) error {
	return s.db.WithContext(ctx).
		Order("decision_id DESC").
		First(decision).Error
}
func (s *decisionStore) GetAllDecisionNoPagination(ctx context.Context) ([]model.Decision, error) {
	var decisions []model.Decision

	if err := s.db.
		Preload("Employees").
		Preload("DecisionType").
		Find(&decisions).Error; err != nil {
		return nil, err
	}
	return decisions, nil
}

func (s *decisionStore) InsertDecisionEmployee(ctx context.Context, decisionID, employeeID string) error {
	link := model.DecisionEmployee{
		DecisionID: decisionID,
		EmployeeID: employeeID,
	}
	return s.db.WithContext(ctx).Create(&link).Error
}
