import api from "@/apis/api.ts";
import type {WorkSchedule} from "@/types/work-schedule.ts";

export const getAllSetupWorkshiftAPI = async () => {
    try {
        const response = await api.get("/employee-workshifts");
        return response.data;
    } catch (error) {
        console.error("Error create contract API:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi lấy cài đặt ca");
    }
}

export const getAllSetupWorkshiftCalendar = async () : Promise<WorkSchedule[]> => {
    try {
        const response = await api.get("/work-schedule");
        return response.data.data;
    } catch (error) {
        console.error("Error getting contract calendar:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi lấy lịch làm việc");
    }
}