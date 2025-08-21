package dto

import "erp/backend/internal/hrm/hr_profile/model"

type EmployeeDTO struct {
	EmployeeID   string  `gorm:"primaryKey;column:employee_id" json:"employee_id"`
	Fullname     string  `gorm:"column:full_name;type:varchar(255);not null" json:"full_name" validate:"required"`
	Birthday     string  `gorm:"column:birthday;type:date" json:"birthday"`
	Gender       string  `gorm:"column:gender;type:varchar(10)" json:"gender" validate:"oneof=Nam Nữ Khác"`
	WorkType     string  `gorm:"column:work_type;type:varchar(255);type:work_type_enum" json:"work_type"`
	PhoneNumber  string  `gorm:"column:phone_number;type:varchar(20)" json:"phone_number" validate:"omitempty,e164"`
	Email        string  `gorm:"column:email;type:varchar(255)" json:"email" validate:"omitempty,email"`
	Address      string  `gorm:"column:address;type:text" json:"address"`
	PositionID   string  `gorm:"column:position_id" json:"position_id"`
	JobTitleID   string  `gorm:"column:job_title_id" json:"job_title_id"`
	Status       string  `gorm:"column:status" json:"status" validate:"required,oneof=active inactive"`
	ManagerID    *string `gorm:"column:manager;foreignKey:EmployeeID;references:employee_id" json:"manager_id"`
	DepartmentID string  `gorm:"column:department_id" json:"department_id"`
	ScheduleID   *int    `gorm:"column:schedule_id" json:"schedule_id"`
	RoleID       string  `gorm:"column:role_id" json:"role_id"`
}

func (dto *EmployeeDTO) ConvertToEmployeeModel() *model.Employee {
	return &model.Employee{
		EmployeeID:   dto.EmployeeID,
		Fullname:     dto.Fullname,
		Birthday:     dto.Birthday,
		Gender:       dto.Gender,
		WorkType:     dto.WorkType,
		PhoneNumber:  dto.PhoneNumber,
		Email:        dto.Email,
		Address:      dto.Address,
		PositionID:   dto.PositionID,
		JobTitleID:   dto.JobTitleID,
		Status:       dto.Status,
		ManagerID:    dto.ManagerID,
		DepartmentID: dto.DepartmentID,
		ScheduleID:   dto.ScheduleID,
	}
}

type EmployeeResponse struct {
	Employee *model.Employee `json:"employee"`
	Role     *model.Role     `json:"role"`
}
