import {useQuery} from "@tanstack/react-query";
import type {ColumnDef} from "@tanstack/react-table";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import type {AttendanceNeedApprovalType} from "@/types/attendance-need-approval-type.ts";
import {getAllAttendanceApprovalAPI} from "@/apis/attendance-approval.api.ts";
import Confirm from "@/components/ui/confirm.tsx";

const ApproveAttendancePage = () => {
    const {data: attendanceApproval, isLoading: pendingAttendances, refetch: refetchAttendances} = useQuery({
        queryKey: ["attendanceApproval"],
        queryFn: getAllAttendanceApprovalAPI
    });

    const columns: ColumnDef<AttendanceNeedApprovalType>[] = [
        {
            accessorKey: "employee.employee_id",
            header: "Mã NV",
        },
        {
            accessorKey: "employee.full_name",
            header: "Tên nhân sự",
        },
        {
            accessorKey: "employee.department.department_name",
            header: "Phòng ban",
        },
        {
            accessorKey: "employee.position.position_name",
            header: "Chức vụ",
        },
        {
            accessorKey: "gps",
            header: "Vị trí GPS",
        },
        {
            accessorKey: "checkin_image_url",
            header: "Ảnh check-in",
            cell: ({row}) => (
                <img className={`w-[100px] h-[100px]`} src={row.original.checkin_image_url} alt={`check-in`}/>
            ),
        },
        {
            accessorKey: "checkout_image_url",
            header: "Ảnh check-out",
            cell: ({row}) => (
                <img className={`w-[100px] h-[100px]`} src={row.original.checkout_image_url} alt={`check-out`}/>
            ),
        },
        {
            accessorKey: "work_type",
            header: "Hình thức",
        },
        {
            accessorKey: "status",
            header: "Trạng thái",
            cell: ({row}) => (
                <p className={`${row.original.status === 'Đã duyệt' ? 'text-green-500' : row.original.status === 'Không duyệt' ? 'text-red-500' : ''}`}>{row.original.status}</p>
            ),
        },
        {
            id: "actions",
            header: "Thao tác",
            cell: ({row}) => {
                const attendanceApproval = row.original;

                return (
                    <div className="flex gap-4">
                        <Confirm message={"Bạn có muốn phê duyệt chấm công này?"}
                                 btn={<Button variant={"outline"}>
                            <SquarePen/>
                        </Button>} onConfirm={function (): void {
                            throw new Error("Function not implemented.");
                        }}/>
                        <ConfirmDelete deleteFn={() => deleteOffice(office.office_id)}/>
                    </div>
                );
            },
        },
    ];

    if (pendingAttendances) {
        return <Loading/>;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={attendanceApproval || []}
                title="Thiết lập chấm công"
                // buttonCreate={<AttendanceApprovalForm refetch={refetchAttendances}/>}
                keyFilter="name"
            />
        </>
    );
};

export default ApproveAttendancePage;
