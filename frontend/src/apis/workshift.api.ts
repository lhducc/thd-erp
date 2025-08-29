import {type Workshift, type WorkShiftRequest} from "@/types/workshift.ts";
import api from "@/apis/api.ts";
import type {EmployeeWorkshift} from "@/types/employee-workshift.ts";

export const getAllWorkshiftApi = async (): Promise<Workshift[]> => {
    try {
        const response = await api.get("/workshifts");
        return response.data.data
    } catch (error) {
        console.error("Error fetching all workshifts API:", error);
        throw error;
    }
};

export const createWorkshiftApi = async (payload: Omit<WorkShiftRequest, 'workshift_id'>): Promise<Workshift> => {
    const response = await api.post("/workshifts", payload);
    return response.data.data
};

export const updateWorkshiftApi = async (id: string, payload: WorkShiftRequest): Promise<Workshift> => {
    const response = await api.put(`/workshifts/${id}`, payload);
    return response.data.data
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
    workshift_id: string,
    date: string // Expected format: "YYYY-MM-DD"
): Promise<EmployeeWorkshift> => {
    try {
        const response = await api.post("/employee-workshifts", {
            employee_id,
            workshift_id,
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

export const getEmployeeWorkshiftsApi = async (month: number, year: number): Promise<EmployeeWorkshift[]> => {
    try {
        const response = await api.get(`/employee-workshifts/personal?month=${month}&year=${year}`);
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
        const response = await api.put(`/employee-workshifts/${employeeWorkshiftId}`, {
            work_shift_id: newWorkshiftId
        });
        return response.data.data;
    } catch (error) {
        console.error("Error updating employee workshift:", error);
        throw error;
    }
};
