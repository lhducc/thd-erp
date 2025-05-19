package model

import "time"

type HierarchyLevel struct {
	ID              string    `gorm:"primaryKey;column:id;type:varchar(50)" json:"id"`
	HierarchyLevel  string    `gorm:"column:hierarchy_level;type:varchar(100)" json:"hierarchy_level"`
	HierarchyNumber int       `gorm:"column:hierarchy_number" json:"hierarchy_number"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`

	// Quan hệ với JobTitle
	JobTitles []JobTitle `gorm:"foreignKey:HierarchyLevelID;references:ID" json:"job_titles,omitempty"`
}

func (HierarchyLevel) TableName() string { return "hierarchylevel" }

type HierarchyLevelCreate struct {
	ID              string    `gorm:"primaryKey;column:id" json:"id"`
	HierarchyLevel  string    `gorm:"column:hierarchy_level" json:"hierarchy_level"`
	HierarchyNumber int       `gorm:"column:hierarchy_number" json:"hierarchy_number"`
	CreatedDate     time.Time `gorm:"column:created_date" json:"created_date"`
}

func (HierarchyLevelCreate) TableName() string { return "hierarchylevel" }
