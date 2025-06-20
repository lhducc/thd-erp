package mock_test

import (
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"github.com/stretchr/testify/mock"
)

type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) GetUserById(id string) (hrmmodel.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(hrmmodel.Employee), args.Error(1)
}
func (m *MockEmployeeRepo) CreateEmployee(employee *hrmmodel.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeRepo) UpdateEmployeeWithAccount(employee *hrmmodel.Employee, accountID int) error {
	args := m.Called(employee, accountID)
	return args.Error(0)
}

func (m *MockEmployeeRepo) GetAllEmployees() ([]hrmmodel.Employee, error) {
	args := m.Called()
	return args.Get(0).([]hrmmodel.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) GetAllEmployeesPagination(page, pageSize int) ([]hrmmodel.Employee, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]hrmmodel.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) GetAllEmployeesByStatus(status string, page, pageSize int) ([]hrmmodel.Employee, error) {
	args := m.Called(status, page, pageSize)
	return args.Get(0).([]hrmmodel.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) UpdateEmployee(id string, updatedEmployee hrmmodel.Employee) error {
	args := m.Called(id, updatedEmployee)
	return args.Error(0)
}

func (m *MockEmployeeRepo) DeleteEmployee(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeRepo) GetLastEmployeeByCode(emp *hrmmodel.Employee) error {
	args := m.Called(emp)
	return args.Error(0)
}
