import type {AttendanceSetting, PayloadAttendanceSetting} from "@/types";
import {useState, useEffect} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import {createAttendanceSettingApi, updateAttendanceSettingApi} from "@/apis/attendance.management.api";
import {toast} from "sonner";
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Plus} from "lucide-react";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form.tsx";
import {Input} from "@/components/ui/input.tsx";
import {useForm} from "react-hook-form";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod";
import {Checkbox} from "@/components/ui/checkbox.tsx";
import {Loader2} from "lucide-react";
import type {Department} from "@/types";
import {getAllDepartmentsApi} from "@/apis/department.api.ts";

// Define form schema
const formSchema = z.object({
    name: z.string().min(1, "Tên hình thức chấm công là bắt buộc"),
    department: z.object({
        department_id: z.string(),
        department_name: z.string()
    }),
    useCamera: z.boolean(),
    gpsEnabled: z.boolean(),
    radius: z.object({
        active: z.boolean(),
        meter: z.number().min(1, "Bán kính phải lớn hơn 0").optional()
    }),
    status: z.boolean()
});

type Props = {
    editBtn?: React.ReactNode;
    data?: AttendanceSetting;
    type?: "edit";
    refetch?: Function;
};

const AttendanceSettingForm = ({ editBtn, data, type, refetch }: Props) => {
    const [open, setOpen] = useState(false);

    const {
        data: departments,
        isPending: pendingGetDepartments,
        refetch: refetchDepartment,
    } = useQuery({
        queryKey: ["departments"],
        queryFn: getAllDepartmentsApi,
    });

    console.log(departments);

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: data?.name || "",
            department: data?.department || { department_id: "", department_name: "" },
            useCamera: data?.useCamera || false,
            gpsEnabled: data?.gpsEnabled || false,
            radius: {
                active: data?.radius?.active || false,
                meter: data?.radius?.meter || 0
            },
            status: data?.status || false
        },
    });

    const radiusActive = form.watch("radius.active");

    useEffect(() => {
        if (!radiusActive) {
            form.setValue("radius.meter", 0);
        }
    }, [radiusActive, form]);

    const { mutateAsync: createAttendanceSetting, isPending: pendingCreateAttendanceSetting } =
        useMutation({
            mutationFn: createAttendanceSettingApi,
            onSuccess: () => {
                refetch && refetch();
                setOpen(false);
                toast.success("Tạo hình thức chấm công thành công");
                form.reset();
            },
            onError: (error) => {
                toast.error(error.message);
            },
        });

    const { mutateAsync: updateAttendanceSetting, isPending: pendingUpdateAttendanceSetting } =
        useMutation({
            mutationFn: ({ id, payload }: { id: string; payload: PayloadAttendanceSetting }) =>
                updateAttendanceSettingApi(id, payload),
            onSuccess: () => {
                refetch && refetch();
                setOpen(false);
                toast.success("Cập nhật hình thức chấm công thành công");
                form.reset();
            },
            onError: (error) => {
                toast.error(error.message);
            },
        });

    const onSubmit = async (values: z.infer<typeof formSchema>) => {
        try {
            const payload: PayloadAttendanceSetting = {
                name: values.name,
                department_id: values.department.department_id,
                useCamera: values.useCamera,
                gpsEnabled: values.gpsEnabled,
                radius: {
                    active: values.radius.active,
                    meter: values.radius.active ? values.radius.meter : 0
                },
                status: values.status
            };

            if (type === "edit" && data) {
                await updateAttendanceSetting({ id: data.id || "", payload });
            } else {
                await createAttendanceSetting(payload);
            }
        } catch (error) {
            console.error("Error submitting form:", error);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger>
                {editBtn ? (
                    editBtn
                ) : (
                    <Button>
                        <Plus />
                        Tạo hình thức chấm công
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className="md:max-w-2xl">
                <DialogHeader>
                    {type === "edit" ? (
                        <DialogTitle>Chỉnh sửa hình thức chấm công</DialogTitle>
                    ) : (
                        <DialogTitle>Tạo hình thức chấm công</DialogTitle>
                    )}
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                        <FormField
                            control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Tên hình thức chấm công</FormLabel>
                                    <FormControl>
                                        <Input placeholder="Nhập tên hình thức chấm công" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="department.department_id"
                            render={({ field }) => (
                                <FormItem className={`flex gap-2`}>
                                    <FormLabel className={`w-[150px]`}>Văn phòng</FormLabel>
                                    <FormControl>
                                        <select
                                            {...field}
                                            className="block w-full rounded-md border border-gray-300 bg-white px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                                            onChange={(e) => {
                                                const selectedDept = departments?.find(d => d.department_id === e.target.value);
                                                form.setValue("department", {
                                                    department_id: e.target.value,
                                                    department_name: selectedDept?.department_name || ""
                                                });
                                            }}
                                        >
                                            <option value="">Chọn phòng ban</option>
                                            {departments?.map((dept) => (
                                                <option key={dept.department_id} value={dept.department_id}>
                                                    {dept.department_name}
                                                </option>
                                            ))}
                                        </select>
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <div className="grid grid-cols-1 gap-4">
                            <FormField
                                control={form.control}
                                name="useCamera"
                                render={({ field }) => (
                                    <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                                        <FormControl>
                                            <Checkbox
                                                checked={field.value}
                                                onCheckedChange={field.onChange}
                                            />
                                        </FormControl>
                                        <div className="space-y-1 leading-none">
                                            <FormLabel>Hình ảnh camera</FormLabel>
                                        </div>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="gpsEnabled"
                                render={({ field }) => (
                                    <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                                        <FormControl>
                                            <Checkbox
                                                checked={field.value}
                                                onCheckedChange={field.onChange}
                                            />
                                        </FormControl>
                                        <div className="space-y-1 leading-none">
                                            <FormLabel>Bản đồ GPS</FormLabel>
                                        </div>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="radius.active"
                                render={({ field }) => (
                                    <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                                        <FormControl>
                                            <Checkbox
                                                checked={field.value}
                                                onCheckedChange={field.onChange}
                                            />
                                        </FormControl>
                                        <div className="space-y-1 leading-none">
                                            <FormLabel>Bán kính cho phép</FormLabel>
                                        </div>
                                    </FormItem>
                                )}
                            />

                            {radiusActive && (
                                <FormField
                                    control={form.control}
                                    name="radius.meter"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Bán kính (mét)</FormLabel>
                                            <FormControl>
                                                <Input
                                                    type="number"
                                                    placeholder="Nhập bán kính cho phép"
                                                    {...field}
                                                    onChange={(e) => field.onChange(parseInt(e.target.value))}
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            )}

                            <FormField
                                control={form.control}
                                name="status"
                                render={({ field }) => (
                                    <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                                        <FormControl>
                                            <Checkbox
                                                checked={field.value}
                                                onCheckedChange={field.onChange}
                                            />
                                        </FormControl>
                                        <div className="space-y-1 leading-none">
                                            <FormLabel>Trạng thái</FormLabel>
                                        </div>
                                    </FormItem>
                                )}
                            />
                        </div>

                        <Button type="submit" disabled={pendingCreateAttendanceSetting || pendingUpdateAttendanceSetting}>
                            {pendingCreateAttendanceSetting || pendingUpdateAttendanceSetting ? (
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

export default AttendanceSettingForm;