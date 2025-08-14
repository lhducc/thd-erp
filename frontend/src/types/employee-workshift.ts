import type {Workshift} from "@/types/workshift.ts";

export type EmployeeWorkshift = {
    id: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string;
    employee_id: string;
    workshift_id: string;
    workshift: Workshift
}