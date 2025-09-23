import {useMutation, useQueryClient} from "@tanstack/react-query";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {z} from "zod";
import type {WorkShift} from "@/types/work-schedule.ts";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {toast} from "sonner";
import {createWorkshiftSchedule, updateWorkshiftSchedule} from "@/apis/work-schedule.api.ts";
import {Link, useNavigate, useParams} from "react-router-dom";
import PATH from "@/constants/Path.ts";
import {useOffice} from "@/query/useOffice.ts";
import {useQueryWorkshift} from "@/query/workshift.query.ts";
import {useWorkScheduleById} from "@/query/useWorkSchedule.ts";
import {useEffect, useState} from "react";
import Loading from "@/components/Loading.tsx";
import {Skeleton} from "@/components/ui/skeleton.tsx";
import {Building2} from "lucide-react";

// Strongly typed enums
export const RepeatTypeEnum = z.enum(["weekly", "monthly"]);
export type RepeatTypeEnum = z.infer<typeof RepeatTypeEnum>;

// Strongly typed weekday map
const WEEKDAY_MAP: Record<number, string> = {
    2: "monday",
    3: "tuesday",
    4: "wednesday",
    5: "thursday",
    6: "friday",
    7: "saturday",
    8: "sunday",
};

// Strongly typed week days
type WeekDay = "monday" | "tuesday" | "wednesday" | "thursday" | "friday" | "saturday" | "sunday";

// Strongly typed day schema với conditional validation
const daySchema = z.object({
    day_of_week: z.number().min(2).max(8),
    enabled: z.boolean(),
    shift_count: z.number().min(1),
    shifts: z.array(z.string())
}).superRefine((data, ctx) => {
    // CHỈ validate khi day được enabled
    if (data.enabled) {
        // Kiểm tra không có shift nào rỗng
        data.shifts.forEach((shift, index) => {
            if (shift === "") {
                ctx.addIssue({
                    code: z.ZodIssueCode.custom,
                    message: "Ca làm việc là bắt buộc",
                    path: [`shifts`, index]
                });
            }
        });

        // Kiểm tra có ít nhất 1 shift không rỗng
        if (data.shifts.length === 0 || data.shifts.every(shift => shift === "")) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: "Ít nhất một ca làm việc phải được chọn",
                path: ['shifts']
            });
        }
    }
});

type Day = z.infer<typeof daySchema>;

// Strongly typed form schema với conditional validation
const auto_schedule = z.object({
    work_schedule_name: z.string().min(1, "Tên lịch làm việc là bắt buộc"),
    office_id: z.string().optional(),
    repeat_type: RepeatTypeEnum,
    repeat_cycle: z.coerce.number().min(1, "Chu kỳ lặp phải lớn hơn 0"),
    effective_date: z.string().min(1, "Ngày hiệu lực là bắt buộc"),
    expiration_date: z.string().optional(),
    days: z.array(daySchema).superRefine((days, ctx) => {
        // CHỈ kiểm tra có ít nhất 1 ngày được enabled và có shift hợp lệ
        const hasValidDay = days.some(day => 
            day.enabled && 
            day.shifts.length > 0 && 
            day.shifts.some(shift => shift !== "")
        );
        
        if (!hasValidDay) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: "Ít nhất một ngày phải được chọn và có ca làm việc",
                path: ['days']
            });
        }
    }),
});

type AutoScheduleFormValues = z.infer<typeof auto_schedule>;

// Strongly typed work schedule payload
type WorkSchedulePayload = Omit<AutoScheduleFormValues, 'days'> & {
    status: string;
    weekdays: {
        week_day: WeekDay;
        workshift_id: string;
    }[];
    effective_date?: string;
    expiration_date?: string;
};

