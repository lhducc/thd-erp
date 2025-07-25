import {useQuery} from "@tanstack/react-query";
import {getAllEmployeesApi, getEmployeeByIdApi} from "@/apis/profile.api.ts";

export const useAllEmployee = () =>
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