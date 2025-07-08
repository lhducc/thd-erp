import {useState, useEffect} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import {createAttendanceSettingApi, updateAttendanceSettingApi} from "@/apis/attendance-management.api.ts";
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
import type {AttendanceSetting} from "@/types/attendance.ts";
import {getAllOfficesApi} from "@/apis/office.api.ts";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";

// Define form schema
const formSchema = z.object({
    attendance_category_name: z.string().min(1, "Tên hình thức chấm công là bắt buộc"),
    office_id: z.string(),
    is_camera: z.boolean(),
    is_gps: z.boolean(),
    is_check_location: z.boolean(),
    scope: z.number().min(1, "Bán kính phải lớn hơn 0").optional(),
    status: z.string()
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
        data: offices,
        isPending: pendingOffices,
        refetch: refetchOffices,
    } = useQuery({
        queryKey: ["offices"],
        queryFn: getAllOfficesApi,
        gcTime: 0,
        staleTime: 0,
    });

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            attendance_category_name: data?.attendance_category_name || "",
            office_id: data?.office_id,
            is_camera: data?.is_camera || false,
            is_gps: data?.is_gps || false,
            is_check_location: data?.is_check_location || false,
            scope: data?.scope || 0,
            status: data?.status || "inactive"
        },
    });

    const radiusActive = form.watch("is_check_location");
    const isGPSActive = form.watch("is_gps");

    useEffect(() => {
        if (!radiusActive) {
            form.setValue("scope", undefined);
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
            mutationFn: ({ id, payload }: { id: string; payload: AttendanceSetting }) =>
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
            const payload: AttendanceSetting = {
                attendance_category_name: values.attendance_category_name,
                office_id: values.office_id,
                is_camera: values.is_camera,
                is_gps: values.is_gps,
                is_check_location: values.is_check_location,
                scope: values.scope,
                status: values.status
            };

            if (type === "edit" && data) {
                await updateAttendanceSetting({ id: data.attendance_category_id || "", payload });
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
                            name="attendance_category_name"
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
                            name="office_id"
                            render={({ field }) => (
                                <FormItem className={`flex gap-2`}>
                                    <FormLabel className={`w-[150px]`}>Văn phòng</FormLabel>
                                    <FormControl>
                                        <select
                                            {...field}
                                            className="block w-full rounded-md border border-gray-300 bg-white px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                                            onChange={(e) => {
                                                form.setValue("office_id",  e.target.value);
                                            }}
                                        >
                                            <option value="">Chọn văn phòng</option>
                                            {offices?.map((dept) => (
                                                <option key={dept.office_id} value={dept.office_id}>
                                                    {dept.office_name}
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
                                name="is_camera"
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
                                name="is_gps"
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
                            {
                                isGPSActive && <FormField
                                    control={form.control}
                                    name="is_check_location"
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
                            }

                            {radiusActive && (
                                <FormField
                                    control={form.control}
                                    name="scope"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Bán kính (mét)</FormLabel>
                                            <FormControl>
                                                <Input
                                                    type="number"
                                                    placeholder="Nhập bán kính cho phép"
                                                    {...field}
                                                    onChange={(e) => {
                                                        const val = e.target.value;
                                                        field.onChange(val === "" ? 0 : parseInt(val));
                                                    }}
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
                                    <FormItem>
                                        <FormLabel>Trạng thái</FormLabel>
                                        <Select onValueChange={field.onChange} value={field.value}>
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Trạng thái"/>
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                <SelectItem value="active">Hoạt động</SelectItem>
                                                <SelectItem value="inactive">Không hoạt động</SelectItem>
                                            </SelectContent>
                                        </Select>
                                        <FormMessage/>
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