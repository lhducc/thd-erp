import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Loader2, Plus} from "lucide-react";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form.tsx";
import {Input} from "@/components/ui/input.tsx";
import {useState} from "react";
import {useMutation} from "@tanstack/react-query";
import {toast} from "sonner";
import {type Workshift, type WorkShiftRequest} from "@/types/workshift.ts";
import {useForm} from "react-hook-form";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";
import {Switch} from "@/components/ui/switch.tsx";
import {TimeOfDayEnum} from "@/types/workshift.ts";
import {createWorkshiftApi, updateWorkshiftApi} from "@/apis/workshift.api.ts";
import {DialogDescription} from "@radix-ui/react-dialog";

type Props = {
    editBtn?: React.ReactNode;
    data?: Workshift;
    type?: "edit" | "view";
    refetch?: () => void;
};

export const formSchema = z.object({
    workshift_name: z.string().min(1, "Tên ca làm việc là bắt buộc"),
    start_time: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)"),
    end_time: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)"),
    checkin_from: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)").optional().or(z.literal("")),
    checkin_to: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)").optional().or(z.literal("")),
    checkout_from: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)").optional().or(z.literal("")),
    checkout_to: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)").optional().or(z.literal("")),
    has_break: z.boolean(),
    break_start: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)").optional().or(z.literal("")),
    break_end: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$/, "Định dạng thời gian không hợp lệ (HH:mm:ss)").optional().or(z.literal("")),
    work_hours: z.number().min(0.1, "Số giờ làm việc phải lớn hơn 0"),
    work_day: z.number().min(0),
    coef_normal_day: z.number().min(0.1, "Hệ số phải lớn hơn 0"),
    coef_weekend: z.number().min(0.1, "Hệ số phải lớn hơn 0"),
    coef_holiday: z.number().min(0.1, "Hệ số phải lớn hơn 0"),
    time_of_day: z.nativeEnum(TimeOfDayEnum),
    effective_date: z.string().min(1, "Ngày hiệu lực là bắt buộc")
        .refine(val => !isNaN(new Date(val).getTime()), { message: "Ngày không hợp lệ" }),
    expiration_date: z.string().optional().or(z.literal(""))
        .refine(val => !val || !isNaN(new Date(val).getTime()), { message: "Ngày không hợp lệ" }),
    created_by: z.string().min(1, "Người tạo là bắt buộc"),
});

