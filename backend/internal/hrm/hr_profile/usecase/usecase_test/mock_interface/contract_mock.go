package mock_interface

import (
	"context"
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/store"
)

type ContractRepo interface {
	GetLastContractByCode(ctx context.Context, contract *hrmmodel.Contract) error
	CreateContract(ctx context.Context, contract *hrmmodel.Contract) error
	GetContract(ctx context.Context, id string) (*hrmmodel.Contract, error)
	DeleteContract(ctx context.Context, id string) error
	UpdateContract(ctx context.Context, id string, data *hrmmodel.ContractCreate) error
	GetAllContract(ctx context.Context) ([]hrmmodel.Contract, error)
	WithTransaction(ctx context.Context, fn func(store.ContractRepo) error) error
}
type EmployeeRepo interface {
	GetUserById(id string) (hrmmodel.Employee, error)
}
