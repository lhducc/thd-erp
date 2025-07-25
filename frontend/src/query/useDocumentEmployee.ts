import {useQuery} from "@tanstack/react-query";
import {getDocumentEmployeeById} from "@/apis/document-employee.api.ts";

export const useDocumentEmployeeById = (employeeId: string) =>
    useQuery({
        queryKey: ["documentEmployee", employeeId],
        queryFn: () => getDocumentEmployeeById(employeeId),
    })