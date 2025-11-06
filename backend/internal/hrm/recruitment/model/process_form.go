package model

import (
	"fmt"
	"gorm.io/gorm"
)

type ProcessForm struct {
	ProcessFormID   string `json:"process_form_id" gorm:"column:process_form_id; primaryKey"`
	ProcessFormName string `json:"process_form_name" gorm:"column:process_form_name; not null" validate:"required"`
	IsDeleted       bool   `json:"is_deleted" gorm:"column:is_deleted; default:false"`

	// Relationships
	Stages []ProcessStage `json:"stages" gorm:"foreignKey:ProcessFormID;references:ProcessFormID"`
}

func (ProcessForm) TableName() string {
	return "process_form"
}

// Tự động tạo ProcessFormID trước khi tạo bản ghi mới theo định dạng PF0001, PF0002, ...
func (pf *ProcessForm) BeforeCreate(tx *gorm.DB) (err error) {
	var last ProcessForm

	// Lấy record cuối cùng để biết ID mới nhất
	if err := tx.Order("process_form_id DESC").First(&last).Error; err == nil {
		// Cắt phần số ra
		var lastNum int
		fmt.Sscanf(last.ProcessFormID, "PF%d", &lastNum)
		pf.ProcessFormID = fmt.Sprintf("PF%04d", lastNum+1)
	} else {
		// Nếu chưa có record nào thì bắt đầu từ 0001
		pf.ProcessFormID = "PF0001"
	}

	// Mặc định is_deleted = false
	pf.IsDeleted = false
	return
}
