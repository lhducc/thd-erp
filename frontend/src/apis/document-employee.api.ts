import api from "@/apis/api.ts";

export const getDocumentEmployeeById = async (id: string) => {
    try {
        const response = await api.get("/employee-document/" + id);
        return response.data.data;
    } catch (error) {
        console.error("Error fetching document employee API:", error);
        throw new Error(error.response.data.message);
    }
}