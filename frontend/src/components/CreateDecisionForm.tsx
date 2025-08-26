'use client';

import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { createDecisionApi } from '@/apis/decision.api';
import { useQuery } from "@tanstack/react-query";
import {useGetAllEmployee} from "@/query/employee.query.ts";
import {getAllDecisionTypeApi} from "@/apis/decision-type.api.ts";
import * as React from "react";

export const formSchema = z.object({
    decision_name: z.string().nonempty('Vui lòng nhập tên quyết định'),
    effective_date: z.string().nonempty('Vui lòng chọn ngày hiệu lực'),
    sign_date: z.string().nonempty('Vui lòng chọn ngày ký'),
    condition: z.string().nonempty('Vui lòng chọn tình trạng'),
    description: z.string().optional(),
    attached_file: z.any().optional(),
    decision_type_id: z.string().nonempty('Vui lòng chọn loại quyết định'),
    employee_ids: z.array(z.string()).nonempty('Vui lòng chọn ít nhất 1 nhân viên'),
});

type Props = {
    open: boolean;
    setOpen: (open: boolean) => void;
};

type Employee = {
    employee_id: string;
    full_name: string;
};

interface MultiSelectEmployeeProps {
    employees: Employee[];
    value: string[];
    onChange: (value: string[]) => void;
    placeholder?: string;
    disabled?: boolean;
}

