import { JSX } from "react";

export type Office = {
  office_id: string;
  office_name: string;
  phone_number: string;
  address: string;
  latitude: number;
  longitude: number;
  created_date: Date;
};

export type PayloadOffice = {
  office_name: string;
  phone_number: string;
  address: string;
  latitude: number;
  longitude: number;
};

export type PayloadDepartment = {
  department_name: string;
  manager: string;
  office_id: string;
};

export type Department = {
  department_id: string;
  department_name: string;
  manager: string;
  created_date: string;
  office_id: string;
  office: null;
};

export type PayloadPosition = {
  position_name: string;
};

export type Position = {
  position_id: string;
  position_name: string;
  created_date: Date;
};

export type PayloadHierarchyLevel = {
  hierarchy_level: string;
  hierarchy_number: number;
};

export type HierarchyLevel = {
  id: string;
  hierarchy_level: string;
  hierarchy_number: number;
  created_date: Date;
};

export type PayloadSignIn = {
  email: string;
  password: string;
};

export type PayloadEmployee = {
  full_name: string;
  birthday: string;
  gender: string;
  phone_number: string;
  email: string;
  work_type: string;
  position_id: string;
  job_title_id: string;
  status: string;
  manager_id: string;
  address: string;
  department_id: string;
  start_date: string;
};

export interface EmployeeFormValues {
  full_name: string;
  birthday: string;
  gender: string;
  phone_number: string;
  email: string;
  work_type: string;
  position_id: string;
  job_title_id: string;
  status: string;
  manager_id: string;
  address: string;
  department_id: string;
}

export type Employee = {
  employee_id: string;
  full_name: string;
  birthday: string;
  gender: string;
  work_type: string;
  phone_number: string;
  email: string;
  account_id: number;
  current_address: string;
  position_id: string;
  job_title_id: string;
  status: string;
  manager_id: string;
  created_date: string;
};

export type HrDocument = {
  employee_id: string;
  full_name: string;
  birthday: string;
  gender: string;
  work_type: string;
  phone_number: string;
  email: string;
  account_id: number;
  position_id: string;
  job_title_id: string;
  status: string;
  manager_id: string;
  created_date: string;
};

export type PayloadDecision = {
  decision_name: string;
  effective_date: Date;
  sign_date: Date;
  content: string;
  condition: string;
  attached_file: File | null;
  employee_ids: string;
  decision_type_id: string;
};

export type Decision = {
  decision_id: string;
  decision_name: string;
  effective_date: string;
  sign_date: string;
  content: string;
  condition: string;
  attached_file: string;
  created_date: string;
  employee_id: string;
  decision_type_id: string;
  decision_type_name: string;
};

export type MenuItem = {
  name: string;
  url?: string;
  icon?: JSX.Element;
  child?: MenuItem[];
};
