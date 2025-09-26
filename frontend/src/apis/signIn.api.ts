import api from "@/apis/api"
import type { PayloadSignIn } from "@/types";
import type { ChangePasswordPayload } from "@/types/changepassword";

export const signInApi = async (payload: PayloadSignIn) => {
  try {
    const response = await api.post("/auth/login", payload);
    return response.data;
  } catch (error: any) {
    console.error("Error delete sign in API:", error);
    throw new Error(error.response.data.message);
  }
}

export const changePasswordApi = async (payload: ChangePasswordPayload) => {
  try {
    const token = localStorage.getItem("access_token");
    
    const response = await api.put("/auth/password", payload, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    return response.data;
  } catch (error: any) {
    console.error("Error change password API:", error);
    throw error; 
  }
};
