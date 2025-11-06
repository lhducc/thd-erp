package model

type ProcessStage struct {
	ProcessStageID int    `json:"process_stage_id" gorm:"column:process_stage_id; primaryKey; autoIncrement"`
	StageName      string `json:"stage_name" gorm:"column:stage_name; not null" validate:"required"`
	StageOrder     int    `json:"stage_order" gorm:"column:stage_order; uniqueIndex:idx_process_form_stage_order"`

	// FK
	ProcessFormID string `json:"process_form_id" gorm:"column:process_form_id; uniqueIndex:idx_process_form_stage_order"`

	// Relationships
	ProcessForm ProcessForm `json:"-" gorm:"foreignKey:ProcessFormID;references:ProcessFormID"`
}

func (ProcessStage) TableName() string {
	return "process_stage"
}
