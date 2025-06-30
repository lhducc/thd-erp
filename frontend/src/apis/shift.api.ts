import api from "@/apis/api.ts";
import type { Shift } from "@/types/shift";

export const getAllShiftsApi = async (): Promise<Shift[]> => {
  const res = await api.get("/shifts");
  return res.data;
};

export const deleteShiftApi = async (id: string): Promise<void> => {
  await api.delete(`/shifts/${id}`);
};

export const createShiftApi = async (data: Partial<Shift>): Promise<Shift> => {
  const res = await api.post("/shifts", data);
  return res.data;
};

export const exportShiftExcelApi = async (): Promise<Blob & { name: string }> => {
  const res = await api.get("/shifts/export", {
    responseType: "blob",
  });
  const blob = new Blob([res.data], { type: res.headers["content-type"] });
  (blob as any).name = "phanca.xlsx";
  return blob as Blob & { name: string };
};
