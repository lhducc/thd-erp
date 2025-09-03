import {useMutation, useQueryClient} from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useEffect, useState } from "react";
import {Link, useNavigate, useParams} from "react-router-dom";
import { toast } from "sonner";

// Components
import { InfoRow } from "@/components/ui/info-row";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Checkbox } from "@/components/ui/checkbox";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import Loading from "@/components/Loading";

// Types & API
import type { Employee } from "@/types/employee";
import {assignWorkshiftSchedule, assignWorkshiftScheduleAuto} from "@/apis/work-schedule.api";
import { useWorkScheduleById } from "@/query/useWorkSchedule";
import { useEmployeeByRoleNameQuery, useGetAllEmployee } from "@/query/employee.query";

// Constants
import PATH from "@/constants/Path";

interface ManagerPermission {
    manager_id: string;
    full_name: string;
    is_reading: boolean;
    is_editing: boolean;
}

interface AssignParams {
    id: string;
    payload: {
        managers: ManagerPermission[];
    };
}

const formSchema = z.object({
    work_schedule_id: z.number().optional(),
    work_schedule_name: z.string().min(1, "Tên lịch làm việc là bắt buộc"),
    office_id: z.string().min(1, "Văn phòng là bắt buộc"),
    repeat_type: z.string().min(1, "Loại lặp lại là bắt buộc"),
    repeat_cycle: z.number().min(1, "Chu kỳ lặp lại là bắt buộc"),
    effective_date: z.string().min(1, "Ngày hiệu lực là bắt buộc"),
    // expiration_date: z.string().min(1, "Ngày hết hiệu lực là bắt buộc"),
    managers: z.array(
        z.object({
            manager_id: z.string(),
            is_reading: z.boolean(),
            is_editing: z.boolean()
        })
    ).min(1, "Chọn ít nhất 1 quản lý"),
});

const SettingWorkScheduleAuto = () => {
    const { id } = useParams<{ id: string }>();
    const workScheduleId = id ? parseInt(id) : 0;

    // Data fetching
    const { data: employees, isLoading: pendingGetEmployees, error } = useGetAllEmployee();
    const { data: managers } = useEmployeeByRoleNameQuery("manager");
    const { data: workSchedule, isLoading: pendingWorkSchedule } = useWorkScheduleById(id);

    // State management
    const [selectedManagers, setSelectedManagers] = useState<ManagerPermission[]>([]);

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            work_schedule_id: workScheduleId,
            work_schedule_name: "",
            office_id: "",
            repeat_type: "",
            repeat_cycle: 1,
            effective_date: "",
            expiration_date: "",
            managers: [],
        },
    });

    useEffect(() => {
        if (workSchedule) {
            form.reset({
                work_schedule_id: workSchedule.work_schedule_id,
                work_schedule_name: workSchedule.work_schedule_name,
                office_id: workSchedule.office_id,
                repeat_type: workSchedule.repeat_type,
                repeat_cycle: workSchedule.repeat_cycle,
                effective_date: workSchedule.effective_date,
                expiration_date: workSchedule.expiration_date,
                managers: workSchedule.managers?.map(m => ({
                    manager_id: m.employee_id, // Sửa key đúng
                    is_reading: m.is_reading,
                    is_editing: m.is_editing
                })) || [],
            });

            const scheduleManagers = workSchedule.managers?.map(m => ({
                manager_id: m.employee_id, // Sửa key đúng
                full_name: m.manager?.full_name || "",
                is_reading: m.is_reading,
                is_editing: m.is_editing
            })) || [];
            setSelectedManagers(scheduleManagers);
        }
    }, [workSchedule, employees, form]);

    const toggleManagerSelection = (managerId: string) => {
        const managerInfo = managers?.find(m => m.employee_id === managerId);
        setSelectedManagers(prev =>
            prev.some(m => m.manager_id === managerId)
                ? prev.filter(m => m.manager_id !== managerId)
                : [
                    ...prev,
                    {
                        manager_id: managerId,
                        full_name: managerInfo?.full_name || "",
                        is_reading: true,
                        is_editing: false
                    }
                ]
        );
    };


    const updateManagerPermission = (managerId: string, permission: 'read' | 'edit', value: boolean) => {
        setSelectedManagers(prev =>
            prev.map(manager =>
                manager.manager_id === managerId
                    ? {
                        ...manager,
                        is_reading: permission === 'read' ? value : manager.is_reading,
                        is_editing: permission === 'edit' ? value : manager.is_editing
                    }
                    : manager
            )
        );
    };

    useEffect(() => {
        form.setValue("managers", selectedManagers, {
            shouldValidate: true
        });
    }, [selectedManagers, form]);

    const availableManagers = managers?.filter(manager =>
        !selectedManagers.some(m => m.manager_id === manager.employee_id)
    ) || [];
    const queryClient = useQueryClient();
