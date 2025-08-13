import {useQuery} from "@tanstack/react-query";
import {getAllContractsApi} from "@/apis/contract.api.ts";

export const contractQuery = () =>
    useQuery({
        queryKey: ["contract"],
        queryFn: getAllContractsApi,
        gcTime: 0,
        staleTime: 0,
    });