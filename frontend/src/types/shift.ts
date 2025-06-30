export type ShiftStatus = "Đang áp dụng" | "Hết áp dụng" | "Chưa áp dụng";

export interface Shift {
  shift_id: string;
  shift_name: string;
  effective_date: string; 
  target: string; 
  office: string; 
  status: ShiftStatus;
  created_at?: string;
  updated_at?: string;
}