import { useEffect, useState } from "react";

interface MonthYearPickerProps {
    value?: string; // "2025-07"
    onChange: (value: string) => void;
    label?: string;
    name?: string;
}

export const MonthYearSelect = ({
                                    value,
                                    onChange,
                                    label = "Chọn tháng",
                                    name,
                                }: MonthYearPickerProps) => {
    const [displayLabel, setDisplayLabel] = useState("");

    useEffect(() => {
        if (value) {
            const [year, month] = value.split("-").map(Number);
            if (!isNaN(year) && !isNaN(month)) {
                setDisplayLabel(`Tháng ${month} năm ${year}`);
            } else {
                setDisplayLabel("");
            }
        } else {
            setDisplayLabel("");
        }
    }, [value]);

    return (
        <div className="flex flex-col gap-1">
            {label && <label className="text-sm font-medium">{label}</label>}
            <input
                type="month"
                name={name}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                className="border rounded-md p-2"
            />
            {displayLabel && (
                <p className="text-sm text-gray-500 italic">{displayLabel}</p>
            )}
        </div>
    );
};
