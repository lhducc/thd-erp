import {useForm} from "react-hook-form";
import {Dialog, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog";
import {Button} from "@/components/ui/button";
import {useQuery} from "@tanstack/react-query";
import {createContractApi, getContractById, updateContractApi} from "@/apis/contract.api.ts";
import type {ContractFormValues, ContractType} from "@/types/contract.ts";
import {useEffect} from "react";
import {getAllContractsTypeApi} from "@/apis/contract-type.api.ts";
import type {Employee} from "@/types/employee.ts";
import {getAllEmployeesApi} from "@/apis/profile.api.ts";

type Props = {
    open: boolean;
    setOpen: (open: boolean) => void;
    contractId: string | null; // if is null ContractFrom will call create
};

export function ContractForm({open, setOpen, contractId}: Props) {
    const {
        register,
        handleSubmit,
        reset,
        watch,
        setValue,
        formState: {errors},
    } = useForm<ContractFormValues>({
        defaultValues: {
            status: "Chưa hiệu lực",
            approve_status: "Chờ duyệt",
            contractType: "",
        },
    });

    const {
        data: contract,
        isPending: pendingContract,
        isError,
    } = useQuery({
        queryKey: ["contract", contractId],
        queryFn: () => contractId ? getContractById(contractId) : null,
        enabled: open && !!contractId,
    });

    // Fetch contract types
    const {
        data: contractTypes = [],
        isPending: pendingContractTypes,
    } = useQuery<ContractType[]>({
        queryKey: ["contractTypes"],
        queryFn: getAllContractsTypeApi,
        enabled: open,
    });

    // Fetch employees
    const {
        data: employees = [],
        isPending: pendingEmployees,
    } = useQuery<Employee[]>({
        queryKey: ["employees"],
        queryFn: () => getAllEmployeesApi(1, 9999),
        enabled: open,
    });

    // Watch employeeId to update department automatically
    const employeeId = watch("employeeId");
    useEffect(() => {
        if (employeeId) {
            const selectedEmployee = employees.find(e => e.employee_id === employeeId);
            if (selectedEmployee) {
                setValue("employeeName", selectedEmployee.full_name);
                setValue("department", selectedEmployee.department?.department_name || "");
            }
        }
    }, [employeeId, employees, setValue]);

    // Reset form with contract data when it's loaded
    useEffect(() => {
        if (contract) {
            console.log(contract);
            reset({
                employeeId: contract.employee.employee_id,
                contractType: contract.contract_type,
                status: contract.condition === "Hiệu lực" ? "Hiệu lực" : "Chưa hiệu lực",
                signingDate: contract.sign_date?.split('T')[0],
                effectiveDate: contract.effective_date?.split('T')[0],
                expirationDate: contract.expired_date?.split('T')[0],
                note: contract.note,
                employeeName: contract.employee?.full_name,
                approve_status: contract.approve_status,
                department: contract.employee?.department?.department_name,
            });
        } else if (contractId === null) {
            reset({
                status: "Chưa hiệu lực",
                approve_status: "Chờ duyệt",
                contractType: "",
            });
        }
    }, [contract, contractId, reset]);

    if (isError) {
        return <div>Error loading contract data</div>;
    }

    const onSubmit = async (data: ContractFormValues) => {
        try {
            if (contractId === null) {
                await createContractApi(data);
                console.log("Đã tạo hợp đồng mới");
            } else {
                await updateContractApi(contractId, data);
                console.log("Đã cập nhật hợp đồng");
            }
            setOpen(false);
            reset();
        } catch (error: any) {
            console.error("Lỗi khi submit form:", error.message);
            alert(`Lỗi: ${error.message}`);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent className="sm:max-w-[90vw] md:max-w-[80vw] lg:max-w-[60vw] xl:max-w-[50vw] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 p-4 md:p-6">
                        <DialogTitle className={`text-lg sm:text-xl md:text-2xl uppercase`}>
                            {contractId === null ? 'Thêm hợp đồng' : 'Thông tin hợp đồng'}
                        </DialogTitle>
                        {contractId && (
                            <div className="border border-gray-500 p-2 md:p-3 rounded-xl md:rounded-2xl w-full md:w-[300px] text-sm md:text-base">
                                {contractId}
                            </div>
                        )}
                    </div>
                </DialogHeader>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="lg:grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Loại hợp đồng</label>
                            {pendingContractTypes ? (
                                <div className="w-full rounded-md border p-2 text-sm h-[43px] animate-pulse bg-gray-100" />
                            ) : (
                                <select
                                    {...register("contractType", {required: "Vui lòng chọn loại hợp đồng"})}
                                    className="w-full rounded-md border p-2 text-sm"
                                >
                                    <option value="">Chọn loại hợp đồng</option>
                                    {contractTypes.map((type) => {
                                        return (
                                        <option selected={contract?.contract_type === type.contract_type_id} key={type.contract_type_id} value={type.contract_type}>
                                            {type.contract_type}
                                        </option>
                                    )})}
                                </select>
                            )}
                            {errors.contractType && (
                                <p className="text-sm text-red-500">{errors.contractType.message}</p>
                            )}
                        </div>

                        {/* Rest of your form remains the same */}
                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Tình trạng</label>
                            <select
                                {...register("status")}
                                className="w-full rounded-md border p-2 text-sm"
                            >
                                <option value="Chưa hiệu lực">Chưa hiệu lực</option>
                                <option value="Hiệu lực">Hiệu lực</option>
                            </select>
                        </div>

                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Ngày ký</label>
                            <input
                                type="date"
                                {...register("signingDate", {required: "Vui lòng chọn ngày ký"})}
                                className="w-full rounded-md border p-2 text-sm"
                            />
                            {errors.signingDate && (
                                <p className="text-sm text-red-500">{errors.signingDate.message}</p>
                            )}
                        </div>

                        <div className="space-y-2 row-span-3">
                            <label className="block text-sm font-medium">Chú thích</label>
                            <textarea
                                {...register("note")}
                                className="w-full rounded-md border p-2 text-sm h-[220px]"
                                rows={3}
                            />
                        </div>

                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Ngày hợp đồng có hiệu lực</label>
                            <input
                                type="date"
                                {...register("effectiveDate", {required: "Vui lòng chọn ngày hiệu lực"})}
                                className="w-full rounded-md border p-2 text-sm"
                            />
                            {errors.effectiveDate && (
                                <p className="text-sm text-red-500">{errors.effectiveDate.message}</p>
                            )}
                        </div>

                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Ngày hết hạn hợp đồng</label>
                            <input
                                type="date"
                                {...register("expirationDate", {required: "Vui lòng chọn ngày hết hạn"})}
                                className="w-full rounded-md border p-2 text-sm"
                            />
                            {errors.expirationDate && (
                                <p className="text-sm text-red-500">{errors.expirationDate.message}</p>
                            )}
                        </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Tên nhân viên</label>
                            {pendingEmployees ? (
                                <div className="w-full rounded-md border p-2 text-sm h-[43px] animate-pulse bg-gray-100" />
                            ) : (
                                <select
                                    {...register("employeeId", {required: "Vui lòng chọn nhân viên"})}
                                    className="w-full rounded-md border p-2 text-sm h-[43px]"
                                >
                                    <option value="">Chọn nhân viên</option>
                                    {employees.map((employee) => {
                                        return (
                                        <option selected={employee.employee_id === contract?.employee.employee_id} key={employee.employee_id} value={employee.employee_id}>
                                            {employee.full_name} - {employee.employee_id}
                                        </option>
                                    )})}
                                </select>
                            )}
                            {errors.employeeId && (
                                <p className="text-sm text-red-500">{errors.employeeId.message}</p>
                            )}
                        </div>

                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Phụ cấp</label>
                            <select
                                {...register("allowanceStatus")}
                                className="w-full rounded-md border p-2 text-sm"
                            >
                                <option value="Chờ duyệt">Chờ duyệt</option>
                                <option value="Đã duyệt">Đã duyệt</option>
                                <option value="Từ chối">Từ chối</option>
                            </select>
                        </div>

                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Phòng ban</label>
                            <input
                                {...register("department", {required: "Vui lòng nhập phòng ban"})}
                                className="w-full rounded-md border p-2 text-sm"
                            />
                            {errors.department && (
                                <p className="text-sm text-red-500">{errors.department.message}</p>
                            )}
                        </div>
                        <div className="space-y-2">
                            <label className="block text-sm font-medium">Trạng thái</label>
                            <select
                                {...register("approve_status")}
                                className="w-full rounded-md border p-2 text-sm"
                            >
                                <option value="Chờ duyệt">Chờ duyệt</option>
                                <option value="Đã duyệt">Đã duyệt</option>
                                <option value="Từ chối">Từ chối</option>
                            </select>
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
                        <Button type="submit" variant="default">
                            Lưu thông tin
                        </Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
}