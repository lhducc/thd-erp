import type {Department} from "@/types/index.ts";

export type Attendance = {

}

export type AttendanceSetting = {
    name: string,
    department: Department,
    useCamera: boolean,
    gpsEnabled: boolean,
    radius: boolean,
    status: boolean,
}