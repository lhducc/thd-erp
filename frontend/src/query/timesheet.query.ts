import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {getAllTimesheets, getEmployeeTimesheetApi, getTimesheetById} from "@/apis/timesheet.api.ts";
import api from "@/apis/api.ts";
import {toast} from "sonner";

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

export const useUpdateTimesheetDetail = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ timesheetDetailId, data }: { timesheetDetailId: number; data: { adjusted_work_day: number } }) => {
            return await api.put(`/timesheet/${timesheetDetailId}`, data);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['timesheet'] });
        }
    });
};

export const useLockTimesheet = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ timesheetDetailId }: { timesheetDetailId: number}) => {
            return await api.put(`/timesheet-list/${timesheetDetailId}/locked`);
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: ['timesheet'] });
        }
    });
};

export const useCalculatorTimesheet = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ timesheetDetailId }: { timesheetDetailId: number}) => {
            return await api.get(`/timesheet/timesheet-list/${timesheetDetailId}`);
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: ['timesheet'] });
        }
    });
};

export const useResetTimesheet = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ timesheetDetailId }: { timesheetDetailId: number}) => {
            return await api.get(`/timesheet/${timesheetDetailId}/reset`);
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: ['timesheet'] });
        }
    });
};

export const useExportTimesheet = () => {
    return useMutation({
        mutationFn: async ({ timesheetDetailId }: { timesheetDetailId: number }) => {
            const res = await api.get(`/timesheet/export/${timesheetDetailId}`, {
                responseType: "arraybuffer",
            });

            return res.data;
        },
        onSuccess: (data) => {
            // Convert to Blob
            const blob = new Blob([data], { type: "application/zip" });
            const url = window.URL.createObjectURL(blob);

            // Create link and trigger download
            const a = document.createElement("a");
            a.href = url;
            a.download = "timesheet_" + Date.now().toLocaleString("vi-VN") + "_.xlsx";
            document.body.appendChild(a);
            a.click();

            // Cleanup
            a.remove();
            window.URL.revokeObjectURL(url);
        },
        onError: (error: any) => {
            toast.error(error.message ?? "Failed to export timesheet");
        },
    });
};
