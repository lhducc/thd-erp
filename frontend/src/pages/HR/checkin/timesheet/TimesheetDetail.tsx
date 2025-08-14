import {useParams} from "react-router-dom";
import {useGetTimesheet} from "@/query/timesheet.query.ts";
import type {ColumnDef} from "@tanstack/react-table";
import type {TimesheetInfor} from "@/types/timesheet.ts";
import DataTable from "@/components/DataTable.tsx";

const TimesheetDetail = () => {
    const { id } = useParams();

    const {data, isLoading}= useGetTimesheet(id || "")

    const columns: ColumnDef<TimesheetInfor>[] = [
        {
            accessorKey: "employee.employee_id",
            header: "Mã nhân viên",
        },
        {
            accessorKey: "employee.full_name",
            header: "Họ và tên",
        },
        {
            accessorKey: "department.department_name",
            header: "Phòng ban",
        },
        {
            accessorKey: "employee.position.position_name",
            header: "Chức vụ",
        },
        {
            accessorKey: "employee.hierarchy_level.hierarchy_level",
            header: "Cấp bậc",
        },
        {
            accessorKey: "total_late_minutes",
            header: "Số phút trễ",
        },
        {
            accessorKey: "total_work_days",
            header: "Tổng công",
        },
        // {
        //     accessorKey: "address",
        //     header: "Thời gian",
        //     cell: ({row}) => {
        //         return (
        //             <>
        //                 {formatDate(row.original.start_date)} - {formatDate(row.original.end_date)}
        //             </>
        //         )
        //     }
        // },
        // {
        //     accessorKey: "is_locked",
        //     header: "Trạng thái ",
        //     cell: ({row}) => {
        //         return (
        //             <>{row.original.is_locked ? "Đang áp dụng" : <p className="text-red-500">Đã chốt công</p>}</>
        //         )
        //     }
        // },
        // {
        //     id: "actions",
        //     header: "Xem chi tiết",
        //     cell: ({ row }) => {
        //         return (
        //             <Link to={`/detail/${row.original.timesheet_list_id}`} className="flex gap-4">
        //                 <EyeIcon />
        //             </Link>
        //         );
        //     },
        // },
    ];

    return (
        <DataTable
            columns={columns}
            data={data || []}
            isLoading={isLoading}
            title="Bảng công"
            buttonCreate={null}
            // keyFilter="time_sheet_list_name"
        />
    );
};

export default TimesheetDetail;
