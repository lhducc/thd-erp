import {useQuery} from "@tanstack/react-query";
import {
    getEmployeeApprovalAttendanceApi
} from "@/pages/HR/checkin/approve-attendance/service/approve-attendance.api.ts";

export const useAttendanceDetail = (id : string) =>
    useQuery({
        queryKey: ["approveAttendance", id],
        queryFn: () => getEmployeeApprovalAttendanceApi(id)
    })