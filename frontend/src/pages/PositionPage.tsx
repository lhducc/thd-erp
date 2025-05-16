import { deletePositionApi, getAllPositionApi } from "@/apis/position.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import CreatePositionForm from "@/components/CreatePositionForm";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import TitleNavLink from "@/components/TitleNavLink";
import { Button } from "@/components/ui/button";
import type { Position } from "@/types";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";

const PositionPage = () => {
  const {
    data: positions,
    isPending: pendingGetPositions,
    refetch: refetchPositions,
  } = useQuery({
    queryKey: ["positions"],
    queryFn: getAllPositionApi,
  });

  const { mutateAsync: deletePosition } = useMutation({
    mutationFn: deletePositionApi,
    onSuccess: () => {
      toast.success("Xóa vị trí thành công");
      refetchPositions();
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  const columns: ColumnDef<Position>[] = [
    {
      accessorKey: "position_id",
      header: "Mã vị trí",
    },
    {
      accessorKey: "position_name",
      header: "Tên vị trí",
    },
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => {
        const position = row.original;

        return (
          <div className="flex gap-4">
            <CreatePositionForm
              editBtn={
                <Button variant="outline">
                  <SquarePen />
                </Button>
              }
              position={position}
              type="edit"
              refetch={refetchPositions}
            />
            <ConfirmDelete
              deleteFn={() => deletePosition(position.position_id)}
            />
          </div>
        );
      },
    },
  ];

  if (pendingGetPositions) {
    return <Loading />;
  }

  return (
    <DataTable
      columns={columns || []}
      data={positions || []}
      title="Chức danh"
      navLink={<TitleNavLink />}
      buttonCreate={<CreatePositionForm refetch={refetchPositions} />}
      keyFilter="position_name"
    />
  );
};

export default PositionPage;
