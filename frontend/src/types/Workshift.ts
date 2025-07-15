export enum WorkDayEnum {
    FullDay = "1",
    HaftDay = "0.5",
    NoWork = "0",
}

export enum TimeOfDayEnum {
    Morning = "Sáng",
    Afternoon = "Trưa",
    Evening = "Tối",
    AllDay = "Cả ngày",
}

export type WorkShift = {
    workshift_id: string;
    workshift_name: string;
    start_time: string;           // dạng 'HH:mm:ss'
    end_time: string;             // dạng 'HH:mm:ss'
    checkin_from?: string | null; // dạng 'HH:mm:ss'
    checkin_to?: string | null;
    checkout_from?: string | null;
    checkout_to?: string | null;
    has_break: boolean;
    break_start?: string | null;
    break_end?: string | null;
    work_hours: string;
    work_day: WorkDayEnum;
    coef_normal_day: number;
    coef_weekend: number;
    coef_holiday: number;
    created_by: string;
    created_date: string;         // ISO date string (e.g., 2025-07-01T08:00:00Z)
    is_deleted: boolean;
    time_of_day: TimeOfDayEnum;
    effective_date: string;       // ISO date string
    expiration_date?: string | null;
    creator?: {
        employee_id: string;
    };
};

// Dạng dữ liệu cho request
export type WorkShiftRequest = {
    workshift_id: string;
    workshift_name: string;
    start_time: string;
    end_time: string;
    checkin_from?: string | null;
    checkin_to?: string | null;
    checkout_from?: string | null;
    checkout_to?: string | null;
    has_break: boolean;
    break_start?: string | null;
    break_end?: string | null;
    work_hours: string;
    work_day: WorkDayEnum;
    coef_normal_day: number;
    coef_weekend: number;
    coef_holiday: number;
    created_by: string;
    time_of_day: TimeOfDayEnum;
    effective_date: string;        // ISO date string
    expiration_date?: string | null;
};