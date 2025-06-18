import type {ContractType} from "@/types/contract.ts";
import api from "@/apis/api.ts";

export const getAllContractsTypeApi = async (): Promise<ContractType[]> => {
    try {
        const response = await api.get("/contracttype");
        return response.data.data;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}

export interface ContractTypeCreate {
    contract_type: string;
    contract_group: string;
    duration: number;
    unit: string;
    working_form: string;
}

export interface ContractType extends ContractTypeCreate {
    contract_type_id: string;
    is_delete: boolean;
    created_date: string;
}

export const createContractTypeApi = async (payload: ContractTypeCreate): Promise<string> => {
    try {
        const response = await api.post("/contracttype", {
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

export const updateContractTypeApi = async (id: string, payload: Partial<ContractTypeCreate>): Promise<string> => {
    try {
        const response = await api.put(`/contracttype/${id}`, payload);
        return response.data.message;
    } catch (error: any) {
        console.error("Error updating contract type:", error);
        throw new Error(error.response.data.message);
    }
}

export const deleteContractTypeApi = async (id: string): Promise<string> => {
    try {
        console.log(id)
        const response = await api.delete(`/contracttype/${id}`);
        console.log(response);
        return response.data.message;
    } catch (error: any) {
        console.error("Error updating contract type:", error);
        throw new Error(error.response.data.message);
    }
}