import {queryOptions, useQuery} from "@tanstack/react-query";
import {getAllWorkScheduleRegister, getWorkScheduleRegisterById} from "@/apis/work-schedule.api.ts";

export const useWorkScheduleRegister = () =>
    useQuery({
        queryKey: ["workScheduleRegister"],
        queryFn: getAllWorkScheduleRegister,
        gcTime: Infinity,
        staleTime: 0,
    });

export const useWorkScheduleRegisterById  = (id?: string) =>
    queryOptions({
        queryKey: ["workScheduleRegisterById", id],
        queryFn: () => getWorkScheduleRegisterById(id),
        enabled: !!id
    })