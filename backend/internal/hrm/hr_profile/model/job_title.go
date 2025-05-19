package model

import "time"

type JobTitle struct {
	JobTitleID       string    `gorm:"primaryKey;column:job_title_id;type:varchar(50)" json:"job_title_id"`
	JobTitle         string    `gorm:"type:varchar(255);not null" json:"job_title"`
	CreatedDate      time.Time `gorm:"autoCreateTime" json:"created_date"`
	HierarchyLevelID string    `gorm:"column:hierarchy_level_id;type:varchar(50)" json:"hierarchy_level_id"`

	// Quan hệ với HierarchyLevel
	HierarchyLevel *HierarchyLevel `gorm:"foreignKey:HierarchyLevelID;references:ID" json:"hierarchy_level,omitempty"`

	// Quan hệ với Employee
	// Employees []Employee `gorm:"foreignKey:JobTitleID;references:JobTitleID" json:"employees,omitempty"`
}

func (JobTitle) TableName() string { return "jobtitle" }

type JobTitleCreate struct {
	JobTitleID       string    `gorm:"primaryKey" json:"job_title_id"`
	JobTitle         string    `gorm:"type:varchar(255);not null" json:"job_title"`
	CreatedDate      time.Time `gorm:"autoCreateTime" json:"created_date"`
	HierarchyLevelID string    `gorm:"column:hierarchy_level_id" json:"hierarchy_level_id"`
}

func (JobTitleCreate) TableName() string { return "jobtitle" }
