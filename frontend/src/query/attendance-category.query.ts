import {useQuery} from "@tanstack/react-query";
import {getAttendanceCategories} from "@/apis/attendance-category.api.ts";

export const useAttendanceCategories = () =>
    useQuery({
        queryKey: ['attendance_category'],
        queryFn: getAttendanceCategories
    })