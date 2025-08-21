import api from "@/apis/api.ts";

export const getDocumentEmployeeById = async (id: string) => {
    const response = await api.get("/employee-document/employee/" + id);
    return response.data.data;
}