package model

import "time"

type Role struct {
	ID          string      `gorm:"primaryKey" json:"id"`
	RoleName    string    `gorm:"" json:"role_name"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
	Accounts    []Account `gorm:"foreignKey:RoleID"`
}

func (Role) TableName() string { return "role" }
