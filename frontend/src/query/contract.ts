import {useQuery} from "@tanstack/react-query";
import {getContractByEmployeeId} from "@/apis/contract.api.ts";

export const useGetContractByEmployeeId = (employeeId: string) =>
    useQuery({
        queryKey: ["contract-employeeId", employeeId],
        queryFn: () => getContractByEmployeeId(employeeId),
    })