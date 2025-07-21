import {useQuery} from "@tanstack/react-query";
import {getAllContractsTypeApi} from "@/apis/contract-type.api.ts";

export const useContractType = () =>
    useQuery({
        queryKey: ["contract_type"],
        queryFn: getAllContractsTypeApi,
        gcTime: 0,
        staleTime: 0,
    });