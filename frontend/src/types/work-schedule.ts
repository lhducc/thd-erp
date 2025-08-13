import type {ManagerPermission} from "@/types/employee.ts";

export type WorkSchedule = {
    work_schedule_id: number;
    work_schedule_name: string;
    office_id: string;
    repeat_type: string;
    repeat_cycle: number;
    effective_date: string; // ISO date string
    expiration_date: string; // ISO date string
    status: string;
    created_at: string; // ISO date string
    is_deleted: boolean;
    managers: ManagerPermission[];
    office: Office;
};

export type WorkShift = {
    workshift_id: string;
    workshift_name: string;
    start_time: string;
    end_time: string;
    checkin_from: string;
    checkin_to: string;
    checkout_from: string;
    checkout_to: string;
    has_break: boolean;
    break_start: string | null;
    break_end: string | null;
    work_hours: number;
    work_day: string;
    coef_normal_day: number;
    coef_weekend: number;
    coef_holiday: number;
    created_by: string;
    created_date: string;
    is_deleted: boolean;
    time_of_day: string;
    effective_date: string;
    expiration_date: string;
};

export type Weekday = {
    work_schedule_id: number;
    week_day: string;
    workshift_id: string;
    order: number;
    work_shift: WorkShift;
};

export type ScheduleAssignment = {
    employee_id: string;
    work_schedule_id: number;
    assigned_at: string;
};

export type Office = {
    office_id: string;
    office_name: string;
};

export type WorkScheduleResponse = {
    work_schedule_id: number;
    work_schedule_name: string;
    office_id: string;
    repeat_type: string;
    repeat_cycle: number;
    effective_date: string;
    expiration_date: string;
    status: string;
    created_at: string;
    is_deleted: boolean;
    managers: ManagerPermission[];
    employees: ScheduleAssignment[];
    weekdays: Weekday[];
    office: Office;
};

export interface WeekdaySelection {
    week_day: string;
    workshift_id: string;
}

export interface WorkScheduleRegisterPayload {
    work_schedule_register_name: string;
    office_id: string;
    effective_date: string;
    expiration_date: string;
    weekdays: WeekdaySelection[];
}
