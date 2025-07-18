import {InfoRow} from "@/components/ui/info-row.tsx";
import {useMutation, useQuery} from "@tanstack/react-query";
import {getAllEmployeesApi} from "@/apis/profile.api.ts";
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
import {Search} from "lucide-react";
import {Input} from "@/components/ui/input.tsx";
import {toast} from "sonner";
import {assignWorkshiftSchedule, getWorkScheduleById} from "@/apis/work-schedule.api.ts";
import {Link, useParams} from "react-router-dom";
import PATH from "@/constants/Path.ts";
import Loading from "@/components/Loading.tsx";

interface AssignParams {
    id: string;
    payload: {
        work_schedule_id: number;
        manager_ids: string[];
        employee_ids: string[];
    };
}

const formSchema = z.object({
    work_schedule_id: z.number().optional(),
    work_schedule_name: z.string().min(1, "Tên lịch làm việc là bắt buộc"),
    office_id: z.string().min(1, "Văn phòng là bắt buộc"),
    repeat_type: z.string().min(1, "Loại lặp lại là bắt buộc"),
    repeat_cycle: z.number().min(1, "Chu kỳ lặp lại là bắt buộc"),
    effective_date: z.string().min(1, "Ngày hiệu lực là bắt buộc"),
    expiration_date: z.string().min(1, "Ngày hết hiệu lực là bắt buộc"),
    manager_ids: z.array(z.string()).min(1, "Chọn ít nhất 1 quản lý"),
    employee_ids: z.array(z.string()).min(1, "Chọn ít nhất 1 nhân viên"),
});

