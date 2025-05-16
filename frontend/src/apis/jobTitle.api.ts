import type { JobTitle, PayloadJobTitle } from "@/types";
import api from "./api";

export const createJobTitleApi = async (payload: PayloadJobTitle) => {
  try {
    const response = await api.post("/jobtitle", payload);
    return response.data;
  } catch (error: any) {
    console.log(error);
    throw new Error(error.response.data.message);
  }
};

export const getJobTitles = async (): Promise<JobTitle[]> => {
  try {
    const response = await api.get("/jobtitle");
    return response.data.data;
  } catch (error: any) {
    console.log(error);
    throw new Error(error.response.data.message);
  }
};

export const updateJobTitleApi = async (
  id: string,
  payload: PayloadJobTitle
) => {
  try {
    const response = await api.put(`/jobtitle/${id}`, payload);
    return response.data;
  } catch (error: any) {
    console.log(error);
    throw new Error(error.response.data.message);
  }
};

export const deleteJobTitleApi = async (id: string) => {
  try {
    const response = await api.delete(`/jobtitle/${id}`);
    return response.data;
  } catch (error: any) {
    console.log(error);
    throw new Error(error.response.data.message);
  }
};
