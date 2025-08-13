import {useMutation, useQuery} from "@tanstack/react-query";
import {deleteAttendanceSettingApi, getAllAttendanceManagementAPI} from "@/apis/attendance-management.api.ts";
import DataTable from "@/components/DataTable.tsx";
import Loading from "@/components/Loading.tsx";
import type {ColumnDef} from "@tanstack/react-table";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import type {AttendanceSetting} from "@/types/attendance.ts";
import AttendanceSettingForm from "@/components/AttendanceSettingForm.tsx";
import {toast} from "sonner";

const AttendanceManagementPage = () => {
    const {data: attendanceSetting, isLoading: pendingAttendances, refetch: refetchAttendances} = useQuery({
        queryKey: ["attendanceSettings"],
        queryFn: getAllAttendanceManagementAPI,
        staleTime: 50000
    });

    const { mutateAsync: deleteAttendanceSetting } = useMutation({
        mutationFn: (id: string) => deleteAttendanceSettingApi(id),
        onSuccess: () => {
            toast.success("Xóa phụ cấp thành công");
            refetchAttendances();
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<AttendanceSetting>[] = [
        {
            accessorKey: "attendance_category_name",
            header: "Hình thức chấm công",
        },
        {
            accessorKey: "Office.office_name",
            header: "Văn phòng",
        },
        {
            accessorKey: "is_gps",
            header: "Định vị GPS",
            cell: ({ row }) => (
                <input
                    type="checkbox"
                    checked={row.original.is_gps}
                    readOnly
                />
            ),
        },
        {
            accessorKey: "is_camera",
            header: "Hình ảnh camera",
            cell: ({ row }) => (
                <input
                    type="checkbox"
                    checked={row.original.is_camera}
                    readOnly
                />
            ),
        },
        {
            accessorKey: "status",
            header: "Trạng thái",
            cell: ({ row }) => (
                <input
                    type="checkbox"
                    checked={row.original.status}
                    readOnly
                />
            ),
        },
        {
            id: "actions",
            header: "Thao tác",
            cell: ({ row }) => {
                const attendanceSetting = row.original;

                return (
                    <div className="flex gap-4">
                        <AttendanceSettingForm
                            editBtn={<Button variant="outline">
                                <SquarePen/>
                            </Button>}
                            data={attendanceSetting}
                            type="edit"
                            refetch={refetchAttendances}/>
                        <ConfirmDelete deleteFn={() => deleteAttendanceSetting(attendanceSetting.attendance_category_id)} />
                    </div>
                );
            },
        },
    ];

    if (pendingAttendances) {
        return <Loading />;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={attendanceSetting || []}
                title="Thiết lập chấm công"
                buttonCreate={<AttendanceSettingForm refetch={refetchAttendances} />}
                keyFilter="attendance_category_name"
            />
        </>
    );
};

export default AttendanceManagementPage;
