'use client';

import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import {Button} from '@/components/ui/button';
import {z} from 'zod';
import {zodResolver} from '@hookform/resolvers/zod';
import {useForm} from 'react-hook-form';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import {Input} from '@/components/ui/input';
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from '@/components/ui/select';
import {updateDecisionApi} from '@/apis/decision.api';
import {getAllDecisionTypeApi} from '@/apis/decision-type.api.ts';
import {useEffect, useState} from 'react';
import {useQuery} from '@tanstack/react-query';
import {format} from 'date-fns';
import {useEmployeeByRoleNameQuery} from "@/query/employee.query.ts";
import * as React from "react";
import type {EmployeeNameAndRole} from "@/types/employee.ts";

// Update the schema to handle array of employee objects
export const formSchema = z.object({
    decision_id: z.string().min(1, 'Vui lòng nhập mã quyết định'),
    decision_name: z.string().min(1, 'Vui lòng nhập tên quyết định'),
    effective_date: z.string().min(1, 'Vui lòng chọn ngày hiệu lực'),
    sign_date: z.string().min(1, 'Vui lòng chọn ngày ký'),
    content: z.string().min(1, 'Vui lòng nhập nội dung'),
    condition: z.string().min(1, 'Vui lòng chọn tình trạng'),
    created_date: z.string().min(1, 'Vui lòng nhập ngày tạo'),
    employees: z.array(z.object({
        employee_id: z.string(),
        fullname: z.string()
    })).min(1, 'Vui lòng chọn ít nhất một nhân viên'),
    decision_type_id: z.string().min(1, 'Vui lòng chọn loại quyết định'),
    attached_file: z.any().optional(),
});

interface MultiSelectEmployeeProps {
    employees: EmployeeNameAndRole[];
    selectedEmployees: {
        employee_id: string;
        fullname: string;
    }[];
    value: string[];
    onChange: (value: string[]) => void;
    placeholder?: string;
    disabled?: boolean;
}

