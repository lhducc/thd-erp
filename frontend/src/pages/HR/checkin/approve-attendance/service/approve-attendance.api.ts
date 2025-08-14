import api from "@/apis/api.ts";

export const getEmployeeApprovalAttendanceApi = async (id: string): Promise<any> => {
    const response = await api.get(`/attendance-record/history-record/${id}`);
    return response.data.data;
}