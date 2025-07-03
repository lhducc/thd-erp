import type {AttendanceSetting} from "@/types/attendance.ts";

export const getAllAttendanceManagementAPI = async (): Promise<AttendanceSetting[]> => {
    return [
        {
            name: 'Tai van phong',
            department: {
                department_id: "PB0001",
                department_name: "Nghiên cứu và Phát triển",
                manager: "",
                created_date: Date.now().toString(),
                office_id: "",
                office: null
            },
            useCamera: true,
            gpsEnabled: true,
            radius: true,
            status: true,
        }
    ]
}

export const createAttendanceSettingApi = async () => {

}

export const updateAttendanceSettingApi = async () => {

}