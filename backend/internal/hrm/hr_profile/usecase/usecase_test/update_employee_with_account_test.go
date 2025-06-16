package usecase_test

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Giả lập FakeEmployeeRepo để test UpdateEmployeeWithAccount
type FakeEmployeeRepoWithUpdate struct {
	UpdatedEmployee  *model.Employee
	UpdatedAccountID int
	Err              error
}

func (r *FakeEmployeeRepoWithUpdate) UpdateEmployeeWithAccount(employee *model.Employee, accountID int) error {
	if r.Err != nil {
		return r.Err
	}
	r.UpdatedEmployee = employee
	r.UpdatedAccountID = accountID
	return nil
}

func TestUpdateEmployeeWithAccount_Success(t *testing.T) {
	repo := &FakeEmployeeRepoWithUpdate{}

	emp := &model.Employee{
		EmployeeID:  "THD001",
		Fullname:    "Nguyễn Văn A",
		PhoneNumber: "0123456789",
	}

	accountID := 123

	err := repo.UpdateEmployeeWithAccount(emp, accountID)

	assert.NoError(t, err)
	assert.Equal(t, emp, repo.UpdatedEmployee)
	assert.Equal(t, accountID, repo.UpdatedAccountID)
}

func TestUpdateEmployeeWithAccount_Failure(t *testing.T) {
	expectedErr := errors.New("cập nhật thất bại")
	repo := &FakeEmployeeRepoWithUpdate{
		Err: expectedErr,
	}

	emp := &model.Employee{
		EmployeeID:  "THD002",
		Fullname:    "Nguyễn Văn B",
		PhoneNumber: "0999999999",
	}

	accountID := 456

	err := repo.UpdateEmployeeWithAccount(emp, accountID)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, repo.UpdatedEmployee)
}
