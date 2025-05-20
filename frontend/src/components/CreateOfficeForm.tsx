import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Loader2, Plus } from "lucide-react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import Map from "./Map";
import { useEffect, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { createOfficeApi, updateOfficeApi } from "@/apis/office.api";
import type { Office, PayloadOffice } from "@/types";
import { toast } from "sonner";

const formSchema = z.object({
  office_name: z.string().nonempty("Tên văn phòng không được để trống"),
  phone_number: z.string().nonempty("Số điện thoại không được để trống"),
  address: z.string().nonempty("Địa chỉ không được để trống"),
});

type Props = {
  editBtn?: React.ReactNode; // Prop để truyền vào nút chỉnh sửa
  office?: Office;
  type?: "edit";
  refetch?: Function;
};

const CreateOfficeForm = ({ editBtn, office, type, refetch }: Props) => {
  const [open, setOpen] = useState(false); // State để quản lý trạng thái mở/đóng của dialog
  const [coordinates, setCoordinates] = useState<{
    lat: number;
    lng: number;
  } | null>(null); // State để lưu tọa độ

  const { mutateAsync: createOffice, isPending: pendingCreateOffice } =
    useMutation({
      mutationFn: createOfficeApi,
      onSuccess: () => {
        refetch && refetch(); // Gọi lại hàm refetch nếu có
        setOpen(false); // Đóng dialog sau khi tạo văn phòng thành công
        toast.success("Tạo văn phòng thành công");
        form.reset();
      },
      onError: (error) => {
        toast.error(error.message);
      },
    });

  const { mutateAsync: updateOffice, isPending: pendingUpdateOffice } =
    useMutation({
      mutationFn: ({ id, payload }: { id: string; payload: PayloadOffice }) =>
        updateOfficeApi(id, payload),
      onSuccess: () => {
        refetch && refetch(); // Gọi lại hàm refetch nếu có
        setOpen(false); // Đóng dialog sau khi cập nhật văn phòng thành công
        toast.success("Cập nhật văn phòng thành công");
        form.reset();
      },
      onError: (error) => {
        toast.error(error.message);
      },
    });

  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      office_name: "",
      phone_number: "",
      address: "",
    },
  });

  useEffect(() => {
    if (office) {
      form.setValue("office_name", office.office_name);
      form.setValue("phone_number", office.phone_number);
      form.setValue("address", office.address);
    }
  }, [office]);

  // 2. Define a submit handler.
  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log("first");
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    if (type === "edit") {
      await updateOffice({
        id: office?.office_id || "",
        payload: {
          ...values,
          latitude: coordinates?.lat || 0,
          longitude: coordinates?.lng || 0,
        },
      });
    } else {
      await createOffice({
        ...values,
        latitude: coordinates?.lat || 0,
        longitude: coordinates?.lng || 0,
      });
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
            Thêm văn phòng
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="md:max-w-2xl">
        <DialogHeader>
          {type === "edit" ? (
            <DialogTitle>Chỉnh sửa văn phòng</DialogTitle>
          ) : (
            <DialogTitle>Thêm văn phòng</DialogTitle>
          )}
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="office_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tên văn phòng</FormLabel>
                  <FormControl>
                    <Input placeholder="" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="phone_number"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Số điện thoại</FormLabel>
                  <FormControl>
                    <Input placeholder="" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="address"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Địa chỉ văn phòng</FormLabel>
                  <FormControl>
                    <Input placeholder="" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Map
              setCoordinates={setCoordinates}
              coordinates={{
                latitude: office?.latitude || 21.0285,
                longitude: office?.longitude || 105.8542,
              }}
            />
            <Button type="submit" disabled={pendingCreateOffice || pendingUpdateOffice}>
              {pendingCreateOffice || pendingUpdateOffice ? (
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

export default CreateOfficeForm;
