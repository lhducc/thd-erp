import type {AttendanceRecord, CreateManualRecord} from "@/types/attendance.ts";
import api from "@/apis/api.ts";

export const getAttendanceRecordByEmployeeId = async (id: string): Promise<AttendanceRecord[]> => {
    const response = await api.get(`/attendance-record/history-record/${id}`);
    return response.data.data;
}

export const createAttendanceRecordByAdminId = async (payload: CreateManualRecord): Promise<AttendanceRecord> => {
    const response = await api.post(`/attendance-record/history-record/manual`, payload);
    return response.data.data;
}

export const createAttendanceRecord = async (payload): Promise<AttendanceRecord> => {
    const response = await api.post('/attendance-record', payload, {
        headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data.data;
}