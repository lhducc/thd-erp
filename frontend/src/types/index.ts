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
  created_date: Date;
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

export type PayloadJobTitle = {
  job_title: string;
  hierarchy_level_id: string;
};

export type JobTitle = {
  job_title_id: string;
  job_title: string;
  created_date: Date;
  hierarchy_level_id: string;
};
