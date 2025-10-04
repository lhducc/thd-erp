import type { ColumnDef } from "@tanstack/react-table";
import more from "../../../assets/more.svg";
import DataTable from "@/components/DataTable.tsx";
import { useGetAllEmployee } from "@/query/employee.query.ts";
import type { Employee } from "@/types/employee.ts";
import { Link } from "react-router-dom";
import NavLinkAttendantHistory from "@/components/ui/NavLinkAttendantHistory";

const ManagementAttendantHistory = () => {
  const { data: employees, isLoading: pendingGetEmployees } =
    useGetAllEmployee();
  const columns: ColumnDef<Employee>[] = [
    {
      accessorKey: "employee_id",
      header: "Mã NV",
    },
    {
      accessorKey: "full_name",
      header: "Tên nhân viên",
    },
    {
      accessorKey: "department.office.office_name",
      header: "Văn phòng",
    },
    {
      accessorKey: "department.department_name",
      header: "Phòng ban",
    },
    {
      accessorKey: "job_title.job_title",
      header: "Chức danh",
    },
    {
      accessorKey: "position.position_name",
      header: "Cấp bậc",
    },
    {
      accessorKey: "is_gps",
      header: "Xem chi tiết",
      cell: ({ row }) => {
        const data = row.original;
        return (
          <div className={`flex justify-center items-center`}>
            <Link to={`${data.employee_id}`}>
              <img className={`w-[25px] h-[25px]`} src={more} alt="" />
            </Link>
          </div>
        );
      },
    },
  ];

  return (
    <>
      <DataTable
        columns={columns}
        data={employees || []}
        isLoading={pendingGetEmployees}
        navLink={<NavLinkAttendantHistory />}
        title="Bảng lịch sử chấm công"
        keyFilter="full_name"
      />
    </>
  );
};

export default ManagementAttendantHistory;
