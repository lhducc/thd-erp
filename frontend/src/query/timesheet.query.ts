import {useQuery} from "@tanstack/react-query";
import {getAllTimesheets, getEmployeeTimesheetApi, getTimesheetById} from "@/apis/timesheet.api.ts";

export const useGetAllTimesheets = (params?: {
    page?: number;
    limit?: number;
    search?: string;
}) => useQuery({
    queryKey: ["timesheet-list", params],
    queryFn: () => getAllTimesheets(params),
});

export const useGetTimesheet = (id: string) =>
    useQuery({
        queryKey: ["timesheet_list_id", id],
        queryFn: () => getTimesheetById(id)
    })

export const useGetPersonalTimesheet = (month: number, year: number) =>
    useQuery({
        queryKey: ["personal-timesheet", month, year],
        queryFn: () => getEmployeeTimesheetApi(month, year)
    })