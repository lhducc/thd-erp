package model

import (
	"errors"
	"time"
)

type EmployeeDocumentType struct {
	ID               string    `gorm:"column:id;primaryKey;type:char(6)"` // Ví dụ định dạng “EDXXXX”
	DocumentGroup    string    `gorm:"column:document_group;type:varchar(100)"`
	DocumentTypeName string    `gorm:"column:document_type_name;type:varchar(100)"`
	Description      string    `gorm:"column:description;type:text"`
	CreatedDate      time.Time `gorm:"column:created_date;type:date"`
}

func (EmployeeDocumentType) TableName() string { return "employeedocumenttype" }

// Validate checks required fields and length for EmployeeDocumentType
func (edt EmployeeDocumentType) Validate() error {
	if edt.ID == "" {
		return errors.New("mã loại tài liệu không được để trống")
	}
	if len(edt.ID) != 6 {
		return errors.New("mã loại tài liệu phải đúng 6 ký tự")
	}
	if edt.DocumentGroup == "" {
		return errors.New("nhóm tài liệu không được để trống")
	}
	if edt.DocumentTypeName == "" {
		return errors.New("tên loại tài liệu không được để trống")
	}
	if edt.CreatedDate.IsZero() {
		return errors.New("ngày tạo không được để trống")
	}
	return nil
}
