package employee_test_test

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type getAllRepo interface {
	GetAllEmployeesPagination(page, pageSize int) ([]model.Employee, error)
}

type fakeGetAllRepo struct {
	Page     int
	PageSize int
	Result   []model.Employee
	Err      error
}

func (f *fakeGetAllRepo) GetAllEmployeesPagination(page, pageSize int) ([]model.Employee, error) {
	f.Page = page
	f.PageSize = pageSize
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Result, nil
}

type testGetAllBiz struct {
	repo getAllRepo
}

func (biz *testGetAllBiz) GetAllEmployees(page, pageSize int) ([]model.Employee, error) {
	employees, err := biz.repo.GetAllEmployeesPagination(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get all employees: %w", err)
	}
	return employees, nil
}

func TestGetAllEmployees_Success(t *testing.T) {
	expected := []model.Employee{
		{EmployeeID: "E1"}, {EmployeeID: "E2"},
	}
	repo := &fakeGetAllRepo{Result: expected}
	biz := &testGetAllBiz{repo: repo}

	res, err := biz.GetAllEmployees(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.Equal(t, 1, repo.Page)
	assert.Equal(t, 10, repo.PageSize)
}

func TestGetAllEmployees_RepoError(t *testing.T) {
	repo := &fakeGetAllRepo{Err: errors.New("query error")}
	biz := &testGetAllBiz{repo: repo}

	_, err := biz.GetAllEmployees(1, 5)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to get all employees: query error")
}
