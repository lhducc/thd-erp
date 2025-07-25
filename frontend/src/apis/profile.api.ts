import type { PayloadEmployee, Employee } from "@/types";
import api from "./api";

export const createEmployeeApi = async (payload: PayloadEmployee) => {
  try {
    const response = await api.post("/employee", payload);
    console.log("Employee created:", response.data); 
    return response.data;
  } catch (error: any) {
    console.error("Error creating Employee API:", error);
    if (error.response) {
      console.error("Response error:", error.response.data); 
      throw new Error(error.response?.data?.message || "Unknown error occurred");
    } else {
      console.error("Request error:", error.message); 
      throw new Error(error.message || "Unknown error occurred");
    }
  }
};

export const getAllEmployeesApi = async (page: number = 2, pageSize: number = 5) => {
  try {
    const response = await api.get(`/employee?page=${page}&pageSize=${pageSize}`);
    if (response.data && response.data.data) {
      return response.data.data.data;
    } else {
      throw new Error("No data found in the response.");
    }
  } catch (error: any) {
    console.error("Error fetching all Employees API:", error);
    if (error.response) {
      console.error("Response error:", error.response.data);
      throw new Error(error.response?.data?.message || "An unknown error occurred while fetching employees.");
    } else {
      console.error("Request error:", error.message);
      throw new Error(error.message || "An unknown error occurred while fetching employees.");
    }
  }
};

export const getEmployeeByIdApi = async (employeeId: string): Promise<any> => {
  try {
    const response = await api.get(`/employee/${employeeId}`);
    return response.data.data;
  } catch (error: any) {
    console.error("Error fetching Employee by ID API:", error);
    if (error.response) {
      console.error("Response error:", error.response.data);
      throw new Error(error.response?.data?.message || "Unknown error occurred");
    } else {
      console.error("Request error:", error.message);
      throw new Error(error.message || "Unknown error occurred");
    }
  }
};

export const updateEmployeeApi = async (employeeId: string, payload: PayloadEmployee) => {
  try {
    const response = await api.put(`/employee/${employeeId}`, payload);
    return response.data;
  } catch (error: any) {
    console.error("Error updating Employee API:", error);
    if (error.response) {
      console.error("Response error:", error.response.data);
      throw new Error(error.response?.data?.message || "Unknown error occurred");
    } else {
      console.error("Request error:", error.message);
      throw new Error(error.message || "Unknown error occurred");
    }
  }
};

export const deleteEmployeeApi = async (employeeId: string) => {
  try {
    const response = await api.delete(`/employee/${employeeId}`);
    return response.data;
  } catch (error: any) {
    console.error("Error deleting Employee API:", error);
    if (error.response) {
      console.error("Response error:", error.response.data);
      throw new Error(error.response?.data?.message || "Unknown error occurred");
    } else {
      console.error("Request error:", error.message);
      throw new Error(error.message || "Unknown error occurred");
    }
  }
};

export const exportEmployeeExcelApi = async (): Promise<File> => {
  try {
    const response = await api.get("employee/export?fields=employee_id,fullname,email,manager,birthday", {
        responseType: "blob",
    });

    const fileName = "employee.xlsx";
    const file = new File([response.data], fileName, { type: response.data.type });
    return file;
  } catch (error: any) {
    console.error("Error exporting Employees to Excel:", error);
    if (error.response) {
      console.error("Response error:", error.response.data);
      throw new Error(error.response?.data?.message || "Unknown error occurred");
    } else {
      console.error("Request error:", error.message);
      throw new Error(error.message || "Unknown error occurred");
    }
  }
};
