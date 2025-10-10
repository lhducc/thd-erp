import type { ColumnDef } from "@tanstack/react-table";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import type { Employee } from "@/types";
import { Link } from "react-router-dom";
import { useGetAllEmployeesActive } from "@/query/employee.query.ts";
import { totalRequestAttendanceApi } from "@/apis/attendance-approval.api.ts";
import { useEffect, useState } from "react";

const ApproveAttendanceOverviewPage = () => {
  const {
    data: employees,
    isLoading: pendingGetEmployees,
    refetch: refetchEmployee,
  } = useGetAllEmployeesActive();
  const TotalRequestsCell = ({ employeeId }: { employeeId: string }) => {
    const [count, setCount] = useState<number | null>(null);

    useEffect(() => {
      let mounted = true;
      totalRequestAttendanceApi(employeeId).then((res) => {
        if (mounted) setCount(res);
      });
      return () => {
        mounted = false;
      };
    }, [employeeId]);

    if (count === null) return <span>Loading...</span>;
    return <span>{count}</span>;
  };

  const columns: ColumnDef<Employee>[] = [
    { accessorKey: "employee_id", header: "Mã nhân viên" },
    {
      accessorKey: "full_name",
      header: "Họ và tên",
      cell: ({ row }: { row: any }) => {
        return (
          <Link to={`/approve/${row.original.employee_id}`}>
            {row.original.full_name}
          </Link>
        );
      },
    },
    {
      accessorKey: "department",
      header: "Phòng ban",
      cell: ({ row }: { row: any }) =>
        row.original.department?.department_name || "N/A",
    },
    {
      accessorKey: "job_title_id",
      header: "Chức vụ",
      cell: ({ row }: { row: any }) =>
        row.original.job_title?.job_title || "N/A",
    },
    {
      accessorKey: "office",
      header: "Văn phòng",
      cell: ({ row }: { row: any }) =>
        row.original.department?.office?.office_name || "N/A",
    },
    {
      id: "request",
      header: "Tổng số yêu cầu trong tháng",
      cell: ({ row }) => (
        <TotalRequestsCell employeeId={row.original.employee_id} />
      ),
    },
  ];

  if (pendingGetEmployees) {
    return <Loading />;
  }

  return (
    <>
      <DataTable
        columns={columns}
        data={employees || []}
        title="Phê duyệt chấm công"
        keyFilter="full_name"
      />
    </>
  );
};

export default ApproveAttendanceOverviewPage;
