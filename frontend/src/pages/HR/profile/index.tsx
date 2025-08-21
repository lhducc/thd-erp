import {useMutation, useQueryClient} from "@tanstack/react-query";
import {toast} from "sonner";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";
import {changeStatusEmployeeApi, deleteEmployeeApi} from "@/apis/profile.api.ts";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import DataTable from "@/components/DataTable.tsx";
import Loading from "@/components/Loading.tsx";
import {lazy, Suspense, useEffect, useRef, useState} from "react";
import CreateEmployeeForm from "@/components/CreateEmployeeForm.tsx";
import type {Employee} from "@/types";
import type {ColumnDef} from "@tanstack/react-table";
import {exportEmployeeExcelApi} from '@/apis/profile.api.ts';
import {Link} from "react-router-dom";
import {useGetAllEmployee} from "@/query/employee.query.ts";

const EditEmployeeForm = lazy(() => import("@/components/EditEmployeeForm"));
interface NavLinkProps {
    selectedFilters: Record<string, string[]>;
    setSelectedFilters: React.Dispatch<React.SetStateAction<Record<string, string[]>>>;
    activeTab: string;
    setActiveTab: React.Dispatch<React.SetStateAction<string>>;
}
interface FilterOption {
    id: string;
    label: string;
    field: keyof Employee;
    options: { value: string; label: string }[];
}
const HRProfilePage = () => {
    const [open, setOpen] = useState(false);
    const [editEmployee, setEditEmployee] = useState<Employee | null>(null);
    const [activeTab, setActiveTab] = useState("hoat-dong");
    const [selectedFilters, setSelectedFilters] = useState<Record<string, string[]>>({});
    const {data: employees, isLoading: pendingGetEmployees, refetch: refetchEmployee} = useGetAllEmployee()
    const {mutate: deleteEmployee} = useMutation({
        mutationFn: deleteEmployeeApi,
        onSuccess: async () => {
            await refetchEmployee();
            toast.success("Xóa nhân viên thành công");
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const queryClient = useQueryClient();

    const { mutate: changeStatus } = useMutation({
        mutationFn: ({ id, status }: { id: string; status: "active" | "inactive" }) =>
            changeStatusEmployeeApi(id, status),
        // Optimistic update
        onMutate: async ({ id, status }) => {
            await queryClient.cancelQueries({queryKey: ["employees"]}); // dừng refetch

            // Lấy snapshot dữ liệu cũ
            const previousEmployees = queryClient.getQueryData<Employee[]>(["employees"]);

            // Cập nhật tạm thời trong cache
            queryClient.setQueryData<Employee[]>(["employees"], old =>
                old?.map(emp =>
                    emp.employee_id === id ? { ...emp, status } : emp
                ) || []
            );

            return { previousEmployees }; // để rollback nếu lỗi
        },
        onError: (err, _, context) => {
            if (context?.previousEmployees) {
                queryClient.setQueryData(["employees"], context.previousEmployees); // rollback
            }
            toast.error("Lỗi cập nhật trạng thái");
        },
        onSuccess: () => {
            toast.success("Cập nhật trạng thái thành công");
        },
        onSettled: () => {
            queryClient.invalidateQueries({queryKey: ["employees"]}); // refetch lại cho chắc
        },
    });


    const columns: ColumnDef<Employee>[] = [
        {accessorKey: "employee_id", header: "Mã nhân viên"},
        {
            accessorKey: "full_name",
            header: "Họ và tên",
            cell: ({row}: { row: any }) => {
                return (
                    <Link to={`/detail_infor/${row.original.employee_id}`}>
                        {row.original.full_name}
                    </Link>
                )
            }
        },
        {
            accessorKey: "department",
            header: "Phòng ban",
            cell: ({row}: { row: any }) => row.original.department?.department_name || "N/A"
        },
        {
            accessorKey: "job_title_id",
            header: "Chức vụ",
            cell: ({row}: { row: any }) => row.original.job_title?.job_title || "N/A"
        },
        {
            accessorKey: "position_id",
            header: "Vị trí",
            cell: ({row}: { row: any }) => row.original.position?.position_name || "N/A"
        },
        {
            accessorKey: "office",
            header: "Văn phòng",
            cell: ({row}: { row: any }) => row.original.department?.office?.office_name || "N/A"
        },
        {
            accessorKey: "phone_number",
            header: "Số điện thoại"
        },
        {
            accessorKey: "manager_id",
            header: "Quản lý trực tiếp",
            cell: ({row}: { row: any }) => row.original.manager?.full_name || "N/A"
        },
        {
            accessorKey: "email",
            header: "Email"
        },
        {
            accessorKey: "status",
            header: "Trạng thái",
            cell: ({ row }: { row: any }) => {
                const status = row.original.status.toLowerCase() as "active" | "inactive";
                const nextStatus = status === "active" ? "inactive" : "active";

                return (
                    <button
                        onClick={() => changeStatus({ id: row.original.employee_id, status: nextStatus })}
                        className={`px-2 py-1 rounded-full text-xs cursor-pointer ${
                            status === "active"
                                ? "bg-green-100 text-green-800"
                                : "bg-red-100 text-red-800"
                        }`}
                    >
                        {status === "active" ? "Hoạt động" : "Vô hiệu hóa"}
                    </button>
                );
            },
        },
        {
            id: "actions",
            header: "Chỉnh sửa",
            cell: ({row}: { row: any }) => {
                const employee = row.original;
                return (
                    <div className="flex gap-4">
                        <Button variant={"outline"} onClick={() => setEditEmployee(employee)}>
                            <SquarePen/>
                        </Button>
                        <ConfirmDelete deleteFn={() => deleteEmployee(employee.employee_id)}/>
                    </div>
                );
            },
        }
    ];

    const ButtonCreate = () => {
        const [isDropdownVisible, setIsDropdownVisible] = useState(false);
        const dropdownRef = useRef<HTMLDivElement | null>(null);
        const [loading, setLoading] = useState(false);

        const toggleDropdown = () => {
            setIsDropdownVisible(!isDropdownVisible);
        };

        const exportFile = async () => {
            try {
                setLoading(true);
                const file = await exportEmployeeExcelApi();
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

        useEffect(() => {
            const handleClickOutside = (event: MouseEvent) => {
                if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                    setIsDropdownVisible(false);
                }
            };

            document.addEventListener("click", handleClickOutside);
            return () => document.removeEventListener("click", handleClickOutside);
        }, []);

        return (
            <div className="button-container flex items-center justify-center space-x-4 text-[17px]">
                <CreateEmployeeForm open={open} setOpen={setOpen}
                                    refetchEmployee={refetchEmployee}/>
                <div ref={dropdownRef} className="dropdown relative">
                    <Button
                        onClick={toggleDropdown}
                        className={`w-[200px] py-5 text-[17px] rounded-[15px]`}
                    >
                        <span className={`mb-1 text-[24px]`}>+</span> Thêm nhân viên
                    </Button>
                    {isDropdownVisible && (
                        <div className="dropdown-menu absolute left-0 w-[200px] text-[13px] bg-white shadow-lg rounded-b-2xl z-10">
                            {["Nhân viên"].map((label, index) => (
                                <button
                                    key={index}
                                    onClick={() => setOpen(true)}
                                    className="dropdown-item text-left p-3 w-[200px] border-b shadow-2xl border-gray-500 cursor-pointer hover:bg-[#f2f2f2] rounded-b-xl transition duration-200"
                                >
                                    <label className="text-[20px]">+ </label> {label}
                                </button>
                            ))}
                        </div>
                    )}
                </div>
                <Button
                    onClick={exportFile}
                    variant={"default"}
                    className={`px-10 py-5 rounded-[15px] text-[17px]`}
                    disabled={loading}
                >
                    {loading ? "Đang xuất..." : "Xuất file"}
                </Button>
            </div>
        );
    };

    const NavLink = ({ selectedFilters, setSelectedFilters, activeTab, setActiveTab }: NavLinkProps) => {
        const [isFilterVisible, setIsFilterVisible] = useState(false);
        const filterRef = useRef<HTMLDivElement | null>(null);
        const [filterOptions, setFilterOptions] = useState<FilterOption[]>([]);

        // Khởi tạo filter options từ dữ liệu employees
        useEffect(() => {
            if (employees && employees.length > 0) {
                const options: FilterOption[] = [
                    {
                        id: 'office',
                        label: 'Văn phòng',
                        field: 'office',
                        options: Array.from(new Set(employees
                            .map(e => e.department?.office?.office_name)
                            .filter(Boolean)))
                            .map(value => ({ value: value as string, label: value as string }))
                    },
                    {
                        id: 'department',
                        label: 'Phòng ban',
                        field: 'department',
                        options: Array.from(new Set(employees
                            .map(e => e.department?.department_name)
                            .filter(Boolean)))
                            .map(value => ({ value: value as string, label: value as string }))
                    },
                    {
                        id: 'job_title',
                        label: 'Chức vụ',
                        field: 'job_title',
                        options: Array.from(new Set(employees
                            .map(e => e.job_title?.job_title)
                            .filter(Boolean)))
                            .map(value => ({ value: value as string, label: value as string }))
                    },
                    {
                        id: 'position',
                        label: 'Vị trí',
                        field: 'position',
                        options: Array.from(new Set(employees
                            .map(e => e.position?.position_name)
                            .filter(Boolean)))
                            .map(value => ({ value: value as string, label: value as string }))
                    }
                ];
                setFilterOptions(options);
            }
        }, [employees]);

        const handleFilterClick = () => {
            setIsFilterVisible(!isFilterVisible);
        };

        const handleFilterChange = (filterId: string, value: string, isChecked: boolean) => {
            setSelectedFilters(prev => {
                const newFilters = { ...prev };
                if (isChecked) {
                    newFilters[filterId] = [...(newFilters[filterId] || []), value];
                } else {
                    newFilters[filterId] = (newFilters[filterId] || []).filter(v => v !== value);
                }
                return newFilters;
            });
        };

        const applyFilters = () => {
            setIsFilterVisible(false);
        };

        return (
            <div className="button-container flex items-center justify-right mb-[10px] mt-[50px] relative">
                <label
                    onClick={() => setActiveTab("hoat-dong")}
                    className={`cursor-pointer text-lg pb-[10px] w-[120px] ${
                        activeTab === "hoat-dong"
                            ? "font-bold border-b-2 border-[#DB3B21] text-[#DB3B21]"
                            : "text-gray-300 border-b-2 border-gray-300"
                    }`}
                >
                    Hoạt động
                </label>
                <label
                    onClick={() => setActiveTab("vo-hieu-hoa")}
                    className={`cursor-pointer text-lg pb-[10px] w-[189px] text-center ${
                        activeTab === "vo-hieu-hoa"
                            ? "font-bold border-b-2 border-[#DB3B21] text-[#DB3B21]"
                            : "text-gray-300 border-b-2 border-gray-300"
                    }`}
                >
                    Đã vô hiệu hóa
                </label>
                <div className="absolute right-0 z-30">
                    <img
                        className="z-50 cursor-pointer"
                        src="/control.png"
                        alt="Control"
                        onClick={handleFilterClick}
                    />
                </div>
                {isFilterVisible && (
                    <div ref={filterRef}
                         className="filter-menu z-1 absolute top-[-20px] right-[-13px] mt-2 w-[250px] bg-white border shadow-lg p-4 rounded-lg">
                        <h3 className="font-bold text-lg mb-4">Lọc</h3>
                        {filterOptions.map((filter) => (
                            <div key={filter.id} className="mb-4">
                                <h4 className="text-sm font-semibold">{filter.label}</h4>
                                <div className="max-h-40 overflow-y-auto">
                                    {filter.options.map((option, idx) => (
                                        <div key={idx} className="flex items-center">
                                            <input
                                                type="checkbox"
                                                id={`${filter.id}-${option.value}`}
                                                checked={selectedFilters[filter.id]?.includes(option.value) || false}
                                                onChange={(e) => handleFilterChange(filter.id, option.value, e.target.checked)}
                                                className="mr-2"
                                            />
                                            <label htmlFor={`${filter.id}-${option.value}`}>{option.label}</label>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ))}
                        <button
                            className="btn-filter bg-[#DB3B21] text-white py-2 px-4 rounded-full w-full"
                            onClick={applyFilters}
                        >
                            Lọc
                        </button>
                    </div>
                )}
            </div>
        );
    };

    // Filter employees based on active tab và các filter được chọn
    const filteredEmployees = employees?.filter((employee: Employee) => {
        // Filter theo trạng thái
        const statusMatch = activeTab === "hoat-dong"
            ? employee.status?.toLowerCase() === "active"
            : employee.status?.toLowerCase() !== "active";

        // Filter theo các điều kiện khác
        const filtersMatch = Object.entries(selectedFilters).every(([field, values]) => {
            if (!values.length) return true;

            switch (field) {
                case 'office':
                    return values.includes(employee.department?.office?.office_name || '');
                case 'department':
                    return values.includes(employee.department?.department_name || '');
                case 'job_title':
                    return values.includes(employee.job_title?.job_title || '');
                case 'position':
                    return values.includes(employee.position?.position_name || '');
                default:
                    return true;
            }
        });

        return statusMatch && filtersMatch;
    }) || [];

    // Truyền setSelectedFilters xuống NavLink component
    const navLink = <NavLink selectedFilters={selectedFilters} setSelectedFilters={setSelectedFilters} setActiveTab={setActiveTab} activeTab={activeTab} />;

    if (pendingGetEmployees) return <Loading/>;
    return (
        <>
            <DataTable
                columns={columns}
                buttonCreate={<ButtonCreate/>}
                data={filteredEmployees}
                navLink={navLink}
                title="Hồ sơ nhân viên"
                keyFilter="employee_id"
            />
            {editEmployee && (
                <Suspense fallback={<Loading />} >
                    <EditEmployeeForm
                        open={Boolean(editEmployee)}
                        setOpen={setEditEmployee}
                        data={editEmployee}
                        refetchEmployee={refetchEmployee}
                    />
                </Suspense>
            )}
        </>
    );
};

export default HRProfilePage;