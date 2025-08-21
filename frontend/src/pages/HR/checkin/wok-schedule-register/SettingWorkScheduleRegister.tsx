import {InfoRow} from "@/components/ui/info-row.tsx";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import type {Employee} from "@/types/employee.ts";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table.tsx";
import {useEffect, useState} from "react";
import {Checkbox} from "@/components/ui/checkbox.tsx";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form.tsx";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {z} from "zod";
import {Button} from "@/components/ui/button.tsx";
import {toast} from "sonner";
import {assignWorkshiftSchedule} from "@/apis/work-schedule.api.ts";
import {Link, useNavigate, useParams} from "react-router-dom";
import PATH from "@/constants/Path.ts";
import Loading from "@/components/Loading.tsx";
import {useWorkScheduleRegisterById} from "@/query/useWorkScheduleRegister.ts";
import {formatDate} from "@/lib/utils.ts";
import {useEmployeeByRoleNameQuery, useGetAllEmployee} from "@/query/employee.query.ts";

// Define the ManagerPermission interface
interface ManagerPermission {
    manager_id: string;
    is_reading: boolean;
    is_editing: boolean;
}

interface AssignParams {
    id: string;
    payload: {
        work_schedule_id: number;
        managers: ManagerPermission[];
    };
}

const formSchema = z.object({
    work_schedule_id: z.number(),
    managers: z.array(
        z.object({
            manager_id: z.string(),
            is_reading: z.boolean(),
            is_editing: z.boolean()
        })
    ).min(1, "Chọn ít nhất 1 quản lý"),
    // employee_ids: z.array(z.string()).min(1, "Chọn ít nhất 1 nhân viên"),
});

