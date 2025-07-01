import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import type {ColumnDef} from "@tanstack/react-table";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import type {WorkShift} from "@/types/Workshift.ts";
import {deleteWorkshiftApi, getAllWorkshiftApi} from "@/apis/workshift.api.ts";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import {toast} from "sonner";
import WorkshiftForm from "@/components/WorkshiftForm.tsx";
import dayjs from "dayjs";

const WorkshiftPage = () => {
    const {data: workshifts, isLoading: isWorkshiftLoading, refetch: refetchWorkshifts} = useQuery({
        queryKey: ["workshift"],
        queryFn: getAllWorkshiftApi,
    })

    const queryClient = useQueryClient();

    const {mutateAsync: deleteWorkshift} = useMutation({
        mutationFn: deleteWorkshiftApi,
        onSuccess: () => {
            toast.success("Xóa phân ca thành công!")
            queryClient.invalidateQueries({
                queryKey: ["workshift"],
            });
        },
        onError: () => {
            toast.error("Lỗi hệ thống!")
        }
    });

    const columns: ColumnDef<WorkShift>[] = [
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
            cell: ({ row }) => {
                const workshift = row.original
                return (
                    <p>
                        {dayjs(workshift.effective_date, "DD/MM/YYYY").format("DD/MM/YYYY")} - {dayjs(workshift.expiration_date, "DD/MM/YYYY").format("DD/MM/YYYY")}
                    </p>
                );
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
            id: "status",
            header: "Trạng thái",
            cell: ({ row }) => {
                const workshift = row.original
                const effectiveDate = dayjs(workshift.effective_date, "DD/MM/YYYY")
                const expirationDate = dayjs(workshift.expiration_date, "DD/MM/YYYY")
                const now = dayjs()

                return (
                    <div className="flex gap-4">
                        {effectiveDate.isBefore(now) && expirationDate.isAfter(now) ? (
                            <p className="text-green-500">Đang áp dụng</p>
                        ) : effectiveDate.isAfter(now) ? (
                            <p className="text-gray-500">Chưa áp dụng</p>
                        ) : (
                            <p className="text-red-500">Hết hạn</p>
                        )}
                    </div>
                );
            },
        },
        {
            id: "actions",
            header: "Thao tác",
            cell: ({ row }) => {
                const workshift = row.original;

                return (
                    <div className="flex gap-4">
                        <WorkshiftForm
                            editBtn={
                                <Button variant="outline">
                                    <SquarePen />
                                </Button>
                            }
                            data={workshift}
                            type="edit"
                            refetch={refetchWorkshifts}
                        />
                        <ConfirmDelete deleteFn={() => deleteWorkshift(workshift.id)} />
                    </div>
                );
            },
        },
    ];

    if (isWorkshiftLoading) {
        return <Loading />;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={workshifts || []}
                title="Quản lý ca làm việc"
                buttonCreate={<WorkshiftForm refetch={refetchWorkshifts} />}
                keyFilter="workshift_name"
            />
        </>
    );
};

export default WorkshiftPage;
