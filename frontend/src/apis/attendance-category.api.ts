import api from "@/apis/api.ts";
import type {AttendanceCategory} from "@/types/attendance-category.ts";

export const getAttendanceCategories = async () : Promise<AttendanceCategory[]> => {
    const response = await api.get(`/attendance-category/office`);
    console.log(response.data.data);
    return response.data.data;
}