const SettingWorkScheduleRegister = () => {
    const {id} = useParams<{ id: string }>();
    const workScheduleId = id ? parseInt(id) : 0;
    const { data: employees } = useGetAllEmployee();
    const { data: managers } = useEmployeeByRoleNameQuery("manager");
    const { data: workSchedule } = useQuery(useWorkScheduleRegisterById(id));

    const [selectedEmployees, setSelectedEmployees] = useState<Employee[]>([]);
    const [selectedManagers, setSelectedManagers] = useState<ManagerPermission[]>([]);

    const navigate = useNavigate();
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            work_schedule_id: workScheduleId,
            managers: [],
            // employee_ids: [],
        },
    });

    // Initialize form with work schedule data
    useEffect(() => {
        if (workSchedule) {
            // Initialize selected managers with permissions
            const initialManagers = workSchedule.managers?.map(manager => ({
                manager_id: manager.employee_id,
                is_reading: manager.is_reading, // Default to true
                is_editing: manager.is_editing // Default to false
            })) || [];

            setSelectedManagers(initialManagers);

            // Initialize selected employees
            const initialEmployees = employees?.filter((e:Employee) =>
                workSchedule.employees?.some(we => we.employee_id === e.employee_id)
            ) || [];
            setSelectedEmployees(initialEmployees);

            form.reset({
                work_schedule_id: workSchedule.work_schedule_id,
                managers: initialManagers,
                // employee_ids: initialEmployees.map(e => e.employee_id),
            });
        }
    }, [workSchedule, employees, form]);

    useEffect(() => {
        form.setValue("managers", selectedManagers, {
            shouldValidate: true
        });
    }, [form, selectedManagers]);

    const availableManagers = managers?.filter(manager =>
        !selectedManagers.some(m => m.manager_id === manager.employee_id)
    ) || [];

    const toggleManagerSelection = (employeeId: string) => {
        const manager = managers?.find(m => m.employee_id === employeeId);
        if (!manager) return;

        setSelectedManagers(prev =>
            prev.some(m => m.manager_id === employeeId)
                ? prev.filter(m => m.manager_id !== employeeId)
                : [...prev, {
                    manager_id: manager.employee_id,
                    is_reading: true,
                    is_editing: false
                }]
        );
    };

    const updateManagerPermission = (employeeId: string, permission: 'read' | 'edit', value: boolean) => {
        setSelectedManagers(prev =>
            prev.map(manager =>
                manager.manager_id === employeeId
                    ? {
                        ...manager,
                        is_reading: permission === 'read' ? value : manager.is_reading,
                        is_editing: permission === 'edit' ? value : manager.is_editing
                    }
                    : manager
            )
        );
    };

    const queryClient = useQueryClient()

    const { mutateAsync: assign, isPending: pendingAssign } = useMutation({
        mutationFn: (params: AssignParams) => assignWorkshiftSchedule(params.id, params.payload.managers),
        onSuccess: async () => {
            toast.success("Cấu hình thành công");
            await queryClient.invalidateQueries({queryKey: ["workScheduleRegisterById", "workScheduleRegister", id]})
            navigate("/setup-work-schedule-register");
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const onSubmit = async (values: z.infer<typeof formSchema>) => {
        await assign({
            id: id!,
            payload: {
                work_schedule_id: values.work_schedule_id,
                managers: selectedManagers,
            }
        });
    };

    if (!employees || !workSchedule) {
        return <Loading />;
    }

    return (
        <div className="bg-gray-50 min-h-screen">
            <div className="container mx-auto py-6">
                <p className="text-3xl font-semibold text-gray-800 mb-5">Cấu hình lịch làm việc</p>

                {/* Office Info Section */}
                <div className="bg-white p-6 rounded-lg shadow-sm mb-8">
                    <div className="border-t-2 border-gray-200 flex flex-col md:flex-row justify-between gap-5">
                        <div className="w-full">
                            <p className="text-xl font-semibold text-gray-600 mb-4">Lịch làm việc</p>
                            <InfoRow label="Văn phòng" value={workSchedule.work_schedule_name || 'N/A'} />
                        </div>
                        <div className="w-full">
                            <InfoRow label="Ngày hiệu lực" value={formatDate(workSchedule.effective_date || '')} />
                            <InfoRow label="Ngày hết hiệu lực" value={formatDate(workSchedule.expiration_date || '')} />
                            <InfoRow label="Tình trạng" value={workSchedule.status || 'Hết hiệu lực'}/>
                        </div>
                    </div>
                </div>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                        {/* Manager Section */}
                        <div className="bg-white p-6 rounded-lg shadow-sm mb-8">
                            <p className="font-semibold text-lg mb-4">Quản lý lịch làm việc</p>

                            {/* Manager selection dropdown */}
                            <div className="mb-8">
                                <FormField
                                    control={form.control}
                                    name="managers"
                                    render={() => (
                                        <FormItem className="w-full">
                                            <FormLabel className="block mb-2 text-sm font-medium text-gray-700">
                                                Quản lý
                                            </FormLabel>
                                            {availableManagers.length > 0 ? (
                                                <Select
                                                    onValueChange={toggleManagerSelection}
                                                    value=""
                                                >
                                                    <FormControl>
                                                        <SelectTrigger className="h-10">
                                                            <SelectValue placeholder="Chọn quản lý"/>
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        {availableManagers.map((manager) => (
                                                            <SelectItem
                                                                key={manager.employee_id}
                                                                value={manager.employee_id}
                                                            >
                                                                {manager.full_name} ({manager.employee_id})
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            ) : (
                                                <p className="text-sm text-gray-500 py-2 px-3">
                                                    Tất cả quản lý đã được chọn
                                                </p>
                                            )}
                                            <FormMessage className="text-red-500 text-xs mt-1"/>
                                        </FormItem>
                                    )}
                                />
                            </div>

                            {selectedManagers.length > 0 && (
                                <div className="mb-8">
                                    <p className="font-medium mb-3 text-gray-700">
                                        Quản lý đã chọn ({selectedManagers.length}):
                                    </p>
                                    <div className="border rounded-lg overflow-auto max-h-[180px]">
                                        <Table>
                                            <TableHeader className="sticky top-0 bg-[#DB3B21]">
                                                <TableRow>
                                                    <TableHead className="text-white">Mã NV</TableHead>
                                                    <TableHead className="text-white">Họ tên</TableHead>
                                                    <TableHead className="text-white">Quyền đọc</TableHead>
                                                    <TableHead className="text-white">Quyền sửa</TableHead>
                                                    <TableHead className="text-white">Thao tác</TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {selectedManagers.map((manager) => {
                                                    const managerInfo = managers?.find(m => m.employee_id === manager.manager_id);
                                                    return (
                                                        <TableRow key={manager.manager_id} className="hover:bg-gray-50">
                                                            <TableCell>{manager.manager_id}</TableCell>
                                                            <TableCell>{managerInfo?.full_name}</TableCell>
                                                            <TableCell>
                                                                <Checkbox
                                                                    checked={manager.is_reading}
                                                                    onCheckedChange={(checked) => {
                                                                        updateManagerPermission(
                                                                            manager.manager_id,
                                                                            'read',
                                                                            !!checked
                                                                        );
                                                                        if (!checked) {
                                                                            updateManagerPermission(
                                                                                manager.manager_id,
                                                                                'edit',
                                                                                false
                                                                            );
                                                                        }
                                                                    }}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                <Checkbox
                                                                    checked={manager.is_editing}
                                                                    onCheckedChange={(checked) =>
                                                                        updateManagerPermission(
                                                                            manager.manager_id,
                                                                            'edit',
                                                                            !!checked
                                                                        )
                                                                    }
                                                                    disabled={!manager.is_reading}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                <Button
                                                                    variant="ghost"
                                                                    className="text-red-500 hover:text-red-700"
                                                                    onClick={() => toggleManagerSelection(manager.manager_id)}
                                                                >
                                                                    Xóa
                                                                </Button>
                                                            </TableCell>
                                                        </TableRow>
                                                    );
                                                })}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Submit Button */}
                        <div className="flex justify-center items-center gap-5">
                            <div className="flex justify-end">
                                <Link to={PATH.WORK_SCHEDULE}>
                                    <Button variant="outline" type="button">Hủy bỏ</Button>
                                </Link>
                            </div>
                            <div className="flex justify-end">
                                <Button type="submit" disabled={pendingAssign}>
                                    {pendingAssign ? <Loading /> : "Lưu cấu hình"}
                                </Button>
                            </div>
                        </div>
                    </form>
                </Form>
            </div>
        </div>
    );
};

export default SettingWorkScheduleRegister;