const SetupWorkScheduleAuto = () => {
    const { id } = useParams<{ id?: string }>();
    const isEditMode = id !== "create" && id !== undefined;

    const { data: offices } = useOffice();
    const { data: workshifts, isPending: pendingWorkshift } = useQueryWorkshift();
    const { data: workSchedule, isPending: pendingWorkSchedule } = useWorkScheduleById(id);

    const queryClient = useQueryClient();
    const navigate = useNavigate();
    const [isSubmitting, setIsSubmitting] = useState(false);

    const { mutateAsync: createWorkSchedule, isPending } = useMutation({
        mutationFn: (data: WorkSchedulePayload) => {
            if (isEditMode) {
                return updateWorkshiftSchedule(id, data);
            } else {
                return createWorkshiftSchedule(data);
            }
        },
        onSuccess: async () => {
            toast.success(isEditMode ? "Cập nhật lịch làm việc thành công" : "Tạo lịch làm việc thành công");
            await queryClient.invalidateQueries({queryKey: ["work-schedule"]});
            await queryClient.invalidateQueries({queryKey: ["work-schedules", id]});
            navigate(PATH.WORK_SCHEDULE);
        },
        onError: (error: Error) => {
            toast.error(error.message);
        },
    });

    const form = useForm<AutoScheduleFormValues>({
        resolver: zodResolver(auto_schedule),
        defaultValues: {
            work_schedule_name: "",
            office_id: "",
            repeat_type: RepeatTypeEnum.enum.weekly,
            repeat_cycle: 1,
            effective_date: new Date().toISOString().split('T')[0],
            expiration_date: "",
            days: Array.from({ length: 7 }, (_, i) => ({
                day_of_week: i + 2,
                enabled: false,
                shift_count: 1,
                shifts: [""],
            })),
        },
    });

    // Debug chi tiết validation errors
    useEffect(() => {
        if (Object.keys(form.formState.errors).length > 0) {
            console.log("=== FORM VALIDATION ERRORS ===");
            console.log("Form Errors:", JSON.stringify(form.formState.errors, null, 2));
            
            // Log chi tiết lỗi của từng day
            form.watch('days').forEach((day, index) => {
                if (form.formState.errors.days?.[index]) {
                    console.log(`Day ${index} errors:`, form.formState.errors.days[index]);
                }
            });
        }
    }, [form.formState.errors]);

    useEffect(() => {
        if (isEditMode && workSchedule && workshifts && offices) {
            const officeId = workSchedule.office_id || workSchedule.office?.office_id || "";

            const days = Array.from({ length: 7 }, (_, i) => {
                const dayOfWeek = i + 2;
                const weekdayKey = WEEKDAY_MAP[dayOfWeek];
                const shiftsForDay = workSchedule.weekdays
                    ?.filter(w => w.week_day === weekdayKey)
                    .map(w => w.work_shift) || [];

                // Đảm bảo shifts luôn có giá trị hợp lệ
                const validShifts = shiftsForDay
                    .filter(shift => shift?.workshift_id)
                    .map(shift => shift.workshift_id);

                return {
                    day_of_week: dayOfWeek,
                    enabled: validShifts.length > 0,
                    shift_count: Math.max(validShifts.length, 1),
                    shifts: validShifts.length > 0 ? validShifts : [""]
                };
            });

            // Reset form với delay nhỏ để đảm bảo component đã render
            setTimeout(() => {
                form.reset({
                    work_schedule_name: workSchedule.work_schedule_name || "",
                    office_id: officeId,
                    repeat_type: (workSchedule.repeat_type as RepeatTypeEnum) || RepeatTypeEnum.enum.weekly,
                    repeat_cycle: workSchedule.repeat_cycle || 1,
                    effective_date: workSchedule.effective_date?.split('T')[0] || new Date().toISOString().split('T')[0],
                    expiration_date: workSchedule?.expiration_date?.split('T')[0] || "",
                    days: days
                });
            }, 150);
        }
    }, [workSchedule, offices, workshifts, isEditMode, form]);

    const isOverlap = (start1: string, end1: string, start2: string, end2: string): boolean => {
        return !(end1 <= start2 || start1 >= end2);
    };

    const handleToggleDay = (index: number) => {
        const currentDays = [...form.getValues('days')];
        const newEnabledState = !currentDays[index].enabled;

        currentDays[index].enabled = newEnabledState;
        
        if (!newEnabledState) {
            // Nếu bỏ tick, reset shifts (không cần validate vì disabled)
            currentDays[index].shifts = [""];
            currentDays[index].shift_count = 1;
        } else {
            // Nếu tick, thêm shift mặc định hợp lệ
            const firstValidShift = workshifts?.find(ws => ws.workshift_id)?.workshift_id;
            currentDays[index].shifts = firstValidShift ? [firstValidShift] : [""];
            currentDays[index].shift_count = 1;
        }

        // Delay validation để UI kịp cập nhật
        setTimeout(() => {
            form.setValue("days", currentDays, { shouldValidate: true });
        }, 50);
    };

    const onSubmit = async (data: AutoScheduleFormValues) => {
        if (isSubmitting) return;
        
        setIsSubmitting(true);
        try {
            // Thêm delay nhỏ trước khi validate
            await new Promise(resolve => setTimeout(resolve, 100));
            
            // Validate form trước khi submit
            const isValid = await form.trigger();
            if (!isValid) {
                console.log("Validation Errors:", form.formState.errors);
                toast.error("Vui lòng kiểm tra lại thông tin các trường");
                return;
            }

            // LỌC CHỈ những ngày được enabled và có shift hợp lệ
            const enabledDays = data.days.filter(day => 
                day.enabled && 
                day.shifts.length > 0 && 
                day.shifts.some(shift => shift !== "")
            );

            if (enabledDays.length === 0) {
                toast.error("Ít nhất một ngày phải được chọn và có ca làm việc");
                return;
            }

            // Tạo weekdays chỉ từ những ngày enabled
            const weekdays = enabledDays
                .flatMap((day) =>
                    day.shifts
                        .filter(shift => shift !== "") // Chỉ lấy shifts không rỗng
                        .map((shiftId) => ({
                            week_day: WEEKDAY_MAP[day.day_of_week] as WeekDay,
                            workshift_id: shiftId,
                        }))
                );

            if (weekdays.length === 0) {
                toast.error("Vui lòng chọn ít nhất một ca làm việc");
                return;
            }

            const finalPayload: WorkSchedulePayload = {
                work_schedule_name: data.work_schedule_name,
                office_id: data.office_id,
                repeat_type: data.repeat_type,
                repeat_cycle: data.repeat_cycle,
                status: "active",
                weekdays,
                effective_date: data.effective_date ? `${data.effective_date}T00:00:00Z` : undefined,
                expiration_date: data.expiration_date ? `${data.expiration_date}T00:00:00Z` : undefined,
            };

            console.log("Submitting payload:", finalPayload);
            await createWorkSchedule(finalPayload);
        } catch (error) {
            console.error("Submit error:", error);
            toast.error("Có lỗi xảy ra khi gửi dữ liệu");
        } finally {
            setIsSubmitting(false);
        }
    };

    if (isEditMode && (pendingWorkSchedule || !offices || pendingWorkshift)) {
        return <Loading />;
    }

    const renderShiftSelect = (day: Day, dayIndex: number, shiftIndex: number) => {
        // Fallback khi workshifts chưa load xong
        if (!workshifts || workshifts.length === 0) {
            return (
                <SelectTrigger className="w-52" disabled>
                    <SelectValue placeholder="Đang tải ca..." />
                </SelectTrigger>
            );
        }

        const selectedShifts = day.shifts
            .filter((_, i) => i !== shiftIndex)
            .map((id) => workshifts?.find((ws) => ws.workshift_id === id))
            .filter(Boolean) as WorkShift[];

        const availableShifts = workshifts?.filter((ws) =>
            !selectedShifts.some((s) =>
                isOverlap(ws.start_time, ws.end_time, s.start_time, s.end_time)
            )
        ) || [];

        return (
            <FormField
                control={form.control}
                name={`days.${dayIndex}.shifts.${shiftIndex}`}
                render={({ field }) => (
                    <FormItem>
                        <Select
                            onValueChange={field.onChange}
                            value={field.value}
                        >
                            <FormControl>
                                <SelectTrigger className="w-52">
                                    <SelectValue placeholder="Chọn ca" />
                                </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                                {availableShifts.map((ws) => (
                                    <SelectItem key={ws.workshift_id} value={ws.workshift_id}>
                                        {ws.workshift_name} ({ws.start_time} - {ws.end_time})
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        <FormMessage />
                    </FormItem>
                )}
            />
        );
    };

    const renderDayRow = (day: Day, index: number) => {
        const dayNames = ["Thứ Hai", "Thứ Ba", "Thứ Tư", "Thứ Năm", "Thứ Sáu", "Thứ Bảy", "Chủ Nhật"];
        
        return (
            <div key={index} className="flex items-center gap-4 border-b pb-2">
                <input
                    type="checkbox"
                    checked={day.enabled}
                    onChange={() => handleToggleDay(index)}
                />
                <span className="w-16">{dayNames[index]}</span>

                {day.enabled && (
                    <>
                        <Input
                            type="number"
                            min="1"
                            className="w-16"
                            value={day.shift_count}
                            onChange={(e) => {
                                const days = [...form.getValues('days')];
                                const newCount = Math.max(1, parseInt(e.target.value) || 1);
                                days[index].shift_count = newCount;

                                // Adjust shifts array length
                                if (newCount > days[index].shifts.length) {
                                    days[index].shifts = [
                                        ...days[index].shifts,
                                        ...Array(newCount - days[index].shifts.length).fill("")
                                    ];
                                } else if (newCount < days[index].shifts.length) {
                                    days[index].shifts = days[index].shifts.slice(0, newCount);
                                }

                                // Delay validation để UI kịp cập nhật
                                setTimeout(() => {
                                    form.setValue('days', days, { shouldValidate: true });
                                }, 50);
                            }}
                        />

                        {Array.from({ length: day.shift_count }).map((_, sIdx) => (
                            renderShiftSelect(day, index, sIdx)
                        ))}
                    </>
                )}
            </div>
        );
    };

    return (
        <div className="space-y-4">
            <h3 className="text-lg font-medium">{isEditMode ? "Chỉnh sửa lịch làm việc" : "Tạo lịch làm việc"}</h3>
            <hr className="border-gray-200" />
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                    <div className="space-y-4 flex gap-5">
                        <div className="space-y-4 w-full">
                            <FormField
                                control={form.control}
                                name="work_schedule_name"
                                render={({ field }) => (
                                    <FormItem className="w-full">
                                        <FormLabel>Tên lịch làm việc</FormLabel>
                                        <FormControl>
                                            <Input {...field} placeholder="Nhập tên lịch làm việc" />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="office_id"
                                render={({ field }) => {
                                    return (
                                        <FormItem>
                                            <FormLabel>Văn phòng</FormLabel>
                                            <Select
                                                onValueChange={field.onChange}
                                                value={field.value || ""}
                                                disabled={isEditMode}
                                                defaultValue={field.value}
                                            >
                                                <FormControl>
                                                    <SelectTrigger className="h-10 w-full">
                                                        <SelectValue placeholder="Chọn văn phòng" />
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    {!offices ? (
                                                        <SelectItem value="loading" disabled>
                                                            <div className="flex items-center gap-2">
                                                                <Skeleton className="h-4 w-4 rounded-full" />
                                                                <span>Đang tải...</span>
                                                            </div>
                                                        </SelectItem>
                                                    ) : (
                                                        offices?.map((item) => (
                                                            <SelectItem key={item.office_id} value={item.office_id}>
                                                                <div className="flex items-center gap-2">
                                                                    <Building2 className="w-4 h-4" />
                                                                    {item.office_name}
                                                                </div>
                                                            </SelectItem>
                                                        ))
                                                    )}
                                                </SelectContent>
                                            </Select>
                                            <FormMessage />
                                        </FormItem>
                                    );
                                }}
                            />

                            <div className="w-full flex gap-5">
                                <FormField
                                    control={form.control}
                                    name="repeat_type"
                                    render={({ field }) => (
                                        <FormItem className="w-full">
                                            <FormLabel>Loại lặp</FormLabel>
                                            <Select onValueChange={field.onChange} value={field.value}>
                                                <FormControl>
                                                    <SelectTrigger>
                                                        <SelectValue placeholder="Chọn loại lặp" />
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    <SelectItem value={RepeatTypeEnum.enum.weekly}>Tuần</SelectItem>
                                                    <SelectItem value={RepeatTypeEnum.enum.monthly}>Tháng</SelectItem>
                                                </SelectContent>
                                            </Select>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <FormField
                                    control={form.control}
                                    name="repeat_cycle"
                                    render={({ field }) => (
                                        <FormItem className="w-full">
                                            <FormLabel>Chu kỳ lặp</FormLabel>
                                            <FormControl>
                                                <Input
                                                    {...field}
                                                    type="number"
                                                    min="1"
                                                    placeholder="Nhập chu kỳ lặp"
                                                    onChange={(e) => field.onChange(parseInt(e.target.value) || 1)}
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </div>

                            <div className="w-full flex gap-5">
                                <FormField
                                    control={form.control}
                                    name="effective_date"
                                    render={({ field }) => (
                                        <FormItem className="w-full">
                                            <FormLabel>Ngày hiệu lực</FormLabel>
                                            <FormControl>
                                                <Input
                                                    {...field}
                                                    type="date"
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <FormField
                                    control={form.control}
                                    name="expiration_date"
                                    render={({ field }) => (
                                        <FormItem className="w-full">
                                            <FormLabel>Ngày hết hiệu lực</FormLabel>
                                            <FormControl>
                                                <Input
                                                    {...field}
                                                    type="date"
                                                    min={form.watch('effective_date')}
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </div>
                        </div>

                        <div className="space-y-2 w-full p-5 bg-white">
                            {form.watch('days').map((day, index) => renderDayRow(day, index))}
                            {/* Hiển thị lỗi tổng cho days */}
                            {form.formState.errors.days && (
                                <p className="text-sm text-destructive">
                                    {form.formState.errors.days.message}
                                </p>
                            )}
                        </div>
                    </div>

                    <div className="flex justify-center items-center gap-5">
                        <div className="flex justify-end">
                            <Link to={PATH.WORK_SCHEDULE}>
                                <Button variant="outline" type="button">Hủy bỏ</Button>
                            </Link>
                        </div>
                        <div className="flex justify-end">
                            {!isPending && !isSubmitting ? (
                                <Button type="submit" disabled={isSubmitting}>
                                    {isEditMode ? "Cập nhật" : "Tạo"} lịch làm việc
                                </Button>
                            ) : <Loading />}
                        </div>
                    </div>
                </form>
            </Form>
        </div>
    );
};

export default SetupWorkScheduleAuto;