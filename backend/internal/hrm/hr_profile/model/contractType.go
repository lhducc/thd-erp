package model

import (
	"errors"
	"time"
)

// ENUM: ContractGroup
type ContractGroup string

const (
	ConfirmTime        ContractGroup = "Hợp đồng xác định thời hạn"
	NoTimeConfirmation ContractGroup = "Hợp đồng không xác định thời hạn"
	Trial              ContractGroup = "Hợp đồng thử việc"
	VocationalTraining ContractGroup = "Hợp đồng đào tạo nghề"
	Service            ContractGroup = "Hợp đồng dịch vụ"
)

// ENUM: Unit
type Unit string

const (
	Year  Unit = "Năm"
	Month Unit = "Tháng"
	Week  Unit = "Tuần"
	Day   Unit = "Ngày"
)

// ENUM: WorkingForm
type WorkingForm string

const (
	FullTime           WorkingForm = "Toàn thời gian"
	PartTime           WorkingForm = "Bán thời gian"
	Collaborator       WorkingForm = "Cộng tác viên"
	Expert             WorkingForm = "Chuyên gia"
	BySession          WorkingForm = "Theo ca"
	ProductContract    WorkingForm = "Khoán sản phẩm"
	WorkSubcontracting WorkingForm = "Khoán công việc"
)

type ContractType struct {
	ContractTypeID   string        `gorm:"column:contract_type_id;primaryKey;type:char(6)" json:"contract_type_id"`
	ContractTypeName string        `gorm:"column:contract_type;type:varchar(50)" json:"contract_type"`
	ContractGroup    ContractGroup `gorm:"column:contract_group;type:contract_group_enum" json:"contract_group"`
	Duration         int           `gorm:"column:duration" json:"duration"`
	Unit             Unit          `gorm:"column:unit;type:unit_enum" json:"unit"`
	IsDelete         bool          `gorm:"column:is_delete" json:"is_delete"`
	WorkingForm      WorkingForm   `gorm:"column:working_type;type:working_type_enum" json:"working_form"`
	CreatedDate      time.Time     `gorm:"column:created_date;type:date" json:"created_date"`
}

func (ContractType) TableName() string {
	return "contracttype"
}

// ValidateContractType checks the required fields and logical consistency of ContractType
func (ct ContractType) ValidateContractType() error {
	if ct.ContractTypeID == "" {
		return errors.New("Mã loại hợp đồng không được để trống")
	}
	if ct.ContractTypeName == "" {
		return errors.New("Tên loại hợp đồng không được để trống")
	}

	if ct.Duration < 0 {
		return errors.New("Thời gian hợp đồng không được âm")
	}

	if ct.CreatedDate.IsZero() {
		return errors.New("Ngày tạo hợp đồng không được bỏ trống")
	}

	// check enum ContractGroup
	validContractGroups := []ContractGroup{
		ConfirmTime, NoTimeConfirmation, Trial, VocationalTraining, Service,
	}
	validGroup := false
	for _, v := range validContractGroups {
		if ct.ContractGroup == v {
			validGroup = true
			break
		}
	}
	if !validGroup {
		return errors.New("Giá trị nhóm hợp đồng không hợp lệ")
	}

	// check enum Unit
	validUnits := []Unit{Year, Month, Week, Day}
	validUnit := false
	for _, v := range validUnits {
		if ct.Unit == v {
			validUnit = true
			break
		}
	}
	if !validUnit {
		return errors.New("Đơn vị thời gian không hợp lệ")
	}

	// check enum WorkingForm
	validWorkingForms := []WorkingForm{
		FullTime, PartTime, Collaborator, Expert, BySession, ProductContract, WorkSubcontracting,
	}
	validWorking := false
	for _, v := range validWorkingForms {
		if ct.WorkingForm == v {
			validWorking = true
			break
		}
	}
	if !validWorking {
		return errors.New("Hình thức làm việc không hợp lệ")
	}

	return nil
}
