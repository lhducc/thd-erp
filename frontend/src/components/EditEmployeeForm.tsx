import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
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
import { updateEmployeeApi } from "@/apis/profile.api";
import { useState } from "react";


const formSchema = z.object({
  full_name: z.string().nonempty("Vui lòng nhập họ và tên"),
  birth_date: z.string().nonempty("Vui lòng nhập ngày sinh"),
  gender: z.string().nonempty("Vui lòng chọn giới tính"),
  phone: z.string().nonempty("Vui lòng nhập số điện thoại"),
  email: z.string().email("Email không hợp lệ").nonempty("Vui lòng nhập email"),
  position: z.string().nonempty("Vui lòng chọn vị trí"),
  current_address: z.string().nonempty("Vui lòng nhập địa chỉ hiện tại"),
  start_date: z.string().nonempty("Vui lòng chọn ngày bắt đầu"),
  office: z.string().nonempty("Vui lòng chọn văn phòng"),
  department: z.string().nonempty("Vui lòng chọn phòng ban"),
  manager: z.string().nonempty("Vui lòng nhập tên quản lý trực tiếp"),
});
export interface Employee {
  employee_id: string;
  full_name: string;
  birthday: string;
  gender: string;
  work_type: string;
  phone_number: string;
  email: string;
  account_id: number;
  position_id: string;
  job_title_id: string;
  status: string;
  manager_id: string;
  created_date: string;
  office?: string;
  department?: string;
}

type Props = {
  open: boolean;
  setOpen: (open: boolean) => void;
  data: Employee;
  refetchEmployee: () => void;
};

const EditEmployeeForm = ({ open, setOpen, data, refetchEmployee }: Props) => {

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      full_name: data.full_name, 
      birth_date: data.birthday, 
      gender: data.gender,
      phone: data.phone_number, 
      email: data.email, 
      position: data.position_id, 
      current_address: data.office || "", 
      start_date: data.created_date, 
      office: data.office || "", 
      department: data.department || "", 
      manager: data.manager_id, 
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      const payload = {
        full_name: values.full_name,
        birthday: values.birth_date,
        gender: values.gender,
        phone_number: values.phone,
        email: values.email,
        work_type: "TTS", 
        position_id: values.position, 
        job_title_id: "CV0002",
        status: "active", 
        manager_id: values.manager,
        address: values.current_address,
        department_id: values.department,
      };
      await updateEmployeeApi(data.employee_id, payload);
      refetchEmployee();
      setOpen(false);
    } catch (error) {
      console.error("Error submitting form:", error);
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className=" min-w-[1261px] p-6">
        <DialogHeader>
          <div className="flex px-[30px] items-center relative">
            <DialogTitle className="text-xl font-semibold text-left">Thêm Nhân Sự</DialogTitle>
            <DialogTitle className="text-xl w-[186px] h-[54px] font-semibold flex items-center justify-center absolute right-[30px] rounded-3xl border-2 border-gray-400">TTS011</DialogTitle>
          </div>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <div className="w-full p-6 flex">
              <div className="w-[50%] mr-[30px] flex items-left">
                <div className="w-full">
                  <FormField
                    control={form.control}
                    name="full_name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="font-medium">Họ và tên</FormLabel>
                        <FormControl>
                          <Input
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="birth_date"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Ngày sinh</FormLabel>
                        <FormControl>
                          <Input
                            type="date"
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="gender"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Giới tính</FormLabel>
                        <FormControl>
                          <select
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          >
                            <option value="">Chọn giới tính</option>
                            <option value="male">Nam</option>
                            <option value="female">Nữ</option>
                          </select>
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="phone"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Số điện thoại</FormLabel>
                        <FormControl>
                          <Input
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Email</FormLabel>
                        <FormControl>
                          <Input
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="position"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Chức vụ</FormLabel>
                        <FormControl>
                          <select
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          >
                            <option value="">Thực tập sinh</option>
                            <option value="developer">Giám đốc</option>
                            <option value="manager">Quản lý</option>
                          </select>
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />
                  <div className="relative w-full mt-[20px]">
                    <Button type="button" onClick={() => setOpen(false)} className="w-[250px] h-[40px] bg-gray-400 hover:bg-gray-500 absolute right-0">
                      Hủy bỏ
                    </Button>                    
                  </div>
                </div>
              </div>
              <div className="w-[50%] ml-[30px] flex items-right">
                <div className="w-full">
                  <FormField
                    control={form.control}
                    name="position"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="font-medium">Vị trí</FormLabel>
                        <FormControl>
                          <select
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          >
                            <option value="">Chọn vị trí</option>
                            <option value="developer">Lập trình viên</option>
                            <option value="manager">Quản lý</option>
                          </select>
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="current_address"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Địa chỉ hiện tại</FormLabel>
                        <FormControl>
                          <Input
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="start_date"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Ngày bắt đầu</FormLabel>
                        <FormControl>
                          <Input
                            type="date"
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="office"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Văn phòng</FormLabel>
                        <FormControl>
                          <select
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          >
                            <option value="">Chọn văn phòng</option>
                            <option value="headquarters">Trụ sở chính</option>
                            <option value="branch">Chi nhánh</option>
                          </select>
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="department"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Phòng ban</FormLabel>
                        <FormControl>
                          <select
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          >
                            <option value="">Chọn phòng ban</option>
                            <option value="it">Công nghệ thông tin</option>
                            <option value="hr">Nhân sự</option>
                          </select>
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="manager"
                    render={({ field }) => (
                      <FormItem className="mt-[20px]">
                        <FormLabel className="font-medium">Quản lý trực tiếp</FormLabel>
                        <FormControl>
                          <Input
                            {...field}
                            className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 text-sm" />
                      </FormItem>
                    )}
                  />
                  <div className="relative w-full mt-[20px]">
                    <Button type="submit" className="w-[250px] h-[40px] bg-[#DB3B21] hover:bg-[#b83a1a] absolute left-0">
                      Lưu thông tin
                    </Button>                    
                  </div>
                </div>
              </div>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};

export default EditEmployeeForm;
