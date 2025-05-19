package usecase

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"time"
)

type ContractRepo interface {
	CreateContract(context context.Context, data *hrmmodel.Contract) error
	GetContract(ctx context.Context, id string) (*hrmmodel.Contract, error)
	GetAllContract(ctx context.Context) ([]hrmmodel.Contract, error)
	UpdateContract(ctx context.Context, id string, data *hrmmodel.ContractCreate) error
	DeleteContract(ctx context.Context, id string) error
	CheckExistName(name string) (bool, error)
	GetLastContractByCode(ctx context.Context, contract *hrmmodel.Contract) error
}

type contractBiz struct {
	repo         ContractRepo
	employeeRepo EmployeeRepo
}

func NewContractBiz(store ContractRepo, employeeRepo EmployeeRepo) *contractBiz {
	return &contractBiz{
		repo:         store,
		employeeRepo: employeeRepo,
	}
}

func (biz *contractBiz) CreateContract(context context.Context, data *hrmmodel.ContractCreate) error {

	code, err := utils.GenerateCode("HD", 6, func() (string, error) {
		var last hrmmodel.Contract
		err := biz.repo.GetLastContractByCode(context, &last)
		if err != nil {
			return "", err
		}
		return last.ContractId, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	contract := &hrmmodel.Contract{
		ContractId:     code,
		EffectiveDate:  data.EffectiveDate,
		ExpiredDate:    data.ExpiredDate,
		SignDate:       data.SignDate,
		Note:           data.Note,
		AttachedFile:   data.AttachedFile,
		Condition:      data.Condition,
		ContractTypeId: data.ContractTypeId,
		EmployeeID:     data.Manager,
		CreatedDate:    time.Now(),
	}
	now := time.Now()
	if contract.SignDate.After(now) {
		return errors.New("Ngày ký không thể trong tương lai")
	}

	if contract.EffectiveDate.Before(contract.SignDate) {
		return errors.New("Ngày hiệu lực không thể trước ngày ký")
	}

	if contract.ExpiredDate.Before(contract.EffectiveDate) {
		return errors.New("Ngày hết hạn phải sau ngày hiệu lực")
	}
	exists, err := biz.employeeRepo.GetUserById(contract.EmployeeID)
	if err != nil {
		return err
	}
	contract.Employee = &exists
	err = biz.repo.CreateContract(context, contract)
	if err != nil {
		return err
	}

	return nil
}

func (biz *contractBiz) GetContract(ctx context.Context, id string) (*hrmmodel.Contract, error) {
	Contract, err := biz.repo.GetContract(ctx, id)
	if err != nil {
		return nil, err
	}
	return Contract, nil
}
func (biz *contractBiz) UpdateContract(ctx context.Context, id string, data *hrmmodel.ContractCreate) error {
	if err := biz.repo.UpdateContract(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *contractBiz) DeleteContract(ctx context.Context, id string) error {
	if err := biz.repo.DeleteContract(ctx, id); err != nil {
		return err
	}
	return nil
}

func (biz *contractBiz) GetAllContract(ctx context.Context) ([]hrmmodel.Contract, error) {
	positions, err := biz.repo.GetAllContract(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (e *contractBiz) ExportContractTest(c context.Context, selectedFields []string) ([]byte, string, error) {
	exporter := utils.NewExcelExporter("Contracts")

	exporter.RegisterField("contract_id", "Mã hợp đồng", "contract_id", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.ContractId
	})

	exporter.RegisterField("effective_date", "Ngày hiệu lực", "effective_date", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.EffectiveDate.Format("2006-01-02")
	})

	exporter.RegisterField("expired_date", "Ngày hết hạn", "expired_date", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.ExpiredDate.Format("2006-01-02")
	})

	exporter.RegisterField("sign_date", "Ngày ký", "sign_date", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.SignDate.Format("2006-01-02")
	})

	exporter.RegisterField("note", "Ghi chú", "note", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.Note
	})

	exporter.RegisterField("attached_file", "Tệp đính kèm", "attached_file", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.AttachedFile
	})

	exporter.RegisterField("condition", "Điều kiện", "condition", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.Condition
	})

	exporter.RegisterField("created_date", "Ngày tạo", "created_date", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.CreatedDate.Format("2006-01-02")
	})

	exporter.RegisterField("contract_type", "Loại hợp đồng", "contract_type", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.ContractTypeId
	})

	exporter.RegisterField("employee_id", "Mã nhân viên", "employee_id", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		return contract.EmployeeID
	})

	exporter.RegisterField("employee_name", "Tên nhân viên", "employee_name", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		if contract.Employee != nil {
			return contract.Employee.Fullname
		}
		return ""
	})

	contracts, err := e.repo.GetAllContract(c)
	if err != nil {
		return nil, "", fmt.Errorf("không thể lấy dữ liệu hợp đồng: %w", err)
	}

	contractPtrs := make([]*hrmmodel.Contract, len(contracts))
	for i := range contracts {
		contractPtrs[i] = &contracts[i]
	}

	return exporter.Export(contractPtrs, selectedFields)
}
