import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { SquarePen } from "lucide-react";
import { getAllEmployeesApi, deleteEmployeeApi } from "@/apis/profile.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import { useEffect, useRef, useState } from "react";
import CreateEmployeeForm from "@/components/CreateEmployeeForm";
import EditEmployeeForm from "@/components/EditEmployeeForm";
import type { Employee } from "@/types";
import type { ColumnDef } from "@tanstack/react-table";
import { exportEmployeeExcelApi } from '@/apis/profile.api';

const EmployeePage = () => {
  const [open, setOpen] = useState(false);
  const [editEmployee, setEditEmployee] = useState<Employee | null>(null);

  const { data: employees, isLoading: pendingGetEmployees, refetch: refetchEmployee } = useQuery({
    queryKey: ["employees"],
    queryFn: getAllEmployeesApi,
  });

  const { mutate: deleteEmployee } = useMutation({
    mutationFn: deleteEmployeeApi,
    onSuccess: () => {
      refetchEmployee();
      toast.success("Xóa nhân viên thành công");
    },
    onError: (error: any) => {
      toast.error(error.message);
    },
  });

  const columns: ColumnDef<Employee>[] = [
    { accessorKey: "employee_id", header: "Mã nhân viên" },
    { accessorKey: "full_name", header: "Họ và tên" },
    {
      accessorKey: "department",
      header: "Phòng ban",
      cell: ({ row }: { row: any }) => row.original.department?.department_name || "N/A"
    },
    { accessorKey: "job_title_id", header: "Chức vụ" },
    { accessorKey: "position_id", header: "Vị trí" },
    {
      accessorKey: "office",
      header: "Văn phòng",
      cell: ({ row }: { row: any }) => row.original.department?.office?.office_name || "N/A"
    },
    { accessorKey: "phone_number", header: "Số điện thoại" },
    { accessorKey: "manager_id", header: "Quản lý trực tiếp" },
    { accessorKey: "email", header: "Email" },
    {
      id: "actions",
      header: "Chỉnh sửa",
      cell: ({ row }: { row: any }) => {
        const employee = row.original;
        return (
          <div className="flex gap-4">
            <Button variant={"outline"} onClick={() => setEditEmployee(employee)}>
              <SquarePen />
            </Button>
            <ConfirmDelete deleteFn={() => deleteEmployee(employee.employee_id)} />
          </div>
        );
      },
    },
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
        <CreateEmployeeForm className="w-[1261px]" open={open} setOpen={setOpen} refetchEmployee={refetchEmployee}/>
        <div ref={dropdownRef} className="dropdown relative">
          <Button
            onClick={toggleDropdown}
              className={`w-[200px] py-5 text-[17px] rounded-[15px]`
            }
           >
            <span className={`mb-1 text-[24px]`}>+</span> Thêm nhân viên
          </Button>
          {isDropdownVisible && (
            <div className="dropdown-menu absolute left-0 w-[200px] text-[13px] bg-white shadow-lg rounded-b-2xl z-10">
              {["Nhân viên chính thức", "Thực tập sinh", "Cộng tác viên"].map((label, index) => (
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
          {loading ? "Đang xuất..." : (
            <>
              Xuất file
            </>
          )}
        </Button>
      </div>
    );
  };

  const NavLink = () => {
    const [activeTab, setActiveTab] = useState("hoat-dong");
    const [isFilterVisible, setIsFilterVisible] = useState(false);
    const filterRef = useRef<HTMLDivElement | null>(null);

    const handleTabClick = (tab: string) => setActiveTab(tab);

    const handleFilterClick = () => {
      setIsFilterVisible(!isFilterVisible);
    };

    return (
      <div className="button-container flex items-center justify-right mb-[10px] mt-[50px] relative">
        <label
          onClick={() => handleTabClick("hoat-dong")}
          className={`cursor-pointer text-lg pb-[10px] w-[120px] ${activeTab === "hoat-dong" ? "font-bold border-b-2 border-[#DB3B21] text-[#DB3B21]" : "text-gray-300 border-b-2 border-gray-300"}`}
        >
          Hoạt động
        </label>
        <label
          onClick={() => handleTabClick("vo-hieu-hoa")}
          className={`cursor-pointer text-lg pb-[10px] w-[189px] text-center ${activeTab === "vo-hieu-hoa" ? "font-bold border-b-2 border-[#DB3B21] text-[#DB3B21]" : "text-gray-300 border-b-2 border-gray-300"}`}
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
          <div ref={filterRef} className="filter-menu z-1 absolute top-[-20px] right-[-13px] mt-2 w-[250px] bg-white border shadow-lg p-4 rounded-lg">
            <h3 className="font-bold text-lg mb-4">Lọc</h3>
            {/* Filter options */}
            {["Văn phòng", "Phòng ban", "Chức vụ", "Vị trí"].map((filter, index) => (
              <div key={index} className="mb-4">
                <h4 className="text-sm font-semibold">{filter}</h4>
                <div>
                  <input type="checkbox" id={`filter-${index}-option1`} /> <label htmlFor={`filter-${index}-option1`}>Option 1</label>
                  <br />
                  <input type="checkbox" id={`filter-${index}-option2`} /> <label htmlFor={`filter-${index}-option2`}>Option 2</label>
                </div>
              </div>
            ))}
            <button className="btn-filter bg-[#DB3B21] text-white py-2 px-4 rounded-full w-full">Lọc</button>
          </div>
        )}
      </div>
    );
  };


  if (pendingGetEmployees) return <Loading />;

  return (
    <>
      <DataTable
        columns={columns}
        buttonCreate={<ButtonCreate />}
        data={employees || []}
        navLink={<NavLink />}
        title="Hồ sơ nhân viên"
        keyFilter="employee_id"
      />
      {editEmployee && (
        <EditEmployeeForm open={Boolean(editEmployee)} setOpen={setEditEmployee} data={editEmployee} refetchEmployee={refetchEmployee}/>
      )}
    </>
  );
};

export default EmployeePage;
