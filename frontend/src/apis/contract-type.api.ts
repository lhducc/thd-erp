import type {ContractType} from "@/types/contract.ts";
import api from "@/apis/api.ts";

export const getAllContractsTypeApi = async (): Promise<ContractType[]> => {
    try {
        const response = await api.get("/contracttype");
        console.log(response);
        console.log(response.data.data);
        return response.data.data;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}