const navigate = useNavigate();
    const { mutateAsync: assign, isPending: pendingAssign } = useMutation({
        mutationFn: (params: AssignParams) => assignWorkshiftScheduleAuto(params.id, params.payload),
        onSuccess: async () => {
            toast.success("Cấu hình thành công");
            form.reset();
            await queryClient.invalidateQueries({queryKey: ["work-schedule", "work-schedules", id]})
            navigate("/setup-work-schedule")
        },
        onError: (error) => {
            toast.error(error.message);
            console.error("Submission error:", error);
        },
    });

    const onSubmit = async () => {
        try {
            await assign({
                id: id!,
                payload: {
                    managers: selectedManagers.map(m => ({
                        manager_id: m.manager_id,
                        is_reading: m.is_reading,
                        is_editing: m.is_editing
                    }))
                }
            });
        } catch (error) {
            console.error("Error during submission:", error);
        }
    };

    if (pendingWorkSchedule) {
        return <Loading />;
    }

    if (error) {
        toast.error(error.message);
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
                            <InfoRow label="Văn phòng" value={workSchedule?.office?.office_name || 'N/A'} />
                        </div>
                        <div className="w-full">
                            <InfoRow label="Ngày hiệu lực" value={workSchedule?.effective_date || 'N/A'} />
                            <InfoRow label="Ngày hết hiệu lực" value={workSchedule?.expiration_date || 'N/A'} />
                            <InfoRow
                                label="Tình trạng"
                                value={workSchedule?.status === "active" ? "Đang hiệu lực" : "Chưa hiệu lực"}
                            />
                        </div>
                    </div>
                </div>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                        {/* Manager Section */}
                        <div className="bg-white p-6 rounded-lg shadow-sm mb-8">
                            <p className="font-semibold text-lg mb-4">Quản lý lịch làm việc</p>

                            <div className="mb-8">
                                <FormField
                                    control={form.control}
                                    name="managers"
                                    render={({ field }) => (
                                        <FormItem className="w-full">
                                            <FormLabel className="block mb-2 text-sm font-medium text-gray-700">
                                                Quản lý
                                            </FormLabel>
                                            <Select
                                                onValueChange={toggleManagerSelection}
                                                value=""
                                            >
                                                <FormControl className="w-full">
                                                    <SelectTrigger className="h-10">
                                                        <SelectValue placeholder="Chọn quản lý"/>
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent className="w-full">
                                                    {availableManagers.map((employee) => (
                                                        <SelectItem
                                                            key={employee.employee_id}
                                                            value={employee.employee_id}
                                                        >
                                                            {employee.full_name} ({employee.employee_id})
                                                        </SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
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
                                            <TableHeader className="bg-[#DB3B21] sticky top-0">
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
                                                    return (
                                                        <TableRow key={manager.manager_id} className="hover:bg-gray-50">
                                                            <TableCell>{manager.manager_id}</TableCell>
                                                            <TableCell>{manager?.full_name}</TableCell>
                                                            <TableCell>
                                                                <Checkbox
                                                                    checked={manager.is_reading}
                                                                    onCheckedChange={(checked) => {
                                                                        updateManagerPermission(manager.manager_id, 'read', !!checked);
                                                                        if (!checked) {
                                                                            updateManagerPermission(manager.manager_id, 'edit', false);
                                                                        }
                                                                    }}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                <Checkbox
                                                                    checked={manager.is_editing}
                                                                    onCheckedChange={(checked) =>
                                                                        updateManagerPermission(manager.manager_id, 'edit', !!checked)
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

export default SettingWorkScheduleAuto;