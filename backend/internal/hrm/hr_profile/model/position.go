package model

import "time"

type Position struct {
	ID          string    `gorm:"primaryKey;column:position_id;type:varchar(50)" json:"position_id"`
	Name        string    `gorm:"column:position_name;type:varchar(100);not null" json:"position_name"`
	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
}

func (Position) TableName() string { return "position" }

