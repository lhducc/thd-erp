import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import DataTable from "@/components/DataTable.tsx";
import Loading from "@/components/Loading.tsx";
import {useMutation, useQuery} from "@tanstack/react-query";
import type {ColumnDef} from "@tanstack/react-table";
import {toast} from "sonner";
import type {Contract} from "@/types/contract.ts";
import {ContractForm} from "@/components/CreateContract.tsx";
import {deleteContractById, getAllContractsApi} from "@/apis/contract.api.ts";
import {useState} from "react";
import {ContractFilter} from "@/components/ContractFilter.tsx";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";

const ContractPage = () => {
    // const [activeTab, setActiveTab] = useState<'approved' | 'pending' | 'rejected'>('approved');
    const [filters, setFilters] = useState({
        department: '',
        contractType: '',
        condition: ''
    });
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
                            // type={activeTab}
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

    const departments = Array.from(new Set(
        contracts?.map(c => c.employee.department.department_name) || []
    ));

    const contractTypes = Array.from(new Set(
        contracts?.map(c => c.contract_type) || []
    ));

    const conditions = Array.from(new Set(
        contracts?.map(c => c.condition) || []
    ));

    // Filter contracts based on active tab and filters
    // const filteredContracts = contracts?.filter(contract => {
    //     // Filter by tab
    //     let tabMatch = false;
    //     switch (activeTab) {
    //         case 'approved':
    //             tabMatch = contract.approve_status === 'Đã duyệt';
    //             break;
    //         case 'pending':
    //             tabMatch = contract.approve_status === 'Chờ duyệt';
    //             break;
    //         case 'rejected':
    //             tabMatch = contract.approve_status === 'Không duyệt';
    //             break;
    //         default:
    //             tabMatch = true;
    //     }
    //
    //     // Filter by department
    //     const departmentMatch = !filters.department ||
    //         contract.employee.department.department_name === filters.department;
    //
    //     // Filter by contract type
    //     const contractTypeMatch = !filters.contractType ||
    //         contract.contract_type === filters.contractType;
    //
    //     // Filter by condition
    //     const conditionMatch = !filters.condition ||
    //         contract.condition === filters.condition;
    //
    //     return tabMatch && departmentMatch && contractTypeMatch && conditionMatch;
    // }) || [];

    // // Count contracts by status
    // const countContractsByStatus = () => {
    //     if (!contracts) return {approved: 0, pending: 0, rejected: 0};
    //
    //     return {
    //         approved: contracts.filter(c => c.approve_status === 'Đã duyệt').length,
    //         pending: contracts.filter(c => c.approve_status === 'Chờ duyệt').length,
    //         rejected: contracts.filter(c => c.approve_status === 'Không duyệt').length,
    //     };
    // };
    //
    // const statusCounts = countContractsByStatus();
    //
    // const navLink = (
    //     <>
    //         <hr className={`mb-10`}/>
    //         <div className={`flex justify-between`}>
    //             <div className="mb-4 border-b border-gray-200 dark:border-gray-700">
    //                 <ul className="flex flex-wrap -mb-px text-sm font-medium text-center" id="default-tab"
    //                     data-tabs-toggle="#default-tab-content" role="tablist">
    //                     <li className="me-2" role="presentation">
    //                         <button
    //                             className={`inline-block p-4 border-b-2 rounded-t-lg ${activeTab === 'approved' ? 'border-[#DB3B21]' : 'hover:text-gray-600 hover:border-gray-300 text-gray-500'}`}
    //                             onClick={() => setActiveTab('approved')}
    //                             type="button"
    //                             role="tab"
    //                         >
    //                             Đã duyệt ({statusCounts.approved})
    //                         </button>
    //                     </li>
    //                     <li className="me-2" role="presentation">
    //                         <button
    //                             className={`inline-block p-4 border-b-2 rounded-t-lg ${activeTab === 'pending' ? 'border-[#DB3B21]' : 'hover:text-gray-600 hover:border-gray-300 text-gray-500'}`}
    //                             onClick={() => setActiveTab('pending')}
    //                             type="button"
    //                             role="tab"
    //                         >
    //                             Chờ duyệt ({statusCounts.pending})
    //                         </button>
    //                     </li>
    //                     <li className="me-2" role="presentation">
    //                         <button
    //                             className={`inline-block p-4 border-b-2 rounded-t-lg ${activeTab === 'rejected' ? 'border-[#DB3B21]' : 'hover:text-gray-600 hover:border-gray-300 text-gray-500'}`}
    //                             onClick={() => setActiveTab('rejected')}
    //                             type="button"
    //                             role="tab"
    //                         >
    //                             Không duyệt ({statusCounts.rejected})
    //                         </button>
    //                     </li>
    //                 </ul>
    //             </div>
    //             <ContractFilter departments={departments} contractTypes={contractTypes} conditions={conditions}
    //                             onFilterChange={setFilters}/>
    //         </div>
    //
    //     </>
    // )

    return (
        <>
            <DataTable
                columns={columns}
                // data={contracts || []}
                data={contracts}
                // navLink={navLink}
                title="Hợp đồng"
                buttonCreate={<ContractForm refetch={refetchContracts}/>}
                keyFilter="contract_type"
            />
        </>
    );
};

export default ContractPage;