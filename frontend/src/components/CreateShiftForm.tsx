'use client';

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { createDecisionApi } from '@/apis/decision.api';
import { useMutation, useQuery } from "@tanstack/react-query";
import { getAllDecisionTypeApi } from '@/apis/decistion-type.api';

export const formSchema = z.object({
  decision_name: z.string().nonempty('Vui lòng nhập tên quyết định'),
  effective_date: z.string().nonempty('Vui lòng chọn ngày hiệu lực'),
  sign_date: z.string().nonempty('Vui lòng chọn ngày ký'),
  condition: z.string().nonempty('Vui lòng chọn tình trạng'),
  description: z.string().nonempty('Vui lòng nhập mô tả'),
  attached_file: z.any().optional(),
  decision_type_id: z.string().nonempty('Vui lòng chọn loại quyết định'),
  employee_ids: z.string().nonempty('Vui lòng nhập id'),
});

type Props = {
  open: boolean;
  setOpen: (open: boolean) => void;
};

const CreateShiftForm = ({ open, setOpen }: Props) => {
  const { data: decisions, isLoading } = useQuery({
    queryKey: ["decisionTypes"],
    queryFn: getAllDecisionTypeApi,
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      decision_name: '',
      effective_date: '',
      sign_date: '',
      condition: '',
      description: '',
      attached_file: undefined,
      employee_ids: '',
      decision_type_id: '',
    },
  });

const onSubmit = async (values: z.infer<typeof formSchema>) => {
  try {
    const formData = new FormData();
    formData.append("decision_name", values.decision_name);
    formData.append("effective_date", values.effective_date);
    formData.append("sign_date", values.sign_date);
    formData.append("content", values.description);
    formData.append("condition", values.condition);
    formData.append("decision_type_id", values.decision_type_id);

    if (values.attached_file?.[0]) {
      formData.append("attached_file", values.attached_file[0]);
    }

    const employeeIds = values.employee_ids
      .split(",")
      .map((id) => id.trim())
      .filter((id) => id !== "");

    employeeIds.forEach((id) => {
      formData.append("employee_ids", id);
    });

    await createDecisionApi(formData);

    setOpen(false);
  } catch (error) {
    console.error("Failed to create decision:", error);
  }
};


  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="w-[1261px] h-[700px] p-6">
        <DialogHeader>
          <div className="flex items-center relative">
            <DialogTitle className="text-xl font-semibold text-left">
              THÊM QUYẾT ĐỊNH
            </DialogTitle>
            <DialogTitle className="text-xl w-[150px] h-[54px] font-semibold flex items-center justify-center absolute right-[30px] rounded-3xl border-2 border-gray-400">
              TL0000001
            </DialogTitle>
          </div>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 flex flex-wrap">
            {/* Left */}
            <div className="w-[50%] relative m-0">
              <div className="w-[96%]">
                <FormField
                  control={form.control}
                  name="decision_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Tên quyết định</FormLabel>
                      <FormControl>
                        <Input {...field} type="text" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] mt-[20px]">
                <FormField
                  control={form.control}
                  name="sign_date"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Ngày ký quyết định</FormLabel>
                      <FormControl>
                        <Input {...field} type="datetime-local" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] mt-[20px]">
                <FormField
                  control={form.control}
                  name="effective_date"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Ngày hiệu lực</FormLabel>
                      <FormControl>
                        <Input {...field} type="datetime-local" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] relative mt-[20px]">
                <FormField
                  control={form.control}
                  name="attached_file"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">File đính kèm</FormLabel>
                      <FormControl>
                        <div>
                          <input
                            id="file-upload"
                            type="file"
                            onChange={(e) => field.onChange(e.target.files)}
                            className="w-full h-[140px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                          <label
                            htmlFor="file-upload"
                            className="absolute p-2 bg-[#EFEFEF] rounded-xl top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 text-center font-medium text-black cursor-pointer"
                          >
                            ĐÍNH KÈM FILE
                          </label>
                        </div>
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <Button
                type="button"
                onClick={() => setOpen(false)}
                className="bg-gray-400 absolute bottom-[-10px] right-[4%] hover:bg-gray-500 w-[120px] h-[40px]"
              >
                Hủy bỏ
              </Button>
            </div>

            {/* Right */}
            <div className="w-[50%] relative">
              <div className="w-[96%]">
                <FormField
                  control={form.control}
                  name="condition"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Tình trạng</FormLabel>
                      <FormControl>
                        <Input {...field} type="text" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] mt-[20px]">
                <FormField
                  control={form.control}
                  name="decision_type_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Loại quyết định</FormLabel>
                      <FormControl>
                        <select {...field} className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]">
                          <option value="" disabled>Chọn loại quyết định</option>
                          {isLoading ? (
                            <option value="">Đang tải...</option>
                          ) : (
                            decisions?.map((item) => (
                              <option key={item.decision_type_id} value={item.decision_type_id}>
                                {item.decision_type}
                              </option>
                            ))
                          )}
                        </select>
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] mt-[20px]">
                <FormField
                  control={form.control}
                  name="employee_ids"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">ID nhân viên</FormLabel>
                      <FormControl>
                        <Input {...field} type="text" placeholder="VD: 1,2,3" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <div className="w-[96%] mt-[20px]">
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Mô tả</FormLabel>
                      <FormControl>
                        <textarea {...field} className="w-full h-[140px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21] resize-none" />
                      </FormControl>
                      <FormMessage className="text-red-500 text-sm" />
                    </FormItem>
                  )}
                />
              </div>

              <Button
                type="submit"
                className="bg-[#DB3B21] absolute bottom-[-10px] hover:bg-[#b83a1a] w-[120px] h-[40px]"
              >
                Lưu thông tin
              </Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default CreateShiftForm;
