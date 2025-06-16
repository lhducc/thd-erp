package usecase_test

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeEmployeeRepo struct {
	Created *model.Employee
	Err     error
}

func (r *FakeEmployeeRepo) CreateEmployee(employee *model.Employee) error {
	if r.Err != nil {
		return r.Err // Không gán nếu có lỗi
	}
	r.Created = employee
	return nil
}

func TestCreateEmployeeRepo_Success(t *testing.T) {
	repo := &FakeEmployeeRepo{}

	emp := &model.Employee{
		Fullname:    "Nguyễn Văn A",
		PhoneNumber: "0123456789",
	}

	err := repo.CreateEmployee(emp)

	assert.NoError(t, err)
	assert.Equal(t, emp.Fullname, repo.Created.Fullname)
	assert.Equal(t, emp.PhoneNumber, repo.Created.PhoneNumber)
}

func TestCreateEmployeeRepo_EmptyFullname(t *testing.T) {
	repo := &FakeEmployeeRepo{}

	emp := &model.Employee{
		Fullname:    "",
		PhoneNumber: "0123456789",
	}

	err := repo.CreateEmployee(emp)

	// Với repo giả, không có validate, ta mong không lỗi
	// Nhưng nếu bạn muốn validate sau này, có thể thay thành assert.Error
	assert.NoError(t, err)
	assert.Equal(t, "", repo.Created.Fullname)
}

func TestCreateEmployeeRepo_InvalidPhoneNumber(t *testing.T) {
	repo := &FakeEmployeeRepo{}

	emp := &model.Employee{
		Fullname:    "Nguyễn Văn B",
		PhoneNumber: "abcde", // không hợp lệ
	}

	err := repo.CreateEmployee(emp)

	assert.NoError(t, err)
	assert.Equal(t, "abcde", repo.Created.PhoneNumber)
}

func TestCreateEmployeeRepo_RepoReturnError(t *testing.T) {
	expectedErr := errors.New("lỗi khi insert vào DB")
	repo := &FakeEmployeeRepo{
		Err: expectedErr,
	}

	emp := &model.Employee{
		Fullname:    "Nguyễn Văn C",
		PhoneNumber: "0988888888",
	}

	err := repo.CreateEmployee(emp)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, repo.Created)
}
