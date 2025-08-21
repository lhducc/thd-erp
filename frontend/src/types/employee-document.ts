import type {Employee} from "@/types/employee.ts";

export interface DocumentType {
    id: string;
    document_group: string;
    document_type_name: string;
    description: string;
    created_date: string; // ISO date string
}

export interface Document {
    document_id: string;
    document_type_id: string;
    employee_id: string;
    effective_date: string; // ISO date string
    expired_date: string;   // ISO date string
    note: string;
    condition: string;
    status: string;
    attached_file: string | null;
    created_date: string; // ISO date string
    document_type: DocumentType;
    employee: Employee;
}