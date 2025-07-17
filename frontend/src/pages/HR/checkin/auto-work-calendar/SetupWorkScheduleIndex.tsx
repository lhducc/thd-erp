import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {deleteOfficeApi} from "@/apis/office.api.ts";
import {toast} from "sonner";
import type {ColumnDef} from "@tanstack/react-table";
import CreateOfficeForm from "@/components/CreateOfficeForm.tsx";
import {Button} from "@/components/ui/button.tsx";
import {EditIcon, SettingsIcon, SquarePen, ViewIcon} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import {getAllSetupWorkshiftCalendar} from "@/apis/setup-workshift.api.ts";
import type {WorkSchedule} from "@/types/work-schedule.ts";
import {Link} from "react-router-dom";
import {deleteWorkshiftSchedule} from "@/apis/work-schedule.api.ts";

const SetupWorkScheduleIndex = () => {
    const {
        data: schedule,
        isPending: pendingSchedule,
        refetch: refetchSchedule,
    } = useQuery({
        queryKey: ["work-schedule"],
        queryFn: getAllSetupWorkshiftCalendar,
        gcTime: 0,
        staleTime: 0,
    });

    const queryClient = useQueryClient();

    const { mutateAsync: deleteSchedule } = useMutation({
        mutationFn: (id: string) => deleteWorkshiftSchedule(id),
        onSuccess: () => {
            toast.success("Xóa lịch làm việc thành công");
            // refetchOffices();
            queryClient.invalidateQueries({
                queryKey: ["work-schedule"],
            });
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
        },
        {
            accessorKey: "status",
            header: "Trạng thái áp dụng",
            cell: ({ row }) => {
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
            cell: ({ row }) => {
                return (
                    <div className="flex gap-4">
                        {/*<CreateOfficeForm*/}
                        {/*    editBtn={*/}
                        {/*        <Button variant="outline">*/}
                        {/*            <SquarePen />*/}
                        {/*        </Button>*/}
                        {/*    }*/}
                        {/*    office={office}*/}
                        {/*    type="edit"*/}
                        {/*    refetch={refetchSchedule}*/}
                        {/*/>*/}
                        <Link to={`${row.original.work_schedule_id}`}>
                            <EditIcon />
                        </Link>
                        <Button>
                            <ViewIcon />
                        </Button>
                        <ConfirmDelete deleteFn={() => deleteSchedule(row.original.work_schedule_id)} />
                        <Link to={`${row.original.work_schedule_id}/setting`}>
                            <SettingsIcon />
                        </Link>
                    </div>
                );
            },
        },
    ];

    if (pendingSchedule) {
        return <Loading />;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={schedule || []}
                title="Lịch làm việc"
                buttonCreate={<CreateOfficeForm refetch={refetchSchedule} />}
                keyFilter="work_schedule_name"
            />
        </>
    );
};

export default SetupWorkScheduleIndex;
