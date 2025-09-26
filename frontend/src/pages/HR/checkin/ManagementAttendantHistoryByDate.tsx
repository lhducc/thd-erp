import type { ColumnDef } from "@tanstack/react-table";
import DataTable from "@/components/DataTable.tsx";
import NavLinkAttendantHistory from "@/components/ui/NavLinkAttendantHistory";
import {useEmployeesByDate } from "@/query/attendance-history-by-date.ts";
import type { AttendanceRecordHistoryByDate } from "@/types/attendance.ts";
import ExportFileDialog from "@/components/CreateExcelFileForm.tsx";
import { Input } from "@/components/ui/input";
import { useState } from "react";

const ManagementAttendantHistory = () => {
  const today = new Date().toLocaleDateString("sv-SE", { timeZone: "Asia/Ho_Chi_Minh" });
  const [selectedDate, setSelectedDate] = useState<string>(today);

  const { data: records, isLoading } = useEmployeesByDate(selectedDate);  

  const columns: ColumnDef<AttendanceRecordHistoryByDate>[] = [
    {
      accessorKey: "employee_id",
      header: "Mã NV",
    },
    {
      accessorKey: "full_name",
      header: "Tên nhân viên",
    },
    {
      accessorKey: "office_name",
      header: "Văn phòng",
    },
    {
      accessorKey: "department_name",
      header: "Phòng ban",
    },
    {
      accessorKey: "timestamp",
      header: "Thời gian chấm công",
      cell: ({ row }) => {
        const record = row.original;
        const ts = record.timestamp;
        const startTime = record.start_time;
        
        if (!ts || !startTime) return "-";

        const display = ts.replace("T", " ").replace("Z", "");

        // Lấy phần thời gian từ timestamp (HH:MM:SS)
        const timestampTime = ts.slice(11, 19); // "08:29:45"
        // Lấy start_time (HH:MM:SS hoặc HH:MM)
        const shiftStartTime = startTime.length === 5 ? startTime + ":00" : startTime; // Đảm bảo cùng format

        // So sánh thời gian
        const isLate = timestampTime > shiftStartTime;

        return (
          <span style={{ color: isLate ? "red" : "inherit", fontWeight: isLate ? "600" : "normal" }}>
            {display}
          </span>
        );
      },
    },

  ];

  return (
    <div className="space-y-4">
      <DataTable
        columns={columns}
        data={records || []}
        isLoading={isLoading}
        navLink={<NavLinkAttendantHistory />}
        title="Dữ liệu theo ngày"
        keyFilter="full_name"
        buttonCreate={
          <div className="flex items-center gap-4">
            <Input
              type="date"
              value={selectedDate}
              onChange={(e) => setSelectedDate(e.target.value)}
            />
            <ExportFileDialog selectedDate={selectedDate}/>
          </div>
        }
      />
    </div>
  );
};

export default ManagementAttendantHistory;
