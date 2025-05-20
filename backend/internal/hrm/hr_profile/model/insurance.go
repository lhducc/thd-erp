package model

import (
	"errors"
	"time"
)

// InsuranceInformation
type Insurance struct {
	ID            string    `gorm:"type:char(6);primaryKey" json:"id"`
	EffectiveDate time.Time `gorm:"type:date;not null" json:"effective_date"`
	CreatedDate   time.Time `gorm:"type:date;not null" json:"created_date"`
	PaymentAmount float64   `gorm:"type:decimal(15,2);not null" json:"payment_amount"`

	Company_SocialInsurance          float64 `gorm:"type:decimal(15,2);not null" json:"company_social_insurance"`           // CTY_BHXH
	Company_AccidentDiseaseInsurance float64 `gorm:"type:decimal(15,2);not null" json:"company_accident_disease_insurance"` // CTY_BH TNLĐ-BNN
	Company_HealthInsurance          float64 `gorm:"type:decimal(15,2);not null" json:"company_health_insurance"`           // CTY_BHYT
	Company_UnemploymentInsurance    float64 `gorm:"type:decimal(15,2);not null" json:"company_unemployment_insurance"`     // CTY_BHTN

	Employee_SocialInsurance          float64 `gorm:"type:decimal(15,2);not null" json:"employee_social_insurance"`           // NV_BHXH
	Employee_AccidentDiseaseInsurance float64 `gorm:"type:decimal(15,2);not null" json:"employee_accident_disease_insurance"` // NV_BH TNLĐ-BNN
	Employee_HealthInsurance          float64 `gorm:"type:decimal(15,2);not null" json:"employee_health_insurance"`           // NV_BHYT
	Employee_UnemploymentInsurance    float64 `gorm:"type:decimal(15,2);not null" json:"employee_unemployment_insurance"`     // NV_BHTN
}

func (Insurance) TableName() string {
	return "insuranceinformation"
}

// ValidateInsurance checks required fields of the insurance
func (i Insurance) ValidateInsurance() error {
	if i.EffectiveDate.IsZero() {
		return errors.New("ngày bắt đầu hiệu lực không được để trống")
	}
	if i.PaymentAmount <= 0 {
		return errors.New("số tiền thanh toán phải lớn hơn 0")
	}
	if i.Company_SocialInsurance <= 0 {
		return errors.New("mức đóng bảo hiểm xã hội của công ty phải lớn hơn 0")
	}
	if i.Employee_SocialInsurance <= 0 {
		return errors.New("mức đóng bảo hiểm xã hội của nhân viên phải lớn hơn 0")
	}
	if i.Company_AccidentDiseaseInsurance < 0 || i.Employee_AccidentDiseaseInsurance < 0 {
		return errors.New("mức bảo hiểm tai nạn và bệnh nghề nghiệp không được âm")
	}
	return nil
}
