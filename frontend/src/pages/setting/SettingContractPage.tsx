import {useMutation, useQuery} from "@tanstack/react-query";
import {deleteOfficeApi} from "@/apis/office.api.ts";
import {toast} from "sonner";
import type {ColumnDef} from "@tanstack/react-table";
import CreateOfficeForm from "@/components/CreateOfficeForm.tsx";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import {getAllContractsTypeApi} from "@/apis/contract-type.api.ts";
import type {ContractType} from "@/types/contract.ts";

const SettingContractPage = () => {
    const {
        data: types,
        isPending: pendingTypes,
        refetch: refetchTypes,
    } = useQuery({
        queryKey: ["contract_type"],
        queryFn: getAllContractsTypeApi,
        gcTime: 0,
        staleTime: 0,
    });

    const { mutateAsync: deleteTypes } = useMutation({
        mutationFn: (id: string) => deleteOfficeApi(id),
        onSuccess: () => {
            toast.success("Xóa loại hợp đồng thành công");
            refetchTypes();
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<ContractType>[] = [
        {
            accessorKey: "contract_type_id",
            header: "Mã",
        },
        {
            accessorKey: "contract_type",
            header: "Tên loại hợp đồng",
        },
        {
            accessorKey: "contract_group",
            header: "Nhóm đơn hợp đồng",
        },
        {
            accessorKey: "duration",
            header: "Thời hạn",
        },
        {
            accessorKey: "unit",
            header: "Đơn vị",
        },
        {
            accessorKey: "working_form",
            header: "Hình thức",
        },
        {
            id: "actions",
            header: "Chỉnh sửa",
            cell: ({ row }) => {
                const data = row.original;

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
                        {/*    refetch={refetchTypes}*/}
                        {/*/>*/}
                        <ConfirmDelete deleteFn={() => deleteTypes(data.contract_type_id)} />
                    </div>
                );
            },
        },
    ];

    if (pendingTypes) {
        return <Loading />;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={types || []}
                title="Loại hợp đồng"
                buttonCreate={<CreateOfficeForm refetch={refetchTypes} />}
                keyFilter="contract_type"
            />
        </>
    );
};

export default SettingContractPage;
