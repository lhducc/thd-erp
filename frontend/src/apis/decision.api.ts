import type { PayloadDecision, Decision } from "@/types";
import api from "./api";

export const createDecisionApi = async (payload: FormData) => {
  try {
      console.log(payload);
    const response = await api.post("/decision", payload);
    return response.data.data;
  } catch (error: any) {
    console.error("Error creating Decision API:", error);
    handleApiError(error);
  }
};

export const getAllDecisionsApi = async (): Promise<Decision[]> => {
  try {
    const response = await api.get(`/decision`);
    console.log(response);
    return response.data?.data?.data || [];
  } catch (error: any) {
    console.error("Error fetching all Decisions API:", error);
    handleApiError(error);
  }
};

export const getDecisionByIdApi = async (DecisionId: string): Promise<Decision> => {
  try {
    const response = await api.get(`/decision/${DecisionId}`);
    return response.data;
  } catch (error: any) {
    console.error("Error fetching Decision by ID API:", error);
    handleApiError(error);
  }
};

export const updateDecisionApi = async (DecisionId: string, payload: PayloadDecision): Promise<Decision> => {
  try {
    // Make the PUT request using the axios instance (api)
    const response = await api.put(`/decision/${DecisionId}`, payload);
    return response.data;
  } catch (error: any) {
    console.error("Error updating Decision API:", error);
    handleApiError(error);  // Handle errors as needed
    throw error;  // Propagate the error
  }
};

export const deleteDecisionApi = async (DecisionId: string): Promise<void> => {
  try {
    const response = await api.delete(`/decision/${DecisionId}`);
    return response.data;
  } catch (error: any) {
    console.error("Error deleting Decision API:", error);
    handleApiError(error);
  }
};

export const exportDecisionExcelApi = async (): Promise<File> => {
  try {
    const response = await api.get("/decisions/export", {
        responseType: "blob",
    });

    const fileName = "contracts.xlsx";
    const file = new File([response.data], fileName, { type: response.data.type });
    return file;
  } catch (error: any) {
    console.error("Error exporting Decisions to Excel:", error);
    handleApiError(error);
  }
};

const handleApiError = (error: any) => {
  if (error.response) {
    console.error("Response error:", error.response.data);
    throw new Error(error.response?.data?.message || "Unknown error occurred");
  } else {
    console.error("Request error:", error.message);
    throw new Error(error.message || "Unknown error occurred");
  }
};

export const createSampleDecisionApi = async (
  decisionData: {
    decision_name: string;
    effective_date: string;
    sign_date: string;
    content: string;
    condition: string;
    attached_file: File | null;
    employee_id: string;
    decision_type_id: string;
  }
): Promise<Decision> => {
  const payload: PayloadDecision = {
    decision_name: decisionData.decision_name,
    effective_date: new Date(decisionData.effective_date).toISOString(),
    sign_date: new Date(decisionData.sign_date).toISOString(),
    content: decisionData.content,
    condition: decisionData.condition,
    attached_file: decisionData.attached_file || null, 
    employee_id: decisionData.employee_id,
    decision_type_id: decisionData.decision_type_id,
  };

  try {
    const response = await api.post("/decision", payload);
    return response.data;
  } catch (error: any) {
    console.error("Error creating sample Decision API:", error);
    handleApiError(error);
    throw error; // optional: rethrow the error if you want to handle it elsewhere
  }
};
