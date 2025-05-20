package model

import (
	"errors"
	"time"
)

// Enum for Condition
type DocumentCondition string

const (
	Expired DocumentCondition = "Hết hạn"
	Active  DocumentCondition = "Đang hiệu lực"
)

// Enum for Status
type DocumentStatus string

const (
	Approved    DocumentStatus = "Duyệt"
	NotApproved DocumentStatus = "Không duyệt"
	Pending     DocumentStatus = "Chờ duyệt"
)

type EmployeeDocument struct {
	DocumentID     string            `gorm:"primaryKey;type:char(6)" json:"document_id"`
	DocumentTypeID string            `gorm:"type:char(6);not null" json:"document_type_id"`
	EmployeeID     string            `gorm:"type:char(8);not null" json:"employee_id"`
	EffectiveDate  time.Time         `gorm:"type:date" json:"effective_date"`
	ExpiredDate    time.Time         `gorm:"type:date" json:"expired_date"`
	Note           string            `gorm:"type:varchar(100);null" json:"note"`
	Condition      DocumentCondition `gorm:"type:document_condition_enum;not null" json:"condition"`
	Status         DocumentStatus    `gorm:"type:document_status_enum;not null" json:"status"`
	AttachedFile   []byte            `gorm:"type:blob;null" json:"attached_file"`
	CreatedDate    time.Time         `gorm:"type:date;default:CURRENT_DATE" json:"created_date"`

	// Foreign key relations
	DocumentType EmployeeDocumentType `gorm:"foreignKey:DocumentTypeID" json:"document_type"`
	Employee     *EmployeeDocument    `gorm:"foreignKey:EmployeeID" json:"employee"`
}

// SetCondition
func (d *EmployeeDocument) SetCondition() {
	currentDate := time.Now()

	if currentDate.After(d.ExpiredDate) {
		d.Condition = Expired
	} else if currentDate.After(d.EffectiveDate) && currentDate.Before(d.ExpiredDate) {
		d.Condition = Active
	}
}

func (EmployeeDocument) TableName() string {
	return "employeedocuments"
}

func (d EmployeeDocument) ValidateEmployeeDocument() error {
	if d.DocumentID == "" {
		return errors.New("mã tài liệu không hợp lệ")
	}
	if d.DocumentTypeID == "" {
		return errors.New("mã loại tài liệu không hợp lệ")
	}
	if d.EmployeeID == "" {
		return errors.New("mã nhân viên không hợp lệ")
	}
	if d.EffectiveDate.After(d.ExpiredDate) {
		return errors.New("ngày hiệu lực phải trước ngày hết hạn")
	}
	if d.Status != Approved && d.Status != NotApproved && d.Status != Pending {
		return errors.New("trạng thái tài liệu không hợp lệ")
	}
	if d.Condition != Active && d.Condition != Expired {
		return errors.New("tình trạng tài liệu không hợp lệ")
	}
	return nil
}
