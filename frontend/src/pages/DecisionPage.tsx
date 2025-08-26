'use client';

import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { SquarePen, Filter } from "lucide-react";
import { useState, useRef, useEffect } from "react";
import export_file from "@/assets/export-file.svg";
import { getAllDecisionsApi, deleteDecisionApi, createSampleDecisionApi, exportDecisionExcelApi } from "@/apis/decision.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import CreateDecisionForm from "@/components/CreateDecisionForm";
import UpdateDecisionForm from "@/components/UpdateDecisionForm";

import type { Decision } from "@/types";
import type { ColumnDef } from "@tanstack/react-table";

interface FilterState {
    decisionTypes: string[];
    conditions: string[];
}

const DecisionPage = () => {
    const [openCreate, setOpenCreate] = useState(false);
    const [openEdit, setOpenEdit] = useState(false);
    const [editDecision, setEditDecision] = useState<Decision | null>(null);
    const [loading, setLoading] = useState(false);
    const [filteredData, setFilteredData] = useState<Decision[]>([]);
    const [filters, setFilters] = useState<FilterState>({
        decisionTypes: [],
        conditions: []
    });

    // Get all unique decision types from the data
    const getUniqueDecisionTypes = (data: Decision[]) => {
        const types = data.map(d => d.decision_type_name).filter(Boolean);
        return Array.from(new Set(types));
    };

    const exportFile = async () => {
        try {
            setLoading(true);
            const file = await exportDecisionExcelApi();
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

    const { data: decisions, isLoading, refetch } = useQuery({
        queryKey: ["decisions"],
        queryFn: getAllDecisionsApi,
    });

    // Apply filters when data or filters change
    useEffect(() => {
        if (decisions) {
            let filtered = [...decisions];

            // Apply decision type filter
            if (filters.decisionTypes.length > 0) {
                filtered = filtered.filter(decision =>
                    decision.decision_type_name &&
                    filters.decisionTypes.includes(decision.decision_type_name)
                );
            }

            // Apply condition filter
            if (filters.conditions.length > 0) {
                filtered = filtered.filter(decision =>
                    decision.condition &&
                    filters.conditions.includes(decision.condition)
                );
            }

            setFilteredData(filtered);
        }
    }, [decisions, filters]);

    const { mutate: deleteDecision } = useMutation({
        mutationFn: deleteDecisionApi,
        onSuccess: () => {
            refetch();
            toast.success("Xóa quyết định thành công");
        },
        onError: (error: any) => {
            toast.error(error.message);
        },
    });

    const { mutate: createDecision } = useMutation({
        mutationFn: createSampleDecisionApi,
        onSuccess: () => {
            refetch();
            toast.success("Tạo quyết định thành công");
            setOpenCreate(false);
        },
        onError: (error: any) => {
            toast.error("Lỗi khi tạo quyết định: " + error.message);
        },
    });

    const handleFilterChange = (type: keyof FilterState, value: string) => {
        setFilters(prev => {
            const currentValues = prev[type];
            if (currentValues.includes(value)) {
                return {
                    ...prev,
                    [type]: currentValues.filter(v => v !== value)
                };
            } else {
                return {
                    ...prev,
                    [type]: [...currentValues, value]
                };
            }
        });
    };

    const clearFilters = () => {
        setFilters({
            decisionTypes: [],
            conditions: []
        });
    };

    const columns: ColumnDef<Decision>[] = [
        { accessorKey: "decision_id", header: "Mã Quyết định" },
        { accessorKey: "decision_name", header: "Tên Quyết định" },
        {
            accessorKey: "decision_type_name",
            header: "Loại Quyết định",
            cell: ({ row }) => row.original.decision_type_name || "-",
        },
        {
            accessorKey: "sign_date",
            header: "Ngày ký",
            cell: ({ row }) => new Date(row.original.sign_date).toLocaleString("vi-VN"),
        },
        {
            accessorKey: "effective_date",
            header: "Hiệu lực từ ngày",
            cell: ({ row }) => new Date(row.original.effective_date).toLocaleString("vi-VN"),
        },
        { accessorKey: "condition", header: "Tình trạng" },
        {
            id: "actions",
            header: "Chỉnh sửa",
            cell: ({ row }) => {
                const decision = row.original;
                return (
                    <div className="flex gap-4">
                        <Button
                            variant="outline"
                            onClick={() => {
                                setEditDecision(decision);
                                setOpenEdit(true);
                            }}
                        >
                            <SquarePen className="w-4 h-4" />
                        </Button>
                        <ConfirmDelete deleteFn={() => deleteDecision(decision.decision_id)} />
                    </div>
                );
            },
        },
    ];

    const ButtonCreate = () => (
        <div className="button-container flex items-center justify-center space-x-4 text-[17px]">
            <CreateDecisionForm open={openCreate} setOpen={setOpenCreate} onSuccess={refetch} />
            <Button
                onClick={() => setOpenCreate(true)}
                variant="default" className={`px-8 py-5 text-[17px] rounded-[15px]`}
            >
                <span className={`mb-1 text-[24px]`}>+</span>
                Thêm quyết định
            </Button>
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
        </div>
    );

    const NavLink = () => {
        const [isFilterVisible, setIsFilterVisible] = useState(false);
        const filterRef = useRef<HTMLDivElement | null>(null);

        // Close filter when clicking outside
        useEffect(() => {
            const handleClickOutside = (event: MouseEvent) => {
                if (filterRef.current && !filterRef.current.contains(event.target as Node)) {
                    setIsFilterVisible(false);
                }
            };

            document.addEventListener('mousedown', handleClickOutside);
            return () => document.removeEventListener('mousedown', handleClickOutside);
        }, []);

        const uniqueDecisionTypes = decisions ? getUniqueDecisionTypes(decisions) : [];

        return (
            <div className="button-container border-t-2 border-gray-300 flex items-center h-[85px] justify-end mb-[10px] pt-[45px] relative">
                <div className="absolute right-0 z-30 flex">
                    <div className="flex items-center bg-gray-100 rounded-xl p-2 w-72 cursor-pointer" onClick={() => setIsFilterVisible(!isFilterVisible)}>
                        <div className="ml-2 p-2 rounded-full bg-white shadow-md flex items-center justify-center">
                            <Filter className="text-gray-500" />
                        </div>
                        <span className="ml-2 text-gray-700">Lọc</span>
                        {(filters.decisionTypes.length > 0 || filters.conditions.length > 0) && (
                            <span className="ml-2 bg-red-500 text-white rounded-full px-2 py-1 text-xs">
                {filters.decisionTypes.length + filters.conditions.length}
              </span>
                        )}
                    </div>
                </div>

                {isFilterVisible && (
                    <div ref={filterRef} className="absolute top-[85px] right-0 w-[300px] bg-white border shadow-lg p-4 rounded-lg z-50">
                        <h3 className="font-bold text-lg mb-4">Lọc</h3>

                        {/* Decision Type Filter */}
                        <div className="mb-4">
                            <h4 className="text-sm font-semibold mb-2">Loại quyết định</h4>
                            <div className="max-h-32 overflow-y-auto">
                                {uniqueDecisionTypes.map((type, i) => (
                                    <div key={i} className="flex items-center mb-1">
                                        <input
                                            type="checkbox"
                                            id={`filter-decision-type-${i}`}
                                            checked={filters.decisionTypes.includes(type)}
                                            onChange={() => handleFilterChange('decisionTypes', type)}
                                            className="mr-2"
                                        />
                                        <label htmlFor={`filter-decision-type-${i}`} className="text-sm">
                                            {type}
                                        </label>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Condition Filter */}
                        <div className="mb-4">
                            <h4 className="text-sm font-semibold mb-2">Tình trạng</h4>
                            <div>
                                {["Đang hiệu lực", "Hết hiệu lực", "Chưa hiệu lực"].map((condition, i) => (
                                    <div key={i} className="flex items-center mb-1">
                                        <input
                                            type="checkbox"
                                            id={`filter-status-${i}`}
                                            checked={filters.conditions.includes(condition)}
                                            onChange={() => handleFilterChange('conditions', condition)}
                                            className="mr-2"
                                        />
                                        <label htmlFor={`filter-status-${i}`} className="text-sm">
                                            {condition}
                                        </label>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Filter Actions */}
                        <div className="flex gap-2">
                            <button
                                onClick={clearFilters}
                                className="flex-1 bg-gray-300 text-gray-700 py-2 px-4 rounded-full hover:bg-gray-400"
                            >
                                Xóa bộ lọc
                            </button>
                            <button
                                onClick={() => setIsFilterVisible(false)}
                                className="flex-1 bg-[#DB3B21] text-white py-2 px-4 rounded-full hover:bg-[#b83a1a]"
                            >
                                Áp dụng
                            </button>
                        </div>
                    </div>
                )}
            </div>
        );
    };

    if (isLoading) return <Loading />;

    return (
        <>
            <DataTable
                columns={columns}
                buttonCreate={<ButtonCreate />}
                data={filteredData}
                navLink={<NavLink />}
                title="QUYẾT ĐỊNH"
                keyFilter="decision_id"
            />
            {editDecision && (
                <UpdateDecisionForm
                    open={openEdit}
                    setOpen={(val) => {
                        setOpenEdit(val);
                        if (!val) setEditDecision(null);
                    }}
                    Decision={{
                        decision_id: editDecision.decision_id,
                        decision_name: editDecision.decision_name,
                        effective_date: new Date(editDecision.effective_date).toISOString().slice(0, 16),
                        sign_date: new Date(editDecision.sign_date).toISOString().slice(0, 16),
                        content: editDecision.content,
                        condition: editDecision.condition,
                        attached_file: editDecision.attached_file || "",
                        created_date: new Date(editDecision.created_date).toISOString().slice(0, 16),
                        employee_id: editDecision.employee_id,
                        decision_type_id: editDecision.decision_type_id,
                        decision_type_name: editDecision.decision_type_name,
                    }}
                />
            )}
        </>
    );
};

export default DecisionPage;