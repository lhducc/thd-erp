import type { Decision, Position } from "@/types/index.ts";
import type { off } from "process";

type Office = {
  office_id: string;
  office_name: string;
  phone_number: string;
  address: string;
  latitude: number;
  longitude: number;
  created_date: string;
};

type Department = {
  department_id: string;
  department_name: string;
  manager: string;
  created_date: string;
  office_id: string;
  office: Office;
};

export type Employee = {
  employee_id: string;
  full_name: string;
  birthday: string;
  gender: string;
  work_type: string;
  phone_number: string;
  email: string;
  address: string;
  account_id: string | null;
  position_id: string;
  role_id: string;
  position: Position;
  job_title_id: string;
  status: "active" | "inactive";
  manager_id: string;
  role: {
    id: string;
    role_name: string;
  };
  schedule_id: number;
  department_id: string;
  created_date: string;
  department: Department;
  Decisions: Decision[];
};

export type EmployeeNameAndRole = {
  employee_id: string;
  full_name: string;
};

export interface ManagerPermission {
  employee_id: string;
  full_name: string;
  is_reading: boolean;
  is_editing: boolean;
  manager: {
    full_name: string;
  };
}

export interface AssignParams {
  id: string;
  payload: {
    work_schedule_id: number;
    managers: ManagerPermission[]; // Changed from manager_ids
    employee_ids: string[];
  };
}

export type ManagerEmployee = {
  employee_id: string;
  full_name: string;
  phone_number: string;
  email: string;
  position_name: string;
  job_title: string;
  department_name: string;
  office_name: string;
};