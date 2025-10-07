import type { Employee, Office } from "@/types/index.ts";

export type AttendanceSetting = {
  attendance_category_id: string;
  attendance_category_name: string;
  office_id: string;
  Office: Office;
  is_camera: boolean;
  is_gps: boolean;
  is_check_location: boolean;
  scope?: number;
  status: string;
};

type CreatedByInfo = {
  employee_id: string;
  full_name: string;
};

export type AttendanceRecord = {
  attendance_record_id: string;
  employee_id: string;
  timestamp: string;
  gps_location: string;
  latitude: number | null;
  longitude: number | null;
  image_name: string;
  image_URL: string;
  office_id: string | null;
  category_id: string | null;
  note_request: string | null;
  note_reject: string | null;
  status: "approved" | "rejected" | "pending";
  created_by: string;
  Employee: Employee;
  Office: Office | null;
  create_by_info: CreatedByInfo;
  AttendanceCategory: AttendanceSetting; // Replace `any` with proper type if available
};

export type CreateManualRecord = {
  employee_id: string;
  timestamp: string;
};

export type AttendanceRecordHistoryByDate = {
  employee_id: string;
  full_name: string;
  office_name: string;
  department_name: string;
  timestamp: string;
  workshift_id: string;
  workshift_name: string;
  start_time: string;
  checkin_to: string;
};
