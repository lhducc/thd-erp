import type { Department, PayloadDepartment } from "@/types";
import api from "./api";

export const createDepartmentApi = async (payload: PayloadDepartment) => {
  try {
    const response = await api.post("/department", payload);
    return response.data;
  } catch (error: any) {
    console.error("Error delete office API:", error);
    throw new Error(error.response.data.message);
  }
};

export const getAllDepartmentsApi = async (): Promise<Department[]> => {
  try {
    const response = await api.get("/department");
    return response.data.data;
  } catch (error: any) {
    console.error("Error fetching all departments API:", error);
    throw new Error(error.response.data.message);
  }
};

export const updateDepartmentApi = async (
  id: string,
  payload: PayloadDepartment
) => {
  try {
    const response = await api.put(`/department/${id}`, payload);
    return response.data;
  } catch (error: any) {
    console.error("Error update department API:", error);
    throw new Error(error.response.data.message);
  }
};

export const deleteDepartmentApi = async (id: string) => {
  try {
    const response = await api.delete(`/department/${id}`);
    return response.data;
  } catch (error: any) {
    console.error("Error delete department API:", error);
    throw new Error(error.response.data.message);
  }
};
