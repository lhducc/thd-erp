package employee_test_test

import (
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type getByStatusRepo interface {
	GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error)
}

type fakeGetByStatusRepo struct {
	Status   string
	Page     int
	PageSize int
	Result   []model.Employee
	Err      error
}

func (f *fakeGetByStatusRepo) GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error) {
	f.Status = status
	f.Page = page
	f.PageSize = pageSize
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Result, nil
}

type testGetByStatusBiz struct {
	repo getByStatusRepo
}

func (biz *testGetByStatusBiz) GetAllEmployeesByStatus(status string, page, pageSize int) ([]model.Employee, error) {
	employees, err := biz.repo.GetAllEmployeesByStatus(status, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees with status: %w", err)
	}
	return employees, nil
}

func TestGetAllEmployeesByStatus_Success(t *testing.T) {
	expected := []model.Employee{{EmployeeID: "EMP1"}}
	repo := &fakeGetByStatusRepo{Result: expected}
	biz := &testGetByStatusBiz{repo: repo}

	res, err := biz.GetAllEmployeesByStatus("active", 2, 20)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	assert.Equal(t, "active", repo.Status)
	assert.Equal(t, 2, repo.Page)
	assert.Equal(t, 20, repo.PageSize)
}

func TestGetAllEmployeesByStatus_RepoError(t *testing.T) {
	repo := &fakeGetByStatusRepo{Err: errors.New("filter error")}
	biz := &testGetByStatusBiz{repo: repo}

	_, err := biz.GetAllEmployeesByStatus("inactive", 1, 10)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to get employees with status: filter error")
}
