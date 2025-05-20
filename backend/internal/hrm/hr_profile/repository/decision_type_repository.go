package repository

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"gorm.io/gorm"
)

// DecisionTypeRepository interface defines methods for decision type-related data operations
type DecisionTypeRepository interface {
	CreateDecisionType(ctx context.Context, decisionType *model.DecisionType) (*model.DecisionType, error)
	UpdateDecisionType(ctx context.Context, decisionTypeID string, decisionType *model.DecisionType) (*model.DecisionType, error)
	DeleteDecisionType(ctx context.Context, decisionTypeID string) error
	GetAllDecisionTypes(ctx context.Context) ([]model.DecisionType, error)
	GetDecisionTypeByID(ctx context.Context, decisionTypeID string) (*model.DecisionType, error)
	GetLastDecisionTypeByCode(ctx context.Context, decisionType *model.DecisionType) error
}

// decisionTypeRepository struct implements the DecisionTypeRepository interface
type decisionTypeRepository struct {
	db *gorm.DB
}

// NewDecisionTypeRepository creates a new instance of decisionTypeRepository
func NewDecisionTypeRepository(db *gorm.DB) DecisionTypeRepository {
	return &decisionTypeRepository{db: db}
}

// CreateDecisionType creates a new decision type in the database
func (r *decisionTypeRepository) CreateDecisionType(ctx context.Context, decisionType *model.DecisionType) (*model.DecisionType, error) {
	if err := r.db.WithContext(ctx).Create(decisionType).Error; err != nil {
		return nil, err
	}
	return decisionType, nil
}

// UpdateDecisionType updates an existing decision type in the database
func (r *decisionTypeRepository) UpdateDecisionType(ctx context.Context, decisionTypeID string, decisionType *model.DecisionType) (*model.DecisionType, error) {
	if err := r.db.WithContext(ctx).Where("decision_type_id = ?", decisionTypeID).Updates(decisionType).Error; err != nil {
		return nil, err
	}
	return decisionType, nil
}

// DeleteDecisionType deletes a decision type by its ID
func (r *decisionTypeRepository) DeleteDecisionType(ctx context.Context, decisionTypeID string) error {
	if err := r.db.WithContext(ctx).Where("decision_type_id = ?", decisionTypeID).Delete(&model.DecisionType{}).Error; err != nil {
		return err
	}
	return nil
}

// GetAllDecisionTypes retrieves all decision types from the database
func (r *decisionTypeRepository) GetAllDecisionTypes(ctx context.Context) ([]model.DecisionType, error) {
	var decisionTypes []model.DecisionType
	if err := r.db.WithContext(ctx).Find(&decisionTypes).Error; err != nil {
		return nil, err
	}
	return decisionTypes, nil
}

// GetDecisionTypeByID retrieves a decision type by its ID
func (r *decisionTypeRepository) GetDecisionTypeByID(ctx context.Context, decisionTypeID string) (*model.DecisionType, error) {
	var decisionType model.DecisionType
	if err := r.db.WithContext(ctx).Where("decision_type_id = ?", decisionTypeID).First(&decisionType).Error; err != nil {
		return nil, err
	}
	return &decisionType, nil
}

func (r *decisionTypeRepository) GetLastDecisionTypeByCode(ctx context.Context, decisionType *model.DecisionType) error {
	return r.db.WithContext(ctx).
		Order("decision_type_id DESC").
		First(decisionType).Error
}
