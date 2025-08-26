import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import {Button} from "./ui/button";
import {Loader2, Plus} from "lucide-react";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod";
import {useForm} from "react-hook-form";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {useMutation, useQuery} from "@tanstack/react-query";
import {
    createDepartmentApi,
    updateDepartmentApi,
} from "@/apis/department.api";
import {getAllOfficesApi} from "@/apis/office.api";
import Loading from "./Loading";
import {toast} from "sonner";
import type {Department, PayloadDepartment} from "@/types";
import {useEffect, useState} from "react";
import {useEmployeeByRoleNameQuery} from "@/query/employee.query.ts";

const formSchema = z.object({
    department_name: z.string().nonempty("Tên phòng ban không được để trống"),
    manager: z.string().optional(),
    office_id: z.string().nonempty("Văn phòng không được để trống"),
});

type Props = {
    editBtn?: React.ReactNode; // Prop để truyền vào nút chỉnh sửa
    department?: Department;
    type?: "edit";
    refetch?: Function;
};

const CreateDepartmentForm = ({
                                  editBtn,
                                  department,
                                  type,
                                  refetch,
                              }: Props) => {
    const [open, setOpen] = useState(false);

    const {data: offices, isPending: pendingOffices} = useQuery({
        queryKey: ["offices"],
        queryFn: getAllOfficesApi,
        gcTime: 0,
        staleTime: 0,
    });

    const {data: managers} = useEmployeeByRoleNameQuery("manager")

    // 1. Define your form.
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            department_name: "",
            manager: "",
            office_id: "",
        },
    });

    useEffect(() => {
        if (department) {
            form.setValue("department_name", department.department_name);
            form.setValue("manager", department.manager);
            form.setValue("office_id", department.office_id);
        }
    }, [department]);

    const {mutateAsync: createDepartment, isPending: pendingCreateDepartment} =
        useMutation({
            mutationFn: createDepartmentApi,
            onSuccess: () => {
                refetch && refetch(); // Gọi lại hàm refetch nếu có
                toast.success("Thêm phòng ban thành công");
                setOpen(false); // Đóng dialog sau khi tạo phòng ban thành công
                form.reset(); // Đặt lại giá trị của form
            },
            onError: (error) => {
                toast.error(error.message);
            },
        });

    const {mutateAsync: updateDepartment, isPending: pendingUpdateDepartment} =
        useMutation({
            mutationFn: ({
                             id,
                             payload,
                         }: {
                id: string;
                payload: PayloadDepartment;
            }) => updateDepartmentApi(id, payload),
            onSuccess: () => {
                refetch && refetch(); // Gọi lại hàm refetch nếu có
                toast.success("Cập nhật phòng ban thành công");
                setOpen(false); // Đóng dialog sau khi cập nhật phòng ban thành công
                form.reset(); // Đặt lại giá trị của form
            },
            onError: (error) => {
                toast.error(error.message);
            },
        });

    // 2. Define a submit handler.
    async function onSubmit(values: z.infer<typeof formSchema>) {
        // Do something with the form values.
        // ✅ This will be type-safe and validated.
        // console.log(values);
        if (type === "edit") {
            await updateDepartment({
                id: department?.department_id || "",
                payload: values,
            });
        } else {
            await createDepartment(values);
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
                        Thêm phòng ban
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className={`md:max-w-2xl`}>
                {pendingOffices ? (
                    <Loading/>
                ) : (
                    <>
                        <DialogHeader>
                            <DialogTitle>
                                {type === "edit" ? "Cập nhật phòng ban" : "Thêm phòng ban"}
                            </DialogTitle>
                        </DialogHeader>
                        <Form {...form}>
                            <form
                                onSubmit={form.handleSubmit(onSubmit)}
                                className="space-y-8"
                            >
                                <FormField
                                    control={form.control}
                                    name="department_name"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Tên phòng ban</FormLabel>
                                            <FormControl>
                                                <Input placeholder="" {...field} />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="manager"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Người quản lý</FormLabel>
                                            <Select
                                                onValueChange={field.onChange}
                                                defaultValue={field.value}
                                            >
                                                <FormControl>
                                                    <SelectTrigger className="w-full">
                                                        <SelectValue placeholder=""/>
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    {managers?.map((manager) => (
                                                        <SelectItem
                                                            value={manager.employee_id}>{manager.full_name}</SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="office_id"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Văn phòng</FormLabel>
                                            <Select
                                                onValueChange={field.onChange}
                                                defaultValue={field.value}
                                            >
                                                <FormControl>
                                                    <SelectTrigger className="w-full">
                                                        <SelectValue placeholder=""/>
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    {offices?.map((office) => (
                                                        <SelectItem
                                                            key={office.office_id}
                                                            value={office.office_id}
                                                        >
                                                            {office.office_name}
                                                        </SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                                <Button type="submit" disabled={pendingCreateDepartment || pendingUpdateDepartment}>
                                    {pendingCreateDepartment || pendingUpdateDepartment ? (
                                        <Loader2 className="animate-spin"/>
                                    ) : type === "edit" ? (
                                        "Cập nhật"
                                    ) : (
                                        "Thêm"
                                    )}
                                </Button>
                            </form>
                        </Form>
                    </>
                )}
            </DialogContent>
        </Dialog>
    );
};

export default CreateDepartmentForm;
