package model

import "time"

type Decision struct {
	DecisionID     string    `gorm:"type:varchar(8);column:decision_id;primaryKey" json:"decision_id"`
	DecisionName   string    `gorm:"type:varchar(200);column:decision_name" json:"decision_name"`
	EffectiveDate  time.Time `gorm:"column:effective_date" json:"effective_date"`
	SignDate       time.Time `gorm:"column:sign_date" json:"sign_date"`
	Content        string    `gorm:"type:text;column:content" json:"content"`
	Condition      string    `gorm:"type:varchar(20);column:condition" json:"condition"`
	AttachedFile   string    `gorm:"column:attached_file" json:"attached_file"`
	CreatedDate    time.Time `gorm:"column:created_date" json:"created_date"`
	EmployeeID     string    `gorm:"type:varchar(8);column:employee_id" json:"employee_id"`
	DecisionTypeID string    `gorm:"type:varchar(6);column:decision_type_id" json:"decision_type_id"`

	// Nếu bạn có struct Employee và DecisionType, bạn có thể dùng các khóa ngoại:
	Employee *Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee,omitempty"`
	// DecisionType   DecisionType `gorm:"foreignKey:DecisionTypeID"`
}

func (Decision) TableName() string {
	return "decision"
}

type DecisionCreate struct {
	DecisionID     string    `gorm:"column:decision_id;primaryKey" json:"decision_id"`
	DecisionName   string    `gorm:"column:decision_name" json:"decision_name"`
	EffectiveDate  time.Time `gorm:"column:effective_date" json:"effective_date"`
	SignDate       time.Time `gorm:"column:sign_date" json:"sign_date"`
	Content        string    `gorm:"column:content" json:"content"`
	Condition      string    `gorm:"column:condition" json:"condition"`
	AttachedFile   string    `gorm:"column:attached_file" json:"attached_file"`
	CreatedDate    time.Time `gorm:"column:created_date" json:"created_date"`
	EmployeeID     string    `gorm:"column:employee_id" json:"employee_id"`
	DecisionTypeID string    `gorm:"column:decision_type_id" json:"decision_type_id"`
}

func (DecisionCreate) TableName() string { return "decision" }
