import type { Office, PayloadOffice } from "@/types";
import api from "./api";

export const createOfficeApi = async (payload: PayloadOffice) => {
  try {
    const response = await api.post("/offices", payload);
    return response.data;
  } catch (error) {
    console.error("Error creating office API:", error);
    throw error;
  }
};

export const getAllOfficesApi = async (): Promise<Office[]> => {
  try {
    const response = await api.get("/offices");
    return response.data.data;
  } catch (error) {
    console.error("Error fetching all offices API:", error);
    throw error;
  }
};

export const updateOfficeApi = async (id: string, payload: PayloadOffice) => {
  try {
    const response = await api.put(`/offices/${id}`, payload);
    return response.data;
  } catch (error) {
    console.error("Error updating office API:", error);
    throw error;
  }
};

export const deleteOffice = async (id: string) => {
  try {
    const response = await api.delete("")
  } catch (error) {
    console.log(error);
    throw error;
  }
};

import type { Office, PayloadOffice } from "@/types";
import api from "./api";

export const createOfficeApi = async (payload: PayloadOffice) => {
  try {
    const response = await api.post("/office", payload);
    return response.data;
  } catch (error: any) {
    console.error("Error creating office API:", error);
    throw new Error(error.response.data.message);
  }
};

export const getAllOfficesApi = async (): Promise<Office[]> => {
  try {
    const response = await api.get("/office");
    return response.data.data;
  } catch (error) {
    console.error("Error fetching all offices API:", error);
    throw error;
  }
};

export const updateOfficeApi = async (id: string, payload: PayloadOffice) => {
  try {
    const response = await api.put(`/office/${id}`, payload);
    return response.data;
  } catch (error: any) {
    console.error("Error update office API:", error);
    throw new Error(error.response.data.message);
  }
};

export const deleteOfficeApi = async (id: string) => {
  try {
    const response = await api.delete(`/office/${id}`);
    return response.data;
  } catch (error: any) {
    console.error("Error delete office API:", error);
    throw new Error(error.response.data.message);
  }
};
