import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { useQuery } from "@tanstack/react-query";
import { getAllOfficesApi } from "@/apis/office.api";
import { getAllWorkshiftApi } from "@/apis/workshift.api";
import WeekScheduleSelector from "@/components/WeekScheduleSelector";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Calendar, Building2 } from "lucide-react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";

const formSchema = z.object({
    name: z.string().min(1, "Tên lịch không được để trống"),
    office: z.string().min(1, "Văn phòng không được để trống"),
    start_date: z.string().min(1, "Ngày bắt đầu không được để trống"),
    end_date: z.string().min(1, "Ngày kết thúc không được để trống"),
});

const WorkshiftSchedule = () => {
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: "",
            office: "",
            start_date: "",
            end_date: ""
        },
    });

    const { data: offices, isPending: pendingOffices } = useQuery({
        queryKey: ["offices"],
        queryFn: getAllOfficesApi,
    });

    const { data: workshifts, isPending: pendingWorkshifts } = useQuery({
        queryKey: ["workshifts"],
        queryFn: getAllWorkshiftApi,
    });

    const onSubmit = (values: z.infer<typeof formSchema>) => {
        console.log("Form submitted", values);
        // Add your form submission logic here
    }

    return (
        <div className="container mx-auto py-6">
            <Card className="border-0 shadow-none">
                <CardHeader className="pb-2">
                    <CardTitle className="text-2xl font-bold">Tạo lịch làm việc</CardTitle>
                    <CardDescription>
                        Thiết lập lịch làm việc mới cho nhân viên
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
                                            render={({ field }) => (
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
                                            )}
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
                                        {/*<h3 className="font-semibold text-lg">Thiết lập lịch làm việc</h3>*/}
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
                                            <WeekScheduleSelector workshifts={workshifts} />
                                        )}
                                    </CardContent>
                                </Card>
                            </div>
                        </CardContent>

                        <div className="flex justify-end gap-3 pt-6">
                            <Button variant="outline" type="button">
                                Hủy bỏ
                            </Button>
                            <Button type="submit">
                                Tạo lịch làm việc
                            </Button>
                        </div>
                    </form>
                </Form>
            </Card>
        </div>
    );
};

export default WorkshiftSchedule;