import api from "@/apis/api.ts";
import type {Role} from "@/types/role.ts";

export const getRolesApi = async (): Promise<Role[]> => {
    const response = await api.get(`/role`);
    return response.data.data;
}