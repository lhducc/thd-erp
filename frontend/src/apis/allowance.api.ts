import api from "@/apis/api.ts";
import type {Allowance} from "@/types/allowance.ts";

export const getAllAllowancesApi = async (): Promise<Allowance[]> => {
    try {
        const response = await api.get("/allowance");
        return response.data.data;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}

export const createAllowanceApi = async (payload: Allowance): Promise<string> => {
    try {
        const response = await api.post("/allowance", {
            ...payload,
            is_delete: false,
            created_date: new Date().toISOString()
        });
        return response.data.message;
    } catch (error: any) {
        console.error("Error creating contract type:", error);
        throw new Error(error.response.data.message);
    }
}

export const updateAllowanceApi = async (id: string, payload: Partial<Allowance>): Promise<string> => {
    try {
        const response = await api.put(`/allowance/${id}`, payload);
        return response.data.message;
    } catch (error: any) {
        console.error("Error updating contract type:", error);
        throw new Error(error.response.data.message);
    }
}

export const deleteAllowanceApi = async (id: string): Promise<string> => {
    try {
        console.log(id)
        const response = await api.delete(`/allowance/${id}`);
        console.log(response);
        return response.data.message;
    } catch (error: any) {
        console.error("Error updating contract type:", error);
        throw new Error(error.response.data.message);
    }
}