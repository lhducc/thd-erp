package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"github.com/stretchr/testify/mock"
)

type MockContractBiz struct {
	mock.Mock
}

func (m *MockContractBiz) CreateContract(ctx context.Context, data *model.ContractCreate) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockContractBiz) GetContract(ctx context.Context, id string) (*model.Contract, error) {
	args := m.Called(ctx, id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.Contract), args.Error(1)
}

func (m *MockContractBiz) GetAllContract(ctx context.Context) ([]model.Contract, error) {
	args := m.Called(ctx)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]model.Contract), args.Error(1)
}

func (m *MockContractBiz) UpdateContract(ctx context.Context, id string, data *model.ContractCreate) error {
	args := m.Called(ctx, id, data)
	return args.Error(0)
}

func (m *MockContractBiz) DeleteContract(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContractBiz) ExportContractTest(ctx context.Context, fields []string) ([]byte, string, error) {
	args := m.Called(ctx, fields)
	var data []byte
	if args.Get(0) != nil {
		data = args.Get(0).([]byte)
	}
	filename := args.String(1)
	return data, filename, args.Error(2)
}
func (m *MockContractBiz) UpdateApproveStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockContractBiz) GetContractByEmployeeID(ctx context.Context, employeeId string) ([]model.ContractBasicInfo, error) {
	args := m.Called(ctx, employeeId)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]model.ContractBasicInfo), args.Error(1)
}
