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
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useEffect, useState } from "react";
import {
  createHierarchyLevelApi,
  updateHierarchyLevelApi,
} from "@/apis/hierarchyLevel.api";
import type { HierarchyLevel, PayloadHierarchyLevel } from "@/types";

const formSchema = z.object({
  hierarchy_level: z.string().nonempty("Vui lòng nhập tên chức vụ"),
  hierarchy_number: z.number().min(0, "Vui lòng nhập số tượng trưng"),
});

type Props = {
  refetch?: Function;
  hierarchyLevel?: HierarchyLevel;
  editBtn?: React.ReactNode;
  type?: "edit";
};

const CreateHierarchyLevelForm = ({
  refetch,
  hierarchyLevel,
  editBtn,
  type,
}: Props) => {
  const [open, setOpen] = useState(false);

  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      hierarchy_level: "",
      hierarchy_number: 0,
    },
  });

  useEffect(() => {
    form.setValue("hierarchy_level", hierarchyLevel?.hierarchy_level || "");
    form.setValue("hierarchy_number", hierarchyLevel?.hierarchy_number || 0);
  }, [hierarchyLevel]);

  const {
    mutateAsync: createHierarchyLevel,
    isPending: pendingCreateHierarchyLevel,
  } = useMutation({
    mutationFn: createHierarchyLevelApi,
    onSuccess: () => {
      toast.success("Thêm cấp bậc thành công");
      form.reset();
      setOpen(false);
      refetch && refetch();
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  const {
    mutateAsync: updateHierarchyLevel,
    isPending: pendingUpdateHierarchyLevel,
  } = useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: PayloadHierarchyLevel;
    }) => updateHierarchyLevelApi(id, payload),
    onSuccess: () => {
      toast.success("Cập nhật cấp bậc thành công");
      form.reset();
      setOpen(false);
      refetch && refetch();
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  // 2. Define a submit handler.
  async function onSubmit(values: z.infer<typeof formSchema>) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    if (type === "edit") {
      await updateHierarchyLevel({
        id: hierarchyLevel?.id || "",
        payload: {
          hierarchy_level: values.hierarchy_level,
          hierarchy_number: values.hierarchy_number,
        },
      });
    } else {
      await createHierarchyLevel(values);
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger>
        {type === "edit" ? (
          editBtn
        ) : (
          <Button>
            <Plus />
            Thêm cấp bậc
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="md:w-[800px] w-[90vw]">
        <DialogHeader>
          <DialogTitle>
            {type === "edit" ? "Cập nhật cấp bậc" : "Thêm cấp bậc"}
          </DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="hierarchy_level"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Cấp bậc</FormLabel>
                  <FormControl>
                    <Input placeholder="" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="hierarchy_number"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Số tượng trưng</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      placeholder=""
                      {...field}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button
              type="submit"
              disabled={
                pendingCreateHierarchyLevel || pendingUpdateHierarchyLevel
              }
            >
              {pendingCreateHierarchyLevel || pendingUpdateHierarchyLevel ? (
                <Loader2 className="animate-spin" />
              ) : type === "edit" ? (
                "Cập nhật"
              ) : (
                "Thêm"
              )}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default CreateHierarchyLevelForm;
