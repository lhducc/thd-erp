import {useQuery} from "@tanstack/react-query";
import {getWorkSchedule, getWorkScheduleById} from "@/apis/work-schedule.api.ts";

export const useWorkSchedule = () =>
    useQuery({
        queryKey: ["work-schedule"],
        queryFn: getWorkSchedule,
        gcTime: 0,
        staleTime: 0,
    });

export const useWorkScheduleById = (id?: string) =>
    useQuery({
        queryKey: ["work-schedule", id],
        queryFn: () => getWorkScheduleById(id!),
        enabled: !!id,
    });