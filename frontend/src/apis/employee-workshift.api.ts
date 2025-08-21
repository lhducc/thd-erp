import type {DocumentType} from "@/types/document-type.ts";
import api from "@/apis/api.ts";
import type {RegisterWorkshiftRequest} from "@/types/employee-workshift.ts";

export const getAllEmployeeWorkshiftsApi = async (currentMonth: number, currentYear: number): Promise<DocumentType[]> => {
    try {
        const response = await api.get(`/employee-workshifts/all?month=${currentMonth}&year=${currentYear}`);
        return response.data.data;
    } catch (error: any) {
        console.error("Error fetching all contracts API:", error);
        throw new Error(error.response.data.message);
    }
}

export const registerManyWorkshiftsApi = async (
    data: RegisterWorkshiftRequest[]
): Promise<void> => {
    try {
        await api.post('/employee-workshifts/register-many', data);
    } catch (error) {
        console.error('Error registering workshifts:', error);
        throw error;
    }
};