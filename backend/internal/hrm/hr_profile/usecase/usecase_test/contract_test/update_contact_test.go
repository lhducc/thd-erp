package contract

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"erp/backend/internal/hrm/hr_profile/usecase/usecase_test/mock_test"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateContract_Success(t *testing.T) {
	mockContractRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockContractRepo, mockEmployeeRepo)

	ctx := context.Background()
	id := "HD001"
	input := &hrmmodel.ContractCreate{
		SignDate:       time.Now().Add(-24 * time.Hour),
		EffectiveDate:  time.Now(),
		ExpiredDate:    time.Now().Add(24 * time.Hour),
		Note:           "Update",
		AttachedFile:   "updated.pdf",
		Condition:      "Updated",
		ContractTypeId: "CT02",
		ApproveStatus:  "1",
		Manager:        "EMP002",
	}

	mockContractRepo.On("UpdateContract", ctx, id, input).Return(nil)

	err := biz.UpdateContract(ctx, id, input)
	assert.NoError(t, err)

	mockContractRepo.AssertExpectations(t)
}

func TestUpdateContract_Error(t *testing.T) {
	mockContractRepo := new(mock_test.MockContractRepo)
	mockEmployeeRepo := new(mock_test.MockEmployeeRepo)
	biz := usecase.NewContractBiz(mockContractRepo, mockEmployeeRepo)

	ctx := context.Background()
	id := "HD001"
	input := &hrmmodel.ContractCreate{
		SignDate:       time.Now().Add(-24 * time.Hour),
		EffectiveDate:  time.Now(),
		ExpiredDate:    time.Now().Add(24 * time.Hour),
		Note:           "Fail update",
		AttachedFile:   "fail.pdf",
		Condition:      "error",
		ContractTypeId: "CT02",
		ApproveStatus:  "1",
		Manager:        "EMP002",
	}

	mockContractRepo.On("UpdateContract", ctx, id, input).Return(fmt.Errorf("update failed"))

	err := biz.UpdateContract(ctx, id, input)
	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")

	mockContractRepo.AssertExpectations(t)
}