const SettingWorkScheduleAuto = () => {
    const {id} = useParams<{ id: string }>();
    const workScheduleId = id ? parseInt(id) : 0;
    const {data: employees, isLoading: pendingGetEmployees, error} = useQuery({
        queryKey: ["employees"],
        queryFn: () => getAllEmployeesApi(1, 9999),
    });

    const {data: workSchedule, isLoading: pendingWorkSchedule} = useQuery({
        queryKey: ["work-schedule", id],
        queryFn: () => getWorkScheduleById(id!),
        enabled: !!id,
    });

    console.log(workSchedule);

    const [selectedEmployees, setSelectedEmployees] = useState<Employee[]>([]);
    const [selectedManagers, setSelectedManagers] = useState<Employee[]>([]);
    const [searchTerm, setSearchTerm] = useState("");
    const [selectedDepartment, setSelectedDepartment] = useState<string>("all");
    const [selectedWorkType, setSelectedWorkType] = useState<string>("all");

    // Get unique departments from employees
    const departments = Array.from(new Set(
        employees?.map(employee => employee.department.department_name) || []
    ));

    const workTypes = Array.from(new Set(
        employees?.map(employee => employee.work_type).filter(Boolean) || []
    ));

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
            manager_ids: [],
            employee_ids: [],
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
                manager_ids: workSchedule.managers?.map(m => m.employee_id),
                employee_ids: workSchedule.employees?.map(e => e.employee_id),
            });

            // Set selected employees and managers
            const scheduleEmployees = employees?.filter(e =>
                workSchedule.employees?.some(we => we.employee_id === e.employee_id)
            ) || [];
            setSelectedEmployees(scheduleEmployees);

            const scheduleManagers = employees?.filter(e =>
                workSchedule.managers?.some(wm => wm.employee_id === e.employee_id)
            ) || [];
            setSelectedManagers(scheduleManagers);
        }
    }, [workSchedule, employees]);

    // Filter employees based on search term
    const filteredEmployees = employees?.filter(employee => {
        const matchesSearch =
            employee.full_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
            employee.employee_id.toLowerCase().includes(searchTerm.toLowerCase());

        const matchesDepartment =
            selectedDepartment === "all" ||
            employee.department.department_name === selectedDepartment;

        const matchesWorkType =
            selectedWorkType === "all" ||
            employee.work_type === selectedWorkType;

        return matchesSearch && matchesDepartment && matchesWorkType;
    }) || [];

    const toggleEmployeeSelection = (employee: Employee) => {
        setSelectedEmployees(prev => {
            const isSelected = prev.some(e => e.employee_id === employee.employee_id);
            if (isSelected) {
                return prev.filter(e => e.employee_id !== employee.employee_id);
            } else {
                return [...prev, employee];
            }
        });
    };

    const toggleAllEmployeesSelection = () => {
        if (selectedEmployees.length === filteredEmployees.length) {
            // Deselect all
            setSelectedEmployees([]);
        } else {
            // Select all filtered employees
            setSelectedEmployees([...filteredEmployees]);
        }
    };

    const toggleManagerSelection = (employeeId: string) => {
        const manager = employees?.find(e => e.employee_id === employeeId);
        if (!manager) return;

        setSelectedManagers(prev => {
            const isSelected = prev.some(m => m.employee_id === employeeId);
            if (isSelected) {
                return prev.filter(m => m.employee_id !== employeeId);
            } else {
                return [...prev, manager];
            }
        });

        // Update form value
        form.setValue(
            "manager_ids",
            selectedManagers.some(m => m.employee_id === employeeId)
                ? selectedManagers.filter(m => m.employee_id !== employeeId).map(m => m.employee_id)
                : [...selectedManagers.map(m => m.employee_id), employeeId]
        );
    };

    useEffect(() => {
        form.setValue(
            "employee_ids",
            selectedEmployees.map(e => e.employee_id),
            {shouldValidate: true}
        );
    }, [selectedEmployees, form]);

    useEffect(() => {
        form.setValue(
            "manager_ids",
            selectedManagers.map(m => m.employee_id),
            {shouldValidate: true}
        );
    }, [selectedManagers, form]);

    const isEmployeeSelected = (employeeId: string) => {
        return selectedEmployees.some(e => e.employee_id === employeeId);
    };

    const availableManagers = employees?.filter(employee =>
        !selectedManagers.some(m => m.employee_id === employee.employee_id)
    ) || [];

    const {mutateAsync: assign, isPending: pendingAssign} =
        useMutation({
            mutationFn: (params: AssignParams) =>
                assignWorkshiftSchedule(params.id, params.payload),
            onSuccess: () => {
                toast.success("Cấu hình thành công");
                form.reset();
            },
            onError: (error) => {
                toast.error(error.message);
            },
        });

    async function onSubmit(values: z.infer<typeof formSchema>) {
        const payload = {
            ...values,
            employee_ids: selectedEmployees.map(e => e.employee_id),
            manager_ids: selectedManagers.map(m => m.employee_id),
        };
        await assign({
            id: id!,
            payload
        });
    }

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
                            <InfoRow label="Văn phòng" value={workSchedule?.office?.office_name} className="mb-2"/>
                        </div>
                        <div className="w-full">
                            <InfoRow label="Ngày hiệu lực" value={workSchedule?.effective_date} className="mb-2"/>
                            <InfoRow label="Ngày hết hiệu lực" value={workSchedule?.expiration_date} className="mb-2"/>
                            <InfoRow label="Tình trạng" value={workSchedule?.status === "active" ? "Đang hiệu lực" : "Chưa hiệu lực"} className="mb-2"/>
                        </div>
                    </div>
                </div>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                        <div className="bg-white p-6 rounded-lg shadow-sm mb-8">
                            <p className="font-semibold text-lg mb-4">Quản lý lịch làm việc</p>

                            {/* Manager selection dropdown */}
                            <div className="mb-8">
                                <FormField
                                    control={form.control}
                                    name="manager_ids"
                                    render={({field}) => (
                                        <FormItem className="w-full">
                                            <FormLabel className="block mb-2 text-sm font-medium text-gray-700">
                                                Quản lý
                                            </FormLabel>
                                            <Select
                                                onValueChange={(value) => {
                                                    toggleManagerSelection(value);
                                                }}
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

                            {/* Selected managers display */}
                            {selectedManagers.length > 0 && (
                                <div className="mb-8">
                                    <p className="font-medium mb-3 text-gray-700">
                                        Quản lý đã chọn ({selectedManagers.length}):
                                    </p>
                                    <div className="border rounded-lg overflow-auto max-h-[180px]">
                                        <Table className="min-w-full">
                                            <TableHeader className="sticky top-0 bg-[#DB3B21]">
                                                <TableRow>
                                                    <TableHead className="px-4 py-3">Mã NV</TableHead>
                                                    <TableHead className="px-4 py-3">Họ và tên</TableHead>
                                                    <TableHead className="px-4 py-3">Phòng ban</TableHead>
                                                    <TableHead className="px-4 py-3">Cấp bậc</TableHead>
                                                    <TableHead className="px-4 py-3">Thao tác</TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {selectedManagers.map((manager) => (
                                                    <TableRow key={manager.employee_id} className="hover:bg-gray-50">
                                                        <TableCell
                                                            className="px-4 py-2">{manager.employee_id}</TableCell>
                                                        <TableCell className="px-4 py-2">{manager.full_name}</TableCell>
                                                        <TableCell
                                                            className="px-4 py-2">{manager.department.department_name}</TableCell>
                                                        <TableCell
                                                            className="px-4 py-2">{/* hierarchy_level */}</TableCell>
                                                        <TableCell className="px-4 py-2">
                                                            <button
                                                                type="button"
                                                                onClick={() => toggleManagerSelection(manager.employee_id)}
                                                                className="text-red-500 hover:text-red-700 text-sm font-medium"
                                                            >
                                                                Xóa
                                                            </button>
                                                        </TableCell>
                                                    </TableRow>
                                                ))}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Employee Section */}
                        <div className="bg-white p-6 rounded-lg shadow-sm">
                            <p className="font-semibold text-lg mb-4">Nhân viên</p>

                            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                                {/* Employee selection table */}
                                <div>
                                    <div className="flex justify-between items-center mb-3">
                                        <p className="font-medium text-gray-700">Chọn nhân viên:</p>
                                        <div className="flex items-center gap-2">
                                            <Button
                                                variant="outline"
                                                size="sm"
                                                onClick={toggleAllEmployeesSelection}
                                                className="text-xs h-8"
                                            >
                                                {selectedEmployees.length === filteredEmployees.length
                                                    ? "Bỏ chọn tất cả"
                                                    : "Chọn tất cả"}
                                            </Button>
                                            <div className="relative">
                                                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-400"/>
                                                <Input
                                                    type="search"
                                                    placeholder="Tìm kiếm..."
                                                    className="pl-8 h-8 w-[150px]"
                                                    value={searchTerm}
                                                    onChange={(e) => setSearchTerm(e.target.value)}
                                                />
                                            </div>
                                        </div>
                                    </div>

                                    {/* Filter controls */}
                                    <div className="flex flex-col sm:flex-row gap-3 mb-3">
                                        <div className="relative flex-1">
                                            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-400"/>
                                            <Input
                                                type="search"
                                                placeholder="Tìm theo tên hoặc mã NV"
                                                className="pl-8 h-8 w-full"
                                                value={searchTerm}
                                                onChange={(e) => setSearchTerm(e.target.value)}
                                            />
                                        </div>
                                        <Select
                                            value={selectedDepartment}
                                            onValueChange={setSelectedDepartment}
                                        >
                                            <SelectTrigger className="h-8">
                                                <SelectValue placeholder="Tất cả phòng ban"/>
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="all">Tất cả phòng ban</SelectItem>
                                                {departments.map(dept => (
                                                    <SelectItem key={dept} value={dept}>
                                                        {dept}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <Select
                                            value={selectedWorkType}
                                            onValueChange={setSelectedWorkType}
                                            className="sm:col-span-1"
                                        >
                                            <SelectTrigger className="h-8">
                                                <SelectValue placeholder="Tất cả loại hình làm việc"/>
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="all">Tất cả loại hình làm việc</SelectItem>
                                                {workTypes.map(type => (
                                                    <SelectItem key={type} value={type}>
                                                        {type}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                    </div>

                                    <div className="border rounded-lg overflow-hidden">
                                        <Table>
                                            <TableHeader className="bg-[#DB3B21]">
                                                <TableRow>
                                                    <TableHead className="w-12">
                                                        <Checkbox
                                                            checked={selectedEmployees.length === filteredEmployees.length && filteredEmployees.length > 0}
                                                            onCheckedChange={toggleAllEmployeesSelection}
                                                            className="ml-2"
                                                        />
                                                    </TableHead>
                                                    <TableHead>Mã NV</TableHead>
                                                    <TableHead>Họ và tên</TableHead>
                                                    <TableHead>Phòng ban</TableHead>
                                                    <TableHead>Cấp bậc</TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {filteredEmployees.length > 0 ? (
                                                    filteredEmployees.map((employee: Employee) => (
                                                        <TableRow key={employee.employee_id}
                                                                  className="hover:bg-gray-50">
                                                            <TableCell>
                                                                <Checkbox
                                                                    checked={isEmployeeSelected(employee.employee_id)}
                                                                    onCheckedChange={() => toggleEmployeeSelection(employee)}
                                                                />
                                                            </TableCell>
                                                            <TableCell>{employee.employee_id}</TableCell>
                                                            <TableCell>{employee.full_name}</TableCell>
                                                            <TableCell>{employee.department.department_name}</TableCell>
                                                            <TableCell>{/* hierarchy_level */}</TableCell>
                                                        </TableRow>
                                                    ))
                                                ) : (
                                                    <TableRow>
                                                        <TableCell colSpan={5}
                                                                   className="text-center py-4 text-gray-500">
                                                            {pendingGetEmployees ? 'Đang tải...' : 'Không tìm thấy nhân viên nào'}
                                                        </TableCell>
                                                    </TableRow>
                                                )}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </div>

                                {/* Selected employees table */}
                                <div>
                                    <p className="font-medium mb-3 text-gray-700">
                                        Nhân viên đã chọn ({selectedEmployees.length}):
                                    </p>
                                    {selectedEmployees.length > 0 ? (
                                        <div className="border rounded-lg overflow-hidden">
                                            <Table>
                                                <TableHeader className="bg-[#DB3B21]">
                                                    <TableRow>
                                                        <TableHead>Mã NV</TableHead>
                                                        <TableHead>Họ và tên</TableHead>
                                                        <TableHead>Phòng ban</TableHead>
                                                        <TableHead>Cấp bậc</TableHead>
                                                    </TableRow>
                                                </TableHeader>
                                                <TableBody>
                                                    {selectedEmployees.map((employee) => (
                                                        <TableRow key={employee.employee_id}
                                                                  className="hover:bg-gray-50">
                                                            <TableCell>{employee.employee_id}</TableCell>
                                                            <TableCell>{employee.full_name}</TableCell>
                                                            <TableCell>{employee.department.department_name}</TableCell>
                                                            <TableCell>{/* hierarchy_level */}</TableCell>
                                                        </TableRow>
                                                    ))}
                                                </TableBody>
                                            </Table>
                                        </div>
                                    ) : (
                                        <div className="border rounded-lg p-4 text-center text-gray-500">
                                            Chưa có nhân viên nào được chọn
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>

                        {/* Submit Button */}
                        <div className="flex justify-center items-center gap-5">
                            <div className="flex justify-end">
                                <Link to={`${PATH.WORK_SCHEDULE}`}>
                                    <Button variant="outline" type="button">Hủy bỏ</Button>
                                </Link>
                            </div>
                            <div className="flex justify-end">
                                <Button type="submit">
                                    {!pendingAssign ? "Lưu cấu hình" : <Loading />}
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