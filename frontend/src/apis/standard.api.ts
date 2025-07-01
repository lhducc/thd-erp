import api from "./api";
import type {PayloadStandard, Standard} from "@/types/standard.ts";

export const createStandardApi = async (payload: PayloadStandard) => {
    try {
        const response = await api.post("/standard", payload);
        return response.data;
    } catch (error: any) {
        console.error("Error creating standard API:", error);
        throw new Error(error.response.data.message);
    }
};

export const getAllStandardsApi = async (): Promise<Standard[]> => {
    try {
        const response = await api.get("/standard");
        return response.data.data;
    } catch (error) {
        console.error("Error fetching all standard API:", error);
        throw error;
    }
};

export const updateStandardApi = async (id: string, payload: PayloadStandard) => {
    try {
        const response = await api.put(`/standard/${id}`, payload);
        return response.data;
    } catch (error: any) {
        console.error("Error update standard API:", error);
        throw new Error(error.response.data.message);
    }
};

export const deleteStandardApi = async (id: string) => {
    try {
        const response = await api.delete(`/standard/${id}`);
        return response.data;
    } catch (error: any) {
        console.error("Error delete standard API:", error);
        throw new Error(error.response.data.message);
    }
};
