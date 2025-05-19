package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type Account struct {
	ID          int       `gorm:"primaryKey;column:id" json:"account_id"`
	LoginMail   string    `gorm:"column:login_mail" json:"login_mail"`
	Password    string    `gorm:"column:password" json:"password"`
	FirstLogin  bool      `gorm:"column:first_login" json:"first_login"`
	RoleID      string    `gorm:"column:role_id" json:"role_id"`
	EmployeeId  string    `gorm:"column:employee_id" json:"employee_id"`
	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`

	Role *Role `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}

func (Account) TableName() string { return "account" }

type Employee struct {
	EmployeeID  string    `gorm:"primaryKey;column:employee_id" json:"employee_id"`
	Fullname    string    `gorm:"column:full_name;type:varchar(255);not null" json:"full_name" validate:"required"`
	Birthday    string    `gorm:"column:birthday;type:date" json:"birthday"`
	Gender      string    `gorm:"column:gender;type:varchar(10)" json:"gender" validate:"oneof=Nam Nữ Khác"`
	WorkType    string    `gorm:"column:work_type;type:varchar(3)" json:"work_type" validate:"oneof=TTS THD CTV"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(20)" json:"phone_number" validate:"omitempty,e164"`
	Email       string    `gorm:"column:email;type:varchar(255)" json:"email" validate:"omitempty,email"`
	AccountID   *int      `gorm:"column:account_id" json:"account_id"`
	PositionID  string    `gorm:"column:position_id" json:"position_id"`
	JobTitleID  string    `gorm:"column:job_title_id" json:"job_title_id"`
	Status      string    `gorm:"column:status" json:"status" validate:"required,oneof=active inactive"`
	ManagerID   string    `gorm:"column:manager" json:"manager_id"`
	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	Account     *Account  `gorm:"foreignKey:AccountID;references:id" json:"account,omitempty"`
	// Position  *Position  gorm:"foreignKey:PositionID;references:position_id" json:"position,omitempty"
	// JobTitle  *JobTitle  gorm:"foreignKey:JobTitleID;references:job_title" json:"job_title,omitempty"
	Manager *Employee `gorm:"foreignKey:ManagerID;references:employee_id" json:"manager,omitempty"`
	// Contracts []Contract gorm:"foreignKey:EmployeeID;references:employee_id" json:"contracts,omitempty"
	// Decisions []Decision gorm:"foreignKey:EmployeeID;references:employee_id" json:"decisions,omitempty"
}

func (Employee) TableName() string { return "employee" }

type EmployeeCreate struct {
	EmployeeID  string `gorm:"primaryKey"`
	Fullname    string `gorm:"type:varchar(255);not null"`
	Birthday    string `gorm:"type:date"`
	Gender      string `gorm:"type:varchar(10)"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	Email       string `gorm:"type:varchar(255)"`
	Status      string `gorm:"column:status" json:"status"`
	LoginID     string
	PositionID  string
	JobTitleID  string
	ManagerID   string
	CreatedDate time.Time `gorm:"autoCreateTime"`
	// Contracts    []Contract `gorm:"foreignKey:Employee_ID"`
	// Decisions    []Decision `gorm:"many2many:employee_decisions;"`
}

func (EmployeeCreate) TableName() string { return "employee" }

type ResetPassword struct {
	AccountID       int    `json:"account_id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	RetypePassword  string `json:"retype_password"`
}

var validate = validator.New()

func ValidateEmployee(emp Employee) error {
	err := validate.Struct(emp)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("Lỗi: field '%s' failed on the '%s' tag", err.Field(), err.Tag())
		}
	}
	return nil
}

func ValidateEmployeeReferences(tx *gorm.DB, e Employee) error {
	if e.PositionID != "" {
		var pos Position
		if err := tx.First(&pos, "position_id = ?", e.PositionID).Error; err != nil {
			return fmt.Errorf("Không tồn tại vị trí này %s: %w", e.PositionID, err)
		}
	}
	if e.JobTitleID != "" {
		var jt JobTitle
		if err := tx.First(&jt, "job_title_id = ?", e.JobTitleID).Error; err != nil {
			return fmt.Errorf("Không tồn tại %s: %w", e.JobTitleID, err)
		}
	}
	if e.ManagerID != "" {
		var mgr Employee
		if err := tx.First(&mgr, "employee_id = ?", e.ManagerID).Error; err != nil {
			return fmt.Errorf("Không tồn tại %s: %w", e.ManagerID, err)
		}
	}
	return nil
}

func UpdateEmployeeFields(existing *Employee, updated Employee) {
	existing.Fullname = updated.Fullname
	existing.Birthday = updated.Birthday
	existing.Gender = updated.Gender
	existing.PhoneNumber = updated.PhoneNumber
	existing.Email = updated.Email
	existing.PositionID = updated.PositionID
	existing.ManagerID = updated.ManagerID
}
