import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "./ui/button";
import { Loader2, Plus } from "lucide-react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
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
import { createPositionApi, updatePositionApi } from "@/apis/position.api";
import { useEffect, useState } from "react";
import type { PayloadPosition, Position } from "@/types";

const formSchema = z.object({
  position_name: z.string().nonempty("Vui lòng nhập tên vị trí"),
});

type Props = {
  refetch?: Function;
  editBtn?: React.ReactNode;
  type?: "edit";
  position?: Position;
};

const CreatePositionForm = ({ refetch, editBtn, type, position }: Props) => {
  const [open, setOpen] = useState(false);

  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      position_name: "",
    },
  });

  useEffect(() => {
    form.setValue("position_name", position?.position_name || "");
  }, [position]);

  const { mutateAsync: createPosition, isPending: pendingCreatePosition } =
    useMutation({
      mutationFn: createPositionApi,
      onSuccess: () => {
        toast.success("Thêm vị trí thành công");
        form.reset();
        setOpen(false);
        refetch && refetch();
      },
      onError: (error) => {
        toast.error(error.message);
      },
    });

  const { mutateAsync: updatePosition, isPending: pendingUpdatePostition } =
    useMutation({
      mutationFn: ({ id, payload }: { id: string; payload: PayloadPosition }) =>
        updatePositionApi(id, payload),
      onSuccess: () => {
        toast.success("Cập nhật vị trí thành công");
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
      await updatePosition({
        id: position?.position_id || "",
        payload: {
          position_name: values.position_name,
        },
      });
    } else {
      await createPosition(values);
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
            Thêm vị trí
          </Button>
        )}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            {type === "edit" ? "Cập nhật vị trí" : "Thêm vị trí"}
          </DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="position_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tên vị trí</FormLabel>
                  <FormControl>
                    <Input placeholder="" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
            <Button
              type="submit"
              disabled={pendingCreatePosition || pendingUpdatePostition}
            >
              {pendingCreatePosition || pendingUpdatePostition ? (
                <Loader2 className="animate-spin" />
              ) : (
                <>{type === "edit" ? "Cập nhật" : "Thêm"}</>
              )}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default CreatePositionForm;
