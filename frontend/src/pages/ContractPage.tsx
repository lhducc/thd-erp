import {useQuery} from "@tanstack/react-query";
import {exportContractsFile, getAllContractsApi} from "@/apis/contract.api.ts";
import DataTable from "@/components/DataTable.tsx";
import Loading from "@/components/Loading.tsx";
import type {ColumnDef} from "@tanstack/react-table";
import {Button} from "@/components/ui/button.tsx";
import type {Contract} from "@/types/contract.ts";
import export_file from "@/assets/export-file.svg";
import {ContractForm} from "@/components/ContractForm.tsx";
import {useState} from "react";
import {ContractActions} from "@/components/ContractActions.tsx";

const ContractPage = () => {
    const [loading, setLoading] = useState(false);
    const [activeTab, setActiveTab] = useState<'approved' | 'pending' | 'rejected'>('approved');
    const [open, setOpen] = useState(false);

    const {
        data: contracts,
        isPending: pendingContracts,
    } = useQuery({
        queryKey: ["contracts"],
        queryFn: getAllContractsApi,
    });

    if (pendingContracts) {
        return <Loading/>;
    }

    // Filter contracts based on active tab
    const filteredContracts = contracts?.filter(contract => {
        switch (activeTab) {
            case 'approved':
                return contract.approve_status === 'Đã duyệt';
            case 'pending':
                return contract.approve_status === 'Chưa duyệt';
            case 'rejected':
                return contract.approve_status === 'Không duyệt';
            default:
                return true;
        }
    }) || [];

    // Count contracts by status
    const countContractsByStatus = () => {
        if (!contracts) return { approved: 0, pending: 0, rejected: 0 };

        return {
            approved: contracts.filter(c => c.approve_status === 'Đã duyệt').length,
            pending: contracts.filter(c => c.approve_status === 'Chưa duyệt').length,
            rejected: contracts.filter(c => c.approve_status === 'Không duyệt').length,
        };
    };

    const statusCounts = countContractsByStatus();

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
                return <ContractActions contractId={contract.contract_id}/>;
            },
        },
    ];

    const exportFile = async () => {
        try {
            setLoading(true);
            const file = await exportContractsFile();
            const url = URL.createObjectURL(file);
            const link = document.createElement("a");
            link.href = url;
            link.download = file.name;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            URL.revokeObjectURL(url);
        } catch (error) {
            console.error("Error exporting file:", error);
        } finally {
            setLoading(false);
        }
    };

    const groupButton = (
        <>
            <Button onClick={() => setOpen(!open)} variant="default" className={`px-8 py-5 text-[17px] rounded-[15px]`}>
                <span className={`mb-1 text-[24px]`}>+</span>
                Thêm hợp đồng
            </Button>
            <ContractForm open={open} setOpen={setOpen} contractId={null}/>
            <Button
                onClick={exportFile}
                variant={"default"}
                className={`px-10 py-5 rounded-[15px] text-[17px]`}
                disabled={loading}
            >
                {loading ? "Đang xuất..." : (
                    <>
                        <img src={export_file} alt="export-file" className="w-[24px]"/>
                        Xuất file
                    </>
                )}
            </Button>
        </>
    )

    const navLink = (
        <>
            <hr className={`mb-10`}/>
            <div className="mb-4 border-b border-gray-200 dark:border-gray-700">
                <ul className="flex flex-wrap -mb-px text-sm font-medium text-center" id="default-tab"
                    data-tabs-toggle="#default-tab-content" role="tablist">
                    <li className="me-2" role="presentation">
                        <button
                            className={`inline-block p-4 border-b-2 rounded-t-lg ${activeTab === 'approved' ? 'border-blue-500 text-blue-600' : 'hover:text-gray-600 hover:border-gray-300'}`}
                            onClick={() => setActiveTab('approved')}
                            type="button"
                            role="tab"
                        >
                            Đã duyệt ({statusCounts.approved})
                        </button>
                    </li>
                    <li className="me-2" role="presentation">
                        <button
                            className={`inline-block p-4 border-b-2 rounded-t-lg ${activeTab === 'pending' ? 'border-blue-500 text-blue-600' : 'hover:text-gray-600 hover:border-gray-300'}`}
                            onClick={() => setActiveTab('pending')}
                            type="button"
                            role="tab"
                        >
                            Chờ duyệt ({statusCounts.pending})
                        </button>
                    </li>
                    <li className="me-2" role="presentation">
                        <button
                            className={`inline-block p-4 border-b-2 rounded-t-lg ${activeTab === 'rejected' ? 'border-blue-500 text-blue-600' : 'hover:text-gray-600 hover:border-gray-300'}`}
                            onClick={() => setActiveTab('rejected')}
                            type="button"
                            role="tab"
                        >
                            Không duyệt ({statusCounts.rejected})
                        </button>
                    </li>
                </ul>
            </div>
        </>
    )

    return (
        <div>
            <DataTable
                columns={columns}
                data={filteredContracts}
                title="Hợp đồng"
                navLink={navLink}
                buttonCreate={groupButton}
                keyFilter="contract_id"
            />
        </div>
    );
};

export default ContractPage;