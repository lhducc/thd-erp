import type { ColumnDef } from "@tanstack/react-table";
import DataTable from "@/components/DataTable.tsx";
import {useGetManagerEmployees } from "@/query/employee.query.ts";
import type { ManagerEmployee } from "@/types/employee.ts";
import { useAuth } from "@/context/AuthContext";

const ManagemenrEmployeePage = () => {

  const { data: employees, isLoading: pendingGetEmployees } = useGetManagerEmployees();


  const columns: ColumnDef<ManagerEmployee>[] = [
    { accessorKey: "employee_id", header: "Mã NV" },
    { accessorKey: "full_name", header: "Tên nhân viên" },
    { accessorKey: "department_name", header: "Phòng ban" },
    { accessorKey: "job_title", header: "Chức vụ" },
    { accessorKey: "position_name", header: "Vị trí" },
    { accessorKey: "office_name", header: "Văn phòng" },
    { accessorKey: "phone_number", header: "Số điện thoại" },
    { accessorKey: "email", header: "Email" },
  ];

  return (
    <>
      <DataTable
        columns={columns}
        data={employees || []}
        isLoading={pendingGetEmployees}
        title="Hồ sơ nhân viên"
        keyFilter="full_name"
      />
    </>
  );
};

export default ManagemenrEmployeePage;
