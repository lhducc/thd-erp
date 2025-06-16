import api from "./api";
import type {Contract, ContractFormValues} from "@/types/contract.ts";

export const createContractApi = async (payload: ContractFormValues) => {
    try {
        const response = await api.post("/contract", payload);
        return response.data;
    } catch (error: any) {
        console.error("Error create contract API:", error);
        throw new Error(error.response?.data?.message || "Lỗi khi tạo hợp đồng");
    }
};


export const updateContractApi = async (id: string, payload: ContractFormValues) => {
    try {
        const response = await api.put(`/contract/${id}`, payload);
        return response.data;
    } catch (error: any) {
        console.error("Error update contract API:", error);
        throw new Error(error.response?.data?.message || "Lỗi khi cập nhật hợp đồng");
    }
};

export const getAllContractsApi = async (): Promise<Contract[]> => {
    try {
        const response = await api.get("/contract");
        console.log(response);
        // return response.data.data;
        const mock: Contract[] = [{
            "contract_id": "C001",
            "effective_date": "2025-01-01T00:00:00Z",
            "expired_date": "2025-12-31T23:59:59Z",
            "sign_date": "2024-12-15T00:00:00Z",
            "note": "Annual employment contract",
            "attached_file": "https://example.com/contracts/C001.pdf",
            "condition": "Full-time employment",
            "created_date": "2024-12-16T10:00:00Z",
            "contract_type": "Employment",
            "approve_status": "Đã duyệt",
            "employee": {
                "employee_id": "E1001",
                "full_name": "Alice Johnson",
                "department": {
                    "department_id": "D001",
                    "department_name": "Human Resources",
                    "office": {
                        "office_id": "O001",
                        "office_name": "Head Office"
                    }
                }
            }
        },
            {
                "contract_id": "C002",
                "effective_date": "2025-02-01T00:00:00Z",
                "expired_date": "2026-01-31T23:59:59Z",
                "sign_date": "2025-01-10T00:00:00Z",
                "note": "Consulting agreement for marketing project",
                "attached_file": "https://example.com/contracts/C002.pdf",
                "condition": "Fixed-term contract",
                "created_date": "2025-01-11T08:30:00Z",
                "contract_type": "Consulting",
                "approve_status": "Chưa duyệt",
                "employee": {
                    "employee_id": "E1002",
                    "full_name": "Bob Smith",
                    "department": {
                        "department_id": "D002",
                        "department_name": "Marketing",
                        "office": {
                            "office_id": "O002",
                            "office_name": "Branch Office"
                        }
                    }
                }
            }, {
                "contract_id": "C002",
                "effective_date": "2025-02-01T00:00:00Z",
                "expired_date": "2026-01-31T23:59:59Z",
                "sign_date": "2025-01-10T00:00:00Z",
                "note": "Consulting agreement for marketing project",
                "attached_file": "https://example.com/contracts/C002.pdf",
                "condition": "Fixed-term contract",
                "created_date": "2025-01-11T08:30:00Z",
                "contract_type": "Consulting",
                "approve_status": "Chưa duyệt",
                "employee": {
                    "employee_id": "E1002",
                    "full_name": "Bob Smith",
                    "department": {
                        "department_id": "D002",
                        "department_name": "Marketing",
                        "office": {
                            "office_id": "O002",
                            "office_name": "Branch Office"
                        }
                    }
                }
            }]
        return mock;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
};

export const getContractById = async (contractId: string): Promise<Contract> => {
    const mock: Contract = {
        "contract_id": contractId,
        "effective_date": "2025-01-01T00:00:00Z",
        "expired_date": "2025-12-31T23:59:59Z",
        "sign_date": "2024-12-15T00:00:00Z",
        "note": "Annual employment contract",
        "attached_file": "https://example.com/contracts/C001.pdf",
        "condition": "Hiệu lực",
        "created_date": "2024-12-16T10:00:00Z",
        "contract_type": "LH0001",
        "approve_status": "Đã duyệt",
        "employee": {
            "employee_id": "THD011",
            "full_name": "Alice Johnson",
            "department": {
                "department_id": "D001",
                "department_name": "Human Resources",
                "office": {
                    "office_id": "O001",
                    "office_name": "Head Office"
                }
            }
        }
    }
    return mock;
}

export const exportContractsFile = async (): Promise<File> => {
    try {
        console.log("heree")
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
