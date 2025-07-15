import type {WorkShift} from "@/types/Workshift.ts";

export type EmployeeWorkshift = {
    ID: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string;
    employee_id: string;
    work_shift_id: string;
    work_shift: WorkShift
}