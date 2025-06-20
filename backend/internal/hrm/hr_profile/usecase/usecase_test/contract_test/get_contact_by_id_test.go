package contract_test

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/internal/hrm/hr_profile/usecase/usecase_test/mock_test"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContract(t *testing.T) {
	mockContractRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockContractRepo, mockEmployeeRepo)

	ctx := context.Background()
	expected := &hrmmodel.Contract{
		ContractId: "HD123",
		Note:       "Test contract",
	}

	mockContractRepo.On("GetContract", ctx, "HD123").Return(expected, nil)

	contract, err := biz.GetContract(ctx, "HD123")

	assert.NoError(t, err)
	assert.Equal(t, expected, contract)
	mockContractRepo.AssertExpectations(t)
}

func TestGetContract_Error(t *testing.T) {
	mockContractRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockContractRepo, mockEmployeeRepo)

	ctx := context.Background()
	mockContractRepo.On("GetContract", ctx, "HD404").Return(nil, fmt.Errorf("not found"))

	contract, err := biz.GetContract(ctx, "HD404")

	assert.Nil(t, contract)
	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	mockContractRepo.AssertExpectations(t)
}
