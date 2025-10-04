import {useQuery} from "@tanstack/react-query";
import {getAllEmployeesApi, getEmployeeByIdApi, getEmployeeByRoleNameApi, getEmployeesByManagerApi} from "@/apis/profile.api.ts";

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

export const useEmployeeByRoleNameQuery = (roleName: string) =>
    useQuery({
        queryKey: ["employeeByRole", roleName],
        queryFn: () => getEmployeeByRoleNameApi(roleName),
    });

export const useGetManagerEmployees = () =>
  useQuery({
    queryKey: ["managerEmployees"],
    queryFn: () => getEmployeesByManagerApi(),
  });