function MultiSelectEmployeeComponent({
                                          employees,
                                          value,
                                          onChange,
                                          placeholder = "Chọn nhân viên...",
                                          disabled,
                                      }: MultiSelectEmployeeProps) {
    const [open, setOpen] = React.useState(false);

    const selectedEmployees = employees.filter((e) =>
        value.includes(e.employee_id)
    );

    const toggleEmployee = (id: string) => {
        if (value.includes(id)) {
            onChange(value.filter((v) => v !== id));
        } else {
            onChange([...value, id]);
        }
    };

    return (
        <div className="w-full relative">
            <div
                className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21] cursor-pointer flex items-center justify-between"
                onClick={() => !disabled && setOpen(!open)}
            >
                <div className="flex flex-wrap gap-1">
                    {selectedEmployees.length > 0 ? (
                        selectedEmployees.map((employee) => (
                            <span
                                key={employee.employee_id}
                                className="bg-blue-100 text-blue-800 text-sm px-2 py-1 rounded"
                            >
                {employee.full_name}
              </span>
                        ))
                    ) : (
                        <span className="text-gray-500">{placeholder}</span>
                    )}
                </div>
                <span className="text-gray-400">▼</span>
            </div>

            {open && (
                <div className="absolute top-full left-0 right-0 bg-white border border-gray-300 rounded-md shadow-lg z-10 mt-1 max-h-60 overflow-y-auto">
                    {employees.map((employee) => (
                        <div
                            key={employee.employee_id}
                            className={`p-2 cursor-pointer hover:bg-gray-100 ${
                                value.includes(employee.employee_id) ? 'bg-blue-50' : ''
                            }`}
                            onClick={() => toggleEmployee(employee.employee_id)}
                        >
                            <label className="flex items-center cursor-pointer">
                                <input
                                    type="checkbox"
                                    checked={value.includes(employee.employee_id)}
                                    onChange={() => toggleEmployee(employee.employee_id)}
                                    className="mr-2"
                                />
                                {employee.full_name}
                            </label>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

const CreateDecisionForm = ({ open, setOpen }: Props) => {
    const { data: decisions, isLoading } = useQuery({
        queryKey: ["decisionTypes"],
        queryFn: getAllDecisionTypeApi,
    });

    const { data: employees, isLoading: isLoadingEmployee } = useGetAllEmployee();

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            decision_name: '',
            effective_date: '',
            sign_date: '',
            condition: "Đang hiệu lực",
            description: '',
            attached_file: undefined,
            employee_ids: [],
            decision_type_id: '',
        },
    });

    const onSubmit = async (values: z.infer<typeof formSchema>) => {
        try {
            const formData = new FormData();
            formData.append("decision_name", values.decision_name);
            formData.append("effective_date", values.effective_date);
            formData.append("sign_date", values.sign_date);
            formData.append("content", values.description || '');
            formData.append("condition", values.condition);
            formData.append("decision_type_id", values.decision_type_id);

            if (values.attached_file?.[0]) {
                formData.append("attached_file", values.attached_file[0]);
            }

            values.employee_ids.forEach((id) => {
                formData.append("employee_ids", id);
            });

            await createDecisionApi(formData);
            form.reset();
            setOpen(false);
        } catch (error) {
            console.error("Failed to create decision:", error);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent className="w-[1261px] h-[700px] p-6 overflow-y-auto">
                <DialogHeader>
                    <div className="flex items-center relative">
                        <DialogTitle className="text-xl font-semibold text-left">
                            THÊM QUYẾT ĐỊNH
                        </DialogTitle>
                    </div>
                </DialogHeader>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 flex flex-wrap">
                        {/* Left Column */}
                        <div className="w-[50%] pr-4 space-y-6">
                            <FormField
                                control={form.control}
                                name="decision_name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Tên quyết định</FormLabel>
                                        <FormControl>
                                            <Input {...field} type="text" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="sign_date"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Ngày ký quyết định</FormLabel>
                                        <FormControl>
                                            <Input {...field} type="datetime-local" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="effective_date"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Ngày hiệu lực</FormLabel>
                                        <FormControl>
                                            <Input {...field} type="datetime-local" className="w-full h-[51px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]" />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="attached_file"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">File đính kèm</FormLabel>
                                        <FormControl>
                                            <div className="relative">
                                                <input
                                                    id="file-upload"
                                                    type="file"
                                                    onChange={(e) => field.onChange(e.target.files)}
                                                    className="w-full h-[140px] p-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21] opacity-0 cursor-pointer"
                                                />
                                                <label
                                                    htmlFor="file-upload"
                                                    className="absolute inset-0 p-2 bg-[#EFEFEF] rounded-md flex items-center justify-center text-center font-medium text-black cursor-pointer border-2 border-dashed border-gray-300"
                                                >
                                                    ĐÍNH KÈM FILE
                                                </label>
                                            </div>
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />
                        </div>

                        {/* Right Column */}
                        <div className="w-[50%] pl-4 space-y-6">
                            <FormField
                                control={form.control}
                                name="condition"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Tình trạng</FormLabel>
                                        <FormControl>
                                            <select {...field} className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]">
                                                <option value="" disabled>Chọn tình trạng</option>
                                                {["Đang hiệu lực", "Hết hiệu lực", "Chưa hiệu lực"].map((item, index) => (
                                                    <option key={index} value={item}>
                                                        {item}
                                                    </option>
                                                ))}
                                            </select>
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="decision_type_id"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Loại quyết định</FormLabel>
                                        <FormControl>
                                            <select {...field} className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]">
                                                <option value="" disabled>Chọn loại quyết định</option>
                                                {isLoading ? (
                                                    <option value="">Đang tải...</option>
                                                ) : (
                                                    decisions?.map((item) => (
                                                        <option key={item.decision_type_id} value={item.decision_type_id}>
                                                            {item.decision_type}
                                                        </option>
                                                    ))
                                                )}
                                            </select>
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="employee_ids"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Nhân viên</FormLabel>
                                        <FormControl>
                                            <MultiSelectEmployeeComponent
                                                employees={employees || []}
                                                value={field.value}
                                                onChange={field.onChange}
                                                disabled={isLoadingEmployee}
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="description"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Mô tả</FormLabel>
                                        <FormControl>
                                            <textarea {...field} className="w-full h-[140px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21] resize-none" />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm" />
                                    </FormItem>
                                )}
                            />
                        </div>

                        {/* Buttons */}
                        <div className="w-full flex justify-end gap-4 pt-6">
                            <Button
                                type="button"
                                onClick={() => setOpen(false)}
                                className="bg-gray-400 hover:bg-gray-500 w-[120px] h-[40px]"
                            >
                                Hủy bỏ
                            </Button>
                            <Button
                                type="submit"
                                className="bg-[#DB3B21] hover:bg-[#b83a1a] w-[120px] h-[40px]"
                            >
                                Lưu thông tin
                            </Button>
                        </div>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
};

export default CreateDecisionForm;