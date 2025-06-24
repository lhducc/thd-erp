package usecase

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/store"
	utils "erp/backend/pkg"
	errpkg "erp/backend/pkg/errors"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type contractBiz struct {
	repo         store.ContractRepo
	employeeRepo EmployeeRepo
}

func NewContractBiz(contractRepo store.ContractRepo, employeeRepo EmployeeRepo) *contractBiz {
	return &contractBiz{
		repo:         contractRepo,
		employeeRepo: employeeRepo,
	}
}

func (biz *contractBiz) CreateContract(ctx context.Context, data *hrmmodel.ContractCreate) error {
	employee, err := biz.employeeRepo.GetUserById(data.Manager)
	if err != nil {
		return fmt.Errorf("không tìm thấy nhân viên: %w", err)
	}

	return biz.repo.WithTransaction(ctx, func(txRepo store.ContractRepo) error {
		code, err := utils.GenerateCode("HD", 6, func() (string, error) {
			var last hrmmodel.Contract
			err := txRepo.GetLastContractByCode(ctx, &last)
			if err != nil {
				return "", err
			}
			return last.ContractId, nil
		})
		if err != nil {
			return fmt.Errorf("không thể tạo mã: %w", err)
		}

		// 4. Validate
		now := time.Now()
		if data.SignDate.After(now) {
			return errors.New("ngày ký không thể trong tương lai")
		}
		if data.EffectiveDate.Before(data.SignDate) {
			return errors.New("ngày hiệu lực không thể trước ngày ký")
		}
		if data.ExpiredDate.Before(data.EffectiveDate) {
			return errors.New("ngày hết hạn phải sau ngày hiệu lực")
		}

		// 5. Mapping
		contract := &hrmmodel.Contract{
			ContractId:     code,
			EffectiveDate:  data.EffectiveDate,
			ExpiredDate:    data.ExpiredDate,
			SignDate:       data.SignDate,
			Note:           data.Note,
			AttachedFile:   data.AttachedFile,
			Condition:      data.Condition,
			ContractTypeId: data.ContractTypeId,
			ApproveStatus:  data.ApproveStatus,
			EmployeeID:     data.Manager,
			CreatedDate:    now,
			Employee:       &employee,
		}

		if err := txRepo.CreateContract(ctx, contract); err != nil {
			return fmt.Errorf("không thể tạo hợp đồng: %w", err)
		}
		unique := map[string]struct{}{}
		for _, aid := range data.AllowanceIDs {
			if _, exists := unique[aid]; exists {
				continue
			}
			unique[aid] = struct{}{}
			ca := &hrmmodel.ContractAllowance{
				ID:          uuid.NewString(),
				ContractID:  contract.ContractId,
				AllowanceID: aid,
			}
			if err := txRepo.CreateContractAllowance(ctx, ca); err != nil {
				return fmt.Errorf("không thể lưu phụ cấp cho hợp đồng: %w", err)
			}
		}
		return nil
	})
}

func (biz *contractBiz) GetContract(ctx context.Context, id string) (*hrmmodel.Contract, error) {
	Contract, err := biz.repo.GetContract(ctx, id)
	if err != nil {
		return nil, err
	}
	return Contract, nil
}

func (biz *contractBiz) GetContractByEmployeeID(ctx context.Context, employeeId string) ([]hrmmodel.ContractBasicInfo, error) {
	contracts, err := biz.repo.GetContractByEmployeeID(ctx, employeeId)
	if err != nil {
		return nil, err
	}
	return hrmmodel.ConvertToBasicContracts(contracts), nil
}

func determineCondition(start, end time.Time) string {
	now := time.Now()
	switch {
	case now.Before(start):
		return "Chưa hiệu lực"
	case now.After(end):
		return "Hết hiệu lực"
	default:
		return "Đang hiệu lực"
	}
}

func isExpired(contract *hrmmodel.Contract) bool {
	return contract.ExpiredDate.Before(time.Now())
}

func (biz *contractBiz) UpdateApproveStatus(ctx context.Context, id string, status string) error {
	contract, err := biz.repo.GetContract(ctx, id)
	if err != nil {
		return err
	}

	if contract.ApproveStatus != "Không duyệt" {
		return fmt.Errorf("Chỉ được yêu cầu duyệt lại khi hợp đồng bị từ chối")
	}

	return biz.repo.UpdateApproveStatus(ctx, id, "Chờ duyệt")
}

func (biz *contractBiz) UpdateContract(ctx context.Context, id string, data *hrmmodel.ContractCreate) error {
	contract, err := biz.repo.GetContract(ctx, id)
	if err != nil {
		return err
	}

	if contract.ApproveStatus == "Đã duyệt" {
		return errpkg.ErrApprovedContractCannotEdit
	}

	if data.Condition == "Thanh lý" {
		if !isExpired(contract) {
			return errpkg.ErrOnlyExpiredContractCanBeLiquidated
		}
		data.Condition = "Thanh lý"
	} else {
		data.Condition = determineCondition(data.EffectiveDate, data.ExpiredDate)
	}

	if contract.ApproveStatus == "Chưa duyệt" {
		data.ApproveStatus = "Đang duyệt"
	} else {
		data.ApproveStatus = contract.ApproveStatus
	}

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

func (biz *contractBiz) GetAllContract(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]hrmmodel.Contract, int64, error) {
	contract, totalRecords, err := biz.repo.GetAllContractPagination(ctx, page, pageSize, filters)
	if err != nil {
		return nil, 0, err
	}
	return contract, totalRecords, nil
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

	exporter.RegisterField("approve_status", "Trạng thái duyệt", "approve_status", func(item interface{}) any {
		contract := item.(*hrmmodel.Contract)
		switch contract.ApproveStatus {
		case strconv.Itoa(0):
			return "Chờ duyệt"
		case strconv.Itoa(1):
			return "Đã duyệt"
		case strconv.Itoa(2):
			return "Từ chối"
		default:
			return "Không xác định"
		}
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
