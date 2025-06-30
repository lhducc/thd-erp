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
import DataTable from "@/components/shift_create";
import Loading from "@/components/Loading";
import CreateShiftForm from "@/components/CreateShiftForm";
import UpdateShiftForm from "@/components/UpdateShiftForm";
import type { Shift } from "@/types";
import type { ColumnDef } from "@tanstack/react-table";
import ShiftScheduleForm from "@/components/ShiftScheduleForm"

const CreateShiftPage = () => {
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
    
  ];

  const NavLink = () => (
    <div className="button-container border-t-2 border-gray-300 flex items-center h-[85px] mb-[10px] pt-[45px] relative">
      <div className="font-bold text-xl">Tạo phân ca lặp</div>
      <div className="absolute right-0 z-30 flex">
        <img
          className="z-50 cursor-pointer"
          src="/control.png"
          alt="Control"
          onClick={() => setIsFilterVisible(!isFilterVisible)}
        />
      </div>
      {isFilterVisible && (
        <div ref={filterRef} className="absolute top-[20px] right-[-13px] w-[250px] bg-white border shadow-lg p-4 rounded-lg z-10">
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
          <button className="btn-filter bg-[#DB3B21] text-white py-2 px-4 rounded-full w-full">Lọc</button>
        </div>
      )}
    </div>
  );

  if (isLoading) return <Loading />;

  return (
    <>
      <DataTable
        columns={columns}
        data={shifts || []}
        buttonCreate={null}
        navLink={<NavLink />}
        title="PHÂN CA"
        keyFilter="shift_id"
      />
      <ShiftScheduleForm/>
    </>
  );
};

export default CreateShiftPage;
