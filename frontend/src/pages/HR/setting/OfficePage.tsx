import { deleteOfficeApi } from "@/apis/office.api.ts";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import CreateOfficeForm from "@/components/CreateOfficeForm.tsx";
import DataTable from "@/components/DataTable.tsx";
import { Button } from "@/components/ui/button.tsx";
import type { Office } from "@/types";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";
import {useOffice} from "@/query/useOffice.ts";

const OfficePage = () => {
  const {data: offices, isPending: pendingOffices, refetch: refetchOffices} = useOffice();

  const { mutateAsync: deleteOffice } = useMutation({
    mutationFn: (id: string) => deleteOfficeApi(id),
    onSuccess: () => {
      refetchOffices();
      toast.success("Xóa văn phòng thành công");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  const columns: ColumnDef<Office>[] = [
    {
      accessorKey: "office_name",
      header: "Tên văn phòng",
    },
    {
      accessorKey: "phone_number",
      header: "Số điện thoại",
    },
    {
      accessorKey: "address",
      header: "Địa chỉ",
    },
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => {
        const office = row.original;

        return (
          <div className="flex gap-4">
            <CreateOfficeForm
              editBtn={
                <Button variant="outline">
                  <SquarePen />
                </Button>
              }
              office={office}
              type="edit"
              refetch={refetchOffices}
            />
            <ConfirmDelete deleteFn={() => deleteOffice(office.office_id)} />
          </div>
        );
      },
    },
  ];

  return (
    <>
      <DataTable
        columns={columns}
        data={offices || []}
        title="Văn phòng"
        isLoading={pendingOffices}
        buttonCreate={<CreateOfficeForm refetch={refetchOffices} />}
        keyFilter="office_name"
      />
    </>
  );
};

export default OfficePage;
