package usecase_test

import (
	"errors"
	"fmt"
	"testing"

	"erp/backend/internal/hrm/hr_profile/model"
	"github.com/stretchr/testify/assert"
)

// Alias interface chỉ cho test
type updateOnlyRepo interface {
	UpdateEmployee(id string, emp model.Employee) error
}

// Usecase tạm chỉ test 1 method
type testEmployeeBiz struct {
	repo updateOnlyRepo
}

func (biz *testEmployeeBiz) UpdateEmployee(id string, emp model.Employee) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}
	if err := biz.repo.UpdateEmployee(id, emp); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

// Fake repo chỉ implement UpdateEmployee
type fakeEmployeeRepo struct {
	UpdateFn   func(id string, emp model.Employee) error
	CalledID   string
	CalledData model.Employee
}

func (f *fakeEmployeeRepo) UpdateEmployee(id string, emp model.Employee) error {
	f.CalledID = id
	f.CalledData = emp
	if f.UpdateFn != nil {
		return f.UpdateFn(id, emp)
	}
	return nil
}

func TestUpdateEmployee_InvalidID(t *testing.T) {
	repo := &fakeEmployeeRepo{}
	biz := &testEmployeeBiz{repo: repo}

	err := biz.UpdateEmployee("", model.Employee{})

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid employee ID")
}

func TestUpdateEmployee_Success(t *testing.T) {
	repo := &fakeEmployeeRepo{}
	biz := &testEmployeeBiz{repo: repo}

	input := model.Employee{
		Fullname:    "Nguyễn Văn A",
		PhoneNumber: "0123456789",
	}

	err := biz.UpdateEmployee("EMP001", input)

	assert.NoError(t, err)
	assert.Equal(t, "EMP001", repo.CalledID)
	assert.Equal(t, input.Fullname, repo.CalledData.Fullname)
	assert.Equal(t, input.PhoneNumber, repo.CalledData.PhoneNumber)
}

func TestUpdateEmployee_RepoError(t *testing.T) {
	repo := &fakeEmployeeRepo{
		UpdateFn: func(id string, emp model.Employee) error {
			return errors.New("db connection failed")
		},
	}
	biz := &testEmployeeBiz{repo: repo}

	err := biz.UpdateEmployee("EMP002", model.Employee{
		Fullname: "Nguyễn Văn B",
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to update employee: db connection failed")
}
