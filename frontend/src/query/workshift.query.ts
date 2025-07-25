import {useQuery} from "@tanstack/react-query";
import {getAllWorkshiftApi} from "@/apis/workshift.api.ts";
import {getAllSetupWorkshiftAPI} from "@/apis/setup-workshift.api.ts";

export const useQueryWorkshift = () =>
    useQuery({
        queryKey: ["workshift"],
        queryFn: getAllWorkshiftApi,
    })

export const useQueryAllSetupWorkshifts = () =>
    useQuery({
        queryKey: ["setupWorkshift"],
        queryFn: getAllSetupWorkshiftAPI,
        gcTime: 0,
        staleTime: 0,
    });