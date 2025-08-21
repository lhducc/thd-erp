import {useMutation} from "@tanstack/react-query";
import {toast} from "sonner";
import type {ColumnDef} from "@tanstack/react-table";
import {Button} from "@/components/ui/button.tsx";
import {EditIcon, SettingsIcon} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import DataTable from "@/components/DataTable.tsx";
import type {WorkSchedule} from "@/types/work-schedule.ts";
import {Link} from "react-router-dom";
import {deleteWorkshiftSchedule} from "@/apis/work-schedule.api.ts";
import {useWorkSchedule} from "@/query/useWorkSchedule.ts";
import PATH from "@/constants/Path.ts";

const WorkScheduleAuto = () => {
    const {
        data: schedule,
        isPending: pendingSchedule,
        refetch: refetchSchedule,
    } = useWorkSchedule();

    const {mutateAsync: deleteSchedule} = useMutation({
        mutationFn: (id: number) => deleteWorkshiftSchedule(id),
        onSuccess: async () => {
            await refetchSchedule();
            toast.success("Xóa lịch làm việc thành công");
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<WorkSchedule>[] = [
        {
            accessorKey: "work_schedule_id",
            header: "STT",
        },
        {
            accessorKey: "work_schedule_name",
            header: "Tên lịch làm việc",
        },
        {
            accessorKey: "office.office_name",
            header: "Văn phòng",
        },
        {
            accessorKey: "manager",
            header: "Quản lý",
            cell: ({row}) => {
                const managers = row.original.managers
                return (
                    <div>
                        {managers.map((manager) => (
                            <p>{manager.manager?.full_name}</p>
                        ))}
                    </div>
                )
            }
        },
        {
            accessorKey: "status",
            header: "Trạng thái áp dụng",
            cell: ({row}) => {
                const workSchedule = row.original;

                const convert = workSchedule.status === "active" ? "Đang áp dụng" : "Không áp dụng"
                const style = workSchedule.status === "active" ? "text-green-500" : "text-red-400";
                return (
                    <p className={style}>
                        {convert}
                    </p>
                )
            }
        },
        {
            id: "actions",
            header: "Thao tác",
            cell: ({row}) => {
                return (
                    <div className="flex gap-4">
                        <Link to={`${row.original.work_schedule_id}`}>
                            <EditIcon/>
                        </Link>
                        <Link to={`${row.original.work_schedule_id}/setting`}>
                            <SettingsIcon/>
                        </Link>
                        <ConfirmDelete deleteFn={() => deleteSchedule(row.original.work_schedule_id)}/>
                    </div>
                );
            },
        },
    ];

    const button = (
        <Button>
            <Link to={`${PATH.WORK_SCHEDULE}/create`}>
                Tạo mới
            </Link>
        </Button>
    )

    return (
        <>
            <DataTable
                columns={columns}
                data={schedule || []}
                isLoading={pendingSchedule}
                title="Lịch làm việc tự động"
                buttonCreate={button}
                keyFilter="work_schedule_name"
            />
        </>
    );
};

export default WorkScheduleAuto;
