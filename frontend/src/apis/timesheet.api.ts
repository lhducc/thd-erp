import api from "@/apis/api.ts";
import type {
  TimesheetInfor,
  TimesheetList,
  TimesheetResponse,
} from "@/types/timesheet.ts";
import type { CreateTimesheet } from "@/types/timesheet.ts";

export const getAllTimesheets = async (params?: {
  page?: number;
  limit?: number;
  search?: string;
}): Promise<TimesheetResponse> => {
  const response = await api.get("/timesheet-list", { params });
  return response.data;
};

export const getTimesheetById = async (
  id: string
): Promise<TimesheetInfor[]> => {
  const response = await api.get(`/timesheet-list/${id}`);
  return response.data.data;
};

export const createTimesheetApi = async (payload: CreateTimesheet) => {
  const response = await api.post(`/timesheet-list`, payload);
  return response.data.data.data;
};

export const getEmployeeTimesheetApi = async (
  month: number,
  year: number
): Promise<TimesheetList[]> => {
  const response = await api.get(
    `/timesheet/personal?month=${month}&year=${year}`
  );
  if (response.data.data && response.data.data.length > 0) {
    const result = await api.get(
      `/timesheet/personal?month=${month}&year=${year}`
    );
  }
  return response.data.data;
};

export const exportTimesheetCheckinCheckoutApi = async (
  id: string
): Promise<Blob> => {
  const response = await api.get(`/timesheet-list/${id}/export`, {
    responseType: "blob",
  });
  return response.data;
};
