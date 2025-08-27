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
	DecisionTypeID string    `gorm:"type:varchar(6);column:decision_type_id" json:"decision_type_id"`

	Employees    []Employee    `gorm:"many2many:decision_employees;joinForeignKey:DecisionID;joinReferences:EmployeeID"`
	DecisionType *DecisionType `gorm:"foreignKey:DecisionTypeID;references:DecisionTypeID" json:"decision_type,omitempty"`
}

func (Decision) TableName() string {
	return "decision"
}

type DecisionEmployee struct {
	DecisionID string `gorm:"type:varchar(8);primaryKey;column:decision_id"`
	EmployeeID string `gorm:"type:varchar(8);primaryKey;column:employee_id"`
}

func (DecisionEmployee) TableName() string {
	return "decision_employees"
}

type DecisionCreate struct {
	DecisionID     string    `json:"decision_id" form:"decision_id"`
	DecisionName   string    `json:"decision_name" form:"decision_name" binding:"required"`
	EffectiveDate  time.Time `json:"effective_date" form:"effective_date" binding:"required"`
	SignDate       time.Time `json:"sign_date" form:"sign_date" binding:"required"`
	Content        string    `json:"content" form:"content" binding:"required"`
	Condition      string    `json:"condition" form:"condition"`
	AttachedFile   string    `json:"-" form:"-"`
	CreatedDate    time.Time `json:"created_date" form:"created_date"`
	EmployeeIDs    []string  `json:"employee_ids" form:"employee_ids"`
	DecisionTypeID string    `json:"decision_type_id" form:"decision_type_id" binding:"required"`
}

type EmployeeShort struct {
	EmployeeID string `json:"employee_id"`
	Fullname   string `json:"fullname"`
}

type DecisionResponse struct {
	DecisionID       string          `json:"decision_id"`
	DecisionName     string          `json:"decision_name"`
	Employees        []EmployeeShort `json:"employees"`
	DecisionTypeID   string          `json:"decision_type_id"`
	DecisionTypeName string          `json:"decision_type_name"`
	EffectiveDate    string          `json:"effective_date"` // YYYY-MM-DD
	SignDate         string          `json:"sign_date"`      // YYYY-MM-DD
	Condition        string          `json:"condition"`
	Content          string          `json:"content"`
	AttachedFile     string          `json:"attached_file"`
	CreatedDate      time.Time       `json:"created_date"`
}

func (DecisionCreate) TableName() string { return "decision" }
