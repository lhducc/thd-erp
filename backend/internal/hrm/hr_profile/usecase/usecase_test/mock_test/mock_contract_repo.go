package mock_test

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/store"
	"github.com/stretchr/testify/mock"
)

type MockContractRepo struct {
	mock.Mock
}

func (m *MockContractRepo) GetLastContractByCode(ctx context.Context, contract *hrmmodel.Contract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}

func (m *MockContractRepo) CreateContract(ctx context.Context, contract *hrmmodel.Contract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}

func (m *MockContractRepo) GetContract(ctx context.Context, id string) (*hrmmodel.Contract, error) {
	args := m.Called(ctx, id)
	// check if nil để tránh panic
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*hrmmodel.Contract), args.Error(1)
}

func (m *MockContractRepo) DeleteContract(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContractRepo) UpdateContract(ctx context.Context, id string, data *hrmmodel.ContractCreate) error {
	args := m.Called(ctx, id, data)
	return args.Error(0)
}

func (m *MockContractRepo) GetAllContract(ctx context.Context) ([]hrmmodel.Contract, error) {
	args := m.Called(ctx)
	return args.Get(0).([]hrmmodel.Contract), args.Error(1)
}

func (m *MockContractRepo) WithTransaction(ctx context.Context, fn func(store.ContractRepo) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockContractRepo) CheckExistName(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}
