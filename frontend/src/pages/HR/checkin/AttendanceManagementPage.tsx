import {useQuery} from "@tanstack/react-query";
import {getAllAttendanceManagementAPI} from "@/apis/attendance.management.api.ts";
import DataTable from "@/components/DataTable.tsx";
import CreateOfficeForm from "@/components/CreateOfficeForm.tsx";
import Loading from "@/components/Loading.tsx";
import type {ColumnDef} from "@tanstack/react-table";
import type {Office} from "@/types";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import type {AttendanceSetting} from "@/types/attendance.ts";
import AttendanceSettingForm from "@/components/AttendanceSettingForm.tsx";

const AttendanceManagementPage = () => {
    const {data: attendanceSetting, isLoading: pendingAttendances, refetch: refetchAttendances} = useQuery({
        queryKey: ["attendanceSettings"],
        queryFn: getAllAttendanceManagementAPI
    });

    const columns: ColumnDef<AttendanceSetting>[] = [
        {
            accessorKey: "name",
            header: "Hình thức chấm công",
        },
        {
            accessorKey: "department.department_name",
            header: "Văn phòng",
        },
        {
            accessorKey: "gpsEnabled",
            header: "Định vị GPS",
            cell: ({ row }) => (
                <input
                    type="checkbox"
                    checked={row.original.gpsEnabled}
                    readOnly
                />
            ),
        },
        {
            accessorKey: "useCamera",
            header: "Hình ảnh camera",
            cell: ({ row }) => (
                <input
                    type="checkbox"
                    checked={row.original.useCamera}
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
                            refetch={refetchAttendances}                       />
                        <ConfirmDelete deleteFn={() => deleteOffice(office.office_id)} />
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
                keyFilter="name"
            />
        </>
    );
};

export default AttendanceManagementPage;
