export type DocumentType = {
    id: string;
    document_group: string;
    document_type_name: string;
    description: string;
    created_date: string;
};

export type DocumentTypeCreate = {
    document_group: string;
    document_type_name: string;
    description?: string;
    created_date?: string;
};