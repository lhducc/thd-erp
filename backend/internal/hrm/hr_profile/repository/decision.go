package repository

import (
	"context"
	model "erp/backend/internal/hrm/hr_profile/model"
	"errors"

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

func (r *decisionStore) GetAllDecision(ctx context.Context) ([]model.Decision, error) {
	var decisions []model.Decision

	if err := r.db.WithContext(ctx).
		Preload("Employees").
		Preload("DecisionType").
		Find(&decisions).Error; err != nil {
		return nil, err
	}

	return decisions, nil
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

func (s *decisionStore) InsertDecisionEmployee(ctx context.Context, decisionID, employeeID string) error {
	link := model.DecisionEmployee{
		DecisionID: decisionID,
		EmployeeID: employeeID,
	}
	return s.db.WithContext(ctx).Create(&link).Error
}
