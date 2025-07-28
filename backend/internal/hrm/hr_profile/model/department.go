package model

import "time"

type Department struct {
	ID          string    `gorm:"type:varchar(6);primaryKey;column:department_id" json:"department_id"`
	Name        string    `gorm:"type:varchar(30);column:department_name" json:"department_name"`
	Manager     string    `gorm:"type:varchar(8);column:manager" json:"manager"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
	OfficeID    string    `gorm:"type:varchar(6);column:office_id" json:"office_id"`
	Office      *Office   `gorm:"foreignKey:OfficeID;references:office_id" json:"office,omitempty"`
}

type DepartmentCreate struct {
	ID          string    `gorm:"primaryKey;column:department_id" json:"department_id"`
	Name        string    `gorm:"column:department_name" json:"department_name"`
	Manager     string    `gorm:"column:manager" json:"manager"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
	OfficeID    string    `gorm:"column:office_id" json:"office_id"`
}

func (Department) TableName() string {
	return "department"
}
func (DepartmentCreate) TableName() string {
	return "department"
}
