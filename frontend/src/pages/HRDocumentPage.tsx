import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { SquarePen } from "lucide-react";
import { getAllEmployeesApi, deleteEmployeeApi } from "@/apis/profile.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import { useState } from "react";
import CreateEmployeeForm from "@/components/CreateEmployeeForm";
import EditEmployeeForm from "@/components/EditEmployeeForm";
import type { Employee } from "@/types";
import type { ColumnDef } from "@tanstack/react-table";

const HRDocumentPage = () => {
  const [open, setOpen] = useState(false);  // State for create employee form
  const [editEmployee, setEditEmployee] = useState<Employee | null>(null);  // State for editing employee

  const {
    data: employees,
    isLoading: pendingGetEmployees,
    refetch: refetchEmployee,
    isError: error,
    error: queryError,
  } = useQuery({
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
    { accessorKey: "department", header: "Phòng ban" },
    { accessorKey: "job_title_id", header: "Chức vụ" },
    { accessorKey: "position_id", header: "Vị trí" },
    { accessorKey: "phone_number", header: "Số điện thoại" },
    { accessorKey: "manager_id", header: "Quản lý trực tiếp" },
    { accessorKey: "office", header: "Văn phòng" },
    { accessorKey: "email", header: "Email" },
    {
      id: "actions",
      header: "Chỉnh sửa",
      cell: ({ row }: { row: any }) => {
        const employee = row.original;
        return (
          <div className="flex gap-4">
            <Button
              variant={"outline"}
              onClick={() => setEditEmployee(employee)} // Set employee data to edit
            >
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
    const handleAddEmployeeClick = () => {
      setIsDropdownVisible(!isDropdownVisible);
    };
    const handleExportClick = () => {
      alert("Exporting file...");
    };

    return (
      <div className="button-container flex items-center justify-center space-x-4 text-[17px]">
        <CreateEmployeeForm className="w-[1261px]" open={open} setOpen={setOpen}></CreateEmployeeForm>
        <div className="dropdown relative">
          <button
            onClick={handleAddEmployeeClick}
            className="btn-add-employee text-left bg-[#DB3B21] w-[200px] text-white p-3 rounded-2xl flex justify-center items-center hover:bg-[#b83a1a] transition duration-200"
          >
            <label htmlFor="" className="text-[20px]">+ </label> Thêm nhân viên
          </button>

          {isDropdownVisible && (
            <div className="dropdown-menu absolute left-0 w-[200px] text-[13px] bg-white shadow-lg rounded-b-2xl z-10">
              <button onClick={() => { setOpen(true); }} className="dropdown-item text-left p-3 w-[200px] border-b shadow-2xl border-gray-500 cursor-pointer hover:bg-[#f2f2f2] rounded-b-xl transition duration-200">
                <label htmlFor="" className="text-[20px]">+ </label> Nhân viên chính thức
              </button>
              <button onClick={() => { setOpen(true); }} className="dropdown-item text-left p-3 w-[200px] border-b shadow-2xl border-gray-500 cursor-pointer hover:bg-[#f2f2f2] rounded-b-xl transition duration-200">
                <label htmlFor="" className="text-[20px]">+ </label> Thực tập sinh
              </button>
              <button onClick={() => { setOpen(true); }} className="dropdown-item text-left p-3 w-[200px] border-b shadow-2xl border-gray-500 cursor-pointer hover:bg-[#f2f2f2] rounded-b-xl transition duration-200">
                <label htmlFor="" className="text-[20px]">+ </label> Cộng tác viên
              </button>
            </div>
          )}
        </div>
        <button
          className="btn-export text-left h-[54px] bg-[#DB3B21] w-[192px] text-white p-3 rounded-2xl flex justify-center items-center hover:bg-[#b83a1a] transition duration-200"
          onClick={handleExportClick}
        >
          <img className="absolute right-0" src="/Vector.png" alt="Control" /> Xuất file
        </button>
      </div>
    );
  };

  const NavLink = () => {
    const [activeTab, setActiveTab] = useState('hoat-dong');
    const [isFilterVisible, setIsFilterVisible] = useState(false);

    const handleTabClick = (tab: string) => {
      setActiveTab(tab);
    };

    const handleFilterClick = () => {
      setIsFilterVisible(!isFilterVisible);
    };

    return (
      <div className="button-container flex items-center justify-right mb-[10px] mt-[50px] relative">
        <label
          onClick={() => handleTabClick('hoat-dong')}
          className={`cursor-pointer text-lg pb-[10px] w-[120px] ${activeTab === 'hoat-dong' ? 'font-bold border-b-2 border-[#DB3B21] text-[#DB3B21]' : 'text-gray-300 border-b-2 border-gray-300'}`}
        >
          Hoạt động
        </label>
        <label
          onClick={() => handleTabClick('vo-hieu-hoa')}
          className={`cursor-pointer text-lg pb-[10px] w-[189px] text-center ${activeTab === 'vo-hieu-hoa' ? 'font-bold border-b-2 border-[#DB3B21] text-[#DB3B21]' : 'text-gray-300 border-b-2 border-gray-300'}`}
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
          <div className="filter-menu z-1 absolute top-[-20px] right-[-13px] mt-2 w-[250px] bg-white border shadow-lg p-4 rounded-lg">
            <h3 className="font-bold text-lg mb-4">Lọc</h3>
            {/* Filter options */}
            <div className="mb-4">
              <h4 className="text-sm font-semibold">Văn phòng</h4>
              <div>
                <input type="checkbox" id="ct-can-tho" /> <label htmlFor="ct-can-tho">Cần Thơ - CT</label>
                <br />
                <input type="checkbox" id="ct-hcm" /> <label htmlFor="ct-hcm">Cần Thơ - HCM</label>
              </div>
            </div>
            <div className="mb-4">
              <h4 className="text-sm font-semibold">Phòng ban</h4>
              <div>
                <input type="checkbox" id="hcns" /> <label htmlFor="hcns">HCNS</label>
                <br />
                <input type="checkbox" id="research" /> <label htmlFor="research">Nghiên cứu & Phát triển</label>
              </div>
            </div>
            <div className="mb-4">
              <h4 className="text-sm font-semibold">Chức vụ</h4>
              <div>
                <input type="checkbox" id="contract3" /> <label htmlFor="contract3">Hợp đồng 3 tháng</label>
                <br />
                <input type="checkbox" id="contract12" /> <label htmlFor="contract12">Hợp đồng 12 tháng</label>
              </div>
            </div>

            {/* Phần lọc "Vị trí" */}
            <div className="mb-4">
              <h4 className="text-sm font-semibold">Vị trí</h4>
              <div>
                <input type="checkbox" id="active" /> <label htmlFor="active">Đang hiệu lực</label>
                <br />
                <input type="checkbox" id="inactive" /> <label htmlFor="inactive">Hết hiệu lực</label>
              </div>
            </div>

            {/* Nút lọc */}
            <button className="btn-filter bg-[#DB3B21] text-white py-2 px-4 rounded-full w-full">
              Lọc
            </button>
          </div>
        )}
      </div>
    );
  };

  if (pendingGetEmployees) {
    return <Loading />;
  }

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
        <EditEmployeeForm
          open={Boolean(editEmployee)}
          setOpen={setEditEmployee}
          data={editEmployee}
        />
      )}
    </>
  );
};

export default HRDocumentPage;
