import {useQuery} from "@tanstack/react-query";
import {getAllDepartmentsApi, getDepartmentByOfficeIdApi} from "@/apis/department.api.ts";
import type {Department} from "@/types";

export const useDepartment = () =>
    useQuery({
        queryKey: ["departments"],
        queryFn: getAllDepartmentsApi,
        staleTime: 1000 * 60 * 5,
    });

export const useGetDepartmentByOfficeId = (officeId: string) =>
    useQuery<Department[]>({
        queryKey: ["departments-office-id", officeId],
        queryFn: () => getDepartmentByOfficeIdApi(officeId),
        staleTime: 0,
        enabled: !!officeId,
    });