import {useMutation, useQuery} from "@tanstack/react-query";
import {deleteOfficeApi} from "@/apis/office.api.ts";
import {toast} from "sonner";
import type {ColumnDef} from "@tanstack/react-table";
import type {Office} from "@/types";
import CreateOfficeForm from "@/components/CreateOfficeForm.tsx";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import {getAllSetupWorkshiftAPI} from "@/apis/setup-workshift.api.ts";

const SetupWorkshift = () => {
    const {
        data: setupWorkshift,
        isPending: pendingSetupWorkshift,
        refetch: refetchSetupWorkshift,
    } = useQuery({
        queryKey: ["setupWorkshift"],
        queryFn: getAllSetupWorkshiftAPI,
        gcTime: 0,
        staleTime: 0,
    });

    const { mutateAsync: deleteOffice } = useMutation({
        mutationFn: (id: string) => deleteOfficeApi(id),
        onSuccess: () => {
            toast.success("Xóa văn phòng thành công");
            refetchSetupWorkshift();
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<Office>[] = [
        {
            accessorKey: "office_name",
            header: "STT",
        },
        {
            accessorKey: "phone_number",
            header: "Tên lịch làm việc",
        },
        {
            accessorKey: "address",
            header: "Văn phòng",
        },
        {
            accessorKey: "address",
            header: "Quản lý",
        },
        {
            accessorKey: "address",
            header: "Trạng thái áp dụng",
        },
        {
            id: "actions",
            header: "Thao tác",
            cell: ({ row }) => {
                const office = row.original;

                return (
                    <div className="flex gap-4">
                        <SetupWorkshiftForm
                            editBtn={
                                <Button variant="outline">
                                    <SquarePen />
                                </Button>
                            }
                            office={office}
                            type="edit"
                            refetch={refetchSetupWorkshift}
                        />
                        <ConfirmDelete deleteFn={() => deleteOffice(office.office_id)} />
                    </div>
                );
            },
        },
    ];

    if (pendingSetupWorkshift) {
        return <Loading />;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={setupWorkshift || []}
                title="Lịch làm việc đăng ký"
                buttonCreate={<CreateOfficeForm refetch={refetchSetupWorkshift} />}
                keyFilter="office_name"
            />
        </>
    );
};

export default SetupWorkshift;
