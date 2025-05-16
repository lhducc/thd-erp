import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "./ui/button";
import { Loader2, Plus } from "lucide-react";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useMutation, useQuery } from "@tanstack/react-query";
import { getAllHierarchyLevelApi } from "@/apis/hierarchyLevel.api";
import Loading from "./Loading";
import { toast } from "sonner";
import { useEffect, useState } from "react";
import { createJobTitleApi, updateJobTitleApi } from "@/apis/jobTitle.api";
import type { JobTitle, PayloadJobTitle } from "@/types";

const formSchema = z.object({
  job_title: z.string().nonempty("Vui lòng nhập tên chức vụ"),
  hierarchy_level_id: z.string().nonempty("Vui lòng nhập cấp bậc"),
});

type Props = {
  editBtn?: React.ReactNode;
  type?: "edit";
  refetch: () => void;
  jobTitle?: JobTitle;
};

const CreateJobTitleForm = ({ editBtn, type, refetch, jobTitle }: Props) => {
  const [open, setOpen] = useState(false);

  const { data: hierarchyLevel, isPending: pendingGetHierarchyLevel } =
    useQuery({
      queryKey: ["hierarchy-level"],
      queryFn: getAllHierarchyLevelApi,
    });

  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      job_title: "",
      hierarchy_level_id: "",
    },
  });

  useEffect(() => {
    if (type === "edit" && jobTitle) {
      form.setValue("job_title", jobTitle.job_title);
      form.setValue("hierarchy_level_id", jobTitle.hierarchy_level_id);
    }
  }, [jobTitle]);

  const { mutateAsync: createJobTitle, isPending: pendingCreateJobTitle } =
    useMutation({
      mutationFn: createJobTitleApi,
      onSuccess: () => {
        form.reset();
        toast.success("Thêm chức vụ thành công");
        setOpen(false);
        refetch();
      },
      onError: (error) => {
        toast.error(error.message);
      },
    });

  const { mutateAsync: updateJobTitle, isPending: pendingUpdateJobTitle } =
    useMutation({
      mutationFn: ({ id, payload }: { id: string; payload: PayloadJobTitle }) =>
        updateJobTitleApi(id, payload),
      onSuccess: () => {
        form.reset();
        toast.success("Cập nhật chức vụ thành công");
        setOpen(false);
        refetch();
      },
      onError: (error) => {
        toast.error(error.message);
      },
    });

  // 2. Define a submit handler.
  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (type === "edit") {
      await updateJobTitle({
        id: jobTitle?.job_title_id || "",
        payload: {
          job_title: values.job_title,
          hierarchy_level_id: values.hierarchy_level_id,
        },
      });
    } else {
      await createJobTitle(values);
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger>
        {editBtn ? (
          editBtn
        ) : (
          <Button>
            <Plus />
            Thêm chức vụ
          </Button>
        )}
      </DialogTrigger>
      <DialogContent>
        {pendingGetHierarchyLevel ? (
          <Loading />
        ) : (
          <>
            <DialogHeader>
              <DialogTitle>
                {type === "edit" ? "Cập nhật chức vụ" : "Thêm chức vụ"}
              </DialogTitle>
            </DialogHeader>
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-8"
              >
                <FormField
                  control={form.control}
                  name="job_title"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Chức vụ</FormLabel>
                      <FormControl>
                        <Input placeholder="" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="hierarchy_level_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Cấp bậc</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <FormControl>
                          <SelectTrigger className="w-full">
                            <SelectValue placeholder="" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {hierarchyLevel?.map((item) => (
                            <SelectItem key={item.id} value={item.id}>
                              {item.hierarchy_level}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  disabled={pendingCreateJobTitle || pendingUpdateJobTitle}
                >
                  {pendingCreateJobTitle || pendingUpdateJobTitle ? (
                    <Loader2 className="animate-spin" />
                  ) : type === "edit" ? (
                    "Cập nhật"
                  ) : (
                    "Thêm"
                  )}
                </Button>
              </form>
            </Form>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
};

export default CreateJobTitleForm;
