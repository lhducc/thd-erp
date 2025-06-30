import api from "./api";
import type {Contract, ContractFormValues} from "@/types/contract.ts";

export const createContractApi = async (payload: ContractFormValues) => {
    try {
        const response = await api.post("/contract", payload);
        return response.data;
    } catch (error: any) {
        console.error("Error create contract API:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi tạo hợp đồng");
    }
};

export const updateContractApi = async (id: string, payload: ContractFormValues) => {
    try {
        const response = await api.put(`/contract/${id}`, payload);
        return response.data;
    } catch (error: any) {
        console.error("Error update contract API:", error);

        throw {
            message: error.response?.data?.message || "Lỗi khi cập nhật hợp đồng",
            status: error.response?.status,
            data: error.response?.data
        };
    }
};

export const reapproveContractApi = async (id: string) => {
    try {
        const response = await api.put(`/contract/${id}/reapprove`);
        return response.data;
    } catch (error: any) {
        console.error("Error update contract API:", error);

        throw {
            message: error.response?.data?.message || "Lỗi khi cập nhật hợp đồng",
            status: error.response?.status,
            data: error.response?.data
        };
    }
};

export const getAllContractsApi = async (): Promise<Contract[]> => {
    try {
        const response = await api.get("/contract?page=1&pageSize=1000");
        return response.data.data.data;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
};

export const getContractById = async (contractId: string): Promise<Contract> => {
    try {
        const response = await api.get(`/contract/${contractId}`);
        return response.data.data;
    } catch (error: any) {
        console.error("Error update contract API:", error);
        throw new Error(error.response?.data?.message || "Lỗi khi cập nhật hợp đồng");
    }
}

export const deleteContractById = async (contractId: string): Promise<Contract> => {
    try {
        const response = await api.delete(`/contract/${contractId}`);
        return response.data.data;
    } catch (error: any) {
        console.error("Error update contract API:", error);
        throw new Error(error.response?.data?.message || "Lỗi khi cập nhật hợp đồng");
    }
}

export const exportContractsFile = async (): Promise<File> => {
    try {
        const response = await api.get("/contract/export?fields=contract_id,effective_date", {
            responseType: "blob",
        });

        const fileName = "contracts.xlsx";
        const file = new File([response.data], fileName, { type: response.data.type });
        return file;
    } catch (error: any) {
        console.error("Error exporting contracts:", error);
        throw error;
    }
};
