import type { HierarchyLevel, PayloadHierarchyLevel } from "@/types";
import api from "./api";

export const createHierarchyLevelApi = async (
  payload: PayloadHierarchyLevel
) => {
  try {
    const response = await api.post("/hierarchylevel", payload);
    return response.data.data;
  } catch (error: any) {
    console.error("Error creating hierarchy level API:", error);
    throw new Error(error.response.data.message);
  }
};

export const getAllHierarchyLevelApi = async (): Promise<HierarchyLevel[]> => {
  try {
    const response = await api.get("/hierarchylevel");
    return response.data.data;
  } catch (error: any) {
    console.error("Error fetching all hierarchy levels API:", error);
    throw new Error(error.response.data.message);
  }
};

export const updateHierarchyLevelApi = async (
  id: string,
  payload: PayloadHierarchyLevel
) => {
  try {
    const response = await api.put(`/hierarchylevel/${id}`, payload);
    return response.data.data;
  } catch (error: any) {
    console.error("Error updating hierarchy level API:", error);
    throw new Error(error.response.data.message);
  }
};

export const deleteHierarchyLevelApi = async (id: string) => {
  try {
    const response = await api.delete(`/hierarchylevel/${id}`);
    return response.data.data;
  } catch (error: any) {
    console.error("Error deleting hierarchy level API:", error);
    throw new Error(error.response.data.message);
  }
};
