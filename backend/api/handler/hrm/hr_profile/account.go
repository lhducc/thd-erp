package handler

//type AccountService interface {
//	GetAccount(ctx context.Context, accountId int64) (*model.Account, error)
//}
//type AccountHandler struct {
//	biz *usecase.EmployeeBiz // or interface if you decouple more
//}
//
//func NewAccountHandler(biz *usecase.EmployeeBiz) *AccountHandler {
//	return &AccountHandler{biz: biz}
//}
//
//// get account by id
//func (h *AccountHandler) GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error) {
//	var (
//		account *model.Account
//		err     error
//	)
//	//check null
//	if accountID == 0 {
//		return nil, fmt.Errorf("thiếu mã tài khoản")
//	}
//	// fetch
//	account, err = h.biz.GetAccount(ctx, accountID)
//	if err != nil {
//		return nil, fmt.Errorf("không thể lấy tài khoản: %w", err)
//	}
//	return account, nil
//}
