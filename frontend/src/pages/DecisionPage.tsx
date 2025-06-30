'use client';

import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { SquarePen, Filter } from "lucide-react";
import { useState, useEffect, useRef } from "react";
import export_file from "@/assets/export-file.svg";
import { getAllDecisionsApi, deleteDecisionApi, createSampleDecisionApi, exportDecisionExcelApi } from "@/apis/decision.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import CreateDecisionForm from "@/components/CreateDecisionForm";
import UpdateDecisionForm from "@/components/UpdateDecisionForm";

import type { Decision } from "@/types";
import type { ColumnDef } from "@tanstack/react-table";

const DecisionPage = () => {
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

  console.log("data", decisions);

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

  const columns: ColumnDef<Decision>[] = [
    { accessorKey: "decision_id", header: "Mã Quyết định" },
    { accessorKey: "decision_name", header: "Tên Quyết định" },
    {
      accessorKey: "decision_type_name",
      header: "Loại Quyết định",
      cell: ({ row }) => row.original.decision_type_name || "-",
    },
    {
      accessorKey: "sign_date",
      header: "Ngày ký",
      cell: ({ row }) => new Date(row.original.sign_date).toLocaleString("vi-VN"),
    },
    {
      accessorKey: "effective_date",
      header: "Hiệu lực từ ngày",
      cell: ({ row }) => new Date(row.original.effective_date).toLocaleString("vi-VN"),
    },
    { accessorKey: "condition", header: "Tình trạng" },
    {
      id: "actions",
      header: "Chỉnh sửa",
      cell: ({ row }) => {
        const decision = row.original;
        return (
          <div className="flex gap-4">
            <Button
              variant="outline"
              onClick={() => {
                setEditDecision(decision);
                setOpenEdit(true);
              }}
            >
              <SquarePen className="w-4 h-4" />
            </Button>
            <ConfirmDelete deleteFn={() => deleteDecision(decision.decision_id)} />
          </div>
        );
      },
    },
  ];

  const ButtonCreate = () => (
    <div className="button-container flex items-center justify-center space-x-4 text-[17px]">
      <CreateDecisionForm open={openCreate} setOpen={setOpenCreate} createDecision={createDecision} />
      <Button
        onClick={() => setOpenCreate(true)}
        variant="default" className={`px-8 py-5 text-[17px] rounded-[15px]`}
      >
        <span className={`mb-1 text-[24px]`}>+</span>
        Thêm hợp đồng
      </Button>
      <Button
          onClick={exportFile}
          variant={"default"}
          className={`px-10 py-5 rounded-[15px] text-[17px]`}
          disabled={loading}
      >
          {loading ? "Đang xuất..." : (
              <>
                  <img src={export_file} alt="export-file" className="w-[24px]"/>
                  Xuất file
              </>
          )}
      </Button>
    </div>
  );

  const Dropdown = () => {
    const [selectedOption, setSelectedOption] = useState<string | null>(null);
    const options = ["Option 1", "Option 2", "Option 3"];
    return (
      <div className="flex items-center bg-gray-100 rounded-xl p-2 w-72">
        <div className="relative w-full">
          <select
            className="w-full p-2 bg-transparent text-gray-600 focus:outline-none rounded-xl appearance-none"
            value={selectedOption || ""}
            onChange={(e) => setSelectedOption(e.target.value)}
          >
            <option value="" disabled>Chọn mục</option>
            {options.map((option, index) => (
              <option key={index} value={option}>{option}</option>
            ))}
          </select>
        </div>
        <div className="ml-2 p-2 rounded-full bg-white shadow-md flex items-center justify-center cursor-pointer">
          <Filter className="text-gray-500" />
        </div>
      </div>
    );
  };

  const NavLink = () => {
    const [isFilterVisible, setIsFilterVisible] = useState(false);
    const filterRef = useRef<HTMLDivElement | null>(null);

    return (
      <div className="button-container border-t-2 border-gray-300 flex items-center h-[85px] justify-end mb-[10px] pt-[45px] relative">
        <div className={`absolute right-[50px] ${isFilterVisible ? "hidden" : ""}`}>
          <Dropdown />
        </div>
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
              <div>
                {["Khen thưởng theo quý", "Khen thưởng theo năm"].map((label, i) => (
                  <div key={i}>
                    <input type="checkbox" id={`filter-decision-type${i}`} />
                    <label htmlFor={`filter-decision-type${i}`}> {label}</label><br />
                  </div>
                ))}
              </div>
            </div>
            <div className="mb-4">
              <h4 className="text-sm font-semibold">Tình trạng</h4>
              <div>
                {["Đang hiệu lực", "Hết hiệu lực", "Chưa hiệu lực"].map((label, i) => (
                  <div key={i}>
                    <input type="checkbox" id={`filter-status${i}`} />
                    <label htmlFor={`filter-status${i}`}> {label}</label><br />
                  </div>
                ))}
              </div>
            </div>
            <button className="btn-filter bg-[#DB3B21] text-white py-2 px-4 rounded-full w-full">Lọc</button>
          </div>
        )}
      </div>
    );
  };

  if (isLoading) return <Loading />;

  return (
    <>
      <DataTable
        columns={columns}
        buttonCreate={<ButtonCreate />}
        data={decisions || []}
        navLink={<NavLink />}
        title="QUYẾT ĐỊNH"
        keyFilter="decision_id"
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

export default DecisionPage;
