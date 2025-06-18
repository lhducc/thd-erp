package model

import "time"

type Allowance struct {
	ID            string    `gorm:"column:id;primaryKey" json:"id"`
	AllowanceName string    `gorm:"column:allowance_name" json:"allowance_name"`
	Tax           bool      `gorm:"column:tax" json:"tax"`
	Amount        float64   `gorm:"column:amount" json:"amount"`
	Unit          string    `gorm:"column:unit" json:"unit"`
	IsDeleted     bool      `gorm:"column:is_deleted" json:"is_deleted"`
	CreatedDate   time.Time `gorm:"column:created_date" json:"created_date"`
}

func (Allowance) TableName() string {
	return "allowance"
}

type AllowanceCreate struct {
	AllowanceName string  `json:"allowance_name"`
	Amount        float64 `json:"amount"`
	Unit          string  `json:"unit"`
	Tax           bool    `json:"tax"`
}
