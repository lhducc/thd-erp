import { deleteJobTitleApi, getJobTitles } from "@/apis/jobTitle.api.ts";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import CreateJobTitleForm from "@/components/CreateJobTitleForm.tsx";
import DataTable from "@/components/DataTable.tsx";
import Loading from "@/components/Loading.tsx";
import TitleNavLink from "@/components/TitleNavLink.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { SquarePen } from "lucide-react";
import { toast } from "sonner";
import type {JobTitle} from "@/types/job-title.ts";

const JobTitlePage = () => {
  const {
    data: jobTitles,
    isPending: pendingJobTitles,
    refetch: refetchJobTitles,
  } = useQuery({
    queryKey: ["jobTitles"],
    queryFn: getJobTitles,
  });

  const { mutateAsync: deleteJobTitle } = useMutation({
    mutationFn: deleteJobTitleApi,
    onSuccess: () => {
      refetchJobTitles();
      toast.success("Xóa chức danh thành công");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  if (pendingJobTitles) {
    return <Loading />;
  }

  const columns: ColumnDef<JobTitle>[] = [
    {
      accessorKey: "job_title",
      header: "Chức danh",
      cell: (info) => info.getValue(),
    },
    {
      accessorKey: "hierarchy_level_id",
      header: "Cấp bậc",
      cell: (info) => info.getValue(),
    },
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => {
        const jobTitle = row.original;

        return (
            <div className="flex gap-4">
              <CreateJobTitleForm
                  editBtn={
                    <Button variant="outline">
                      <SquarePen />
                    </Button>
                  }
                  type="edit"
                  refetch={refetchJobTitles}
                  jobTitle={jobTitle}
              />
              <ConfirmDelete
                  deleteFn={() => deleteJobTitle(jobTitle.job_title_id)}
              />
            </div>
        );
      },
    },
  ];

  return (
      <DataTable
          columns={columns}
          data={jobTitles || []}
          title="Chức danh"
          navLink={<TitleNavLink />}
          buttonCreate={<CreateJobTitleForm refetch={refetchJobTitles} />}
          keyFilter="job_title"
      />
  );
};

export default JobTitlePage;