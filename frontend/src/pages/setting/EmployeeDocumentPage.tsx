import {useMutation, useQuery} from "@tanstack/react-query";
import {toast} from "sonner";
import type {ColumnDef} from "@tanstack/react-table";
import CreateOfficeForm from "@/components/CreateOfficeForm.tsx";
import {Button} from "@/components/ui/button.tsx";
import {SquarePen} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import {deleteDocumentTypesApi, getAllDocumentTypesApi} from "@/apis/document-type.api.ts";
import type {DocumentType} from "@/types/document-type.ts";
import CreateDocumentTypeForm from "@/components/CreateDocumentTypeForm.tsx";

const EmployeeDocumentPage = () => {
    const {
        data: types,
        isPending: pendingTypes,
        refetch: refetchTypes,
    } = useQuery({
        queryKey: ["documentType"],
        queryFn: getAllDocumentTypesApi,
        gcTime: 0,
        staleTime: 0,
    });

    const { mutateAsync: deleteDocumentTypes } = useMutation({
        mutationFn: (id: string) => deleteDocumentTypesApi(id),
        onSuccess: () => {
            toast.success("Xóa loại tài liệu thành công");
            refetchTypes();
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const columns: ColumnDef<DocumentType>[] = [
        {
            accessorKey: "id",
            header: "Mã",
        },
        {
            accessorKey: "document_type_name",
            header: "Tên loại tài liệu",
        },
        {
            accessorKey: "document_group",
            header: "Nhóm tài liệu",
        },
        {
            accessorKey: "description",
            header: "Mô tả",
        },
        {
            id: "actions",
            header: "Chỉnh sửa",
            cell: ({ row }) => {
                const type = row.original;

                return (
                    <div className="flex gap-4">
                        <CreateDocumentTypeForm
                            editBtn={
                                <Button variant="outline">
                                    <SquarePen />
                                </Button>
                            }
                            data={type}
                            type="edit"
                            refetch={refetchTypes}
                        />
                        <ConfirmDelete deleteFn={() => deleteDocumentTypes(type.id)} />
                    </div>
                );
            },
        },
    ];

    if (pendingTypes) {
        return <Loading />;
    }

    return (
        <>
            <DataTable
                columns={columns}
                data={types || []}
                title="Tài liệu nhân sự"
                buttonCreate={<CreateDocumentTypeForm refetch={refetchTypes} />}
                keyFilter="document_type_name"
            />
        </>
    );
};

export default EmployeeDocumentPage;
