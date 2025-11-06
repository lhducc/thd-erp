package model

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type JobPosting struct {
	JobPostingID   string  `json:"job_posting_id" gorm:"column:job_posting_id; primaryKey"`
	Title          string  `json:"title" gorm:"column:title; not null" validate:"required"`
	JobDescription string  `json:"job_description" gorm:"column:job_description; not null" validate:"required"`
	HiringQuantity int     `json:"hiring_quantity" gorm:"column:hiring_quantity" validate:"gt=0"`
	Location       string  `json:"location" gorm:"column:location not null" validate:"required"`
	Salary         float64 `json:"salary" gorm:"column:salary type:numeric(12,2)"`
	IsActive       bool    `json:"is_active" gorm:"column:is_active default:true"`
	IsDeleted      bool    `json:"is_deleted" gorm:"column:is_deleted default:false"`

	CreatedDate time.Time `json:"created_date" gorm:"column:created_date autoCreateTime"`
	ExpiredDate time.Time `json:"expired_date" gorm:"column:expired_date"`

	// FK
	ProcessFormID string `json:"process_form_id" gorm:"column:process_form_id"`
	CreatedBy     string `json:"created_by" gorm:"column:created_by"`
	DepartmentID  string `json:"department_id" gorm:"column:department_id"`
	PositionID    string `json:"position_id" gorm:"column:position_id"`

	// Relationships
	ProcessForm ProcessForm      `json:"process_form" gorm:"foreignKey:ProcessFormID;references:ProcessFormID"`
	CreatorInfo model.Employee   `json:"creator_info" gorm:"foreignKey:CreatedBy;references:EmployeeID"`
	Department  model.Department `json:"department" gorm:"foreignKey:DepartmentID;references:ID"`
	Position    model.Position   `json:"position" gorm:"foreignKey:PositionID;references:ID"`
}

func (JobPosting) TableName() string {
	return "job_posting"
}

func (jp *JobPosting) BeforeCreate(tx *gorm.DB) (err error) {
	var last JobPosting
	// Lấy record cuối cùng để biết ID mới nhất
	if err := tx.Order("job_posting_id DESC").First(&last).Error; err == nil {
		// Cắt phần số ra
		var lastNum int
		fmt.Sscanf(last.JobPostingID, "JP%d", &lastNum)
		jp.JobPostingID = fmt.Sprintf("JP%04d", lastNum+1)
	} else {
		// Nếu chưa có record nào thì bắt đầu từ 0001
		jp.JobPostingID = "JP0001"
	}

	jp.IsDeleted = false
	return
}
