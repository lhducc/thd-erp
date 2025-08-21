import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import {Button} from "@/components/ui/button";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod";
import {useForm} from "react-hook-form";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {updateEmployeeApi} from "@/apis/profile.api";
import {usePosition} from "@/query/usePosition.ts";
import Asterisk from "@/components/ui/Asterisk.tsx";
import {useDepartment} from "@/query/useDepartment.ts";
import {useOffice} from "@/query/useOffice.ts";
import type {Employee} from "@/types/employee.ts";
import {useJobTitle} from "@/query/useJobTitle.ts";
import {useState} from "react";
import Loading from "@/components/Loading.tsx";
import {toast} from "sonner";
import {useGetRoles} from "@/query/role.query.ts";
import {useWorkSchedule} from "@/query/useWorkSchedule.ts";
import {useWorkScheduleRegister} from "@/query/useWorkScheduleRegister.ts";
import {useEmployeeByRoleNameQuery} from "@/query/employee.query.ts";

const employee = z.object({
    full_name: z.string().nonempty("Vui lòng nhập họ và tên"),
    birth_date: z.string().nonempty("Vui lòng nhập ngày sinh"),
    gender: z.string().nonempty("Vui lòng chọn giới tính"),
    phone: z.string().nonempty("Vui lòng nhập số điện thoại"),
    email: z.string().email("Email không hợp lệ").nonempty("Vui lòng nhập email"),
    position: z.string().nonempty("Vui lòng chọn vị trí"),
    job_title_id: z.string().nonempty("Vui lòng chọn chức vụ"),
    work_type: z.string().nonempty("Vui lòng chọn hình thức làm việc"),
    current_address: z.string().nonempty("Vui lòng nhập địa chỉ hiện tại"),
    start_date: z.string().nonempty("Vui lòng chọn ngày bắt đầu"),
    office: z.string().nonempty("Vui lòng chọn văn phòng"),
    role_id: z.string().min(1, "Vui lòng chọn phân quyền"),
    department: z.string().nonempty("Vui lòng chọn phòng ban"),
    manager: z.string().optional(),
    schedule_id: z.string().optional(),
});

type Props = {
    open: boolean;
    setOpen: (open: boolean) => void;
    data: Employee;
    refetchEmployee: () => void;
};

