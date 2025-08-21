import {useMemo, Suspense, lazy} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import type {ColumnDef} from "@tanstack/react-table";
import {Button} from "@/components/ui/button.tsx";
import type {Workshift} from "@/types/workshift.ts";
import {deleteWorkshiftApi, getAllWorkshiftApi} from "@/apis/workshift.api.ts";
import DataTable from "@/components/DataTable.tsx";
import {toast} from "sonner";
import dayjs from "dayjs";
import {ActionButtons} from "@/pages/HR/checkin/workshift/_components/ActionButtons.tsx";

const WorkshiftForm = lazy(() => import("./_components/WorkshiftForm"))

const DateRangeDisplay = ({startDate, endDate}: {startDate: string, endDate: string}) => {
    return (
        <p>
            {dayjs(startDate, "DD/MM/YYYY").format("DD/MM/YYYY")} - {dayjs(endDate, "DD/MM/YYYY").format("DD/MM/YYYY")}
        </p>
    );
};

const Index = () => {
    const {data: workshifts, refetch: refetchWorkshifts} = useQuery({
        queryKey: ["workshift"],
        queryFn: getAllWorkshiftApi,
    });

    const {mutateAsync: deleteWorkshift} = useMutation({
        mutationFn: (id: string) => deleteWorkshiftApi(id),
        onSuccess: async () => {
            await refetchWorkshifts();
            toast.success("Xóa phân ca thành công!");
        },
        onError: () => {
            toast.error("Lỗi hệ thống!")
        }
    });

    const columns : ColumnDef<Workshift>[] = [
        {
            accessorKey: "workshift_name",
            header: "Tên ca làm việc",
        },
        {
            accessorKey: "workshift_id",
            header: "Mã ca",
        },
        {
            accessorKey: "effective_date",
            header: "Thời gian áp dụng",
            cell: ({row}) => {
                const workshift = row.original;
                return <DateRangeDisplay
                    startDate={workshift.effective_date}
                    endDate={workshift.expiration_date || ""}
                />;
            },
        },
        {
            accessorKey: "start_time",
            header: "Giờ bắt đầu",
        },
        {
            accessorKey: "end_time",
            header: "Giờ kết thúc",
        },
        {
            id: "actions",
            header: "Thao tác",
            cell: ({row}) => {
                const workshift = row.original;
                return (
                    <ActionButtons
                        workshift={workshift}
                        refetchWorkshifts={refetchWorkshifts}
                        deleteWorkshift={deleteWorkshift}
                    />
                );
            },
        },
    ];

    return (
        <DataTable
            columns={columns}
            data={workshifts || []}
            isLoading={!workshifts}
            title="Quản lý ca mẫu"
            buttonCreate={
                <Suspense fallback={<Button disabled>Đang tải...</Button>}>
                    <WorkshiftForm refetch={refetchWorkshifts}/>
                </Suspense>
            }
            keyFilter="workshift_name"
        />
    );
};

export default Index;