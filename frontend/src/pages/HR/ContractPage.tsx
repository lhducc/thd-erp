import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import DataTable from "@/components/DataTable.tsx";
import Loading from "@/components/Loading.tsx";
import {useMutation, useQuery} from "@tanstack/react-query";
import type {ColumnDef} from "@tanstack/react-table";
import {toast} from "sonner";
import type {Contract} from "@/types/contract.ts";
import {ContractForm} from "@/components/CreateContract.tsx";
import {deleteContractById, getAllContractsApi} from "@/apis/contract.api.ts";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";

const ContractPage = () => {
    const {
        data: contracts,
        isPending: pendingContracts,
        refetch: refetchContracts,
    } = useQuery({
        queryKey: ["contract"],
        queryFn: getAllContractsApi,
        gcTime: 0,
        staleTime: 0,
    });

    const {mutateAsync: deleteContract} = useMutation({
        mutationFn: (id: string) => deleteContractById(id),
        onSuccess: () => {
            toast.success("Xóa hợp đồng thành công");
            refetchContracts();
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<Contract>[] = [
        {
            accessorKey: "contract_id",
            header: "Mã hợp đồng",
        },
        {
            accessorKey: "employee.full_name",
            header: "Tên nhân sự",
        },
        {
            accessorKey: "employee.department.department_name",
            header: "Phòng ban",
        },
        {
            accessorKey: "contract_type",
            header: "Loại hợp đồng",
        },
        {
            accessorKey: "sign_date",
            header: "Ngày ký",
            cell: ({getValue}) => new Date(getValue<string>()).toLocaleDateString(),
        },
        {
            accessorKey: "effective_date",
            header: "Hiệu lực từ ngày",
            cell: ({getValue}) => new Date(getValue<string>()).toLocaleDateString(),
        },
        {
            accessorKey: "expired_date",
            header: "Ngày hết hạn",
            cell: ({getValue}) => new Date(getValue<string>()).toLocaleDateString(),
        },
        {
            accessorKey: "condition",
            header: "Tình trạng",
        },
        {
            id: "actions",
            header: "Thao tác",
            cell: ({row}) => {
                const contract = row.original;

                return (
                    <div className="flex gap-4">
                        <ContractForm
                            editBtn={
                                <Button variant="outline">
                                    <SquarePen />
                                </Button>
                            }
                            data={contract}
                            type="pending"
                            refetch={refetchContracts}
                        />
                        <ConfirmDelete deleteFn={() => deleteContract(contract.contract_id)}/>
                    </div>
                );
            },
        },
    ];

    if (pendingContracts) {
        return <Loading/>;
    }

    return (
        <>
            <DataTable
                columns={columns}
                // data={contracts || []}
                data={contracts || []}
                // navLink={navLink}
                title="Hợp đồng"
                buttonCreate={<ContractForm refetch={refetchContracts}/>}
                keyFilter="contract_type"
            />
        </>
    );
};

export default ContractPage;