const EditEmployeeForm = ({open, setOpen, data, refetchEmployee}: Props) => {
    const [isLoading, setIsLoading] = useState(false);
    const {data: workschedulesAuto} = useWorkSchedule()
    const {data: workschedulesRegister} = useWorkScheduleRegister()
    const [workScheduleType, setWorkScheduleType] = useState<boolean>(workschedulesAuto?.includes(data.schedule_id) || false) // True là tự động, False là đăng ký
    function toYMD(date: string | Date): string {
        return new Date(date).toISOString().split("T")[0];
    }

    const workSchedules = workScheduleType ? workschedulesAuto : workschedulesRegister
    const form = useForm<z.infer<typeof employee>>({
        resolver: zodResolver(employee),
        defaultValues: {
            full_name: data.full_name,
            birth_date: toYMD(data.birthday),
            gender: data.gender,
            phone: data.phone_number,
            email: data.email,
            position: data.position_id,
            current_address: data.address || "",
            job_title_id: data.job_title_id || "",
            start_date: toYMD(data.created_date),
            office: data.department.office.office_id || "",
            department: data.department.department_id || "",
            manager: data.manager_id || "",
            role_id: data.role.id,
            work_type: data.work_type,
            schedule_id: data.schedule_id?.toString(),
        },
    });

    const {data: positions} = usePosition()
    const {data: departments} = useDepartment()
    const {data: offices} = useOffice()
    const {data: jobTitle} = useJobTitle()
    const {data: roles} = useGetRoles()
    const {data: managers} = useEmployeeByRoleNameQuery("manager")

    async function onSubmit(values: z.infer<typeof form>) {
        setIsLoading(true);
        try {
            const payload = {
                full_name: values.full_name,
                birthday: values.birth_date,
                gender: values.gender,
                phone_number: values.phone,
                email: values.email,
                work_type: values.work_type,
                position_id: values.position,
                job_title_id: values.job_title_id,
                status: "active",
                manager_id: values.manager === "" ? null : values.manager,
                address: values.current_address,
                department_id: values.department,
                schedule_id: Number(values.schedule_id)
            };
            await updateEmployeeApi(data.employee_id, payload);
            refetchEmployee();
            setIsLoading(false);
            toast.success("Cập nhật nhân viên thành công")
            setOpen(false);
        } catch (error) {
            console.error("Error submitting form:", error);
        }
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent className="w-[1000px] py-10 max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <div className="flex px-[30px] items-center relative">
                        <DialogTitle className="text-xl font-semibold text-left">Chỉnh sửa nhân Sự</DialogTitle>
                        <DialogTitle
                            className="text-3xl w-[186px] h-[54px] font-semibold flex items-center justify-center absolute right-[30px] rounded-3xl border-2 border-gray-400">{data.employee_id}</DialogTitle>
                    </div>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                        <div className="w-full p-6 flex items-center justify-center gap-5">
                            <div className="flex flex-col gap-6 w-full">
                                <div className="w-full mt-9">
                                    <FormField
                                        control={form.control}
                                        name="full_name"
                                        render={({field}) => (
                                            <FormItem>
                                                <FormLabel className="font-medium">Họ và tên</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    />
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="birth_date"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Ngày sinh</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        type="date"
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    />
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="gender"
                                        render={({field}) => (
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
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="phone"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Số điện thoại</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    />
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="email"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Email</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    />
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="job_title_id"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Chức vụ <Asterisk/> </FormLabel>
                                                <FormControl>
                                                    <select
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    >
                                                        <option value="">Chọn chức vụ</option>
                                                        {jobTitle?.map((item) => {
                                                            return <option
                                                                value={item.job_title_id}>{item.job_title} - {item.hierarchy_level.hierarchy_level}</option>
                                                        })}
                                                    </select>
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />
                                    <FormItem className="mt-[20px]">
                                        <FormLabel className="font-medium">Loại lịch làm việc</FormLabel>
                                        <div className="flex items-center space-x-4">
                                            <div className="flex items-center">
                                                <input
                                                    id="register-schedule"
                                                    type="radio"
                                                    checked={!workScheduleType}
                                                    onChange={() => setWorkScheduleType(false)}
                                                    className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500"
                                                />
                                                <label htmlFor="register-schedule"
                                                       className="ms-2 text-sm font-medium text-gray-900">
                                                    Lịch làm việc đăng ký
                                                </label>
                                            </div>
                                            <div className="flex items-center">
                                                <input
                                                    id="auto-schedule"
                                                    type="radio"
                                                    checked={workScheduleType}
                                                    onChange={() => setWorkScheduleType(true)}
                                                    className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500"
                                                />
                                                <label htmlFor="auto-schedule"
                                                       className="ms-2 text-sm font-medium text-gray-900">
                                                    Lịch làm việc tự động
                                                </label>
                                            </div>
                                        </div>
                                    </FormItem>

                                    <FormField
                                        control={form.control}
                                        name="schedule_id"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                {/*<FormLabel className="font-medium">Lịch làm việc</FormLabel>*/}
                                                <FormControl>
                                                    <select
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    >
                                                        <option value="">Chọn lịch làm việc <Asterisk/></option>
                                                        {workSchedules?.map((item) => {
                                                            return <option
                                                                value={item.work_schedule_id}>{item.work_schedule_name}</option>
                                                        })}
                                                    </select>
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />
                                </div>
                            </div>
                            <div className="flex flex-col gap-6 w-full">
                                <div className="w-full">
                                    <FormField
                                        control={form.control}
                                        name="position"
                                        render={({field}) => (
                                            <FormItem>
                                                <FormLabel className="font-medium">Vị trí</FormLabel>
                                                <FormControl>
                                                    <select
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    >
                                                        <option value="">Chọn vị trí</option>
                                                        {positions?.map((position) => {
                                                            return <option
                                                                value={position.position_id}>{position.position_name}</option>
                                                        })}
                                                    </select>
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="current_address"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Địa chỉ hiện tại</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    />
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="start_date"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Ngày bắt đầu</FormLabel>
                                                <FormControl>
                                                    <Input
                                                        type="date"
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    />
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="office"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Văn phòng <Asterisk/></FormLabel>
                                                <select
                                                    {...field}
                                                    className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                >
                                                    <option value="">Chọn văn phòng <Asterisk/></option>
                                                    {offices?.map((office) => {
                                                        return <option
                                                            value={office.office_id}>{office.office_name}</option>
                                                    })}
                                                </select>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="department"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Phòng ban</FormLabel>
                                                <FormControl>
                                                    <select
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    >
                                                        <option value="">Chọn phòng ban</option>
                                                        {departments?.map((department) => {
                                                            return <option
                                                                value={department.department_id}>{department.department_name}</option>
                                                        })}
                                                    </select>
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={form.control}
                                        name="manager"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Quản lý trực tiếp</FormLabel>
                                                <FormControl>
                                                    <select
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    >
                                                        <option value="">Chọn quản lý trực tiếp</option>
                                                        {managers?.map((manager) => {
                                                            return <option
                                                                value={manager.employee_id}>{manager.full_name}</option>
                                                        })}
                                                    </select>
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={form.control}
                                        name="role_id"
                                        render={({field}) => (
                                            <FormItem className="mt-[20px]">
                                                <FormLabel className="font-medium">Phân
                                                    quyền <Asterisk/></FormLabel>
                                                <FormControl>
                                                    <select
                                                        {...field}
                                                        className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                    >
                                                        <option value="">Chọn phân quyền</option>
                                                        {roles?.map((role) => (
                                                            <option value={role.id}>{role.role_name}</option>
                                                        ))}
                                                    </select>
                                                </FormControl>
                                                <FormMessage className="text-red-500 text-sm"/>
                                            </FormItem>
                                        )}
                                    />
                                </div>
                            </div>
                        </div>
                        <div className="flex items-center justify-center w-full gap-5">
                            <Button type="button" onClick={() => setOpen(false)}
                                    className="w-[250px] h-[40px] bg-gray-400 hover:bg-gray-500">
                                Hủy bỏ
                            </Button>
                            <Button type="submit"
                                    className="w-[250px] h-[40px] bg-[#DB3B21] hover:bg-[#b83a1a]">
                                {!isLoading ? "Lưu thông tin" : <Loading/>}
                            </Button>
                        </div>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
};

export default EditEmployeeForm;
