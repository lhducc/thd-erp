package employee_test_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type deleteOnlyRepo interface {
	DeleteEmployee(id string) error
}
type testDeleteEmployeeBiz struct {
	repo deleteOnlyRepo
}

func (biz *testDeleteEmployeeBiz) DeleteEmployee(id string) error {
	if id == "" {
		return errors.New("invalid employee ID")
	}
	if err := biz.repo.DeleteEmployee(id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}
	return nil
}

type fakeDeleteEmployeeRepo struct {
	CalledID string
	Err      error
}

func (f *fakeDeleteEmployeeRepo) DeleteEmployee(id string) error {
	f.CalledID = id
	return f.Err
}
func TestDeleteEmployee_InvalidID(t *testing.T) {
	repo := &fakeDeleteEmployeeRepo{}
	biz := &testDeleteEmployeeBiz{repo: repo}

	err := biz.DeleteEmployee("")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid employee ID")
}

func TestDeleteEmployee_Success(t *testing.T) {
	repo := &fakeDeleteEmployeeRepo{}
	biz := &testDeleteEmployeeBiz{repo: repo}

	err := biz.DeleteEmployee("EMP001")

	assert.NoError(t, err)
	assert.Equal(t, "EMP001", repo.CalledID)
}

func TestDeleteEmployee_RepoError(t *testing.T) {
	repo := &fakeDeleteEmployeeRepo{
		Err: errors.New("record not found"),
	}
	biz := &testDeleteEmployeeBiz{repo: repo}

	err := biz.DeleteEmployee("EMP002")

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to delete employee: record not found")
}
