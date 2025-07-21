import {useQuery} from "@tanstack/react-query";
import {getAllEmployeesApi} from "@/apis/profile.api.ts";

export const useAllEmployee = () =>
    useQuery({
        queryKey: ["employees"],
        queryFn: () => getAllEmployeesApi(1, 9999),
    });