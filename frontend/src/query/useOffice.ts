import {useQuery} from "@tanstack/react-query";
import {getAllOfficesApi} from "@/apis/office.api.ts";

export const useOffice = () =>
    useQuery({
        queryKey: ["office"],
        queryFn: getAllOfficesApi,
    });
