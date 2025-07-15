import api from "@/apis/api.ts";

export const getAllSetupWorkshiftAPI = async () => {
    try {
        const response = await api.get("/employee-workshifts");
        return response.data;
    } catch (error: any) {
        console.error("Error create contract API:", error.message);
        throw new Error(error.response?.data?.error || "Lỗi khi tạo hợp đồng");
    }
}