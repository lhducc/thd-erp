package model

import "time"

type Contract struct {
	ContractId     string    `gorm:"type:varchar(8);primaryKey;column:contract_id" json:"contract_id"`
	EffectiveDate  time.Time `gorm:"column:effective_date" json:"effective_date"`
	ExpiredDate    time.Time `gorm:"column:expired_date" json:"expired_date"`
	SignDate       time.Time `gorm:"column:sign_date" json:"sign_date"`
	Note           string    `gorm:"note" json:"note"`
	AttachedFile   string    `gorm:"attached_file" json:"attached_file"`
	Condition      string    `gorm:"type:varchar(30);column:condition" json:"condition"`
	CreatedDate    time.Time `gorm:"column:created_date" json:"created_date"`
	ContractTypeId string    `gorm:"type:varchar(6);column:contract_type_id" json:"contract_type"`
	EmployeeID     string    `gorm:"type:varchar(8);column:employee_id" json:"employee_id"`
	Employee       *Employee `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee_info"`
}

func (Contract) TableName() string { return "contract" }

type ContractCreate struct {
	ContractId     string    `gorm:"primaryKey;column:contract_id" json:"contract_id"`
	EffectiveDate  time.Time `gorm:"column:effective_date" json:"effective_date"`
	ExpiredDate    time.Time `gorm:"column:expired_date" json:"expired_date"`
	SignDate       time.Time `gorm:"column:sign_date" json:"sign_date"`
	Note           string    `gorm:"note" json:"note"`
	AttachedFile   string    `gorm:"attached_file" json:"attached_file"`
	Condition      string    `gorm:"condition" json:"condition"`
	CreatedDate    time.Time `gorm:"column:created_date" json:"created_date"`
	ContractTypeId string    `gorm:"column:contract_type_id" json:"contract_type"`
	Manager        string    `gorm:"column:employee_id" json:"employee_id"`
}

func (ContractCreate) TableName() string { return "contract" }
