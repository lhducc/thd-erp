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
import { useForm, Controller } from 'react-hook-form';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { updateDecisionApi } from '@/apis/decision.api';
import { getAllDecisionTypeApi } from '@/apis/decistion-type.api';
import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';

export const formSchema = z.object({
  decision_id: z.string().nonempty('Vui lòng nhập mã quyết định'),
  decision_name: z.string().nonempty('Vui lòng nhập tên quyết định'),
  effective_date: z.string().nonempty('Vui lòng chọn ngày hiệu lực'),
  sign_date: z.string().nonempty('Vui lòng chọn ngày ký'),
  content: z.string().nonempty('Vui lòng nhập nội dung'),
  condition: z.string().nonempty('Vui lòng nhập tình trạng'),
  created_date: z.string().nonempty('Vui lòng nhập ngày tạo'),
  employee_id: z.string().nonempty('Vui lòng nhập mã nhân viên'),
  decision_type_id: z.string().nonempty('Vui lòng nhập mã loại quyết định'),
  attached_file: z.string().optional(), // For file
});
type FormSchemaType = z.infer<typeof formSchema>;

type Props = {
  open: boolean;
  setOpen: (open: boolean) => void;
  Decision: FormSchemaType;
};

const UpdateDecisionForm = ({ open, setOpen, Decision }: Props) => {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  const { data: decisions, isLoading, isError, error } = useQuery({
    queryKey: ["decisionTypes"],
    queryFn: getAllDecisionTypeApi,
  });

  const onSubmit = async (values: FormSchemaType) => {
    try {
      const formData: FormSchemaType = {
        decision_name: values.decision_name,
        effective_date: new Date(values.effective_date).toISOString(),
        sign_date: new Date(values.sign_date).toISOString(),
        content: values.content,
        condition: values.condition,
        decision_type_id: values.decision_type_id,  // Use decision_type_id here
        employee_id: 'THD002', // Adjust as needed
        attached_file: values.attached_file ? values.attached_file[0] : null,
      };

      await updateDecisionApi(Decision.decision_id, formData);
      setOpen(false);
    } catch (error) {
      console.error('Failed to update decision:', error);
    }
  };

  useEffect(() => {
    if (Decision) {
      form.reset({
        decision_id: Decision.decision_id,
        decision_name: Decision.decision_name,
        effective_date: Decision.effective_date,
        sign_date: Decision.sign_date,
        content: Decision.content,
        condition: Decision.condition,
        created_date: Decision.created_date,
        employee_id: Decision.employee_id,
        decision_type_id: Decision.decision_type_id, // Set default decision_type_id
      });
    }
  }, [Decision, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="w-[1261px] h-[700px] p-6">
        <DialogHeader>
          <div className="flex items-center relative">
            <DialogTitle className="text-xl font-semibold text-left">
              CẬP NHẬT QUYẾT ĐỊNH
            </DialogTitle>
            <DialogTitle className="text-xl w-[150px] h-[54px] font-semibold flex items-center justify-center absolute right-[30px] rounded-3xl border-2 border-gray-400">
              {Decision.decision_id}
            </DialogTitle>
          </div>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 flex flex-wrap">
            <div className="w-[50%] relative m-0">
              <div className="w-[96%]">
                <FormField
                  control={form.control}
                  name="decision_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Tên quyết định</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          type="text"
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
                  name="sign_date"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Ngày ký quyết định</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          type="datetime-local"
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
                  name="effective_date"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="font-medium">Ngày quyết định có hiệu lực</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          type="datetime-local"
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

            <div className="w-[50%] flex justify-end relative">
              <div className="w-[96%]">
                <div className="w-full">
                  <FormField
                    control={form.control}
                    name="condition"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="font-medium">Tình trạng</FormLabel>
                        <FormControl>
                          <Input
                            {...field}
                            type="text"
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />
                </div>

                <div className="w-full mt-[20px]">
                  <FormField
                    control={form.control}
                    name="decision_type_id" // Changed to decision_type_id
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="font-medium">Loại quyết định</FormLabel>
                        <FormControl>
                          <Controller
                            name="decision_type_id"
                            control={form.control}
                            render={({ field }) => (
                              <Select
                                {...field}
                                defaultValue={field.value}
                                className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                              >
                                <SelectTrigger className="w-full h-[51px]">
                                  <SelectValue placeholder="Chọn loại quyết định" />
                                </SelectTrigger>
                                <SelectContent>
                                  {isLoading ? (
                                    <SelectItem value="">Loading...</SelectItem>
                                  ) : decisions?.map((item) => {
                                      return (
                                        <SelectItem key={item.decision_type_id} value={item.decision_type_id}>
                                          {item.decision_type}
                                        </SelectItem>
                                      );
                                    })
                                  }
                                </SelectContent>
                              </Select>
                            )}
                          />
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
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default UpdateDecisionForm;
