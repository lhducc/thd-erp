import api from "@/apis/api.ts";

export const getAllShiftAllocations = async (
    month: number,
    year: number,
    page: number = 1,
    pageSize: number = 100,
): Promise<any> => {
    const response = await api.get(`/shift-allocation?month=${month}&year=${year}&page=${page}&page_size=${pageSize}`);
    return response.data.data;
};