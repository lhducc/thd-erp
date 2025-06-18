package model

import (
	"errors"
	"time"
)

type AllowanceUnit string

const (
	UnitVND AllowanceUnit = "VND"
	UnitUSD AllowanceUnit = "USD"
	UnitEUR AllowanceUnit = "EUR"
	UnitJPY AllowanceUnit = "JPY"
	UnitGBP AllowanceUnit = "GBP"
)

type Allowance struct {
	ID            string        `gorm:"column:id;primaryKey" json:"id"`
	AllowanceName string        `gorm:"column:allowance_name" json:"allowance_name"`
	Tax           bool          `gorm:"column:tax" json:"tax"`
	Amount        float64       `gorm:"column:amount" json:"amount"`
	Unit          AllowanceUnit `gorm:"column:unit" json:"unit"`
	IsDeleted     bool          `gorm:"column:is_deleted" json:"is_deleted"`
	CreatedDate   time.Time     `gorm:"column:created_date" json:"created_date"`
}

func (Allowance) TableName() string {
	return "allowance"
}

type AllowanceCreate struct {
	AllowanceName string        `json:"allowance_name"`
	Amount        float64       `json:"amount"`
	Unit          AllowanceUnit `json:"unit"`
	Tax           bool          `json:"tax"`
}

func (AllowanceCreate) TableName() string {
	return "allowance"
}

func IsValidUnit(unit AllowanceUnit) bool {
	switch unit {
	case UnitVND, UnitUSD, UnitEUR, UnitJPY, UnitGBP:
		return true
	default:
		return false
	}
}

func (a *Allowance) Validate() error {
	if a.ID == "" {
		return errors.New("ID không được để trống")
	}
	if a.AllowanceName == "" {
		return errors.New("Tên phụ cấp không được để trống")
	}
	if a.Amount < 0 {
		return errors.New("Số tiền phụ cấp không được âm")
	}
	if !IsValidUnit(a.Unit) {
		return errors.New("Đơn vị tiền tệ không hợp lệ")
	}
	if a.CreatedDate.After(time.Now()) {
		return errors.New("Ngày tạo không hợp lệ")
	}
	return nil
}
