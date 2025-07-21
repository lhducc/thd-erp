import api from "@/apis/api.ts";
import type {WorkSchedule, WorkScheduleRegisterPayload, WorkScheduleResponse} from "@/types/work-schedule.ts";
import type {WorkScheduleRegister} from "@/types/work-schedule-register.ts";

export const createWorkshiftSchedule = async (data: any): Promise<void> => {
    try {
        const response = await api.post(`/work-schedule`, data);
        return response.data;
    } catch (error) {
        console.error("Error create work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}

export const getWorkSchedule = async () : Promise<WorkSchedule[]> => {
    try {
        const response = await api.get("/work-schedule");
        console.log(response.data.data)
        return response.data.data;
    } catch (error) {
        console.error("Error getting contract calendar:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi lấy lịch làm việc");
    }
}

export const getAllWorkScheduleRegister = async () : Promise<WorkScheduleRegister[]> => {
    try {
        const response = await api.get("/work-schedule-register");
        return response.data.data;
    } catch (error) {
        console.error("Error getting contract calendar:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi lấy lịch làm việc");
    }
}

export const getWorkScheduleById = async (id: string): Promise<WorkScheduleResponse> => {
    try {
        const response = await api.get(`/work-schedule/${id}`);
        return response.data.data;
    } catch (error) {
        console.error("Error fetching work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}

export const getWorkScheduleRegisterById = async (id?: string) => {
    try {
        const response = await api.get(`/work-schedule-register/${id}`);
        return response.data.data;
    } catch (error) {
        console.error("Error fetching work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}


export const deleteWorkshiftSchedule = async (id: number): Promise<void> => {
    try {
        const response = await api.delete(`/work-schedule/${id}`);
        return response.data;
    } catch (error) {
        console.error("Error assign work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}

export const updateWorkshiftSchedule = async (id: number, data: WorkScheduleRegisterPayload): Promise<void> => {
    try {
        const response = await api.post(`/work-schedule-schedule/${id}`, data);
        return response.data;
    } catch (error) {
        console.error("Error assign work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}
export const updateWorkScheduleRegisterApi = async (id: number, data: WorkScheduleRegisterPayload): Promise<void> => {
    try {
        const response = await api.put(`/work-schedule-schedule/${id}`, data);
        return response.data;
    } catch (error) {
        console.error("Error assign work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}

export const assignWorkshiftSchedule = async (id: string, payload:any) => {
    try {
        const response = await api.post(`/work-schedule/assign/${id}`, payload);
        return response.data;
    } catch (error) {
        console.error("Error assign work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}

export const registerWorkScheduleRegisterApi = async (payload: WorkScheduleRegisterPayload) => {
    try {
        const response = await api.post(`/work-schedule-register`, payload);
        return response.data;
    } catch (error) {
        console.error("Error assign work schedule API:", error);
        throw new Error(error.response.data.message);
    }
};