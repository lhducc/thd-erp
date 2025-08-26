import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import DataTable from "@/components/DataTable.tsx";
import Loading from "@/components/Loading.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";
import {deleteDecisionTypeApi, getAllDecisionTypeApi} from "@/apis/decistion-type.api.ts";
import type {DecisionType} from "@/types/decistion-type.ts";
import CreateDecisionType from "@/components/CreateDecisionType.tsx";

const SettingDecisionPage = () => {
    const {
        data: decisionTypes,
        isPending: pendingDecisionTypes,
        refetch: refetchDecisionTypes,
    } = useQuery({
        queryKey: ["decisionTypes"],
        queryFn: getAllDecisionTypeApi,
        gcTime: 0,
        staleTime: 0,
    });

    const { mutateAsync: deleteOffice } = useMutation({
        mutationFn: (id: string) => deleteDecisionTypeApi(id),
        onSuccess: () => {
            toast.success("Xóa văn phòng thành công");
            refetchDecisionTypes();
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<DecisionType>[] = [
        {
            accessorKey: "decision_type_id",
            header: "Mã",
        },
        {
            accessorKey: "decision_type",
            header: "Tên quyết định",
        },
        {
            accessorKey: "description",
            header: "Mô tả",
        },
        {
            id: "actions",
            header: "Chỉnh sửa",
            cell: ({ row }) => {
                const decisionType = row.original;

                return (
                    <div className="flex gap-4">
                        <CreateDecisionType
                            editBtn={
                                <Button variant="outline">
                                    <SquarePen />
                                </Button>
                            }
                            data={decisionType}
                            type="edit"
                            refetch={refetchDecisionTypes}
                        />
                        <ConfirmDelete deleteFn={() => deleteOffice(decisionType.decision_type_id)} />
                    </div>
                );
            },
        },
    ];

    if (pendingDecisionTypes) {
        return <Loading />;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={decisionTypes || []}
                title="Loại quyết định"
                buttonCreate={<CreateDecisionType refetch={refetchDecisionTypes} />}
                keyFilter="decision_type"
            />
        </>
    );
};

export default SettingDecisionPage;
