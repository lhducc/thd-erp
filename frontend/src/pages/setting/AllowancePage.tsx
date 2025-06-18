import {useMutation, useQuery} from "@tanstack/react-query";
import {toast} from "sonner";
import type {ColumnDef} from "@tanstack/react-table";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import {SquarePen} from "lucide-react";
import { Button } from "@/components/ui/button";
import {Link, useLocation} from "react-router-dom";
import {deleteAllowanceApi, getAllAllowancesApi} from "@/apis/allowance.api.ts";
import type {Allowance} from "@/types/allowance.ts";
import {useState} from "react";
import CreateAllowance from "@/components/CreateAllowance.tsx";

const AllowancePage = () => {
    const location = useLocation().pathname;
    const {
        data: types,
        isPending: pendingTypes,
        refetch: refetchTypes,
    } = useQuery({
        queryKey: ["allowance"],
        queryFn: getAllAllowancesApi,
        gcTime: 0,
        staleTime: 0,
    });

    const { mutateAsync: deleteTypes } = useMutation({
        mutationFn: (id: string) => deleteAllowanceApi(id),
        onSuccess: () => {
            toast.success("Xóa phụ cấp thành công");
            refetchTypes();
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const [columns, setColumns] = useState<ColumnDef<Allowance>[]>([
        {
            accessorKey: "id",
            header: "Mã",
        },
        {
            accessorKey: "amount",
            header: "Số tiền",
        },
        {
            accessorKey: "unit",
            header: "Đơn vị",
        },
        {
            accessorKey: "tax",
            header: "Chịu thuế",
        },
        {
            id: "actions",
            header: "Chỉnh sửa",
            cell: ({row}) => {
                const data = row.original;

                return (
                    <div className="flex gap-4">
                        <CreateAllowance
                            editBtn={
                                <Button variant="outline">
                                    <SquarePen/>
                                </Button>
                            }
                            data={data}
                            type="edit"
                            refetch={refetchTypes}
                        />
                        <ConfirmDelete deleteFn={() => deleteTypes(data.id)}/>
                    </div>
                );
            },
        },
    ]);

    if (pendingTypes) {
        return <Loading />;
    }

    const navLink = (
        <>
            <hr className={`mb-10`}/>
            <div className="mb-4 border-b border-gray-200 dark:border-gray-700">
                <ul className="flex flex-wrap -mb-px text-sm font-medium text-center" id="default-tab"
                    data-tabs-toggle="#default-tab-content" role="tablist">
                    <li className="me-2" role="presentation">
                        <Link
                            className={`inline-block p-4 border-b-2 rounded-t-lg ${location === '/setting/contract' ? 'border-[#DB3B21]' : 'hover:text-gray-600 hover:border-gray-300'}`}
                            type="button"
                            role="tab"
                            to="/setting/contract"
                        >
                            Loại hợp đồng
                        </Link>
                    </li>
                    <li className="me-2" role="presentation">
                        <Link
                            className={`inline-block p-4 border-b-2 rounded-t-lg ${location === '/setting/allowance' ? 'border-[#DB3B21]' : 'hover:text-gray-600 hover:border-gray-300'}`}
                            type="button"
                            role="tab"
                            to="/setting/allowance"
                        >
                            Phụ cấp
                        </Link>
                    </li>
                </ul>
            </div>
        </>
    )

    return (
        <>
            <DataTable
                columns={columns}
                data={types || []}
                title="Loại hợp đồng"
                navLink={navLink}
                buttonCreate={<CreateAllowance refetch={refetchTypes} />}
                keyFilter="contract_type"
            />
        </>
    );
};

export default AllowancePage;
