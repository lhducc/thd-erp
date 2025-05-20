import {
  deleteHierarchyLevelApi,
  getAllHierarchyLevelApi,
} from "@/apis/hierarchyLevel.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import CreateHierarchyLevelForm from "@/components/CreateHierarchyLevelForm";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import TitleNavLink from "@/components/TitleNavLink";
import { Button } from "@/components/ui/button";
import type { HierarchyLevel } from "@/types";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";

const HierarchyLevelPage = () => {
  const {
    data: hierarchyLevel,
    isPending: pendingGetHierarchyLevel,
    refetch: refetchHierarchyLevel,
  } = useQuery({
    queryKey: ["hierarchy-level"],
    queryFn: getAllHierarchyLevelApi,
  });

  const { mutateAsync: deleteHierarchyLevel } = useMutation({
    mutationFn: deleteHierarchyLevelApi,
    onSuccess: () => {
      toast.success("Xóa thành công");
      refetchHierarchyLevel();
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  if (pendingGetHierarchyLevel) {
    return <Loading />;
  }

  const columns: ColumnDef<HierarchyLevel>[] = [
    {
      accessorKey: "id",
      header: "Mã cấp bậc",
    },
    {
      accessorKey: "hierarchy_level",
      header: "Tên cấp bậc",
    },
    {
      accessorKey: "hierarchy_number",
      header: "Cấp bậc",
    },
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => {
        const hierarchyLevel = row.original;

        return (
          <div className="flex gap-2">
            <CreateHierarchyLevelForm
              hierarchyLevel={hierarchyLevel}
              refetch={refetchHierarchyLevel}
              editBtn={
                <Button variant="outline">
                  <SquarePen />
                </Button>
              }
              type="edit"
            />
            <ConfirmDelete
              deleteFn={() => deleteHierarchyLevel(hierarchyLevel.id)}
            />
          </div>
        );
      },
    },
  ];

  return (
    <DataTable
      columns={columns || []}
      data={hierarchyLevel || []}
      title="Chức danh"
      navLink={<TitleNavLink />}
      buttonCreate={
        <CreateHierarchyLevelForm refetch={refetchHierarchyLevel} />
      }
      keyFilter="hierarchy_level"
    />
  );
};

export default HierarchyLevelPage;
