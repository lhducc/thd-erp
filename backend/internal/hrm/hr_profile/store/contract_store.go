package store

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
)

type ContractRepo interface {
	CreateContract(context.Context, *model.Contract) error
	GetContract(context.Context, string) (*model.Contract, error)
	GetAllContractPagination(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Contract, int64, error)
	GetAllContract(ctx context.Context) ([]model.Contract, error)
	UpdateContract(context.Context, string, *model.ContractCreate) error
	DeleteContract(context.Context, string) error
	CheckExistName(name string) (bool, error)
	GetLastContractByCode(context.Context, *model.Contract) error
	WithTransaction(ctx context.Context, fn func(txRepo ContractRepo) error) error
	CreateContractAllowance(ctx context.Context, ca *model.ContractAllowance) error
	UpdateApproveStatus(ctx context.Context, id string, status string) error
	GetContractByEmployeeID(ctx context.Context, id string) ([]model.Contract, error)
}