const WorkshiftForm = ({editBtn, data, type, refetch}: Props) => {
    const [open, setOpen] = useState(false);

    const formatDateForInput = (dateString?: string) => {
        if (!dateString) return '';
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return '';
        return date.toISOString().slice(0, 16); // Format: YYYY-MM-DDTHH:mm
    };

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            workshift_name: data?.workshift_name || "",
            start_time: data?.start_time || "08:00:00",
            end_time: data?.end_time || "17:00:00",
            checkin_from: data?.checkin_from || "",
            checkin_to: data?.checkin_to || "",
            checkout_from: data?.checkout_from || "",
            checkout_to: data?.checkout_to || "",
            has_break: data?.has_break || false,
            break_start: data?.break_start || "",
            break_end: data?.break_end || "",
            work_hours: data?.work_hours || 8,
            work_day: data?.work_day || 1,
            coef_normal_day: data?.coef_normal_day || 1.0,
            coef_weekend: data?.coef_weekend || 1.5,
            coef_holiday: data?.coef_holiday || 2.0,
            time_of_day: data?.time_of_day || TimeOfDayEnum.Morning,
            effective_date: data?.effective_date ? formatDateForInput(data.effective_date) : formatDateForInput(new Date().toISOString()),
            expiration_date: data?.expiration_date ? formatDateForInput(data.expiration_date) : "",
            created_by: data?.created_by || "current_user_id", // Replace with actual user ID
        },
    });

    const {mutateAsync: createWorkshift, isPending: pendingCreateWorkshift} =
        useMutation({
            mutationFn: createWorkshiftApi,
            onSuccess: () => {
                refetch?.();
                setOpen(false);
                toast.success("Tạo ca làm việc thành công");
                form.reset();
            },
            onError: (error) => {
                toast.error(error.message);
            },
        });

    const {mutateAsync: updateWorkshift, isPending: pendingUpdateWorkshift} =
        useMutation({
            mutationFn: ({id, payload}: { id: string; payload: WorkShiftRequest }) =>
                updateWorkshiftApi(id, payload),
            onSuccess: () => {
                refetch?.();
                setOpen(false);
                toast.success("Cập nhật ca làm việc thành công");
                form.reset();
            },
            onError: (error) => {
                toast.error(error.message);
            },
        });

    function onSubmit(values: z.infer<typeof formSchema>) {
        if (type === "view") return; // Prevent submission in view mode

        const formatForApi = (dateString: string) => {
            if (!dateString) return null;
            const date = new Date(dateString);
            return date.toISOString();
        };

        const payload: WorkShiftRequest = {
            ...values,
            work_day: Number(values.work_day) || 1,
            checkin_from: values.checkin_from || null,
            checkin_to: values.checkin_to || null,
            checkout_from: values.checkout_from || null,
            checkout_to: values.checkout_to || null,
            break_start: values.has_break ? (values.break_start || null) : null,
            break_end: values.has_break ? (values.break_end || null) : null,
            effective_date: formatForApi(values.effective_date) as string,
            expiration_date: values.expiration_date ? formatForApi(values.expiration_date) : null,
        };

        if (type === "edit" && data) {
            return updateWorkshift({id: data.workshift_id, payload});
        } else {
            return createWorkshift(payload);
        }
    }

    // Helper component for read-only fields
    const ReadOnlyField = ({ label, value }: { label: string; value: string | number | boolean | undefined }) => (
        <div className="space-y-2">
            <label className="text-sm font-medium leading-none">{label}</label>
            <div className="p-2 border rounded-md bg-gray-50">
                {typeof value === 'boolean' ? (value ? 'Có' : 'Không') : (value || '-')}
            </div>
        </div>
    );

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                {editBtn ? (
                    editBtn
                ) : (
                    <Button>
                        <Plus/>
                        Tạo ca làm việc
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className="md:w-[1335px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    {type === "edit" ? (
                        <DialogTitle>Chỉnh sửa ca làm việc</DialogTitle>
                    ) : type === "view" ? (
                        <DialogTitle>Xem chi tiết ca</DialogTitle>
                    ) : (
                        <DialogTitle>Tạo ca làm việc</DialogTitle>
                    )}
                    <DialogDescription>
                        This is a description of the dialog content.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        {/* Basic Information */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="md:grid-cols-1 md:gap-4 space-y-4">
                                {type === "view" ? (
                                    <ReadOnlyField label="Tên ca" value={form.watch('workshift_name')} />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="workshift_name"
                                        render={({field}) => (
                                            <FormItem>
                                                <FormLabel>Tên ca</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="Nhập tên ca làm việc" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}

                                {type === "view" ? (
                                    <ReadOnlyField
                                        label="Buổi làm việc"
                                        value={form.watch('time_of_day') === TimeOfDayEnum.Morning ? 'Buổi sáng' :
                                            form.watch('time_of_day') === TimeOfDayEnum.Afternoon ? 'Buổi chiều' :
                                                form.watch('time_of_day') === TimeOfDayEnum.Evening ? 'Buổi tối' : 'Cả ngày'}
                                    />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="time_of_day"
                                        render={({field}) => (
                                            <FormItem>
                                                <FormLabel>Buổi làm việc</FormLabel>
                                                <Select onValueChange={field.onChange} value={field.value}>
                                                    <FormControl>
                                                        <SelectTrigger>
                                                            <SelectValue placeholder="Chọn buổi làm việc"/>
                                                        </SelectTrigger>
                                                    </FormControl>
                                                    <SelectContent>
                                                        <SelectItem value={TimeOfDayEnum.Morning}>Buổi sáng</SelectItem>
                                                        <SelectItem value={TimeOfDayEnum.Afternoon}>Buổi chiều</SelectItem>
                                                        <SelectItem value={TimeOfDayEnum.Evening}>Buổi tối</SelectItem>
                                                        <SelectItem value={TimeOfDayEnum.AllDay}>Cả ngày</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}
                            </div>

                            <div className={`flex flex-col gap-2`}>
                                {
                                    type === "view" ? (
                                        <div>
                                            <p>Mã ca</p>
                                            <p className={`border border-blue-300 rounded-lg p-2 w-fit min-w-20`}>{data?.workshift_id}</p>
                                        </div>
                                    ) : null
                                }

                                <div>
                                    Ngày có hiệu lực
                                    {/* Date Settings */}
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        {type === "view" ? (
                                            <>
                                                <ReadOnlyField
                                                    label="Từ"
                                                    value={form.watch('effective_date') ? new Date(form.watch('effective_date')).toLocaleString() : ''}
                                                />
                                                <ReadOnlyField
                                                    label="Đến"
                                                    value={form.watch('expiration_date') ? new Date(form.watch('expiration_date')).toLocaleString() : 'Không có'}
                                                />
                                            </>
                                        ) : (
                                            <>
                                                <FormField
                                                    control={form.control}
                                                    name="effective_date"
                                                    render={({field}) => (
                                                        <FormItem className={`flex justify-center gap-3 h-fit`}>
                                                            <FormLabel>Từ:</FormLabel>
                                                            <FormControl>
                                                                <Input
                                                                    type="datetime-local"
                                                                    {...field}
                                                                />
                                                            </FormControl>
                                                            <FormMessage/>
                                                        </FormItem>
                                                    )}
                                                />

                                                <FormField
                                                    control={form.control}
                                                    name="expiration_date"
                                                    render={({field}) => (
                                                        <FormItem className={`flex justify-center gap-3 h-fit`}>
                                                            <FormLabel>Đến: </FormLabel>
                                                            <FormControl>
                                                                <Input
                                                                    type="datetime-local"
                                                                    {...field}
                                                                    value={field.value || ''}
                                                                />
                                                            </FormControl>
                                                            <FormMessage/>
                                                        </FormItem>
                                                    )}
                                                />
                                            </>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Time Settings */}
                        <div className="grid md:grid-cols-1 grid-cols-2 gap-4 justify-between mt-5">
                            <div className="flex justify-between gap-3 md:flex-row flex-col">
                                {type === "view" ? (
                                    <ReadOnlyField label="Giờ bắt đầu" value={form.watch('start_time')} />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="start_time"
                                        render={({field}) => (
                                            <FormItem className={`w-full`}>
                                                <FormLabel>Giờ bắt đầu</FormLabel>
                                                <FormControl>
                                                    <Input type="time" step="1" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}

                                {type === "view" ? (
                                    <ReadOnlyField label="Check-in từ" value={form.watch('checkin_from') || 'Không có'} />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="checkin_from"
                                        render={({field}) => (
                                            <FormItem className={`w-full`}>
                                                <FormLabel>Check-in từ (tùy chọn)</FormLabel>
                                                <FormControl>
                                                    <Input type="time" step="1" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}

                                {type === "view" ? (
                                    <ReadOnlyField label="Check-out từ" value={form.watch('checkout_from') || 'Không có'} />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="checkout_from"
                                        render={({field}) => (
                                            <FormItem className={`w-full`}>
                                                <FormLabel>Check-out từ (tùy chọn)</FormLabel>
                                                <FormControl>
                                                    <Input type="time" step="1" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}
                            </div>

                            <div className="flex justify-between gap-3 md:flex-row flex-col">
                                {type === "view" ? (
                                    <ReadOnlyField label="Giờ kết thúc" value={form.watch('end_time')} />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="end_time"
                                        render={({field}) => (
                                            <FormItem className={`w-full`}>
                                                <FormLabel>Giờ kết thúc</FormLabel>
                                                <FormControl>
                                                    <Input className={`w-full`} type="time" step="1" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}

                                {type === "view" ? (
                                    <ReadOnlyField label="Check-in đến" value={form.watch('checkin_to') || 'Không có'} />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="checkin_to"
                                        render={({field}) => (
                                            <FormItem className={`w-full`}>
                                                <FormLabel>Check-in đến (tùy chọn)</FormLabel>
                                                <FormControl>
                                                    <Input type="time" step="1" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}

                                {type === "view" ? (
                                    <ReadOnlyField label="Check-out đến" value={form.watch('checkout_to') || 'Không có'} />
                                ) : (
                                    <FormField
                                        control={form.control}
                                        name="checkout_to"
                                        render={({field}) => (
                                            <FormItem className={`w-full`}>
                                                <FormLabel>Check-out đến (tùy chọn)</FormLabel>
                                                <FormControl>
                                                    <Input type="time" step="1" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                )}
                            </div>
                        </div>

                        {/* Break Settings */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 my-5">
                            {type === "view" ? (
                                <ReadOnlyField label="Có nghỉ giữa ca" value={form.watch('has_break')} />
                            ) : (
                                <FormField
                                    control={form.control}
                                    name="has_break"
                                    render={({field}) => (
                                        <FormItem
                                            className="flex flex-row items-center justify-between rounded-lg border p-4">
                                            <div className="space-y-0.5">
                                                <FormLabel>Có nghỉ giữa ca</FormLabel>
                                            </div>
                                            <FormControl>
                                                <Switch
                                                    checked={field.value}
                                                    onCheckedChange={field.onChange}
                                                />
                                            </FormControl>
                                        </FormItem>
                                    )}
                                />
                            )}

                            <div className="flex justify-between gap-3">
                                {form.watch('has_break') && (
                                    <>
                                        {type === "view" ? (
                                            <>
                                                <ReadOnlyField
                                                    label="Giờ bắt đầu nghỉ"
                                                    value={form.watch('break_start') || 'Không có'}
                                                />
                                                <ReadOnlyField
                                                    label="Giờ kết thúc nghỉ"
                                                    value={form.watch('break_end') || 'Không có'}
                                                />
                                            </>
                                        ) : (
                                            <>
                                                <FormField
                                                    control={form.control}
                                                    name="break_start"
                                                    render={({field}) => (
                                                        <FormItem>
                                                            <FormLabel>Giờ bắt đầu nghỉ</FormLabel>
                                                            <FormControl>
                                                                <Input type="time" step="1" {...field} />
                                                            </FormControl>
                                                            <FormMessage/>
                                                        </FormItem>
                                                    )}
                                                />

                                                <FormField
                                                    control={form.control}
                                                    name="break_end"
                                                    render={({field}) => (
                                                        <FormItem>
                                                            <FormLabel>Giờ kết thúc nghỉ</FormLabel>
                                                            <FormControl>
                                                                <Input type="time" step="1" {...field} />
                                                            </FormControl>
                                                            <FormMessage/>
                                                        </FormItem>
                                                    )}
                                                />
                                            </>
                                        )}
                                    </>
                                )}
                            </div>
                        </div>

                        <div className="grid md:grid-cols-2 grid-cols-1 gap-10">
                            {type === "view" ? (
                                <ReadOnlyField label="Ngày công" value={form.watch('work_day')} />
                            ) : (
                                <FormField
                                    control={form.control}
                                    name="work_day"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Ngày công</FormLabel>
                                            <Select onValueChange={field.onChange} value={field.value}>
                                                <FormControl>
                                                    <SelectTrigger>
                                                        <SelectValue placeholder="Ngày công"/>
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    <SelectItem value={1}>1</SelectItem>
                                                    <SelectItem value={0.5}>0.5</SelectItem>
                                                    <SelectItem value={0}>0</SelectItem>
                                                </SelectContent>
                                            </Select>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            )}

                            {type === "view" ? (
                                <ReadOnlyField label="Giờ công" value={form.watch('work_hours')} />
                            ) : (
                                <FormField
                                    control={form.control}
                                    name="work_hours"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Giờ công</FormLabel>
                                            <FormControl>
                                                <Input
                                                    type="number"
                                                    step="0.1"
                                                    placeholder="Nhập số giờ làm việc"
                                                    {...field}
                                                    onChange={(e) => field.onChange(parseFloat(e.target.value))}
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            )}

                            {type === "view" ? (
                                <ReadOnlyField label="Hệ số ngày thường" value={form.watch('coef_normal_day')} />
                            ) : (
                                <FormField
                                    control={form.control}
                                    name="coef_normal_day"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Hệ số ngày thường</FormLabel>
                                            <FormControl>
                                                <Input
                                                    type="number"
                                                    step="0.1"
                                                    placeholder="Hệ số ngày thường"
                                                    {...field}
                                                    onChange={(e) => field.onChange(parseFloat(e.target.value))}
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            )}

                            {type === "view" ? (
                                <ReadOnlyField label="Hệ số cuối tuần" value={form.watch('coef_weekend')} />
                            ) : (
                                <FormField
                                    control={form.control}
                                    name="coef_weekend"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Hệ số cuối tuần</FormLabel>
                                            <FormControl>
                                                <Input
                                                    type="number"
                                                    step="0.1"
                                                    placeholder="Hệ số cuối tuần"
                                                    {...field}
                                                    onChange={(e) => field.onChange(parseFloat(e.target.value))}
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            )}

                            {type === "view" ? (
                                <ReadOnlyField label="Hệ số ngày lễ" value={form.watch('coef_holiday')} />
                            ) : (
                                <FormField
                                    control={form.control}
                                    name="coef_holiday"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Hệ số ngày lễ</FormLabel>
                                            <FormControl>
                                                <Input
                                                    type="number"
                                                    step="0.1"
                                                    placeholder="Hệ số ngày lễ"
                                                    {...field}
                                                    onChange={(e) => field.onChange(parseFloat(e.target.value))}
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            )}
                        </div>

                        {/* Hidden field for created_by */}
                        <FormField
                            control={form.control}
                            name="created_by"
                            render={({field}) => (
                                <FormItem className="hidden">
                                    <FormControl>
                                        <Input type="hidden" {...field} />
                                    </FormControl>
                                </FormItem>
                            )}
                        />

                        <div className="flex justify-center gap-5">
                            <Button className={`bg-[#EFEFEF] text-black mt-2`} onClick={() => setOpen(false)}>
                                {type === "view" ? "Đóng" : "Hủy bỏ"}
                            </Button>
                            {type !== "view" && (
                                <Button className="mt-2" type="submit" disabled={pendingCreateWorkshift || pendingUpdateWorkshift}>
                                    {pendingCreateWorkshift || pendingUpdateWorkshift ? (
                                        <Loader2 className="animate-spin"/>
                                    ) : type === "edit" ? (
                                        "Cập nhật"
                                    ) : (
                                        "Thêm"
                                    )}
                                </Button>
                            )}
                        </div>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
};

export default WorkshiftForm;