import {
  deleteDepartmentApi,
} from "@/apis/department.api.ts";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import CreateDepartmentForm from "@/components/CreateDepartmentForm.tsx";
import DataTable from "@/components/DataTable.tsx";
import { Button } from "@/components/ui/button.tsx";
import type { Department } from "@/types";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";
import {useDepartment} from "@/query/useDepartment.ts";

const DepartmentPage = () => {
  const {
    data: departments,
    isPending: pendingGetDepartments,
    refetch: refetchDepartment,
  } = useDepartment();

  const queryClient = useQueryClient();

  const { mutate: deleteDepartment } = useMutation({
    mutationFn: deleteDepartmentApi,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['departments'] });
      toast.success("Xóa phòng ban thành công");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
  const columns: ColumnDef<Department>[] = [
    {
      accessorKey: "department_name",
      header: "Tên phòng ban",
    },
    {
      accessorKey: "manager",
      header: "Người quản lý",
    },
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => {
        const department = row.original;

        return (
          <div className="flex gap-4">
            <CreateDepartmentForm
              editBtn={
                <Button variant={"outline"}>
                  <SquarePen />
                </Button>
              }
              department={department}
              type="edit"
              refetch={refetchDepartment}
            />
            <ConfirmDelete
              deleteFn={() => deleteDepartment(department.department_id)}
            />
          </div>
        );
      },
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={departments || []}
      title="Phòng ban"
      isLoading={pendingGetDepartments}
      buttonCreate={<CreateDepartmentForm refetch={refetchDepartment} />}
      keyFilter="department_name"
    />
  );
};

export default DepartmentPage;
