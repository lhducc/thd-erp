import {useWorkScheduleRegister} from "@/query/useWorkScheduleRegister.ts";
import {useMutation} from "@tanstack/react-query";
import {deleteWorkshiftSchedule} from "@/apis/work-schedule.api.ts";
import {toast} from "sonner";
import type {ColumnDef} from "@tanstack/react-table";
import {Link} from "react-router-dom";
import {EditIcon, SettingsIcon} from "lucide-react";
import {Button} from "@/components/ui/button.tsx";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import DataTable from "@/components/DataTable.tsx";
import type {WorkScheduleRegister} from "@/types/work-schedule-register.ts";
import PATH from "@/constants/Path.ts";

const WorkScheduleRegisterPage = () => {
    const {data: schedule, isPending: pendingSchedule, refetch: refetchSchedule} = useWorkScheduleRegister()

    const {mutateAsync: deleteSchedule} = useMutation({
        mutationFn: (id: number) => deleteWorkshiftSchedule(id),
        onSuccess: () => {
            refetchSchedule();
            toast.success("Xóa lịch làm việc thành công");
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<WorkScheduleRegister>[] = [
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
            header: "Quản lý",
            cell: ({ row }) => {
                const managers = row.original.managers
                return (
                    <div className={`flex flex-col space-y-2`}>
                        {managers.map((manager) => (
                            <p>{manager.manager.full_name}</p>
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
            <Link to={`${PATH.WORK_SCHEDULE_REGISTER}/create`}>
                Tạo mới
            </Link>
        </Button>
    );

    return (
        <>
            <DataTable
                columns={columns}
                data={schedule || []}
                isLoading={pendingSchedule}
                title="Lịch làm việc đăng ký"
                buttonCreate={button}
                keyFilter="work_schedule_name"
            />
        </>
    );
};

export default WorkScheduleRegisterPage;
