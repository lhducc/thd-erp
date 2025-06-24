package errors

import "github.com/pkg/errors"

const (
	MsgCreatedSuccess = "Tạo thành công"
	MsgListData       = "Danh sách dữ liệu"
	MsgUpdateSuccess  = "Cập nhật thành công"
	MsgDeleteSuccess  = "Xóa dữ liệu thành công"
)

var (
	ErrApprovedContractCannotEdit         = errors.New("Hợp đồng đã được duyệt, không được chỉnh sửa")
	ErrOnlyExpiredContractCanBeLiquidated = errors.New("Chỉ có thể thanh lý hợp đồng đã hết hiệu lực")
)
