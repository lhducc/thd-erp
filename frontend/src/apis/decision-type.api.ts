import api from "@/apis/api.ts";
import type {DecisionType} from "@/types/decistion-type.ts";

export const getAllDecisionTypeApi = async (): Promise<DecisionType[]> => {
    try {
        const response = await api.get("/decisiontype");
        return response.data.data;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}

export const createDecisionTypeApi = async (payload: DecisionType): Promise<string> => {
    try {
        const response = await api.post("/decisiontype", {
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

export const updateDecisionTypeApi = async (id: string, payload: Partial<DecisionType>): Promise<string> => {
    try {
        const response = await api.put(`/decisiontype/${id}`, payload);
        return response.data.message;
    } catch (error: any) {
        console.error("Error updating contract type:", error);
        throw new Error(error.response.data.message);
    }
}

export const deleteDecisionTypeApi = async (id: string): Promise<string> => {
    try {
        console.log(id)
        const response = await api.delete(`/decisiontype/${id}`);
        console.log(response);
        return response.data.message;
    } catch (error: any) {
        console.error("Error updating contract type:", error);
        throw new Error(error.response.data.message);
    }
}