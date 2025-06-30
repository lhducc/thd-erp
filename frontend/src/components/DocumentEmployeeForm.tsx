'use client';

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
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

const formSchema = z.object({
  document_type: z.string().nonempty("Vui lòng chọn loại tài liệu"),
  effective_date: z.string().nonempty("Vui lòng chọn ngày hiệu lực"),
  expiration_date: z.string().nonempty("Vui lòng chọn ngày hết hạn"),
  status: z.string().nonempty("Vui lòng chọn tình trạng"),
  description: z.string().nonempty("Vui lòng nhập mô tả"),
});

type Props = {
  open: boolean;
  setOpen: (open: boolean) => void;
};

const DocumentEmployeeForm = ({ open, setOpen }: Props) => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      document_type: "",
      effective_date: "",
      expiration_date: "",
      status: "Chưa hiệu lực",
      description: "",
    },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    console.log(values);
    setOpen(false);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="min-w-[600px] p-6">
        <DialogHeader>
          <div className="flex items-center relative">
            <DialogTitle className="text-xl font-semibold text-left">Tài Liệu Nhân Sự</DialogTitle>
            <DialogTitle className="text-xl w-[150px] h-[54px] font-semibold flex items-center justify-center absolute right-[30px] rounded-3xl border-2 border-gray-400">
              TL0000001
            </DialogTitle>
          </div>
        </DialogHeader>
        <Form select {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 flex flex-wrap justify-between">
            <div className="w-[50%] relative">
              <div className="w-[96%]">
                <FormField
                  control={form.control}
                  name="document_type"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Loại tài liệu</FormLabel>
                      <FormControl>
                        <select
                          {...field}
                          className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                        >
                          <option value="">Chọn loại tài liệu</option>
                          <option value="hsk">Chứng chỉ HSK6</option>
                          <option value="other">Khác</option>
                        </select>
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] relative mt-[20px]">
                <FormField
                  control={form.control}
                  name="effective_date"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Ngày hiệu lực</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          type="date"
                          className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                        />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] mt-[20px]">
                <FormField
                  control={form.control}
                  name="expiration_date"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Ngày hết hạn</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          type="date"
                          className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                        />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>
              <div className="w-[96%] relative mt-[20px]">
                <FormField
                  control={form.control}
                  name="file_attachment"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">File đính kèm</FormLabel>
                      <FormControl>
                        <input
                          {...field}
                          type="file"
                          className="w-full h-[155px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                        />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
                  <div className="absolute p-2 bg-[#EFEFEF] rounded-2xl top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 text-center font-medium text-[#DB3B21] cursor-pointer">
                    ĐÍNH KÈM FILE
                  </div>
              </div>
              <Button
                type="button"
                onClick={() => setOpen(false)}
                className="bg-gray-400 mt-[20px] absolute right-[4%] hover:bg-gray-500 w-[120px] h-[40px]"
              >
                Hủy bỏ
              </Button>
            </div>
            <div className="w-[50%] flex justify-end">
              <div className="w-[96%]">
              <div className="w-full">
                <FormField
                  control={form.control}
                  name="status"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Tình trạng</FormLabel>
                      <FormControl>
                        <select
                          {...field}
                          className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                        >
                          <option value="Chưa hiệu lực">Chưa hiệu lực</option>
                          <option value="Đang hiệu lực">Đang hiệu lực</option>
                          <option value="Đã duyệt">Đã duyệt</option>
                        </select>
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-full mt-[20px]">
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Mô tả / Thông tin quyết định</FormLabel>
                      <FormControl>
                        <textarea
                          {...field}
                          className="w-full h-[340px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          placeholder="Nhập mô tả"
                        />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>
              <Button
                type="submit"
                className="bg-[#DB3B21] mt-[20px] hover:bg-[#b83a1a] w-[120px] h-[40px]"
              >
                Lưu thông tin
              </Button>
            </div>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default DocumentEmployeeForm;
