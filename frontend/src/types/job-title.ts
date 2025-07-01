export type JobTitle = {
    job_title_id: string;
    job_title: string;
    created_date: Date;
    hierarchy_level_id: string;
};

export type PayloadJobTitle = {
    job_title: string;
    hierarchy_level_id: string;
};