import type {JobTitle} from "@/types/job-title.ts";

export type Standard = {
    id: number;
    name: string;
    effective_date: Date;
    amount: number;
    time_work_type: "Theo ca" | "Theo ngày";
    flex_work_type: boolean;
    job_title: JobTitle;
}

export type PayloadStandard = {
    name: string;
    effective_date: Date;
    amount: number;
    time_work_type: "Theo ca" | "Theo ngày";
    flex_work_type: boolean;
    job_title: string;
}