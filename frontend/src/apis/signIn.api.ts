import api from "@/apis/api"
import type { PayloadSignIn } from "@/types";

export const signInApi = async (payload: PayloadSignIn) => {
  try {
    const response = await api.post("/auth/login", payload);
    return response.data;
  } catch (error: any) {
    console.error("Error delete sign in API:", error);
    throw new Error(error.response.data.message);
  }
}