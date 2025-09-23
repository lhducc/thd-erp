import { useQuery } from "@tanstack/react-query";

import type { AttendanceRecordHistoryByDate } from "@/types/attendance";
import {getEmployeesByDateApi, getEmployeesByMonthApi } from "@/apis/attendance-record.api";



export const useEmployeesByDate = (date: string | null) => {
  return useQuery<AttendanceRecordHistoryByDate[]>({
    queryKey: ["employeesByDate", date],
    queryFn: () => getEmployeesByDateApi(date!),
    enabled: !!date, 
  });
};

export const useEmployeesByMonth = (month: string | null) => {
  return useQuery<AttendanceRecordHistoryByDate[]>({
    queryKey: ["employeesByMonth", month],
    queryFn: () => getEmployeesByMonthApi(month!),
    enabled: !!month,
  });
};