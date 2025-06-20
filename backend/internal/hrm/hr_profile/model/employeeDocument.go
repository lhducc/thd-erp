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
	DocumentID     string            `gorm:"column:document_id;primaryKey;type:char(8)" json:"document_id"`
	DocumentTypeID string            `gorm:"column:document_type_id;type:char(8);not null" json:"document_type_id"`
	EmployeeID     string            `gorm:"column:employee_id;type:char(6);not null" json:"employee_id"`
	EffectiveDate  time.Time         `gorm:"column:effective_date;type:date" json:"effective_date"`
	ExpiredDate    time.Time         `gorm:"column:expired_date;type:date" json:"expired_date"`
	Note           string            `gorm:"column:note;type:varchar(100)" json:"note"`
	Condition      DocumentCondition `gorm:"column:condition;type:document_condition_enum;not null" json:"condition"`
	Status         DocumentStatus    `gorm:"column:status;type:document_status_enum;not null" json:"status"`
	AttachedFile   []byte            `gorm:"column:attached_file;type:bytea" json:"attached_file"`
	CreatedDate    time.Time         `gorm:"column:created_date;type:date;default:CURRENT_DATE" json:"created_date"`

	// Foreign key relations
	DocumentType *EmployeeDocumentType `gorm:"foreignKey:DocumentTypeID;references:ID" json:"document_type,omitempty"`
	Employee     *Employee             `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee"`
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
	return "employeedocument"
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

type EmployeeDocumentResponse struct {
	DocumentID    string            `json:"document_id"`
	EmployeeName  string            `json:"employee_name"`
	DocumentType  string            `json:"document_type"`
	EffectiveDate time.Time         `json:"effective_date"`
	ExpiredDate   time.Time         `json:"expired_date"`
	Condition     DocumentCondition `json:"condition"`
	EmployeeID    string            `json:"employee_id"`
}
