import {useQuery} from "@tanstack/react-query";
import {getRolesApi} from "@/apis/role.api.ts";

export const useGetRoles = () =>
    useQuery({
        queryKey: ["roles"],
        queryFn: getRolesApi,
    });