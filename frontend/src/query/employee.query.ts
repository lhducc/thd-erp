import { useQuery } from "@tanstack/react-query";
import {
  getAllEmployeesActiveApi,
  getAllEmployeesApi,
  getEmployeeByIdApi,
  getEmployeeByRoleNameApi,
  getEmployeesByManagerApi,
} from "@/apis/profile.api.ts";

export const useGetAllEmployee = () =>
  useQuery({
    queryKey: ["employees"],
    queryFn: () => getAllEmployeesApi(1, 9999),
  });

export const useGetEmployeeById = (employeeId: string) =>
  useQuery({
    queryKey: ["employeeId", employeeId],
    queryFn: () => getEmployeeByIdApi(employeeId),
    enabled: !!employeeId,
  });

export const useEmployeeByRoleNameQuery = (roleNames: string | string[]) =>
  useQuery({
    queryKey: ["employeeByRole", roleNames],
    queryFn: async () => {
      const roles = Array.isArray(roleNames) ? roleNames : [roleNames];
      const results = await Promise.all(
        roles.map((role) => getEmployeeByRoleNameApi(role))
      );
      return results.flat();
    },
  });

export const useGetManagerEmployees = () =>
  useQuery({
    queryKey: ["managerEmployees"],
    queryFn: () => getEmployeesByManagerApi(),
  });

export const useGetAllEmployeesActive = () =>
  useQuery({
    queryKey: ["employeesActive"],
    queryFn: () => getAllEmployeesActiveApi(1, 9999),
  });
