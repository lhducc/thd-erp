import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import {Loader2, Plus} from "lucide-react";
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
import {useMutation, useQuery} from "@tanstack/react-query";
import {toast} from "sonner";
import type {DocumentType, DocumentTypeCreate} from "@/types/document-type";
import {
    createDocumentTypeApi,
    updateDocumentTypeApi,
} from "@/apis/document-type.api";
import {getAllDocumentsGroupApi} from "@/apis/document-group.api.ts";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

const formSchema = z.object({
    document_group: z.string().nonempty("Nhóm tài liệu không được để trống"),
    document_type_name: z.string().nonempty("Tên loại tài liệu không được để trống"),
    description: z.string().optional(),
});

type Props = {
    editBtn?: React.ReactNode;
    data?: DocumentType;
    type?: "edit";
    refetch?: () => void;
};

const CreateDocumentTypeForm = ({editBtn, data, type, refetch}: Props) => {
    const [open, setOpen] = useState(false);

    const {
        data: groups,
        isPending: pendingGroups,
    } = useQuery({
        queryKey: ["documentGroup"],
        queryFn: getAllDocumentsGroupApi,
    });

    const {mutateAsync: createDocumentType, isPending: pendingCreate} = useMutation({
        mutationFn: (payload: DocumentTypeCreate) => createDocumentTypeApi(payload),
        onSuccess: () => {
            refetch?.();
            setOpen(false);
            toast.success("Tạo loại tài liệu thành công");
            form.reset();
        },
        onError: (error: any) => {
            toast.error(error.message);
        },
    });

    const {mutateAsync: updateDocumentType, isPending: pendingUpdate} = useMutation({
        mutationFn: ({id, payload}: { id: string; payload: DocumentTypeCreate }) =>
            updateDocumentTypeApi(id, payload),
        onSuccess: () => {
            refetch?.();
            setOpen(false);
            toast.success("Cập nhật loại tài liệu thành công");
            form.reset();
        },
        onError: (error: any) => {
            toast.error(error.message);
        },
    });

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            document_group: "",
            document_type_name: "",
            description: "",
        },
    });

    useEffect(() => {
        if (data && groups) {
            const matchedGroup = groups.find(g => g.name === data.document_group);
            if (matchedGroup) {
                form.setValue("document_group", matchedGroup.name);
            }
            form.setValue("document_type_name", data.document_type_name);
            form.setValue("description", data.description);
        }
    }, [data, groups]);

    async function onSubmit(values: z.infer<typeof formSchema>) {
        if (type === "edit") {
            await updateDocumentType({
                id: data?.id || "",
                payload: values,
            });
        } else {
            await createDocumentType({
                ...values,
                created_date: new Date().toISOString(),
            });
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
                        Thêm loại tài liệu
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className="sm:max-w-[90vw] md:max-w-2xl lg:max-w-4xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 p-4 sm:p-5">
                        {type === "edit" ? (
                            <DialogTitle className="text-lg sm:text-xl lg:text-2xl uppercase">
                                Chỉnh sửa loại tài liệu
                            </DialogTitle>
                        ) : (
                            <DialogTitle className="text-lg sm:text-xl lg:text-2xl uppercase">
                                Thêm loại tài liệu nhân sự
                            </DialogTitle>
                        )}
                        {data?.id && (
                            <div className="border border-gray-500 px-3 py-1 rounded-lg text-sm">
                                {data.id}
                            </div>
                        )}
                    </div>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                        <div className="flex flex-col lg:flex-row gap-6">
                            <div className="flex-1 space-y-4">
                                <FormField
                                    control={form.control}
                                    name="document_type_name"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Tên loại tài liệu</FormLabel>
                                            <FormControl>
                                                <Input
                                                    placeholder="Nhập tên loại tài liệu"
                                                    {...field}
                                                    className="w-full"
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="document_group"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Nhóm tài liệu</FormLabel>
                                            <Select onValueChange={field.onChange} value={field.value}>
                                                <FormControl>
                                                    <SelectTrigger className="w-full">
                                                        <SelectValue placeholder="Chọn nhóm tài liệu"/>
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent className="bg-white max-h-60 overflow-auto">
                                                    {groups?.map((group) => (
                                                        <SelectItem
                                                            key={group.id}
                                                            value={group.name}
                                                            className="hover:bg-gray-100 cursor-pointer"
                                                        >
                                                            {group.name}
                                                        </SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            </div>

                            <div className="flex-1">
                                <FormField
                                    control={form.control}
                                    name="description"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Mô tả</FormLabel>
                                            <FormControl>
                                                <textarea
                                                    placeholder="Nhập mô tả (nếu có)"
                                                    className="w-full min-h-[150px] p-3 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                                                    {...field}
                                                />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            </div>
                        </div>
                        <div className="flex flex-col sm:flex-row gap-3 justify-center pt-4">
                            <Button
                                type="button"
                                variant="outline"
                                className="w-full sm:w-[200px]"
                                onClick={() => setOpen(false)}
                            >
                                Hủy
                            </Button>
                            <Button
                                type="submit"
                                className="w-full sm:w-[200px] bg-[#DB3B21] hover:bg-[#DB3B21]/90"
                                disabled={pendingCreate || pendingUpdate || pendingGroups}
                            >
                                {pendingCreate || pendingUpdate ? (
                                    <Loader2 className="animate-spin mr-2 h-4 w-4"/>
                                ) : null}
                                {type === "edit" ? "Cập nhật" : "Lưu thông tin"}
                            </Button>
                        </div>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
};

export default CreateDocumentTypeForm;