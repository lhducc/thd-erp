import type {AttendanceRecord, AttendanceRecordHistoryByDate, CreateManualRecord} from "@/types/attendance.ts";
import api from "@/apis/api.ts";

export const getAttendanceRecordByEmployeeId = async (id: string): Promise<AttendanceRecord[]> => {
    const response = await api.get(`/attendance-record/history-record/${id}`);
    return response.data.data;
}

export const getAttendanceRecordPersonalApi = async (page: number, limit: number): Promise<AttendanceRecord[]> => {
    const response = await api.get(`/attendance-record/personal?page=${page}&limit=${limit}`);
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

export const exportAttendanceExcelApi = async (params: {date?: string;}): Promise<Blob> => {
  const response = await api.get("/attendance-record/attendance/export", {
    params,
    responseType: "blob",
  });
  return response.data;
};


export const getEmployeesByDateApi = async (date: string): Promise<AttendanceRecordHistoryByDate[]> => {
  const q = encodeURIComponent(date); 
  const response = await api.get(`/attendance-record/history/by-date?date=${q}`);
  return response.data.data;
};

export const getEmployeesByMonthApi = async (month: string): Promise<AttendanceRecordHistoryByDate[]> => {
  const response = await api.get(`/attendance-record/attendance?month=${month}`);
  return response.data.data;
};
