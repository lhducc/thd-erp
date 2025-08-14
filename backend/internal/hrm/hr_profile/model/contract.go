package model

import "time"

type Contract struct {
	ContractId     string        `gorm:"type:varchar(8);primaryKey;column:contract_id" json:"contract_id"`
	EffectiveDate  time.Time     `gorm:"column:effective_date" json:"effective_date"`
	ExpiredDate    time.Time     `gorm:"column:expired_date" json:"expired_date"`
	SignDate       time.Time     `gorm:"column:sign_date" json:"sign_date"`
	Note           string        `gorm:"note" json:"note"`
	AttachedFile   string        `gorm:"attached_file" json:"attached_file"`
	Condition      string        `gorm:"type:condition_enum;column:condition" json:"condition"`
	CreatedDate    time.Time     `gorm:"column:created_date" json:"created_date"`
	IsDeleted      bool          `gorm:"column:is_deleted" json:"is_deleted"`
	ContractTypeId string        `gorm:"type:varchar(6);column:contract_type_id" json:"contract_type"`
	ContractType   *ContractType `gorm:"foreignKey:ContractTypeId;references:ContractTypeID" json:"contract_type_info,omitempty"`
	ApproveStatus  string        `gorm:"type:approve_status_enum;column:approve_status" json:"approve_status"`
	EmployeeID     string        `gorm:"type:varchar(8);column:employee_id" json:"employee_id"`
	
	Employee       *Employee     `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"employee_info"`
	Allowances     []*Allowance  `gorm:"many2many:contract_allowances;joinForeignKey:ContractID;joinReferences:AllowanceID" json:"allowances,omitempty"`
}

func (Contract) TableName() string { return "contract" }

type ContractAllowance struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ContractID  string `gorm:"column:contract_id" json:"contract_id"`
	AllowanceID string `gorm:"column:allowance_id" json:"allowance_id"`
}

func (ContractAllowance) TableName() string {
	return "contract_allowances"
}

type ContractCreate struct {
	ContractId     string    `gorm:"primaryKey;column:contract_id" json:"contract_id"`
	EffectiveDate  time.Time `gorm:"column:effective_date" json:"effective_date"`
	ExpiredDate    time.Time `gorm:"column:expired_date" json:"expired_date"`
	SignDate       time.Time `gorm:"column:sign_date" json:"sign_date"`
	Note           string    `gorm:"note" json:"note"`
	AttachedFile   string    `gorm:"attached_file" json:"attached_file"`
	Condition      string    `gorm:"condition" json:"condition"`
	CreatedDate    time.Time `gorm:"column:created_date" json:"created_date"`
	ContractTypeId string    `gorm:"column:contract_type_id" json:"contract_type"`
	ApproveStatus  string    `gorm:"column:approve_status" json:"approve_status"`
	Manager        string    `gorm:"column:employee_id" json:"employee_id"`

	AllowanceIDs []string `json:"allowance_ids"`
}

func (ContractCreate) TableName() string { return "contract" }

type ContractResponse struct {
	ContractID    string         `json:"contract_id"`
	EffectiveDate time.Time      `json:"effective_date"`
	ExpiredDate   time.Time      `json:"expired_date"`
	SignDate      time.Time      `json:"sign_date"`
	Note          string         `json:"note"`
	AttachedFile  string         `json:"attached_file"`
	Condition     string         `json:"condition"`
	CreatedDate   time.Time      `json:"created_date"`
	ContractType  string         `json:"contract_type"`
	ApproveStatus string         `gorm:"column:approve_status" json:"approve_status"`
	Employee      EmployeeSimple `json:"employee"`
	Allowances    []*Allowance   `json:"allowances"`
}

type EmployeeSimple struct {
	EmployeeID string           `json:"employee_id"`
	FullName   string           `json:"full_name"`
	Department DepartmentSimple `json:"department"`
}

type DepartmentSimple struct {
	DepartmentID   string       `json:"department_id"`
	DepartmentName string       `json:"department_name"`
	Office         OfficeSimple `json:"office"`
}

type OfficeSimple struct {
	OfficeID   string `json:"office_id"`
	OfficeName string `json:"office_name"`
}

type ContractBasicInfo struct {
	ContractID   string    `json:"contract_id"`
	ContractName string    `json:"contract_name"`
	Condition    string    `json:"condition"`
	CreatedDate  time.Time `json:"created_date"`
}

func ConvertToBasicContracts(contracts []Contract) []ContractBasicInfo {
	var result []ContractBasicInfo
	for _, c := range contracts {
		result = append(result, ContractBasicInfo{
			ContractID:   c.ContractId,
			ContractName: c.ContractType.ContractTypeName,
			Condition:    c.Condition,
			CreatedDate:  c.CreatedDate,
		})
	}
	return result
}
