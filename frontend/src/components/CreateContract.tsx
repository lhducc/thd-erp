import {useForm} from "react-hook-form";
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import {Button} from "@/components/ui/button";
import {useMutation, useQuery} from "@tanstack/react-query";
import {createContractApi, reapproveContractApi, updateContractApi} from "@/apis/contract.api.ts";
import type {Contract, ContractFormValues} from "@/types/contract.ts";
import {useEffect, useState} from "react";
import {getAllContractsTypeApi} from "@/apis/contract-type.api.ts";
import {getAllEmployeesApi} from "@/apis/profile.api.ts";
import {getAllAllowancesApi} from "@/apis/allowance.api.ts";
import Loading from "@/components/Loading.tsx";
import {toast} from "sonner";
import {z} from "zod";
import {zodResolver} from "@hookform/resolvers/zod";
import {Form, FormControl, FormField, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {Loader2, Plus} from "lucide-react";

const contractFormSchema = z.object({
    contract_type: z.string().min(1, "Vui lòng chọn loại hợp đồng"),
    // con: z.string(),
    sign_date: z.string().min(1, "Vui lòng chọn ngày ký"),
    effective_date: z.string().min(1, "Vui lòng chọn ngày hiệu lực"),
    expired_date: z.string().min(1, "Vui lòng chọn ngày hết hạn"),
    note: z.string().optional(),
    employee_id: z.string().min(1, "Vui lòng chọn nhân viên"),
    department: z.string().min(1, "Phòng ban không được để trống"),
    allowance_ids: z.array(z.string()).optional(),
    approve_status: z.string(),
});

type Props = {
    editBtn?: React.ReactNode;
    data?: Contract;
    type?: "pending" | "approved" | "rejected" ;
    refetch?: () => void;
};

export function ContractForm({editBtn, data, type, refetch}: Props) {
    const [open, setOpen] = useState(false);
    const form = useForm<ContractFormValues>({
        resolver: zodResolver(contractFormSchema),
        defaultValues: {
            contract_id: "",
            condition: "Chưa hiệu lực",
            approve_status: "Chờ duyệt",
            contract_type: "",
            employee_id: "",
            department: "",
            sign_date: "",
            effective_date: "",
            expired_date: "",
            note: "",
            allowance_ids: [],
        },
    });

    // Fetch contract types
    const {data: contractTypes = [], isPending: pendingContractTypes} = useQuery({
        queryKey: ["contractTypes"],
        queryFn: getAllContractsTypeApi,
        enabled: open,
    });

    // Fetch allowances
    const {data: allowances = [], isPending: pendingAllowances} = useQuery({
        queryKey: ["allowances"],
        queryFn: getAllAllowancesApi,
        enabled: open,
    });

    // Fetch employees
    const {data: employees = [], isPending: pendingEmployees} = useQuery({
        queryKey: ["employees"],
        queryFn: () => getAllEmployeesApi(1, 9999),
        enabled: open,
    });

    const {mutateAsync: createContract, isPending: pendingCreate} = useMutation({
        mutationFn: createContractApi,
        onSuccess: () => {
            toast.success("Tạo hợp đồng thành công");
            setOpen(false);
            refetch?.();
        },
        onError: (error) => {
            console.log(error);
            toast.error(error.message);
        },
    });

    const {mutateAsync: reapproveContract, isPending: pendingReapprove} = useMutation({
        mutationFn: reapproveContractApi,
        onSuccess: () => {
            toast.success("Yêu cầu duyệt lại thành công");
            setOpen(false);
            refetch?.();
        },
        onError: (error) => {
            console.log(error);
            toast.error(error.message);
        },
    });

    const {mutateAsync: updateContract, isPending: pendingUpdate} = useMutation({
        mutationFn: (payload: { id: string; data: Partial<ContractFormValues> }) =>
            updateContractApi(payload.id, payload.data),
        onSuccess: () => {
            toast.success("Cập nhật hợp đồng thành công");
            setOpen(false);
            refetch?.();
        },
        onError: (error: any) => {
            if (error.status === 500) {
                toast.error(error.data.error);
            } else {
                toast.error(error.message || "Có lỗi xảy ra khi cập nhật hợp đồng");
            }

            console.log("Error details:", {
                status: error.status,
                data: error.data
            });
        },
    });

    // Watch employeeId to update department automatically
    const employeeId = form.watch("employee_id");
    useEffect(() => {
        if (employeeId) {
            const selectedEmployee = employees?.find(e => e.employee_id === employeeId);
            if (selectedEmployee) {
                form.setValue("department", selectedEmployee.department?.department_name || "");
            }
        }
    }, [employeeId, employees, form]);

    // Reset form when opening/closing or when data changes
    useEffect(() => {
        if (!open) {
            form.reset();
            return;
        }

        if (data) {
            form.reset({
                contract_id: data.contract_id || "",
                condition: data.condition || "Chưa hiệu lực",
                approve_status: data.approve_status || "Chờ duyệt",
                employee_id: data.employee?.employee_id || "",
                contract_type: data.contract_type || "",
                sign_date: formatDateForInput(data.sign_date) || "",
                effective_date: formatDateForInput(data.effective_date) || "",
                expired_date: formatDateForInput(data.expired_date) || "",
                note: data.note || "",
                department: data.employee?.department?.department_name || "",
                allowance_ids: data.allowances?.map(a => a.id) || [],
            });
        }
    }, [open, data, form]);

    const formatDateForInput = (dateString: string) => {
        if (!dateString) return "";
        const date = new Date(dateString);
        const offset = date.getTimezoneOffset() * 60000;
        const localDate = new Date(date.getTime() - offset);
        return localDate.toISOString().split('T')[0];
    };

    const onSubmit = async (values: ContractFormValues) => {
        try {
            const payload: ContractFormValues = {
                ...values,
                sign_date: `${values.sign_date}T00:00:00+07:00`,
                effective_date: `${values.effective_date}T00:00:00+07:00`,
                expired_date: `${values.expired_date}T00:00:00+07:00`
            };

            if (type === "pending" && data?.contract_id) {
                await updateContract({id: data.contract_id, data: payload});
            } else if (type === "rejected" && data?.contract_id) {
                await reapproveContract(data.contract_id);
            }
            else {
                await createContract(payload);
            }
        } catch (error) {
            console.error("Error submitting form:", error);
        }
    };

    const handleAddAllowance = (allowanceId: string) => {
        const currentAllowances = form.watch("allowance_ids") || [];
        if (!currentAllowances.includes(allowanceId)) {
            form.setValue("allowance_ids", [...currentAllowances, allowanceId]);
        }
    };

    const handleRemoveAllowance = (allowanceId: string) => {
        const currentAllowances = form.watch("allowance_ids") || [];
        form.setValue("allowance_ids", currentAllowances.filter(id => id !== allowanceId));
    };

    // const isLoading = pendingContract;

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                {editBtn ? (
                    editBtn
                ) : (
                    <Button>
                        <Plus className="mr-2 h-4 w-4"/>
                        Thêm hợp đồng
                    </Button>
                )}
            </DialogTrigger>
            <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <div className={`flex justify-between md:flex-row flex-col p-5 items-center gap-2`}>
                    <DialogTitle className="text-xl">
                        {type === "pending" ? 'Chỉnh sửa hợp đồng' : 'Thêm hợp đồng mới'}
                    </DialogTitle>
                        {
                            data?.contract_id ? (
                                    <div className="border p-3 w-fit rounded-lg border-black">
                                        {data?.contract_id}
                                    </div>
                                ) :
                                null
                        }
                    </div>
                </DialogHeader>

                {pendingAllowances || pendingContractTypes || pendingEmployees ? (
                    <Loading/>
                ) : (
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Contract Type */}
                            <FormField
                                control={form.control}
                                name="contract_type"
                                render={({field}) => {
                                    return (
                                        <FormItem>
                                            <FormLabel>Loại hợp đồng</FormLabel>
                                            <Select onValueChange={field.onChange} value={field.value}>
                                                <FormControl  className={`w-full col-span-2 ${type !== "approved" ? "border-blue-300" : "border-gray-300"}`}>
                                                    <SelectTrigger>
                                                        <SelectValue placeholder="Chọn loại hợp đồng"/>
                                                    </SelectTrigger>
                                                </FormControl>
                                                <SelectContent>
                                                    {contractTypes.map((type) => (
                                                        <SelectItem key={type.contract_type}
                                                                    value={type.contract_type_id}>
                                                            {type.contract_type}
                                                        </SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                            <FormMessage/>
                                        </FormItem>
                                    )
                                }}
                            />

                            <div className="md:grid-cols-2 md:col-span-1 col-span-2 gap-4">
                                Tình trạng
                                <div className={`border border-gray-300 text-gray-500 p-2 rounded-lg`}>
                                    {
                                        data?.condition ?
                                            <p>{data.condition}</p>
                                            :
                                            <p>Chưa hiệu lực</p>
                                    }
                                </div>
                            </div>

                            {/* Sign Date */}
                            <FormField
                                control={form.control}
                                name="sign_date"
                                render={({field}) => (
                                    <FormItem className={`rounded-lg col-span-2 md:col-span-1`}>
                                        <FormLabel>Ngày ký</FormLabel>
                                        <FormControl>
                                            <Input className={`rounded-lg h-[50px] w-full col-span-2`} type="date" {...field} />
                                        </FormControl>
                                        <FormMessage/>
                                    </FormItem>
                                )}
                            />

                            {/* Note */}
                            <div className="md:row-span-3 col-span-2 md:col-span-1">
                                <FormField
                                    control={form.control}
                                    name="note"
                                    render={({field}) => (
                                        <FormItem>
                                            <FormLabel>Chú thích</FormLabel>
                                            <FormControl>
                          <textarea
                              {...field}
                              value={field.value || ""}
                              className={`w-full rounded-md border p-2 min-h-[250px] ${type !== "approved" ? "border-blue-300" : ""}`}
                          />
                                            </FormControl>
                                            <FormMessage/>
                                        </FormItem>
                                    )}
                                />
                            </div>

                            {/* Effective Date */}
                            <FormField
                                control={form.control}
                                name="effective_date"
                                render={({field}) => (
                                    <FormItem className={`rounded-lg col-span-2 md:col-span-1`}>
                                        <FormLabel>Ngày hiệu lực</FormLabel>
                                        <FormControl>
                                            <Input className={`rounded-lg h-[50px]`} type="date" {...field} />
                                        </FormControl>
                                        <FormMessage/>
                                    </FormItem>
                                )}
                            />

                            {/* Expired Date */}
                            <FormField
                                control={form.control}
                                name="expired_date"
                                render={({field}) => (
                                    <FormItem className={`rounded-lg col-span-2 md:col-span-1`}>
                                        <FormLabel>Ngày hết hạn</FormLabel>
                                        <FormControl>
                                            <Input className={`rounded-lg h-[50px]`} type="date" {...field} />
                                        </FormControl>
                                        <FormMessage/>
                                    </FormItem>
                                )}
                            />

                            {/* Employee */}
                            <FormField
                                control={form.control}
                                name="employee_id"
                                render={({field}) => (
                                    <FormItem className={`rounded-lg col-span-2 md:col-span-1`}>
                                        <FormLabel>Tên nhân viên</FormLabel>
                                        <Select onValueChange={field.onChange} value={field.value}>
                                            <FormControl className={`w-full`}>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Chọn nhân viên"/>
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                {employees?.map((employee) => (
                                                    <SelectItem key={employee.employee_id}
                                                                value={employee.employee_id}>
                                                        {employee.full_name}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                        <FormMessage/>
                                    </FormItem>
                                )}
                            />

                            <div className={`md:col-span-1 col-span-2`}>
                                Trạng thái
                                <div className={`border border-gray-300 rounded-lg p-2 text-gray-500`}>
                                    {
                                        data?.contract_id ?
                                            <p>{data?.approve_status}</p>
                                            : <p>Chờ duyệt</p>
                                    }
                                </div>
                            </div>

                            {/* Department */}
                            <FormField
                                control={form.control}
                                name="department"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel>Phòng ban</FormLabel>
                                        <FormControl>
                                            <Input className={`cursor-default rounded-lg h-[50px]`} {...field} readOnly/>
                                        </FormControl>
                                        <FormMessage/>
                                    </FormItem>
                                )}
                            />

                            {/* Allowances */}
                            <div className="col-span-2 md:col-span-1 space-y-2">
                                <FormLabel>Phụ cấp</FormLabel>
                                <div className="flex gap-2">
                                    <Select
                                        onValueChange={(value) => {
                                            if (value) handleAddAllowance(value);
                                        }}
                                    >
                                        <SelectTrigger className="flex-1">
                                            <SelectValue placeholder="Chọn phụ cấp"/>
                                        </SelectTrigger>
                                        <SelectContent>
                                            {allowances
                                                .filter(a => !form.watch("allowance_ids")?.includes(a.id))
                                                .map((a) => (
                                                    <SelectItem key={a.id} value={a.id}>
                                                        {a.allowance_name}
                                                    </SelectItem>
                                                ))}
                                        </SelectContent>
                                    </Select>
                                </div>

                                <div className="mt-2 space-y-2">
                                    {form.watch("allowance_ids")?.map(allowanceId => {
                                        const allowance = allowances.find(a => a.id === allowanceId);
                                        return allowance ? (
                                            <div key={allowanceId}
                                                 className="flex items-center justify-between p-2 border rounded">
                                                <span>{allowance.allowance_name}</span>
                                                <Button
                                                    type="button"
                                                    variant="ghost"
                                                    size="sm"
                                                    className={`text-red-500`}
                                                    onClick={() => handleRemoveAllowance(allowanceId)}
                                                >
                                                    Xóa
                                                </Button>
                                            </div>
                                        ) : null;
                                    })}
                                </div>
                            </div>
                        </div>

                        <div className="flex justify-end gap-2 pt-4">
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => setOpen(false)}
                            >
                                Hủy bỏ
                            </Button>
                            <Button type="submit"
                                    className={`${type === "approved" ? "hidden" : ""}`}
                                    disabled={pendingCreate || pendingUpdate || pendingReapprove}>
                                {pendingCreate || pendingUpdate || pendingReapprove ? (
                                    <Loader2 className="animate-spin mr-2 h-4 w-4"/>
                                ) : null}
                                {type === "pending" ? "Cập nhật" : type === "rejected" ? "Yêu cầu duyệt lại" : "Tạo hợp đồng"}
                            </Button>
                        </div>
                    </form>
                </Form>
                )}
            </DialogContent>
        </Dialog>
    );
}