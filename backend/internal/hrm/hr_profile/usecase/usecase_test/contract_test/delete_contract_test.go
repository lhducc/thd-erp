package contract

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/internal/hrm/hr_profile/usecase/usecase_test/mock_test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteContract_Success(t *testing.T) {
	mockRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo) // vẫn cần truyền vào constructor
	biz := usecase.NewContractBiz(mockRepo, mockEmployeeRepo)

	ctx := context.Background()
	contractID := "HD000001"

	mockRepo.On("DeleteContract", ctx, contractID).Return(nil)

	err := biz.DeleteContract(ctx, contractID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteContract_Error(t *testing.T) {
	mockRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockRepo, mockEmployeeRepo)

	ctx := context.Background()
	contractID := "HD_NOT_FOUND"

	expectedErr := errors.New("contract not found")
	mockRepo.On("DeleteContract", ctx, contractID).Return(expectedErr)

	err := biz.DeleteContract(ctx, contractID)
	assert.EqualError(t, err, expectedErr.Error())
	mockRepo.AssertExpectations(t)
}
