"use client";

import * as React from "react";

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

export function MultiSelectEmployee({
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
        <div className="w-full">

        </div>
    );
}