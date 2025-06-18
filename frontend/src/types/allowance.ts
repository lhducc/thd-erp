export type Allowance = {
    id: string;
    allowance_name: string;
    tax: boolean;
    amount: number;
    unit: string;
    is_deleted: boolean;
    created_date: string; // ISO 8601 date string
};
