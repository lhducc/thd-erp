import { useState, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Calendar, Building2 } from "lucide-react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";
import {registerWorkScheduleRegisterApi, updateWorkScheduleRegisterApi} from "@/apis/work-schedule.api";
import type {WeekdaySelection} from "@/types/work-schedule.ts";
import {useNavigate, useParams} from "react-router-dom";
import {useOffice} from "@/query/useOffice.ts";
import {useQueryWorkshift} from "@/query/workshift.query.ts";
import {useWorkScheduleRegisterById} from "@/query/useWorkScheduleRegister.ts";
import {WorkshiftScheduler} from "@/components/WeekScheduleSelector.tsx";
import PATH from "@/constants/Path.ts";
import Loading from "@/components/Loading.tsx";

const formSchema = z.object({
    name: z.string().min(1, "Tên lịch không được để trống"),
    office: z.string().min(1, "Văn phòng không được để trống"),
    start_date: z.string().min(1, "Ngày bắt đầu không được để trống"),
    end_date: z.string().optional(),
});

const SetupWorkScheduleRegister = () => {
    const { id } = useParams();
    const isEditMode = id !== "create";
    const [weekdays, setWeekdays] = useState<WeekdaySelection[]>([]);
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: "",
            office: "",
            start_date: "",
            end_date: ""
        },
    });
    const { data: workSchedule, isLoading: pendingWorkSchedule } = useWorkScheduleRegisterById(id);
    const { data: offices, isLoading: pendingOffices } = useOffice()
    const { data: workshifts, isLoading: pendingWorkshifts } = useQueryWorkshift()
    const navigate = useNavigate();

    // Set form values when workSchedule data is loaded (edit mode)
    useEffect(() => {
        if (isEditMode && workSchedule && offices && !pendingOffices) {
            form.reset({
                name: workSchedule.work_schedule_name,
                office: workSchedule.office_id,
                start_date: workSchedule.effective_date?.split('T')[0],
                end_date: workSchedule.expiration_date?.split('T')[0],
            });

            // Convert weekdays data to WeekdaySelection format
            const initialWeekdays = workSchedule?.weekdays?.map(day => ({
                week_day: day.week_day,
                workshift_id: day.workshift_id,
                order: day.order
            }));
            setWeekdays(initialWeekdays);
        }
    }, [workSchedule, isEditMode, form, offices, pendingOffices]);

    const queryClient = useQueryClient();
    const registerMutation = useMutation({
        mutationFn: isEditMode ?
            (data: any) => updateWorkScheduleRegisterApi(id, data) :
            registerWorkScheduleRegisterApi,
        onSuccess: async () => {
            toast.success(isEditMode ? "Cập nhật lịch làm việc thành công" : "Đăng ký lịch làm việc thành công");
            if (!isEditMode) {
                form.reset();
                setWeekdays([]);
            }
            queryClient.removeQueries({queryKey: ["workScheduleRegisterById", id]})
            queryClient.removeQueries({ queryKey: ["workScheduleRegister"] })
            navigate("/setup-work-schedule-register")
        },
        onError: (error) => {
            toast.error(isEditMode ? "Cập nhật lịch làm việc thất bại" : "Đăng ký lịch làm việc thất bại");
            console.error("Error:", error);
        },
    });

    useEffect(() => {
        return () => {
            queryClient.removeQueries({
                queryKey: ["workScheduleRegisterById", id],
                exact: true
            });
        };
    }, [id]);

    const onSubmit = (values: z.infer<typeof formSchema>) => {
        if (weekdays.length === 0) {
            toast.error("Vui lòng chọn ít nhất một ca làm việc");
            return;
        }

        const payload = {
            work_schedule_name: values.name,
            office_id: values.office,
            effective_date: new Date(values.start_date).toISOString(),
            expiration_date: values?.end_date ? new Date(values?.end_date).toISOString() : null,
            weekdays: weekdays.map(day => ({
                week_day: day.week_day,
                workshift_id: day.workshift_id,
                order: day.order || 0
            })),
        };
        registerMutation.mutate(payload);
    };

    const handleWeekdaySelectionChange = (selections: WeekdaySelection[]) => {
        setWeekdays(selections);
    };
    useEffect(() => {
        if (isEditMode && workSchedule && offices) {
            const officeId = workSchedule.office?.office_id || workSchedule.office_id || "";

            // Set giá trị sau một khoảng delay nhỏ để đảm bảo Select component đã render
            setTimeout(() => {
                form.setValue("office", officeId);
            }, 100);
        }
    }, [workSchedule, offices, isEditMode, form]);

    if (pendingWorkSchedule && isEditMode && pendingOffices && pendingWorkshifts) {
        return (
            <Loading />
        );
    }

    return (
        <div className="container mx-auto py-6">
            <Card className="border-0 shadow-none">
                <CardHeader className="pb-2">
                    <CardTitle className="text-2xl font-bold">
                        {isEditMode ? "Chỉnh sửa lịch làm việc" : "Tạo lịch làm việc"}
                    </CardTitle>
                    <CardDescription>
                        {isEditMode ?
                            "Cập nhật thông tin lịch làm việc" :
                            "Thiết lập lịch làm việc mới cho nhân viên"}
                    </CardDescription>
                </CardHeader>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)}>
                        <CardContent className="grid grid-cols-1 lg:grid-cols-1 gap-8 p-0">
                            {/* Left Column - Form Inputs */}
                            <div className="space-y-6">
                                <Card className="border-0 shadow-none">
                                    <CardHeader className="pb-4">
                                        <div className="flex items-center gap-2">
                                            <Calendar className="w-5 h-5 text-primary" />
                                            <h3 className="font-semibold text-lg">Thông tin cơ bản</h3>
                                        </div>
                                    </CardHeader>
                                    <CardContent className="space-y-4">
                                        <FormField
                                            control={form.control}
                                            name="name"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Tên lịch làm việc</FormLabel>
                                                    <FormControl>
                                                        <Input
                                                            {...field}
                                                            placeholder="Nhập tên lịch làm việc"
                                                            className="h-10"
                                                        />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />

                                        <FormField
                                            control={form.control}
                                            name="office"
                                            render={({ field }) => {
                                                return(
                                                <FormItem>
                                                    <FormLabel>Văn phòng</FormLabel>
                                                    <Select onValueChange={field.onChange} value={field.value}>
                                                        <FormControl>
                                                            <SelectTrigger className="h-10 w-full">
                                                                <SelectValue placeholder="Chọn văn phòng" />
                                                            </SelectTrigger>
                                                        </FormControl>
                                                        <SelectContent>
                                                            {pendingOffices ? (
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
                                                )}}
                                        />

                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                            <FormField
                                                control={form.control}
                                                name="start_date"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormLabel>Ngày hiệu lực</FormLabel>
                                                        <FormControl>
                                                            <Input
                                                                {...field}
                                                                type="date"
                                                                className="h-10"
                                                            />
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />

                                            <FormField
                                                control={form.control}
                                                name="end_date"
                                                render={({ field }) => (
                                                    <FormItem>
                                                        <FormLabel>Ngày hết hiệu lực</FormLabel>
                                                        <FormControl>
                                                            <Input
                                                                {...field}
                                                                type="date"
                                                                className="h-10"
                                                            />
                                                        </FormControl>
                                                        <FormMessage />
                                                    </FormItem>
                                                )}
                                            />
                                        </div>
                                    </CardContent>
                                </Card>
                            </div>

                            {/* Right Column - Workshift Selector */}
                            <div className="space-y-2">
                                <Card className="border-0 shadow-none">
                                    <CardHeader className="pb-4">
                                        <CardDescription>
                                            Chọn ngày và ca làm việc sẽ áp dụng
                                        </CardDescription>
                                    </CardHeader>
                                    <CardContent>
                                        {pendingWorkshifts ? (
                                            <div className="space-y-4">
                                                <Skeleton className="h-[300px] w-full rounded-lg" />
                                            </div>
                                        ) : (
                                            <WorkshiftScheduler
                                                workshifts={workshifts || []}
                                                onSelectionChange={handleWeekdaySelectionChange}
                                                initialSelections={weekdays}
                                            />
                                        )}
                                    </CardContent>
                                </Card>
                            </div>
                        </CardContent>

                        <div className="flex justify-end gap-3 pt-6">
                            <Button
                                variant="outline"
                                type="button"
                                onClick={() => {
                                    navigate(PATH.WORK_SCHEDULE_REGISTER)
                                }}
                            >
                                Hủy bỏ
                            </Button>
                            <Button
                                type="submit"
                                disabled={registerMutation.isPending}
                            >
                                {registerMutation.isPending ? "Đang xử lý..." :
                                    isEditMode ? "Cập nhật lịch làm việc" : "Tạo lịch làm việc"}
                            </Button>
                        </div>
                    </form>
                </Form>
            </Card>
        </div>
    );
};

export default SetupWorkScheduleRegister;