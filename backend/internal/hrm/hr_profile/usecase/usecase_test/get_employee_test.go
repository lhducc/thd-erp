package usecase_test

import (
	"errors"
	"fmt"
	"testing"

	"erp/backend/internal/hrm/hr_profile/model"
	"github.com/stretchr/testify/assert"
)

type getUserByIDRepo interface {
	GetUserById(id string) (model.Employee, error)
}

type fakeGetUserByIDRepo struct {
	ExpectedID string
	Result     model.Employee
	Err        error
}

func (f *fakeGetUserByIDRepo) GetUserById(id string) (model.Employee, error) {
	if f.Err != nil {
		return model.Employee{}, f.Err
	}
	f.ExpectedID = id
	return f.Result, nil
}

type testGetUserByIDBiz struct {
	repo getUserByIDRepo
}

func (biz *testGetUserByIDBiz) GetUserById(id string) (model.Employee, error) {
	if id == "" {
		return model.Employee{}, errors.New("invalid employee ID")
	}

	employee, err := biz.repo.GetUserById(id)
	if err != nil {
		return model.Employee{}, fmt.Errorf("failed to get employee: %w", err)
	}

	return employee, nil
}

func TestGetUserById_InvalidID(t *testing.T) {
	biz := &testGetUserByIDBiz{repo: &fakeGetUserByIDRepo{}}

	_, err := biz.GetUserById("")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid employee ID")
}

func TestGetUserById_Success(t *testing.T) {
	expected := model.Employee{
		EmployeeID: "EMP123",
		Fullname:   "Nguyễn Văn A",
	}
	repo := &fakeGetUserByIDRepo{
		Result: expected,
	}
	biz := &testGetUserByIDBiz{repo: repo}

	result, err := biz.GetUserById("EMP123")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, "EMP123", repo.ExpectedID)
}

func TestGetUserById_RepoError(t *testing.T) {
	repo := &fakeGetUserByIDRepo{
		Err: errors.New("db error"),
	}
	biz := &testGetUserByIDBiz{repo: repo}

	_, err := biz.GetUserById("EMP124")

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to get employee: db error")
}
