import api from "@/apis/api.ts";
import type {DocumentType, DocumentTypeCreate} from "@/types/document-type.ts";

export const getAllDocumentTypesApi = async (): Promise<DocumentType[]> => {
    try {
        const response = await api.get("/documenttype");
        return response.data.data;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}

export const createDocumentTypeApi = async (payload: DocumentTypeCreate): Promise<string> => {
    try {
        const response = await api.post(`/documenttype`, payload);
        return response.data.message;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}

export const updateDocumentTypeApi = async (id: string, payload: DocumentTypeCreate): Promise<string> => {
    try {
        const response = await api.put(`/documenttype/${id}`, payload);
        return response.data.message;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}

export const deleteDocumentTypesApi = async (id: string): Promise<string> => {
    try {
        const response = await api.delete(`/documenttype/${id}`);
        return response.data.message;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}