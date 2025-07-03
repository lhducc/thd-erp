import {TimeOfDayEnum, WorkDayEnum, type WorkShift, type WorkShiftRequest} from "@/types/Workshift.ts";
import api from "@/apis/api.ts";

// // Mock data storage
// const mockWorkShifts: WorkShift[] = [
//     {
//         workshift_id: "1",
//         workshift_name: "Morning Shift",
//         start_time: "08:00:00",
//         end_time: "12:00:00",
//         checkin_from: "07:30:00",
//         checkin_to: "08:30:00",
//         checkout_from: "11:30:00",
//         checkout_to: "12:30:00",
//         has_break: false,
//         break_start: null,
//         break_end: null,
//         work_hours: 4,
//         work_day: WorkDayEnum.Weekday,
//         coef_normal_day: 1.0,
//         coef_weekend: 1.5,
//         coef_holiday: 2.0,
//         created_by: "user1",
//         created_date: "2025-01-01T00:00:00Z",
//         is_deleted: false,
//         time_of_day: TimeOfDayEnum.Morning,
//         effective_date: "2025-01-01",
//         expiration_date: null,
//         creator: {
//             employee_id: "emp1"
//         }
//     },
//     {
//         workshift_id: "2",
//         workshift_name: "Afternoon Shift",
//         start_time: "13:00:00",
//         end_time: "17:00:00",
//         checkin_from: "12:30:00",
//         checkin_to: "13:30:00",
//         checkout_from: "16:30:00",
//         checkout_to: "17:30:00",
//         has_break: true,
//         break_start: "15:00:00",
//         break_end: "15:30:00",
//         work_hours: 4,
//         work_day: WorkDayEnum.Weekday,
//         coef_normal_day: 1.0,
//         coef_weekend: 1.5,
//         coef_holiday: 2.0,
//         created_by: "user1",
//         created_date: "2025-01-01T00:00:00Z",
//         is_deleted: false,
//         time_of_day: TimeOfDayEnum.Afternoon,
//         effective_date: "2025-01-01",
//         expiration_date: null,
//         creator: {
//             employee_id: "emp1"
//         }
//     },
//     {
//         workshift_id: "3",
//         workshift_name: "Weekend Shift",
//         start_time: "09:00:00",
//         end_time: "15:00:00",
//         checkin_from: "08:30:00",
//         checkin_to: "09:30:00",
//         checkout_from: "14:30:00",
//         checkout_to: "15:30:00",
//         has_break: true,
//         break_start: "12:00:00",
//         break_end: "12:30:00",
//         work_hours: 5.5,
//         work_day: WorkDayEnum.Weekend,
//         coef_normal_day: 1.0,
//         coef_weekend: 1.5,
//         coef_holiday: 2.0,
//         created_by: "user2",
//         created_date: "2025-01-15T00:00:00Z",
//         is_deleted: false,
//         time_of_day: TimeOfDayEnum.Morning,
//         effective_date: "2025-01-15",
//         expiration_date: "2025-08-15",
//         creator: {
//             employee_id: "emp2"
//         }
//     }
// ];

export const getAllWorkshiftApi = async (): Promise<WorkShift[]> => {
    try {
        const response = await api.get("/workshifts");
        return response.data.data
    } catch (error) {
        console.error("Error fetching all workshifts API:", error);
        throw error;
    }
};

export const createWorkshiftApi = async (payload: Omit<WorkShiftRequest, 'workshift_id'>): Promise<WorkShift> => {
    try {
        const response = await api.post("/workshifts", payload);
        return response.data.data
    } catch (error) {
        console.error("Error creating workshift API:", error);
        throw error;
    }
};

export const updateWorkshiftApi = async (id: string, payload: WorkShiftRequest): Promise<WorkShift> => {
    try {
        const response = await api.put(`/workshifts/${id}`, payload);
        return response.data.data
    } catch (error) {
        console.error("Error updating workshift API:", error);
        throw error;
    }
};

export const deleteWorkshiftApi = async (id: string): Promise<void> => {
    try {
        const response = await api.delete(`/workshifts/${id}`);
        return response.data.data
    } catch (error) {
        console.error("Error deleting workshift API:", error);
        throw error;
    }
};