import api from "@/apis/api.ts";

export const createWorkshiftSchedule = async (data: any): Promise<void> => {
    try {
        const response = await api.post(`/work-schedule`, data);
        return response.data;
    } catch (error) {
        console.error("Error create work schedule API:", error);
        throw new Error(error.response.data.message);
    }
}

export const deleteWorkshiftSchedule = async (id: string): Promise<void> => {
    try {
        const response = await api.delete(`/work-schedule/${id}`);
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