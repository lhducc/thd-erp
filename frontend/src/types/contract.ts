export interface Contract {
    contract_id: string;
    effective_date: string; // ISO date string
    expired_date: string;   // ISO date string
    sign_date: string;      // ISO date string
    approve_status: "Chờ duyệt" | "Đã duyệt" | "Không duyệt";
    note: string;
    attached_file: string;  // URL string
    condition: string;
    created_date: string;   // ISO date string
    contract_type: string;
    employee: ContractEmployee;
    allowances: []
}

interface ContractEmployee {
    employee_id: string;
    full_name: string;
    department: Department;
}

interface Department {
    department_id: string;
    department_name: string;
    office: Office;
}

interface Office {
    office_id: string;
    office_name: string;
}

export interface ContractFormValues {
    contract_id: string;
    employee_id: string;
    contract_type: string;
    condition: string;
    sign_date: string;
    effective_date: string;
    expired_date: string;
    employee_name: string;
    department: string;
    status: "Chưa hiệu lực" | "Hiệu lực";
    note: string;
    approve_status: "Chờ duyệt" | "Đã duyệt" | "Không duyệt";
    allowance_ids: [];
}

export interface ContractType {
    contract_type_id: string;
    contract_type: string;
    contract_group: string;
    duration: number;
    unit: string;
    is_delete: boolean;
    working_form: string;
    created_date: string; // ISO 8601 date string
};