function MultiSelectEmployeeComponent({
                                          employees,
                                          selectedEmployees,
                                          value = [],
                                          onChange,
                                          placeholder = "Chọn nhân viên...",
                                          disabled,
                                      }: MultiSelectEmployeeProps) {
    const [open, setOpen] = React.useState(false);
    const dropdownRef = React.useRef<HTMLDivElement>(null);
    console.log(selectedEmployees)
    React.useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setOpen(false);
            }
        }

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const toggleEmployee = (id: string) => {
        if (value.includes(id)) {
            onChange(value.filter((v) => v !== id));
        } else {
            onChange([...value, id]);
        }
    };

    return (
        <div className="w-full relative" ref={dropdownRef}>
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

type FormSchemaType = z.infer<typeof formSchema>;

type DecisionType = {
    decision_type_id: string;
    decision_type: string;
};

type Props = {
    open: boolean;
    setOpen: (open: boolean) => void;
    Decision: FormSchemaType;
    onSuccess?: () => void;
};

const UpdateDecisionForm = ({open, setOpen, Decision, onSuccess}: Props) => {
    const [filePreview, setFilePreview] = useState<string | null>(null);

    const form = useForm<FormSchemaType>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            decision_id: '',
            decision_name: '',
            effective_date: '',
            sign_date: '',
            content: '',
            condition: '',
            created_date: '',
            employees: [],
            decision_type_id: '',
            attached_file: undefined,
        },
    });

    const {data: decisionTypes, isLoading} = useQuery<DecisionType[]>({
        queryKey: ["decisionTypes"],
        queryFn: getAllDecisionTypeApi,
    });
    const {data: employees, isLoading: isLoadingEmployee} = useEmployeeByRoleNameQuery("manager");

    const onSubmit = async (values: FormSchemaType) => {
        try {
            const formData = new FormData();

            // Append all fields to formData
            Object.keys(values).forEach(key => {
                if (key === 'attached_file' && values[key]?.[0]) {
                    formData.append(key, values[key][0]);
                } else if (key === 'employees') {
                    // Convert array of employee objects to array of employee IDs
                    const employeeIds = values.employees.map(emp => emp.employee_id);
                    formData.append('employees', JSON.stringify(employeeIds));
                } else if (key !== 'attached_file') {
                    formData.append(key, values[key as keyof FormSchemaType] as string);
                }
            });

            await updateDecisionApi(Decision.decision_id, formData);
            setOpen(false);
            onSuccess?.();
        } catch (error) {
            console.error('Failed to update decision:', error);
        }
    };

    useEffect(() => {
        if (Decision) {
            // Format dates for datetime-local input
            const formatDateForInput = (dateString: string) => {
                try {
                    const date = new Date(dateString);
                    return format(date, "yyyy-MM-dd'T'HH:mm");
                } catch {
                    return dateString;
                }
            };

            form.reset({
                ...Decision,
                effective_date: formatDateForInput(Decision.effective_date),
                sign_date: formatDateForInput(Decision.sign_date),
                created_date: formatDateForInput(Decision.created_date),
                attached_file: undefined,
            });

            // If there's an existing file, show its name
            if (Decision.attached_file) {
                setFilePreview(Decision.attached_file);
            }
        }
    }, [Decision, form]);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>, onChange: (value: any) => void) => {
        const file = e.target.files?.[0];
        if (file) {
            onChange(e.target.files);
            setFilePreview(file.name);
        }
    };

    const handleClose = () => {
        form.reset();
        setFilePreview(null);
        setOpen(false);
    };

    // Get selected employee IDs for the multi-select component
    const selectedEmployeeIds = form.watch('employees')?.map(emp => emp.employee_id) || [];

    // Handle employee selection change
    const handleEmployeeChange = (employeeIds: string[]) => {
        const selectedEmployees = employees?.filter(emp =>
            employeeIds.includes(emp.employee_id)
        ) || [];

        form.setValue('employees', selectedEmployees);
    };

    // Get currently selected employees for display
    const currentSelectedEmployees = form.watch('employees') || [];

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogContent className="w-[1261px] max-h-[90vh] overflow-y-auto p-6">
                <DialogHeader>
                    <div className="flex items-center justify-between">
                        <DialogTitle className="text-xl font-semibold">
                            CẬP NHẬT QUYẾT ĐỊNH
                        </DialogTitle>
                        <div
                            className="text-xl w-[150px] h-[54px] font-semibold flex items-center justify-center rounded-3xl border-2 border-gray-400">
                            {Decision.decision_id}
                        </div>
                    </div>
                </DialogHeader>

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 flex flex-wrap">
                        {/* Left Column */}
                        <div className="w-[50%] pr-4 space-y-6">
                            <FormField
                                control={form.control}
                                name="decision_name"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Tên quyết định</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                type="text"
                                                className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                                placeholder="Nhập tên quyết định"
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm"/>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="sign_date"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Ngày ký quyết định</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                type="datetime-local"
                                                className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm"/>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="effective_date"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Ngày quyết định có hiệu lực</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                type="datetime-local"
                                                className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]"
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm"/>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="attached_file"
                                render={({field: {onChange, value, ...field}}) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">File đính kèm</FormLabel>
                                        <FormControl>
                                            <div className="relative">
                                                <input
                                                    id="file-upload"
                                                    type="file"
                                                    {...field}
                                                    onChange={(e) => handleFileChange(e, onChange)}
                                                    className="w-full h-[140px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21] opacity-0 absolute z-10 cursor-pointer"
                                                />
                                                <label
                                                    htmlFor="file-upload"
                                                    className="w-full h-[140px] p-2 bg-[#EFEFEF] rounded-md flex flex-col items-center justify-center text-center font-medium text-black cursor-pointer border-2 border-dashed border-gray-300"
                                                >
                                                    {filePreview ? (
                                                        <span className="text-sm text-gray-700">{filePreview}</span>
                                                    ) : (
                                                        <>
                                                            <span className="block mb-2">ĐÍNH KÈM FILE</span>
                                                            <span
                                                                className="text-xs text-gray-500">Click để chọn file</span>
                                                        </>
                                                    )}
                                                </label>
                                            </div>
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm"/>
                                    </FormItem>
                                )}
                            />
                        </div>

                        {/* Right Column */}
                        <div className="w-[50%] pl-4 space-y-6">
                            <FormField
                                control={form.control}
                                name="condition"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Tình trạng</FormLabel>
                                        <FormControl>
                                            <Select
                                                value={field.value}
                                                onValueChange={field.onChange}
                                            >
                                                <SelectTrigger
                                                    className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]">
                                                    <SelectValue placeholder="Chọn tình trạng"/>
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="Đang hiệu lực">Đang hiệu lực</SelectItem>
                                                    <SelectItem value="Hết hiệu lực">Hết hiệu lực</SelectItem>
                                                    <SelectItem value="Chưa hiệu lực">Chưa hiệu lực</SelectItem>
                                                </SelectContent>
                                            </Select>
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm"/>
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="decision_type_id"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Loại quyết định</FormLabel>
                                        <FormControl>
                                            <Select
                                                value={field.value}
                                                onValueChange={field.onChange}
                                                disabled={isLoading}
                                            >
                                                <SelectTrigger
                                                    className="w-full h-[51px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21]">
                                                    <SelectValue placeholder="Chọn loại quyết định"/>
                                                </SelectTrigger>
                                                <SelectContent>
                                                    {isLoading ? (
                                                        <SelectItem value="loading">Đang tải...</SelectItem>
                                                    ) : (
                                                        decisionTypes?.map((item) => (
                                                            <SelectItem key={item.decision_type_id}
                                                                        value={item.decision_type_id}>
                                                                {item.decision_type}
                                                            </SelectItem>
                                                        ))
                                                    )}
                                                </SelectContent>
                                            </Select>
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm"/>
                                    </FormItem>
                                )}
                            />

                            <FormItem>
                                <FormLabel className="font-medium">Nhân viên</FormLabel>
                                <FormControl>
                                    <MultiSelectEmployeeComponent
                                        employees={employees || []}
                                        selectedEmployees={currentSelectedEmployees}
                                        value={selectedEmployeeIds}
                                        onChange={handleEmployeeChange}
                                        disabled={isLoadingEmployee}
                                    />
                                </FormControl>
                                <FormMessage className="text-red-500 text-sm"/>
                            </FormItem>

                            <FormField
                                control={form.control}
                                name="content"
                                render={({field}) => (
                                    <FormItem>
                                        <FormLabel className="font-medium">Nội dung</FormLabel>
                                        <FormControl>
                                            <textarea
                                                {...field}
                                                className="w-full h-[140px] p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-[#DB3B21] resize-none"
                                                placeholder="Nhập nội dung quyết định"
                                            />
                                        </FormControl>
                                        <FormMessage className="text-red-500 text-sm"/>
                                    </FormItem>
                                )}
                            />
                        </div>

                        {/* Buttons */}
                        <div className="w-full flex justify-end gap-4 pt-6">
                            <Button
                                type="button"
                                onClick={handleClose}
                                className="bg-gray-400 hover:bg-gray-500 w-[120px] h-[40px]"
                            >
                                Hủy bỏ
                            </Button>
                            <Button
                                type="submit"
                                className="bg-[#DB3B21] hover:bg-[#b83a1a] w-[120px] h-[40px]"
                                disabled={form.formState.isSubmitting}
                            >
                                {form.formState.isSubmitting ? 'Đang xử lý...' : 'Lưu thông tin'}
                            </Button>
                        </div>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
};

export default UpdateDecisionForm;