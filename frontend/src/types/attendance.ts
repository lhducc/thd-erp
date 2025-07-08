import type {Office} from "@/types/index.ts";

export type Attendance = {

}

export type AttendanceSetting = {
    attendance_category_id: string,
    attendance_category_name: string,
    office_id: string,
    Office: Office,
    is_camera: boolean,
    is_gps: boolean,
    is_check_location: boolean,
    scope?: number,
    status: string,
}