import api from "./api";
import type {Contract, ContractFormValues} from "@/types/contract.ts";

export const createContractApi = async (payload: ContractFormValues) => {
    const formData = new FormData();

    // Thêm các trường thông thường
    formData.append('contract_type', payload.contract_type);
    formData.append('employee_id', payload.employee_id);
    formData.append('sign_date', payload.sign_date);
    formData.append('effective_date', payload.effective_date);
    formData.append('expired_date', payload.expired_date);
    formData.append('approve_status', payload.approve_status || 'Chờ duyệt');
    formData.append('condition', payload.condition || 'Chưa hiệu lực');

    if (payload.note) {
        formData.append('note', payload.note);
    }

    if (payload.department) {
        formData.append('department', payload.department);
    }

    // Thêm allowance_ids dưới dạng mảng
    if (payload.allowance_ids && payload.allowance_ids.length > 0) {
        payload.allowance_ids.forEach((id, index) => {
            formData.append(`allowance_ids[${index}]`, id);
        });
    }

    // Thêm file nếu có
    if (payload.attached_file) {
        formData.append('attached_file', payload.attached_file);
    }

    const response = await api.post("/contract", formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
    return response.data;
};

export const updateContractApi = async (id: string, payload: Partial<ContractFormValues>) => {
    const formData = new FormData();

    // Thêm các trường có giá trị
    if (payload.contract_type) formData.append('contract_type', payload.contract_type);
    if (payload.employee_id) formData.append('employee_id', payload.employee_id);
    if (payload.sign_date) formData.append('sign_date', payload.sign_date);
    if (payload.effective_date) formData.append('effective_date', payload.effective_date);
    if (payload.expired_date) formData.append('expired_date', payload.expired_date);
    if (payload.approve_status) formData.append('approve_status', payload.approve_status);
    if (payload.condition) formData.append('condition', payload.condition);
    if (payload.note !== undefined) formData.append('note', payload.note || '');
    if (payload.department) formData.append('department', payload.department);

    // Thêm allowance_ids
    if (payload.allowance_ids) {
        payload.allowance_ids.forEach((id, index) => {
            formData.append(`allowance_ids[${index}]`, id);
        });
    }

    // Thêm file nếu có
    if (payload.attached_file) {
        formData.append('attached_file', payload.attached_file);
    }

    const response = await api.put(`/contract/${id}`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
    return response.data;
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

export const getContractByEmployeeId = async (employeeId: string): Promise<Contract[]> => {
    const response = await api.get(`/contract/employee/${employeeId}`);
    return response.data.data;
}

export const exportContractsFile = async (): Promise<File> => {
    try {
        const response = await api.get("/contract/export?fields=contract_id,effective_date", {
            responseType: "blob",
        });

        const fileName = "contracts.xlsx";
        const file = new File([response.data], fileName, {type: response.data.type});
        return file;
    } catch (error: any) {
        console.error("Error exporting contracts:", error);
        throw error;
    }
};
