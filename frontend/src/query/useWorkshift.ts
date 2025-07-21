import {useQuery} from "@tanstack/react-query";
import {getAllWorkshiftApi} from "@/apis/workshift.api.ts";

export const useWorkshift = () =>
    useQuery({
        queryKey: ["workshift"],
        queryFn: getAllWorkshiftApi,
    })