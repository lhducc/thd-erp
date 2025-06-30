'use client';

import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { SquarePen } from "lucide-react";
import { useState } from "react";
import export_file from "@/assets/export-file.svg";
import {
  getAllDecisionsApi,
  deleteDecisionApi,
  createSampleDecisionApi,
  exportDecisionExcelApi,
} from "@/apis/decision.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import DataTable from "@/components/ShiftTable";
import Loading from "@/components/Loading";
import CreateDecisionForm from "@/components/CreateDecisionForm";
import UpdateDecisionForm from "@/components/UpdateDecisionForm";
import type { ColumnDef } from "@tanstack/react-table";

const weekdays = ["T2", "T3", "T4", "T5", "T6", "T7", "CN"];

export interface Decision {
  decision_id: string;
  employee_id: string;
  employee_name: string;
  department: string;
  expected_workdays: number;
  sign_date: string;
  effective_date: string;
  created_date: string;
  content: string;
  condition: string;
  attached_file?: string;
  decision_type_id: string;
  decision_type_name: string;
  shift_days?: Record<string, string>;
}

const mockDecisions: Decision[] = [
  {
    decision_id: "D001",
    employee_name: "Nguyễn Văn A",
    department: "HCNS",
    expected_workdays: 24,
    sign_date: "2025-06-15T10:00:00Z",
    shift_days: {
      T2: "F", T3: "F", T4: "F", T5: "F", T6: "C", T7: "N", CN: "F",
    },
    effective_date: "2025-04-15T00:00:00Z",
    content: "Phân ca tháng 6",
    condition: "Hoạt động",
    attached_file: "",
    created_date: "2025-04-10T00:00:00Z",
    employee_id: "E001",
    decision_type_id: "DT001",
    decision_type_name: "Phân ca lặp",
  },
  {
    decision_id: "D002",
    employee_name: "Nguyễn Văn B",
    department: "Kế Toán",
    expected_workdays: 22,
    sign_date: "2025-06-10T15:00:00Z",
    shift_days: {
      T2: "C", T3: "C", T4: "F", T5: "N", T6: "F", T7: "F", CN: "N",
    },
    effective_date: "2025-04-15T00:00:00Z",
    content: "Phân ca tháng 6",
    condition: "Hoạt động",
    attached_file: "",
    created_date: "2025-04-12T00:00:00Z",
    employee_id: "E002",
    decision_type_id: "DT001",
    decision_type_name: "Phân ca lặp",
  },
];

const ShiftDetailPage = () => {
  const [openCreate, setOpenCreate] = useState(false);
  const [openEdit, setOpenEdit] = useState(false);
  const [editDecision, setEditDecision] = useState<Decision | null>(null);
  const [loading, setLoading] = useState(false);

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

  const shiftColumns: ColumnDef<Decision>[] = weekdays.map((day) => ({
    id: `shift_${day}`,
    header: day,
    cell: ({ row }) => <div className="text-center">{row.original.shift_days?.[day] || ""}</div>,
    meta: { className: "w-[60px] text-center" },
  }));

  const columns: ColumnDef<Decision>[] = [
    {
      accessorKey: "employee_name",
      header: "Tên nhân sự",
      cell: ({ getValue }) => <div className="truncate">{getValue<string>()}</div>,
      meta: { className: "w-[120px]" },
    },
    {
      accessorKey: "department",
      header: "Phòng ban",
      cell: ({ getValue }) => <div className="truncate">{getValue<string>()}</div>,
      meta: { className: "w-[100px]" },
    },
    {
      accessorKey: "expected_workdays",
      header: "Số ngày công dự kiến",
      cell: ({ getValue }) => <div className="text-center">{getValue<number>()}</div>,
      meta: { className: "w-[80px]" },
    },
    ...shiftColumns,
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => {
        const decision = row.original;
        return (
          <div className="flex gap-2 justify-center">
            <Button variant="outline" size="icon" onClick={() => {
              setEditDecision(decision);
              setOpenEdit(true);
            }}>
              <SquarePen className="w-4 h-4" />
            </Button>
            <ConfirmDelete deleteFn={() => deleteDecision(decision.decision_id)} />
          </div>
        );
      },
      meta: { className: "w-[100px] text-center" },
    },
    {
      accessorKey: "sign_date",
      header: "Trạng thái",
      cell: ({ row }) => {
        const date = new Date(row.original.sign_date);
        return (
          <div className="text-center">
            {date.toLocaleTimeString("vi-VN", { hour: "2-digit", minute: "2-digit" })} {date.toLocaleDateString("vi-VN")}
          </div>
        );
      },
      meta: { className: "w-[140px] text-center" },
    },
  ];

  const ButtonCreate = () => (
    <div className="button-container flex items-center justify-center space-x-4 text-[17px]">
      <CreateDecisionForm open={openCreate} setOpen={setOpenCreate} createDecision={createDecision} />
      <Button onClick={() => setOpenCreate(true)} variant="default" className="px-8 py-5 text-[17px] rounded-[15px]">
        <span className="mb-1 text-[24px]">+</span> Yêu cầu xác nhận công
      </Button>
      <Button onClick={exportFile} variant="default" className="px-10 py-5 rounded-[15px] text-[17px]" disabled={loading}>
        {loading ? "Đang xuất..." : (<><img src={export_file} alt="export-file" className="w-[24px]" /> Xuất file</>)}
      </Button>
    </div>
  );

  const NavLink = () => (
    <div className="p-4 w-full">
      <div className="flex items-center mb-4">
        <h2 className="text-lg font-semibold">Phân ca lặp NVCT THD HCM</h2>
        <span className="px-3 py-1 text-sm font-semibold text-green-700 border border-green-500 rounded-full ml-2">ACTIVE</span>
      </div>
      <div className="grid grid-cols-2 gap-4 text-sm">
        <div className="space-y-2">
          <div className="flex border-b pb-1 relative"><span className="text-gray-700 font-bold">Ngày hiệu lực</span><span className="text-black absolute right-[50%]">15/04/2025</span></div>
          <div className="flex border-b pb-1 relative"><span className="text-gray-700 font-bold">Hết hiệu lực</span><span className="text-black absolute right-[50%]">15/06/2025</span></div>
          <div className="flex border-b pb-1 relative"><span className="text-gray-700 font-bold">Lặp lại theo</span><span className="text-black absolute right-[50%]">Tháng</span></div>
        </div>
        <div className="space-y-2">
          <div className="flex border-b pb-1 relative"><span className="text-gray-700 font-bold">Đối tượng áp dụng</span><span className="text-black absolute right-[50%]">Nhân viên chính thức</span></div>
          <div className="flex border-b pb-1 relative"><span className="text-gray-700 font-bold">Văn phòng</span><span className="text-black absolute right-[50%]">VP Hồ Bá Kiện - HCM</span></div>
        </div>
      </div>
    </div>
  );

  if (isLoading) return <Loading />;

  return (
    <>
      <DataTable
        columns={columns}
        buttonCreate={<ButtonCreate />}
        data={mockDecisions}
        navLink={<NavLink />}
        title="PHÂN CA"
        keyFilter=""
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

export default ShiftDetailPage;
