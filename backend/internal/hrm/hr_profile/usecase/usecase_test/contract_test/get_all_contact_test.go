package contract

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/internal/hrm/hr_profile/usecase/usecase_test/mock_test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllContract_Success(t *testing.T) {
	mockRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockRepo, mockEmployeeRepo)

	ctx := context.Background()

	expectedContracts := []hrmmodel.Contract{
		{ContractId: "HD001", Note: "Hợp đồng 1"},
		{ContractId: "HD002", Note: "Hợp đồng 2"},
	}

	mockRepo.On("GetAllContract", ctx).Return(expectedContracts, nil)

	result, err := biz.GetAllContract(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedContracts, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllContract_Error(t *testing.T) {
	mockRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockRepo, mockEmployeeRepo)

	ctx := context.Background()

	expectedErr := errors.New("cannot fetch contracts")
	mockRepo.On("GetAllContract", ctx).Return([]hrmmodel.Contract(nil), expectedErr)

	result, err := biz.GetAllContract(ctx)
	assert.Nil(t, result)
	assert.EqualError(t, err, expectedErr.Error())
	mockRepo.AssertExpectations(t)
}
