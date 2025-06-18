import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Loader2, Plus } from "lucide-react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useEffect, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import type { ContractType } from "@/types/contract";
import {
    type ContractTypeCreate,
    createContractTypeApi,
    updateContractTypeApi,
} from "@/apis/contract-type.api";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

// Define enum values for TypeScript
const ContractGroup = {
    ConfirmTime: "Hợp đồng xác định thời hạn",
    NoTimeConfirmation: "Hợp đồng không xác định thời hạn",
    Trial: "Hợp đồng thử việc",
    VocationalTraining: "Hợp đồng đào tạo nghề",
    Service: "Hợp đồng dịch vụ",
} as const;

const Unit = {
    Year: "Năm",
    Month: "Tháng",
    Week: "Tuần",
    Day: "Ngày",
} as const;

const WorkingForm = {
    FullTime: "Toàn thời gian",
    PartTime: "Bán thời gian",
    Collaborator: "Cộng tác viên",
    Expert: "Chuyên gia",
    BySession: "Theo ca",
    ProductContract: "Khoán sản phẩm",
    WorkSubcontracting: "Khoán công việc",
} as const;

const formSchema = z.object({
    contract_type: z.string().min(1, "Tên loại hợp đồng không được để trống"),
    contract_group: z.string().min(1, "Nhóm hợp đồng không được để trống"),
    duration: z.number().min(0, "Thời hạn không được âm"),
    unit: z.string().min(1, "Đơn vị không được để trống"),
    working_form: z.string().min(1, "Hình thức làm việc không được để trống"),
});

type Props = {
    editBtn?: React.ReactNode;
    data?: ContractType;
    type?: "edit";
    refetch?: () => void;
};

const CreateContractType = ({ editBtn, data, type, refetch }: Props) => {
    const [open, setOpen] = useState(false);

    const { mutateAsync: createContractType, isPending: pendingCreate } = useMutation({
        mutationFn: (payload: ContractTypeCreate) => createContractTypeApi(payload),
        onSuccess: () => {
            refetch?.();
            setOpen(false);
            toast.success("Tạo loại hợp đồng thành công");
            form.reset();
        },
        onError: (error: any) => {
            toast.error(error.message);
        },
    });

    const { mutateAsync: updateContractType, isPending: pendingUpdate } = useMutation({
        mutationFn: ({ id, payload }: { id: string; payload: Partial<ContractTypeCreate> }) =>
            updateContractTypeApi(id, payload),
        onSuccess: () => {
            refetch?.();
            setOpen(false);
            toast.success("Cập nhật loại hợp đồng thành công");
            form.reset();
        },
        onError: (error: any) => {
            toast.error(error.message);
        },
    });

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            contract_type: "",
            contract_group: "",
            duration: 0,
            unit: "",
            working_form: "",
        },
    });

    useEffect(() => {
        if (data) {
            form.reset({
                contract_type: data.contract_type,
                contract_group: data.contract_group,
                duration: data.duration,
                unit: data.unit,
                working_form: data.working_form,
            });
        }
    }, [data]);

    async function onSubmit(values: z.infer<typeof formSchema>) {
        try {
            if (type === "edit") {
                await updateContractType({
                    id: data?.contract_type_id || "",
                    payload: values,
                });
            } else {
                await createContractType({
                    ...values,
                    // created_date: new Date().toISOString(),
                    // is_delete: false,
                });
            }
        } catch (error) {
            console.error("Error submitting form:", error);
        }
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger>
                {editBtn ? (
                    editBtn
                ) : (
                    <Button>
                        <Plus />
                        Thêm loại hợp đồng
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className="md:max-w-2xl">
                <DialogHeader>
                    {type === "edit" ? (
                        <DialogTitle>Chỉnh sửa loại hợp đồng</DialogTitle>
                    ) : (
                        <DialogTitle>Thêm loại hợp đồng</DialogTitle>
                    )}
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <div className="grid grid-cols-2 gap-4">
                            <FormField
                                control={form.control}
                                name="contract_type"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Tên loại hợp đồng</FormLabel>
                                        <FormControl>
                                            <Input placeholder="Nhập tên loại hợp đồng" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="contract_group"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Nhóm hợp đồng</FormLabel>
                                        <Select onValueChange={field.onChange} value={field.value}>
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Chọn nhóm hợp đồng" />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {Object.values(ContractGroup).map((group) => (
                                                    <SelectItem key={group} value={group}>
                                                        {group}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        <div className="grid grid-cols-3 gap-4">
                            <FormField
                                control={form.control}
                                name="duration"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Thời hạn</FormLabel>
                                        <FormControl>
                                            <Input
                                                type="number"
                                                min="0"
                                                placeholder="Nhập thời hạn"
                                                {...field}
                                                onChange={(e) => field.onChange(parseInt(e.target.value) || 0)}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="unit"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Đơn vị</FormLabel>
                                        <Select onValueChange={field.onChange} value={field.value}>
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Chọn đơn vị" />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {Object.values(Unit).map((unit) => (
                                                    <SelectItem key={unit} value={unit}>
                                                        {unit}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="working_form"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Hình thức làm việc</FormLabel>
                                        <Select onValueChange={field.onChange} value={field.value}>
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Chọn hình thức" />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {Object.values(WorkingForm).map((form) => (
                                                    <SelectItem key={form} value={form}>
                                                        {form}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        <Button
                            type="submit"
                            className="w-full mt-4"
                            disabled={pendingCreate || pendingUpdate}
                        >
                            {pendingCreate || pendingUpdate ? (
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

export default CreateContractType;