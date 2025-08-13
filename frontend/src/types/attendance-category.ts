import type {Office} from "@/types/index.ts";

export type AttendanceCategory = {
    attendance_category_id: string;
    attendance_category_name: string;
    is_gps: boolean;
    is_camera: boolean;
    is_check_location: boolean;
    scope: number;
    status: "active" | "inactive"; // có thể mở rộng thêm nếu cần
    is_deleted: boolean;
    office_id: string;
    created_by: string;
    created_at: string; // ISO date string
    auto_approve: boolean;
    Office: Office;
};