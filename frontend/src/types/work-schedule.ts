import type {Office} from "@/types/index.ts";

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
    office: Office;
};