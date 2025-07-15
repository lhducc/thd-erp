package model

import "time"

type Office struct {
	ID          string    `gorm:"type:varchar(8);primaryKey;column:office_id;unique;" json:"office_id"`
	Name        string    `gorm:"type:varchar(100);column:office_name" json:"office_name"`
	PhoneNumber string    `gorm:"type:varchar(12);column:phone_number" json:"phone_number"`
	Address     string    `gorm:"type:text;column:address" json:"address"`
	Latitude    *float64  `gorm:"type:float;column:latitude" json:"latitude"`
	Longitude   *float64  `gorm:"type:float;column:longitude" json:"longitude"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
}

func (Office) TableName() string {
	return "office"
}

type OfficeCreate struct {
	ID          string    `gorm:"primaryKey;column:office_id" json:"office_id"`
	Name        string    `gorm:"column:office_name" json:"office_name"`
	PhoneNumber string    `gorm:"column:phone_number" json:"phone_number"`
	Address     string    `gorm:"column:address" json:"address"`
	Latitude    float64   `gorm:"column:latitude" json:"latitude"`
	Longitude   float64   `gorm:"column:longitude" json:"longitude"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
}

func (OfficeCreate) TableName() string {
	return "office"
}
