import type { PayloadPosition, Position } from "@/types";
import api from "./api";

export const createPositionApi = async (payload: PayloadPosition) => {
  try {
    const response = await api.post("/position", payload);
    return response.data.data;
  } catch (error: any) {
    console.error("Error creating position API:", error);
    throw new Error(error.response.data.message);
  }
};

export const getAllPositionApi = async (): Promise<Position[]> => {
  try {
    const response = await api.get("/position");
    return response.data.data.data;
  } catch (error: any) {
    console.error("Error fetching all positions API:", error);
    throw new Error(error.response.data.message);
  }
};

export const updatePositionApi = async (
  id: string,
  payload: PayloadPosition
) => {
  try {
    const response = await api.put(`/position/${id}`, payload);
    return response.data.data;
  } catch (error: any) {
    console.error("Error updating position API:", error);
    throw new Error(error.response.data.message);
  }
};

export const deletePositionApi = async (id: string) => {
  try {
    const response = await api.delete(`/position/${id}`);
    return response.data.data;
  } catch (error: any) {
    console.error("Error deleting position API:", error);
    throw new Error(error.response.data.message);
  }
};
