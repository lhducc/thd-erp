import type {ColumnDef} from "@tanstack/react-table";
import type {Document} from "@/types/employee-document.ts";
import {formatDate} from "@/lib/utils.ts";
import type {Contract} from "@/types/contract.ts";

export const documentColumn : ColumnDef < Document >[]  = [
    {accessorKey: "document_type.document_type_name", header: "Loại tài liệu"},
    {accessorKey: "status", header: "Trạng thái"},
    {
        accessorKey: "condition",
        header: "Tình trạng",
    },
    {
        accessorKey: "expired_date",
        header: "Ngày hết hạn",
        cell: ({row}) => formatDate(row.original.expired_date),
    },
    {
        accessorKey: "sign_date",
        header: "Chỉnh sửa",
        // cell: ({row}) => formatDate(row.original.,
    }
]

export const contractColumns : ColumnDef < Contract >[]  = [
    {accessorKey: "contract_id", header: "Mã hợp đồng"},
    {accessorKey: "contract_name", header: "Tên hợp đồng"},
    {
        accessorKey: "condition",
        header: "Tình trạng",
    },
    {
        accessorKey: "created_date",
        header: "Ngày tạo",
        cell: ({row}) => formatDate(row.original.created_date),
    }
]