import type { PayloadEmployee } from "@/types";
import api from "./api";
import type {Employee} from "@/types/employee.ts";

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

export const getAllEmployeesApi = async (page: number = 1, pageSize: number = 5) : Promise<Employee[]> => {
  try {
    const response = await api.get(`/employee?page=${page}&pageSize=${pageSize}`);
    if (response.data && response.data.data) {
      return response.data.data;
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

// API lấy thông tin nhân viên theo ID
export const getEmployeeByIdApi = async (employeeId: string): Promise<Employee> => {
  try {
    const response = await api.get(`/employee/${employeeId}`);
    return response.data;
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

// API cập nhật nhân viên
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

// API xóa nhân viên
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

// API xuất dữ liệu nhân viên sang Excel
export const exportEmployeeExcelApi = async () => {
  try {
    const response = await api.get("/employee/export", { responseType: "blob" });
    const file = new Blob([response.data], { type: "application/vnd.ms-excel" });
    const fileURL = URL.createObjectURL(file);
    const link = document.createElement("a");
    link.href = fileURL;
    link.download = "employees.xlsx";
    link.click();
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
