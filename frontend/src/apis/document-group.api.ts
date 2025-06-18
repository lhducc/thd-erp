import type {DocumentGroup} from "@/types/document-group.ts";

export const getAllDocumentsGroupApi = async (): Promise<DocumentGroup[]> => {
    try {
        const mock: DocumentGroup[] = [
            {
                id: "1",
                name: "Loại chứng chỉ",
            },
            {
                id: "2",
                name: "Loại chứng từ",
            },
            {
                id: "3",
                name: "Hồ sơ nhân sự"
            }
        ]
        return mock;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}