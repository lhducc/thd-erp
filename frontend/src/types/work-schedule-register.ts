import type {WorkShift} from "@/types/Workshift.ts";
import type {Employee, Office} from "@/types/index.ts";

export type WorkScheduleRegister = {
    work_schedule_register_id: number;
    work_schedule_register_name: string;
    office_id: string;
    effective_date: string; // ISO date string
    expiration_date: string; // ISO date string
    status: string;
    created_at: string; // ISO date string
    is_deleted: boolean;
    managers: WorkScheduleRegisterManager[];
    employees: WorkScheduleRegisterEmployee[];
    weekdays: WorkScheduleRegisterShift[];
    office: Office;
};

type WorkScheduleRegisterManager = {
    work_schedule_register_manager_id: number;
    employee_id: string;
    work_schedule_register_id: number;
    is_reading: boolean;
    is_editing: boolean;
    manager: Employee;
};

type WorkScheduleRegisterEmployee = {
    work_schedule_register_employee_id: number;
    employee_id: string;
    work_schedule_register_id: number;
    assigned_at: string; // ISO date string
    employee: Employee;
};

type WorkScheduleRegisterShift = {
    work_schedule_register_shift_id: number;
    work_schedule_register_id: number;
    week_day: string;
    workshift_id: string;
    order: number;
    work_shift: WorkShift;
};