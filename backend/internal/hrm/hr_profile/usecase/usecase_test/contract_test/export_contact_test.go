package contract

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/internal/hrm/hr_profile/usecase/usecase_test/mock_test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExportContractTest(t *testing.T) {
	mockContractRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockContractRepo, mockEmployeeRepo)

	ctx := context.Background()

	contracts := []model.Contract{
		{
			ContractId:     "HD0001",
			EffectiveDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpiredDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			SignDate:       time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			Note:           "Test note",
			AttachedFile:   "file.pdf",
			Condition:      "None",
			CreatedDate:    time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			ContractTypeId: "CT01",
			ApproveStatus:  "1",
			EmployeeID:     "EMP001",
			Employee:       &model.Employee{Fullname: "Nguyễn Văn A"},
		},
	}

	// Setup mock
	mockContractRepo.On("GetAllContract", ctx).Return(contracts, nil)

	selectedFields := []string{
		"contract_id", "effective_date", "expired_date", "sign_date",
		"note", "attached_file", "condition", "created_date",
		"contract_type", "approve_status", "employee_id", "employee_name",
	}

	// Run
	fileData, fileName, err := biz.ExportContractTest(ctx, selectedFields)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fileData)
	assert.NotEmpty(t, fileName)

	mockContractRepo.AssertExpectations(t)
}
func TestExportContractTest_Error(t *testing.T) {
	mockContractRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockContractRepo, mockEmployeeRepo)

	ctx := context.Background()

	mockContractRepo.On("GetAllContract", ctx).Return(nil, assert.AnError)

	selectedFields := []string{"contract_id"}

	fileData, fileName, err := biz.ExportContractTest(ctx, selectedFields)

	assert.Error(t, err)
	assert.Nil(t, fileData)
	assert.Empty(t, fileName)

	mockContractRepo.AssertExpectations(t)
}
