import {useQuery} from "@tanstack/react-query";
import {getAllWorkScheduleRegister, getWorkScheduleRegisterById} from "@/apis/work-schedule.api.ts";

export const useWorkScheduleRegister = () =>
    useQuery({
        queryKey: ["workScheduleRegister"],
        queryFn: getAllWorkScheduleRegister,
        staleTime: 0,
    });

export const useWorkScheduleRegisterById  = (id?: string) =>
    useQuery({
        queryKey: ["workScheduleRegisterById", id],
        queryFn: () => getWorkScheduleRegisterById(id),
        staleTime: 0,
        enabled: !!id
    })