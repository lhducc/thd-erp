import {useQuery} from "@tanstack/react-query";
import {getAllDepartmentsApi} from "@/apis/department.api.ts";

export const useDepartment = () =>
    useQuery({
        queryKey: ["departments"],
        queryFn: getAllDepartmentsApi,
        staleTime: 1000 * 60 * 5,
    });