import type {Employee} from "@/types/employee";

export type AttendanceNeedApprovalType = {
    employee: Employee;
    gps: string;
    checkin_image_url: string;
    checkout_image_url: string;
    work_type: string;
    status: string;
}