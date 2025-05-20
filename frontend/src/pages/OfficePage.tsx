import { deleteOfficeApi, getAllOfficesApi } from "@/apis/office.api";
import ConfirmDelete from "@/components/ConfirmDelete";
import CreateOfficeForm from "@/components/CreateOfficeForm";
import DataTable from "@/components/DataTable";
import Loading from "@/components/Loading";
import { Button } from "@/components/ui/button";
import type { Office } from "@/types";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";

const OfficePage = () => {
  const {
    data: offices,
    isPending: pendingOffices,
    refetch: refetchOffices,
  } = useQuery({
    queryKey: ["offices"],
    queryFn: getAllOfficesApi,
    gcTime: 0,
    staleTime: 0,
  });

  const { mutateAsync: deleteOffice } = useMutation({
    mutationFn: (id: string) => deleteOfficeApi(id),
    onSuccess: () => {
      toast.success("Xóa văn phòng thành công");
      refetchOffices();
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

  if (pendingOffices) {
    return <Loading />;
  }

  return (
    <>
      <DataTable
        columns={columns}
        data={offices || []}
        title="Văn phòng"
        buttonCreate={<CreateOfficeForm refetch={refetchOffices} />}
        keyFilter="office_name"
      />
    </>
  );
};

export default OfficePage;
