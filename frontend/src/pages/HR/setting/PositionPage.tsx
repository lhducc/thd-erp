import { deletePositionApi } from "@/apis/position.api.ts";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import CreatePositionForm from "@/components/CreatePositionForm.tsx";
import DataTable from "@/components/DataTable.tsx";
import TitleNavLink from "@/components/TitleNavLink.tsx";
import { Button } from "@/components/ui/button.tsx";
import type { Position } from "@/types";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";
import {usePosition} from "@/query/usePosition.ts";

const PositionPage = () => {
  const {
    data: positions,
    isPending: pendingGetPositions,
    refetch: refetchPositions,
  } = usePosition()

  const { mutateAsync: deletePosition } = useMutation({
    mutationFn: deletePositionApi,
    onSuccess: () => {
      refetchPositions();
      toast.success("Xóa vị trí thành công");
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

  return (
    <DataTable
      columns={columns || []}
      data={positions || []}
      title="Chức danh"
      isLoading={pendingGetPositions}
      navLink={<TitleNavLink />}
      buttonCreate={<CreatePositionForm refetch={refetchPositions} />}
      keyFilter="position_name"
    />
  );
};

export default PositionPage;
