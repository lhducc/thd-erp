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
    job_title_id: string;
    status: 'active' | 'inactive'; // có thể mở rộng nếu có thêm trạng thái
    manager_id: string;
    department_id: string;
    created_date: string;
    department: Department;
    Decisions: any; // có thể thay `any` bằng kiểu cụ thể nếu bạn có cấu trúc cho `Decisions`
};
