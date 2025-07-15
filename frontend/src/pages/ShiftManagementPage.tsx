'use client';

import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { SquarePen, Shield } from "lucide-react";
import { useState, useRef } from "react";
import {
  getAllShiftsApi,
  deleteShiftApi,
  createShiftApi,
  exportShiftExcelApi,
} from "@/apis/shift.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import CreateShiftForm from "@/components/CreateShiftForm";
import UpdateShiftForm from "@/components/UpdateShiftForm";
import type { ColumnDef } from "@tanstack/react-table";
import type {Shift} from "@/types/shift.ts";

const ShiftScheduleManagement = () => {
  const [openCreate, setOpenCreate] = useState(false);
  const [openEdit, setOpenEdit] = useState(false);
  const [editShift, setEditShift] = useState<Shift | null>(null);
  const [loading, setLoading] = useState(false);
  const [isFilterVisible, setIsFilterVisible] = useState(false);
  const filterRef = useRef<HTMLDivElement | null>(null);

  const { data: shifts, isLoading, refetch } = useQuery({
    queryKey: ["shifts"],
    queryFn: getAllShiftsApi,
  });

  const { mutate: deleteShift } = useMutation({
    mutationFn: deleteShiftApi,
    onSuccess: () => {
      refetch();
      toast.success("Xóa phân ca thành công");
    },
    onError: (error: any) => {
      toast.error(error.message);
    },
  });

  const { mutate: createShift } = useMutation({
    mutationFn: createShiftApi,
    onSuccess: () => {
      refetch();
      toast.success("Tạo phân ca thành công");
      setOpenCreate(false);
    },
    onError: (error: any) => {
      toast.error("Lỗi khi tạo phân ca: " + error.message);
    },
  });

  const exportFile = async () => {
    try {
      setLoading(true);
      const file = await exportShiftExcelApi();
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

  const columns: ColumnDef<Shift>[] = [
    { accessorKey: "shift_name", header: "Tên phân ca lặp" },
    { accessorKey: "effective_date", header: "Ngày hiệu lực" },
    { accessorKey: "target", header: "Đối tượng áp dụng" },
    { accessorKey: "office", header: "Văn phòng" },
    { accessorKey: "status", header: "Tình trạng" },
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => {
        const shift = row.original;
        return (
          <div className="flex gap-4">
            <Button variant="outline" onClick={() => { setEditShift(shift); setOpenEdit(true); }}>
              <SquarePen className="w-4 h-4" />
            </Button>
            <Button variant="ghost" onClick={() => toast("Tính năng đổi trạng thái chưa được xử lý")}>
              <Shield className="w-4 h-4" />
            </Button>
            <ConfirmDelete deleteFn={() => deleteShift(shift.shift_id)} />
          </div>
        );
      }
    }
  ];

  const ButtonCreate = () => (
    <div className="relative flex items-center justify-center space-x-4 text-[17px]">
      <CreateShiftForm open={openCreate} setOpen={setOpenCreate} createShift={createShift} />
      <Button
        onClick={() => setOpenCreate(true)}
        variant="default"
        className="px-8 py-5 text-[17px] rounded-[15px] mr-[60px]"
      >
        <span className="mb-1 text-[24px]">+</span> Tạo phân ca lặp
      </Button>

      <div className="absolute right-0 top-0 z-30 flex">
        <img
          className="cursor-pointer"
          src="/control.png"
          alt="Control"
          onClick={() => setIsFilterVisible(!isFilterVisible)}
        />
        
      {isFilterVisible && (
        <div
          ref={filterRef}
          className="absolute top-[60px] right-[-13px] w-[250px] bg-white border shadow-lg p-4 rounded-lg z-10"
        >
          <h3 className="font-bold text-lg mb-4">Lọc</h3>
          <div className="mb-4">
            <h4 className="text-sm font-semibold">Loại quyết định</h4>
            {["Khen thưởng theo quý", "Khen thưởng theo năm"].map((label, i) => (
              <div key={i}>
                <input type="checkbox" id={`filter-decision-type${i}`} />
                <label htmlFor={`filter-decision-type${i}`}> {label}</label><br />
              </div>
            ))}
          </div>
          <div className="mb-4">
            <h4 className="text-sm font-semibold">Tình trạng</h4>
            {["Đang hiệu lực", "Hết hiệu lực", "Chưa hiệu lực"].map((label, i) => (
              <div key={i}>
                <input type="checkbox" id={`filter-status${i}`} />
                <label htmlFor={`filter-status${i}`}> {label}</label><br />
              </div>
            ))}
          </div>
          <button className="bg-[#DB3B21] text-white py-2 px-4 rounded-full w-full">Lọc</button>
        </div>
      )}
      </div>

    </div>
  );

  if (isLoading) return <Loading />;

  return (
    <div>
      <DataTable
        columns={columns}
        data={shifts || []}
        buttonCreate={<ButtonCreate />}
        navLink={null}
        title="PHÂN CA"
        keyFilter="shift_id"
      />
      {editShift && (
        <UpdateShiftForm
          open={openEdit}
          setOpen={(val) => {
            setOpenEdit(val);
            if (!val) setEditShift(null);
          }}
          shift={editShift}
        />
      )}
    </div>
  );
};

export default ShiftScheduleManagement;
