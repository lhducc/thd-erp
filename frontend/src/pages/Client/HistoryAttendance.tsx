import {useGetAttendanceRecordPersonal} from "@/query/attendance-record.query.ts";
import type {ColumnDef} from "@tanstack/react-table";
import type {AttendanceRecord} from "@/types/attendance.ts";
import {Button} from "@/components/ui/button.tsx";
import {Eye} from "lucide-react";
import {formatDate, getTimeFromTimestamp} from "@/lib/utils.ts";
import DataTable from "@/components/DataTable.tsx";
import {useState} from "react";
import ImageDialog from "@/pages/Client/attendance-report/components/ImageDialog.tsx";

const HistoryAttendance = () => {
    // TODO: backend không pagination api này
    const [pageIndex, setPageIndex] = useState(0);
    const [pageSize, setPageSize] = useState(9999);
    const [selectedImage, setSelectedImage] = useState<string | null>(null);
    const [dialogOpen, setDialogOpen] = useState(false);

    // Sử dụng pageIndex và pageSize trong query
    const {data, isLoading} = useGetAttendanceRecordPersonal(
        pageIndex + 1, // Nếu API của bạn bắt đầu từ page 1
        pageSize
    );

    // Hàm mở dialog với hình ảnh
    const handleOpenImage = (imageUrl: string) => {
        setSelectedImage(imageUrl);
        setDialogOpen(true);
    };

    const columns: ColumnDef<AttendanceRecord>[] = [
        {
            accessorKey: "timestamp",
            header: "Ngày",
            cell: ({row}) => {
                return <p>{formatDate(row.original.timestamp)}</p>
            }
        },
        {
            header: "Giờ",
            cell: ({row}) => {
                return <p>{getTimeFromTimestamp(row.original.timestamp)}</p>
            }
        },
        {
            accessorKey: "AttendanceCategory.attendance_category_name",
            header: "Hình thức chấm công",
        },
        {
            accessorKey: "is_gps",
            header: "Vị trí",
            cell: ({row}) => (
                <p>{row.original.longitude} - {row.original.latitude}</p>
            ),
        },
        {
            accessorKey: "image_URL",
            header: "Xem hình ảnh chấm công",
            cell: ({row}) => (
                row.original.image_URL ? (
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleOpenImage(row.original.image_URL)}
                    >
                        <Eye className="h-4 w-4 mr-2" />
                        Xem hình
                    </Button>
                ) : (
                    <p>Không có hình ảnh</p>
                )
            ),
        },
        {
            accessorKey: "create_by_info",
            header: "Người tạo",
        },
        {
            accessorKey: "status",
            header: "Trạng thái",
            cell: ({row}) => {
                let status : string = row.original.status;
                switch (status) {
                    case "approved":
                        status = "Đã phê duyệt"
                        break;
                    case "rejected":
                        status = "Từ chối"
                        break;
                    case "pending":
                        status = "Chờ phê duyệt"
                        break;
                }
                return (
                    <p>{status}</p>
                )
            },
        },
    ];

    return (
        <div>
            <DataTable
                title="Lịch sử chấm công cá nhân"
                columns={columns}
                data={data || []}
                isLoading={isLoading}
                keyFilter="timestamp"
            />

            {/* Dialog hiển thị hình ảnh */}
            <ImageDialog
                open={dialogOpen}
                onOpenChange={setDialogOpen}
                imageUrl={selectedImage || ""}
            />
        </div>
    );
};

export default HistoryAttendance;