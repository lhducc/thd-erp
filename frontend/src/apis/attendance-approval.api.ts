import api from "@/apis/api.ts";

export const totalRequestAttendanceApi = async (id: string): Promise<number> => {
    const response = await api.get(`/attendance-record/employee/${id}/total-request`);
    return response.data.data;
}

export const createAttendanceApprovalApi = async () => {
    // Mock implementation
    return { success: true };
}

export const updateAttendanceApprovalApi = async () => {
    // Mock implementation
    return { success: true };
}