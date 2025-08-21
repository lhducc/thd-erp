import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import {Building2, Loader2, Plus} from "lucide-react";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod";
import {useForm} from "react-hook-form";
import {Button} from "@/components/ui/button";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {useEffect, useState} from "react";
import {useMutation} from "@tanstack/react-query";
import {toast} from "sonner";
import {createDecisionTypeApi, updateDecisionTypeApi} from "@/apis/decistion-type.api.ts";
import type {DecisionType} from "@/types/decistion-type.ts";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";
import {Skeleton} from "@/components/ui/skeleton.tsx";
import Asterisk from "@/components/ui/Asterisk.tsx";

const formSchema = z.object({
    decision_type: z.string().min(1, "Tên loại quyết định không được để trống"),
    // decision_group: z.string().min(1, "Nhóm quyết định không được để trống"),
    description: z.string().optional(),
});

type Props = {
    editBtn?: React.ReactNode;
    data?: DecisionType;
    type?: "edit";
    refetch?: () => void;
    buttonText?: string;
};

// const decistionGroup = [
//     "Hình thức khen thưởng",
//     "Hình thức kỷ luật",
//     "Lý do điều chuyển",
//     "Lý do bổ nhiệm",
//     "Lý do miễn nhiệm",
//     "Lý do chấm dứt HĐLĐ"
// ]

const CreateDecisionType = ({editBtn, data, type, refetch, buttonText = "Thêm loại quyết định"}: Props) => {
    const [open, setOpen] = useState(false);

    const {mutateAsync: createDecisionType, isPending: pendingCreate} = useMutation({
        mutationFn: (payload: Omit<DecisionType, "decision_type_id" | "is_delete" | "created_date">) =>
            createDecisionTypeApi({
                ...payload,
                is_delete: false,
                created_date: new Date().toISOString(),
            }),
        onSuccess: () => {
            refetch?.();
            setOpen(false);
            toast.success("Tạo loại quyết định thành công");
            form.reset();
        },
        onError: (error: any) => {
            toast.error(error.message);
        },
    });

    const {mutateAsync: updateDecisionType, isPending: pendingUpdate} = useMutation({
        mutationFn: ({id, payload}: { id: string; payload: Partial<DecisionType> }) =>
            updateDecisionTypeApi(id, payload),
        onSuccess: () => {
            refetch?.();
            setOpen(false);
            toast.success("Cập nhật loại quyết định thành công");
            form.reset();
        },
        onError: (error: any) => {
            toast.error(error.message);
        },
    });

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            decision_type: "",
            // decision_group: "",
            description: "",
        },
    });

    useEffect(() => {
        if (data) {
            form.reset({
                decision_type: data.decision_type,
                // decision_group: data.decision_group,
                description: data.description,
            });
        }
    }, [data]);

    async function onSubmit(values: z.infer<typeof formSchema>) {
        try {
            if (type === "edit") {
                await updateDecisionType({
                    id: data?.decision_type_id || "",
                    payload: values,
                });
            } else {
                await createDecisionType(values);
            }
        } catch (error) {
            console.error("Error submitting form:", error);
        }
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                {editBtn ? (
                    editBtn
                ) : (
                    <Button className="w-full sm:w-auto">
                        <Plus className="mr-2 h-4 w-4"/>
                        {buttonText}
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className="w-[95vw] max-w-2xl sm:max-w-3xl max-h-[90vh]">
                <DialogHeader>
                    <div className={`flex justify-between md:flex-row flex-col p-5 items-center gap-2`}>
                        <DialogTitle className="text-xl sm:text-2xl">
                            {type === "edit" ? "Chỉnh sửa loại quyết định" : "Thêm loại quyết định"}
                        </DialogTitle>
                        {
                            data?.decision_type_id ? (
                                    <div className="border p-3 w-fit rounded-lg border-black">
                                        {data?.decision_type_id}
                                    </div>
                                ) :
                                null
                        }
                    </div>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <div className="flex flex-col lg:flex-row gap-4">
                            <div className="space-y-4 flex-1">
                                <FormField
                                    control={form.control}
                                    name="decision_type"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Tên loại quyết định <Asterisk /></FormLabel>
                                            <FormControl>
                                                <Input
                                                    className="border-primary border rounded-lg h-12"
                                                    placeholder="Nhập tên loại quyết định"
                                                    {...field}
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                                {/*<FormField*/}
                                {/*    control={form.control}*/}
                                {/*    name="decision_group"*/}
                                {/*    render={({field}) => (*/}
                                {/*        <FormItem>*/}
                                {/*            <FormLabel>Nhóm quyết định <Asterisk /></FormLabel>*/}
                                {/*            <Select onValueChange={field.onChange} value={field.value}>*/}
                                {/*                <FormControl>*/}
                                {/*                    <SelectTrigger className="h-10 w-full">*/}
                                {/*                        <SelectValue placeholder="Chọn nhóm quyết định" />*/}
                                {/*                    </SelectTrigger>*/}
                                {/*                </FormControl>*/}
                                {/*                <SelectContent>*/}

                                {/*                    {decistionGroup?.map((item, index) => (*/}
                                {/*                        <SelectItem key={index} value={item}>*/}
                                {/*                            <div className="flex items-center gap-2">*/}
                                {/*                                {item}*/}
                                {/*                            </div>*/}
                                {/*                        </SelectItem>*/}
                                {/*                    ))}*/}

                                {/*                </SelectContent>*/}
                                {/*            </Select>*/}
                                {/*            <FormMessage />*/}
                                {/*        </FormItem>*/}
                                {/*    )}*/}
                                {/*/>*/}
                            </div>

                            <div className="flex-1">
                                <FormField
                                    control={form.control}
                                    name="description"
                                    render={({field}) => (
                                        <FormItem className="h-full flex flex-col">
                                            <FormLabel>Mô tả</FormLabel>
                                            <FormControl>
                                                <textarea
                                                    className="flex-1 w-full min-h-[150px] sm:min-h-[200px] border-primary rounded-lg p-3 border focus:ring-2 focus:ring-primary focus:border-transparent"
                                                    placeholder="Nhập mô tả (nếu có)"
                                                    {...field}
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            </div>
                        </div>

                        <div className="flex flex-col sm:flex-row gap-3 pt-4 max-w-[600px] mx-auto">
                            <Button
                                type="button"
                                variant="outline"
                                className="md:w-[300px] w-full"
                                onClick={() => setOpen(false)}
                            >
                                Hủy
                            </Button>
                            <Button
                                type="submit"
                                className="md:w-[300px] w-full"
                                disabled={pendingCreate || pendingUpdate}
                            >
                                {pendingCreate || pendingUpdate ? (
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                ) : null}
                                {type === "edit" ? "Cập nhật" : "Thêm"}
                            </Button>
                        </div>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
};

export default CreateDecisionType;