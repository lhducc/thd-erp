import type {AttendanceSetting} from "@/types/attendance.ts";
import api from "@/apis/api.ts";

export const getAllAttendanceManagementAPI = async (): Promise<AttendanceSetting[]> => {
    try {
        const response = await api.get("/attendance-category");
        return response.data.data;
    } catch (error) {
        console.error("Error create contract API:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi tạo hợp đồng");
    }
}

export const createAttendanceSettingApi = async (payload: AttendanceSetting) => {
    try {
        const response = await api.post("/attendance-category", payload);
        return response.data.data;
    } catch (error) {
        console.error("Error create contract API:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi tạo hợp đồng");
    }
}

export const updateAttendanceSettingApi = async (id: string, payload: AttendanceSetting) => {
    try {
        const response = await api.put(`/attendance-category/${id}`, payload);
        console.log(response.data);
        return response.data.data;
    } catch (error) {
        console.error("Error create contract API:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi tạo hợp đồng");
    }
}

export const deleteAttendanceSettingApi = async (id: string) => {
    try {
        const response = await api.delete(`/attendance-category/${id}`);
        console.log(response.data);
        return response.data.data;
    } catch (error) {
        console.error("Error create contract API:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi tạo hợp đồng");
    }
}