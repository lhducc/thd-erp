import type {PayloadEmployee} from "@/types";
import api from "./api";
import type {Employee, EmployeeNameAndRole} from "@/types/employee.ts";
import {unwrap} from "@/lib/utils.ts";

export const createEmployeeApi = async (payload: PayloadEmployee) => {
    const response = await api.post("/employee", payload);
    return response.data;
};

export const getAllEmployeesApi = async (page: number = 2, pageSize: number = 5) => {
    const response = await api.get(`/employee?page=${page}&pageSize=${pageSize}`);
    if (response.data && response.data.data) {
        const list = response.data.data.data
        const result: Employee[] = [];
        list.forEach((employee) => {
            const role = employee.role;
            const Employee: Employee = employee.employee;
            Employee.role = role;
            delete employee.role;
            result.push(unwrap(employee));
        })
        return result;
    } else {
        throw new Error("No data found in the response.");
    }
};

export const getEmployeeByIdApi = async (id: string): Promise<Employee> => {
    const response = await api.get(`/employee/${id}`);
    return response.data.data.employee;
};

export const getEmployeePersonalApi = async (): Promise<Employee> => {
    const response = await api.get(`/employee/personal`);
    return response.data.data.employee;
};

export const updateEmployeeApi = async (employeeId: string, payload: PayloadEmployee) => {
    const response = await api.put(`/employee/${employeeId}`, payload);
    return response.data;
};

export const deleteEmployeeApi = async (employeeId: string) => {
    const response = await api.delete(`/employee/${employeeId}`);
    return response.data;
};

export const exportEmployeeExcelApi = async (): Promise<File> => {
    const response = await api.get("employee/export?fields=employee_id,fullname,email,manager,birthday", {
        responseType: "blob",
    });
    const fileName = "employee.xlsx";
    const file = new File([response.data], fileName, {type: response.data.type});
    return file;
};

export const getEmployeeByRoleNameApi = async (roleName: string): Promise<EmployeeNameAndRole[]> => {
    const response = await api.get(`/employee/user?roleID=${roleName}`);
    return response.data.data;
}

export const changeStatusEmployeeApi = async (
    id: string,
    status: "active" | "inactive"
): Promise<string> => {
    const response = await api.put(`/employee/${id}/status?status=${status}`);
    if (response.status === 200) {
        return status;
    }
    throw new Error("Không thể thay đổi trạng thái");
};

