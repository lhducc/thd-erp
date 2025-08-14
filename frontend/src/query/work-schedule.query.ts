import {useQuery} from "@tanstack/react-query";
import {getWorkScheduleRegisterById} from "@/apis/work-schedule.api.ts";

export const useGetWorkScheduleRegisterById = (id : string) =>
    useQuery({
        queryKey: ["work-schedule-id", id],
        queryFn: () => getWorkScheduleRegisterById(id)
    })