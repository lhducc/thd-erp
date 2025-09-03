package contract_test

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/store"
	usecase "erp/backend/internal/hrm/hr_profile/usecase"
	mock_test "erp/backend/internal/hrm/hr_profile/usecase/usecase_test/mock_test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateContract(t *testing.T) {
	mockContractRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)

	biz := usecase.NewContractBiz(mockContractRepo, mockEmployeeRepo)

	ctx := context.Background()

	input := &hrmmodel.ContractCreate{
		Manager:        "EMP001",
		SignDate:       time.Now().Add(-24 * time.Hour),
		EffectiveDate:  time.Now(),
		ExpiredDate:    time.Now().Add(24 * time.Hour),
		Note:           "Test",
		AttachedFile:   "file.pdf",
		Condition:      "None",
		ContractTypeId: "CT01",
		ApproveStatus:  "0",
	}

	// Expect: lấy thông tin nhân viên
	mockEmployeeRepo.On("GetUserById", "EMP001").
		Return(hrmmodel.Employee{EmployeeID: "EMP001", Fullname: "Nguyễn Văn A"}, nil)

	// Expect: gọi trong GenerateCode
	mockContractRepo.On("GetLastContractByCode", mock.Anything, mock.AnythingOfType("*model.Contract")).
		Return(nil).
		Run(func(args mock.Arguments) {
			contract := args.Get(1).(*hrmmodel.Contract)
			contract.ContractId = "HD000005"
		})

	// Expect: tạo hợp đồng
	mockContractRepo.On("CreateContract", mock.Anything, mock.AnythingOfType("*model.Contract")).
		Return(nil)

	// Expect: transaction wrapper
	mockContractRepo.On("WithTransaction", mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(store.ContractRepo) error)
			_ = txFunc(mockContractRepo)
		})

	err := biz.CreateContract(ctx, input)
	assert.NoError(t, err)

	mockContractRepo.AssertExpectations(t)
	mockEmployeeRepo.AssertExpectations(t)
}
