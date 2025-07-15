import {type WorkShift, type WorkShiftRequest} from "@/types/Workshift.ts";
import api from "@/apis/api.ts";
import type {EmployeeWorkshift} from "@/types/employee-workshift.ts";

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

export const registerEmployeeWorkshiftApi = async (
    employee_id: string,
    work_shift_id: string,
    date: string // Expected format: "YYYY-MM-DD"
): Promise<EmployeeWorkshift> => {
    try {
        const response = await api.post("/employee-workshifts", {
            employee_id,
            work_shift_id,
            date: new Date(date).toISOString() // Convert to ISO string
        });
        return response.data.data;
    } catch (error) {
        console.error("Error registering employee workshift:", error);
        throw error;
    }
};

export const deleteEmployeeWorkshiftApi = async (id: number): Promise<void> => {
    try {
        await api.delete(`/employee-workshifts/${id}`);
    } catch (error) {
        console.error("Error deleting employee workshift:", error);
        throw error;
    }
};

export const getEmployeeWorkshiftsApi = async (employee_id: string): Promise<EmployeeWorkshift[]> => {
    try {
        const response = await api.get(`/employee-workshifts/${employee_id}`);
        return response.data.data;
    } catch (error) {
        console.error("Error fetching employee workshifts:", error);
        throw error;
    }
};

export const getAllEmployeeWorkshiftsApi = async (): Promise<EmployeeWorkshift[]> => {
    try {
        const response = await api.get("/employee-workshifts");
        return response.data.data;
    } catch (error) {
        console.error("Error fetching all employee workshifts:", error);
        throw error;
    }
};

export const updateEmployeeWorkshiftApi = async (
    employeeWorkshiftId: string,
    newWorkshiftId: string
): Promise<EmployeeWorkshift> => {
    try {
        console.log(employeeWorkshiftId, newWorkshiftId);
        const response = await api.put(`/employee-workshifts/${employeeWorkshiftId}`, {
            work_shift_id: newWorkshiftId
        });
        return response.data.data;
    } catch (error) {
        console.error("Error updating employee workshift:", error);
        throw error;
    }
};