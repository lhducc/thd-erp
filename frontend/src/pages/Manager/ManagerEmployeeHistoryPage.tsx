import { useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import DataTable from "@/components/DataTable";
import { Input } from "@/components/ui/input";
import type { AttendanceRecordHistoryByDate } from "@/types/attendance";
import { useEmployeeHistoryByManager } from "@/query/attendance-history-by-date";

const ManagerEmployeeHistoryPage = () => {
  const today = new Date().toLocaleDateString("sv-SE", {
    timeZone: "Asia/Ho_Chi_Minh",
  });

  const [selectedDate, setSelectedDate] = useState<string>(today);

  const { data: records, isLoading } =
    useEmployeeHistoryByManager(selectedDate);

  const columns: ColumnDef<AttendanceRecordHistoryByDate>[] = [
    { accessorKey: "employee_id", header: "Mã NV" },
    { accessorKey: "full_name", header: "Tên nhân viên" },
    { accessorKey: "office_name", header: "Văn phòng" },
    { accessorKey: "department_name", header: "Phòng ban" },
    {
      accessorKey: "timestamp",
      header: "Thời gian chấm công",
      cell: ({ row }) => {
        const record = row.original;
        const ts = record.timestamp;
        const startTime = record.start_time;

        if (!ts || !startTime) return "-";

        const date = new Date(ts);
        const display = date
          .toLocaleString("sv-SE", {
            timeZone: "Asia/Ho_Chi_Minh",
            year: "numeric",
            month: "2-digit",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit",
            second: "2-digit",
          })
          .replace("T", " ");

        const [dayPart] = display.split(" ");
        const shiftStartTime =
          startTime.length === 5 ? startTime + ":00" : startTime;
        const shiftStartDate = new Date(`${dayPart}T${shiftStartTime}+07:00`);
        const shiftStartWithGrace = new Date(
          shiftStartDate.getTime() + 3 * 60 * 1000
        );

        const isLate = date > shiftStartWithGrace;

        return (
          <span
            style={{
              color: isLate ? "red" : "inherit",
              fontWeight: isLate ? "600" : "normal",
            }}
          >
            {display}
          </span>
        );
      },
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={records || []}
      isLoading={isLoading}
      title="Dữ liệu theo ngày"
      keyFilter="full_name"
      buttonCreate={
        <Input
          type="date"
          value={selectedDate}
          onChange={(e) => setSelectedDate(e.target.value)}
          className="w-[230px]"
        />
      }
    />
  );
};

export default ManagerEmployeeHistoryPage;
