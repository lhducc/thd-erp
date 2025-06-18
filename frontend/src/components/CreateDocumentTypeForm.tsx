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
            // Tìm group có name khớp với data.document_group
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
            <DialogTrigger>
                {editBtn ? (
                    editBtn
                ) : (
                    <Button>
                        <Plus/>
                        Thêm loại tài liệu
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className="md:max-w-5xl">
                <DialogHeader>
                    <div className={`flex justify-between items-center p-5`}>
                        {type === "edit" ? (
                            <DialogTitle className={`text-[21px] uppercase`}>Chỉnh sửa loại tài liệu</DialogTitle>
                        ) : (
                            <DialogTitle className={`text-[21px] uppercase`}>Thêm loại tài liệu nhân sự</DialogTitle>
                        )}
                        {
                            data?.id !== null ?
                                <div className={`border border-gray-500 p-2 rounded-lg`}>
                                    {data?.id}
                                </div>
                                : null
                        }
                    </div>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                        <div className="flex gap-4">
                            <div className="flex flex-col gap-5">
                                <FormField
                                    control={form.control}
                                    name="document_type_name"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Tên loại tài liệu</FormLabel>
                                            <FormControl>
                                                <Input className={`md:w-sm`}
                                                       placeholder="Nhập tên loại tài liệu" {...field} />
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
                                                    <SelectTrigger className="md:w-sm">
                                                        <SelectValue placeholder="Chọn nhóm tài liệu"/>
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent className="bg-white">
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

                            <FormField
                                control={form.control}
                                name="description"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel>Mô tả</FormLabel>
                                        <FormControl>
                                        <textarea
                                            placeholder="Nhập mô tả (nếu có)"
                                            className="min-h-[100px] md:w-xl border rounded-2xl w-[300px] h-[400px] p-3"
                                            {...field}
                                        />
                                        </FormControl>
                                        <FormMessage/>
                                    </FormItem>
                                )}
                            />
                        </div>
                        <div className="flex gap-5 justify-center">

                            <button
                                className={`bg-[#EFEFEF] text-black hover:bg-[#EFEFEF] md:w-[200px] rounded-lg py-2`}>
                                Hủy
                            </button>
                            <button className={`md:w-[200px] rounded-lg py-2 bg-[#DB3B21] text-white `} type="submit"
                                    disabled={pendingCreate || pendingUpdate || pendingGroups}>
                                {pendingCreate || pendingUpdate ? (
                                    <div className="flex justify-center">
                                        <Loader2 className="animate-spin"/>
                                    </div>
                                ) : type === "edit" ? (
                                    "Cập nhật"
                                ) : (
                                    "Lưu thông tin"
                                )}
                            </button>
                        </div>

                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
};

export default CreateDocumentTypeForm;