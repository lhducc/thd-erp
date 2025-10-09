package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type Account struct {
	ID          int64     `gorm:"primaryKey;column:id;autoincrement" json:"account_id"`
	LoginMail   string    `gorm:"column:login_mail;index" json:"login_mail"`
	Password    string    `gorm:"column:password" json:"password"`
	FirstLogin  bool      `gorm:"column:first_login" json:"first_login"`
	RoleID      string    `gorm:"column:role_id" json:"role_id"`
	EmployeeId  string    `gorm:"column:employee_id;index" json:"employee_id"`
	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`

	Role     *Role     `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	Employee *Employee `gorm:"foreignKey:EmployeeId;references:EmployeeID" json:"employee,omitempty"`
}

func (Account) TableName() string { return "account" }

type Employee struct {
	EmployeeID   string    `gorm:"primaryKey;column:employee_id" json:"employee_id"`
	Fullname     string    `gorm:"column:full_name;type:varchar(255);not null" json:"full_name" validate:"required"`
	Birthday     string    `gorm:"column:birthday;type:date" json:"birthday"`
	Gender       string    `gorm:"column:gender;type:varchar(10)" json:"gender" validate:"oneof=Nam Nữ Khác"`
	WorkType     string    `gorm:"column:work_type;type:varchar(255);type:work_type_enum" json:"work_type"`
	PhoneNumber  string    `gorm:"column:phone_number;type:varchar(20)" json:"phone_number"`
	Email        string    `gorm:"column:email;type:varchar(255)" json:"email" validate:"omitempty,email"`
	Address      string    `gorm:"column:address;type:text" json:"address"`
	AccountID    *int64    `gorm:"column:account_id" json:"account_id"`
	PositionID   string    `gorm:"column:position_id" json:"position_id"`
	JobTitleID   string    `gorm:"column:job_title_id" json:"job_title_id"`
	Status       string    `gorm:"column:status" json:"status" validate:"required,oneof=active inactive"`
	ManagerID    *string   `gorm:"column:manager" json:"manager_id"`
	DepartmentID string    `gorm:"column:department_id" json:"department_id"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	ScheduleID   *int      `gorm:"column:schedule_id" json:"schedule_id"`

	Account    *Account         `gorm:"foreignKey:AccountID;references:id" json:"-"`
	Position   *Position        `gorm:"foreignKey:PositionID;references:position_id" json:"position,omitempty"`
	JobTitle   *JobTitle        `gorm:"foreignKey:JobTitleID;references:job_title_id" json:"job_title,omitempty"`
	Manager    *ManagerResponse `gorm:"foreignKey:ManagerID;references:employee_id" json:"manager,omitempty"`
	Department *Department      `gorm:"foreignKey:DepartmentID;references:department_id" json:"department,omitempty"`

	Contracts []Contract `gorm:"foreignKey:EmployeeID;references:EmployeeID" json:"contracts,omitempty"`
	Decisions []Decision `gorm:"many2many:decision_employees;joinForeignKey:EmployeeID;joinReferences:DecisionID" json:"decisions,omitempty"`
}

type ManagerResponse struct {
	EmployeeID string `gorm:"primaryKey;column:employee_id" json:"employee_id"`
	Fullname   string `gorm:"column:full_name;type:varchar(255);not null" json:"full_name" validate:"required"`
}
type JobTitleResponse struct {
	JobTitleID string `gorm:"primaryKey;column:job_title_id;type:varchar(50)" json:"job_title_id"`
	JobTitle   string `gorm:"type:varchar(255);not null" json:"job_title"`
}

type EmployeeInforResponse struct {
	EmployeeID   string `gorm:"column:employee_id" json:"employee_id"`
	Fullname     string `gorm:"column:full_name;" json:"full_name" validate:"required"`
	PositionID   string `gorm:"column:position_id" json:"-"`
	JobTitleID   string `gorm:"column:job_title_id" json:"-"`
	DepartmentID string `gorm:"column:department_id" json:"-"`
	WorkType     string `gorm:"column:work_type" json:"work_type"`

	JobTitle       *JobTitle       `gorm:"foreignKey:JobTitleID;references:job_title_id" json:"-"`
	HierarchyLevel *HierarchyLevel `gorm:"-" json:"hierarchy_level,omitempty"`
	Position       *Position       `gorm:"foreignKey:PositionID;references:position_id" json:"position,omitempty"`
	Department     *Department     `gorm:"column:department_id" json:"department,omitempty"`
}

func (JobTitleResponse) TableName() string { return "jobtitle" }

func (ManagerResponse) TableName() string { return "employee" }

func (Employee) TableName() string { return "employee" }

func (EmployeeInforResponse) TableName() string { return "employee" }

type ResetPassword struct {
	AccountID       int    `json:"account_id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	RetypePassword  string `json:"retype_password"`
}

var validate = validator.New()

func ValidateEmployee(emp *Employee) error {
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
	if e.ManagerID != nil {
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
