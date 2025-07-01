import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { SquarePen } from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete";
import DataTable from "@/components/DataTableDesision";
import Loading from "@/components/Loading";
import { useEffect, useRef, useState } from "react";
import CreateDecisionForm from "@/components/CreateDecisionForm";
import type { Decision } from "@/types";
import type { ColumnDef } from "@tanstack/react-table";
import { TableCell } from "@/components/ui/table";
import {deleteDecisionApi, getAllDecisionsApi} from "@/apis/decision.api.ts";

const DecisionPage = () => {
  const [open, setOpen] = useState(false);
  const [editDecision, setEditDecision] = useState<Decision | null>(null);

  const { data: decisions, isLoading: pendingGetDecisions, refetch: refetchDecision } = useQuery({
    queryKey: ["decisions"],
    queryFn: getAllDecisionsApi,
  });

  const { mutate: deleteDecision } = useMutation({
    mutationFn: deleteDecisionApi,
    onSuccess: () => {
      refetchDecision();
      toast.success("Xóa quyết định thành công");
    },
    onError: (error: any) => {
      toast.error(error.message);
    },
  });

  const columns: ColumnDef<Decision>[] = [
    { accessorKey: "decision_id", header: "Mã" },
    { accessorKey: "decision_name", header: "Tên loại quyết định" },
    { accessorKey: "department", header: "Nhóm đơn quyết định" },
    {
      id: "description",
      header: "Mô tả",
      cell: ({ row }: { row: any }) => (
        <TableCell style={{ width: "300px" }}>
          {row.original.description}
        </TableCell>
      ),
    },
    {
      id: "actions",
      header: "Chỉnh sửa",
      cell: ({ row }: { row: any }) => {
        const decision = row.original;
        return (
          <div className="flex gap-4">
            <Button variant={"outline"} onClick={() => setEditDecision(decision)}>
              <SquarePen />
            </Button>
            <ConfirmDelete deleteFn={() => deleteDecision(decision.decision_id)} />
          </div>
        );
      },
    },
  ];

  const ButtonCreate = () => {
    const [isDropdownVisible, setIsDropdownVisible] = useState(false);
    const dropdownRef = useRef<HTMLDivElement | null>(null);

    const toggleDropdown = () => {
      setIsDropdownVisible(!isDropdownVisible);
    };

    useEffect(() => {
      const handleClickOutside = (event: MouseEvent) => {
        if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
          setIsDropdownVisible(false);
        }
      };

      document.addEventListener("click", handleClickOutside);
      return () => document.removeEventListener("click", handleClickOutside);
    }, []);

    return (
      <div className="button-container flex items-center justify-center space-x-4 text-[17px]">
        <CreateDecisionForm className="w-[1261px]" open={open} setOpen={setOpen} />
        <div ref={dropdownRef} className="dropdown relative">
          <button
            onClick={toggleDropdown}
            className="btn-add-Decision text-left bg-[#DB3B21] w-[200px] text-white p-3 rounded-2xl flex justify-center items-center hover:bg-[#b83a1a] transition duration-200"
          >
            <label htmlFor="" className="text-[20px]">+ </label> Thêm quyết định
          </button>
          {isDropdownVisible && (
            <div className="dropdown-menu absolute left-0 w-[200px] text-[13px] bg-white shadow-lg rounded-b-2xl z-10">
              {["Quyết định thăng chức", "Quyết định khen thưởng", "Quyết định kỷ luật"].map((label, index) => (
                <button
                  key={index}
                  onClick={() => setOpen(true)}
                  className="dropdown-item text-left p-3 w-[200px] border-b shadow-2xl border-gray-500 cursor-pointer hover:bg-[#f2f2f2] rounded-b-xl transition duration-200"
                >
                  <label className="text-[20px]">+ </label> {label}
                </button>
              ))}
            </div>
          )}
        </div>
        <button
          className="btn-export text-left h-[54px] bg-[#DB3B21] w-[192px] text-white p-3 rounded-2xl flex justify-center items-center hover:bg-[#b83a1a] transition duration-200"
          onClick={() => alert("Exporting file...")}
        >
          <img className="absolute right-0" src="/Vector.png" alt="Control" /> Xuất file
        </button>
      </div>
    );
  };

  const NavLink = () => {
    const [activeTab, setActiveTab] = useState("hoat-dong");

    const handleTabClick = (tab: string) => setActiveTab(tab);

    return (
      <div className="button-container flex items-center justify-right mb-[10px] mt-[50px] relative">
        <Button><label>+ </label> Thêm loại quyết định</Button>
      </div>
    );
  };

  if (pendingGetDecisions) return <Loading />;

  return (
    <>
      <DataTable
        columns={columns}
        buttonCreate={<ButtonCreate />}
        data={decisions || []}
        navLink={<NavLink />}
        title="Danh sách quyết định"
        keyFilter="decision_id"
      />
    </>
  );
};

export default DecisionPage;
