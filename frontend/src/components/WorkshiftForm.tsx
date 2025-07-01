import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Loader2, Plus} from "lucide-react";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form.tsx";
import {Input} from "@/components/ui/input.tsx";
import {useState} from "react";
import {useMutation} from "@tanstack/react-query";
import {createWorkshiftApi, updateWorkshiftApi} from "@/apis/workshift.api.ts";
import {toast} from "sonner";
import type {WorkShift, WorkShiftRequest} from "@/types/Workshift.ts";
import {useForm} from "react-hook-form";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";
import {Switch} from "@/components/ui/switch.tsx";
import {WorkDayEnum, TimeOfDayEnum} from "@/types/Workshift.ts";

type Props = {
    editBtn?: React.ReactNode;
    data?: WorkShift;
    type?: "edit";
    refetch?: Promise<void>;
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
    work_day: z.nativeEnum(WorkDayEnum),
    coef_normal_day: z.number().min(0.1, "Hệ số phải lớn hơn 0"),
    coef_weekend: z.number().min(0.1, "Hệ số phải lớn hơn 0"),
    coef_holiday: z.number().min(0.1, "Hệ số phải lớn hơn 0"),
    time_of_day: z.nativeEnum(TimeOfDayEnum),
    effective_date: z.string().min(1, "Ngày hiệu lực là bắt buộc"),
    expiration_date: z.string().optional().or(z.literal("")),
    created_by: z.string().min(1, "Người tạo là bắt buộc"),
});

const WorkshiftForm = ({editBtn, data, type, refetch}: Props) => {
    const [open, setOpen] = useState(false);

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
            work_day: data?.work_day || WorkDayEnum.Weekday,
            coef_normal_day: data?.coef_normal_day || 1.0,
            coef_weekend: data?.coef_weekend || 1.5,
            coef_holiday: data?.coef_holiday || 2.0,
            time_of_day: data?.time_of_day || TimeOfDayEnum.Morning,
            effective_date: data?.effective_date || new Date().toISOString().split('T')[0],
            expiration_date: data?.expiration_date || "",
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
        const payload: WorkShiftRequest = {
            ...values,
            checkin_from: values.checkin_from || null,
            checkin_to: values.checkin_to || null,
            checkout_from: values.checkout_from || null,
            checkout_to: values.checkout_to || null,
            break_start: values.has_break ? (values.break_start || null) : null,
            break_end: values.has_break ? (values.break_end || null) : null,
            expiration_date: values.expiration_date || null,
        };

        if (type === "edit" && data) {
            return updateWorkshift({id: data.workshift_id, payload});
        } else {
            return createWorkshift(payload);
        }
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger>
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
                    ) : (
                        <DialogTitle>Tạo ca làm việc</DialogTitle>
                    )}
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        {/* Basic Information */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="md:grid-cols-1 md:gap-4 space-y-4">
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
                                                </SelectContent>
                                            </Select>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            </div>

                            <div>
                                Ngày có hiệu lực
                                {/* Date Settings */}
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <FormField
                                        control={form.control}
                                        name="effective_date"
                                        render={({field}) => (
                                            <FormItem className={`flex justify-center gap-3 h-fit`}>
                                                <FormLabel>Từ:</FormLabel>
                                                <FormControl>
                                                    <Input type="date" {...field} />
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
                                                    <Input type="date" {...field} />
                                                </FormControl>
                                                <FormMessage/>
                                            </FormItem>
                                        )}
                                    />
                                </div>
                            </div>
                        </div>

                        {/* Time Settings */}
                        <div className="grid md:grid-cols-1 grid-cols-2 gap-4 justify-between">
                            <div className="flex justify-between gap-3">
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
                            </div>

                            <div className="flex justify-between gap-3">
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
                            </div>

                        </div>

                        {/* Break Settings */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
                            <div className="flex justify-between gap-3">
                                {form.watch('has_break') && (
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
                            </div>
                        </div>
                        <div className="grid md:grid-cols-2 grid-cols-1 gap-10">
                            <FormField
                                control={form.control}
                                name="work_day"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel>Ngày công</FormLabel>
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

                        <Button className={`bg-[#EFEFEF] text-black`} onClick={() => setOpen(false)} >
                            Hủy bỏ
                        </Button>
                        <Button type="submit" disabled={pendingCreateWorkshift || pendingUpdateWorkshift}>
                            {pendingCreateWorkshift || pendingUpdateWorkshift ? (
                                <Loader2 className="animate-spin"/>
                            ) : type === "edit" ? (
                                "Cập nhật"
                            ) : (
                                "Thêm"
                            )}
                        </Button>
                        </div>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
        ;
};

export default WorkshiftForm;