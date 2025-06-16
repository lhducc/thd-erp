export type Employee = {
    employee_id: string;
    full_name: string;
    birthday: string; // ISO 8601 date string
    gender: string;
    work_type: string;
    phone_number: string;
    email: string;
    account_id: number;
    address: string;
    position_id: Position;
    job_title_id: JobTitle;
    department: Department;
    status: string;
    manager: Manager;
    created_date: string; // ISO 8601 date string with time
};

type Position = {
    position_id: string;
    position_name: string;
};

type JobTitle = {
    job_title_id: string;
    job_title: string;
};

type Department = {
    department_id: string;
    department_name: string;
    office: Office;
};

type Office = {
    office_id: string;
    office_name: string;
};

type Manager = {
    employee_id: string;
    full_name: string;
};
