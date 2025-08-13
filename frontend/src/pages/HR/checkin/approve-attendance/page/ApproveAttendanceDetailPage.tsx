import {useState} from "react";
import {useParams} from "react-router-dom";
import {useGetEmployeeById} from "@/query/employee.query.ts";
import Loading from "@/components/Loading.tsx";
import type {ColumnDef} from "@tanstack/react-table";
import type {AttendanceRecord} from "@/types/attendance.ts";
import {formatDate, getTimeFromTimestamp} from "@/lib/utils.ts";
import DataTable from "@/components/DataTable.tsx";
import {useAttendanceDetail} from "@/pages/HR/checkin/approve-attendance/hooks/useAttendanceDetail.ts";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog.tsx";
import {Tabs, TabsList, TabsTrigger} from "@/components/ui/tabs.tsx";
import {MapContainer, TileLayer, Marker, Popup} from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import {Button} from "@/components/ui/button.tsx";
import {useUpdateAttendanceStatus} from "@/query/attendance-record.query.ts";

// Fix for default marker icons in Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon-2x.png',
    iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
    shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
});

const ApproveAttendanceDetailPage = () => {
    const {id} = useParams<{ id: string }>();
    const [selectedItem, setSelectedItem] = useState<AttendanceRecord>();
    const [selectedImage, setSelectedImage] = useState<string | null>(null);
    const [selectedLocation, setSelectedLocation] = useState<{ lat: number, lng: number } | null>(null);
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const [rejectReason, setRejectReason] = useState<string>("");
    const [showRejectDialog, setShowRejectDialog] = useState<boolean>(false);

    const {data: employee, isPending, isError} = useGetEmployeeById(id ?? "");
    const {
        data: attendant,
        isPending: isPendingAttendant,
        isError: isErrorAttendant,
        refetch: refreshAttendance
    } = useAttendanceDetail(id || "");

    const {mutate: updateStatus} = useUpdateAttendanceStatus();

    if (isPending) return <Loading/>;
    if (isError || isErrorAttendant) return <p>Không thể tải dữ liệu. Vui lòng thử lại.</p>;

    // Filter records based on selected tab
    const filteredRecords = attendant?.filter(record => {
        if (statusFilter === "all") return true;
        return record.status === statusFilter;
    }) || [];

    const handleApprove = () => {
        if (!selectedItem) return;

        updateStatus({
            attendance_record_id: selectedItem.attendance_record_id,
            status: "approved",
            note_reject: ""
        }, {
            onSuccess: () => {
                refreshAttendance();
                setSelectedImage(null);
            }
        });
    };

    const handleReject = () => {
        if (!selectedItem) return;

        updateStatus({
            attendance_record_id: selectedItem.attendance_record_id,
            status: "rejected",
            note_reject: rejectReason
        }, {
            onSuccess: () => {
                refreshAttendance();
                setSelectedImage(null);
                setShowRejectDialog(false);
                setRejectReason("");
            }
        });
    };

    const InfoSection = (
        <div className="p-6 rounded-md">
            <h2 className="text-xl font-semibold mb-4">Thông tin nhân viên</h2>
            <div className={`flex gap-5 mb-10`}>
                <div className="space-y-2 w-full">
                    {[
                        {label: "Mã nhân viên", value: employee.employee_id},
                        {label: "Họ và tên", value: employee.full_name},
                        {label: "Văn phòng", value: employee.department?.office?.office_name},
                    ].map(({label, value}, index) => (
                        <div key={index} className="flex justify-between w-full border-b-1 border-gray-400 py-2">
                            <p className="text-gray-500 text-sm">{label}:</p>
                            <p className="font-medium text-left w-1/2">{value}</p>
                        </div>
                    ))}
                </div>
                <div className="space-y-2 w-full">
                    {[
                        {label: "Phòng ban", value: employee.department?.department_name},
                        {label: "Chức vụ", value: employee.job_title?.job_title},
                        {label: "Cấp bậc", value: employee.position?.position_name},
                    ].map(({label, value}, index) => (
                        <div key={index} className="flex justify-between w-full border-b-1 border-gray-400 py-2">
                            <p className="text-gray-500 text-sm">{label}:</p>
                            <p className="font-medium text-left w-1/2">{value}</p>
                        </div>
                    ))}
                </div>
            </div>
            <Tabs value={statusFilter} onValueChange={setStatusFilter} className="mb-4">
                <TabsList>
                    <TabsTrigger value="all">Tất cả</TabsTrigger>
                    <TabsTrigger value="pending">Chờ phê duyệt</TabsTrigger>
                    <TabsTrigger value="approved">Đã phê duyệt</TabsTrigger>
                    <TabsTrigger value="rejected">Từ chối</TabsTrigger>
                </TabsList>
            </Tabs>
        </div>
    );

    const columns: ColumnDef<AttendanceRecord>[] = [
        {
            accessorKey: "AttendanceCategory.attendance_category_name",
            header: "Hình thức chấm công",
        },
        {
            accessorKey: "gps_location",
            header: "Vị trí GPS",
            cell: ({row}) => {
                if (!row.original.latitude || !row.original.longitude) return null;
                return (
                    <button
                        onClick={() => setSelectedLocation({
                            lat: row.original.latitude!,
                            lng: row.original.longitude!
                        })}
                        className="text-blue-500 hover:underline"
                    >
                        Xem bản đồ
                    </button>
                )
            }
        },
        {
            accessorKey: "timestamp",
            header: "Thời gian yêu cầu",
            cell: ({row}) => {
                return (
                    <div className={`flex gap-3`}>
                        <p>{getTimeFromTimestamp(row.original.timestamp)}</p>
                        <p>{formatDate(row.original.timestamp)}</p>
                    </div>
                )
            }
        },
        {
            accessorKey: "note_request",
            header: "Ghi chú",
        },
        {
            accessorKey: "status",
            header: "Trạng thái",
            cell: ({row}) => {
                return (
                    <div
                        className={`flex gap-3 ${row.original.status === 'pending' ? "" : row.original.status === 'approved' ? 'text-green-500' : 'text-red-500'}`}>
                        {row.original.status === 'pending' ? "Chờ phê duyệt" : row.original.status === 'approved' ? 'Đã phê duyệt' : 'Từ chối'}
                    </div>
                )
            }
        },
        {
            id: "image",
            header: "Hình ảnh",
            cell: ({row}) => {
                if (!row.original.image_URL) return null;
                return (
                    <button
                        onClick={() => {
                            setSelectedImage(row.original.image_URL);
                            setSelectedLocation({
                                lat: row.original.latitude! || 0,
                                lng: row.original.longitude! || 0
                            });
                            setSelectedItem(row.original);
                        }}
                        className="text-blue-500 hover:underline"
                    >
                        Xem hình
                    </button>
                )
            }
        },
    ];

    return (
        <>
            <DataTable
                columns={columns}
                data={filteredRecords}
                isLoading={isPendingAttendant}
                navLink={InfoSection}
                title="Phê duyệt chấm công"
                keyFilter="timestamp"
            />

            {/* Image Preview Dialog */}
            <Dialog open={!!selectedImage} onOpenChange={(open) => {
                if (!open) {
                    setSelectedImage(null);
                    setSelectedItem(undefined);
                }
            }}>
                <DialogContent className="sm:max-w-3xl">
                    <DialogHeader>
                        <DialogTitle>Phê duyệt chấm công</DialogTitle>
                    </DialogHeader>
                    <div className="flex justify-center gap-3">
                        <div className="min-h-[20vh] min-w-[20vw] object-contain">
                            <p>Hình thức chấm công: {selectedItem?.AttendanceCategory?.attendance_category_name}</p>
                            {selectedImage && (
                                <img
                                    src={selectedImage}
                                    alt="Attendance proof"
                                    className="w-full h-full rounded-lg"
                                />
                            )}
                        </div>
                        <div className="min-h-[20vh] min-w-[20vw] object-contain">
                            <p>Vị trí GPS: {selectedItem?.latitude} - {selectedItem?.longitude}</p>
                            {selectedLocation && (
                                <MapContainer
                                    center={[selectedLocation?.lat, selectedLocation?.lng]}
                                    zoom={15}
                                    style={{height: '100%', width: '100%', borderRadius: '0.5rem'}}
                                >
                                    <TileLayer
                                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                                    />
                                    <Marker position={[selectedLocation?.lat, selectedLocation?.lng]}>
                                        <Popup>
                                            Vị trí chấm công <br/>
                                            Kinh độ: {selectedLocation?.lng?.toFixed(6) || 0} <br/>
                                            Vĩ độ: {selectedLocation?.lat?.toFixed(6) || 0}
                                        </Popup>
                                    </Marker>
                                </MapContainer>
                            )}
                        </div>
                    </div>
                    <div
                        className={`flex gap-3 mt-5 ${selectedItem?.status === 'pending' ? "" : selectedItem?.status === 'approved' ? 'text-green-500' : 'text-red-500'}`}>
                        Trạng thái: {selectedItem?.status === 'pending' ? "Chờ phê duyệt" : selectedItem?.status === 'approved' ? 'Đã phê duyệt' : 'Từ chối'}
                    </div>
                    <div>
                        <p>Ghi chú của nhân viên: {selectedItem?.note_request || "Không có ghi chú"}</p>
                    </div>
                    {selectedItem?.status === 'pending' && (
                        <div className={`w-full flex items-center gap-3 justify-center mt-4`}>
                            <button className={`rounded-xl px-5 cursor-pointer py-3 font-semibold text-white bg-red-700`}
                                onClick={() => setShowRejectDialog(true)}
                            >
                                Từ chối
                            </button>
                            <button className={`rounded-xl px-5 cursor-pointer py-3 font-semibold text-white bg-green-700`}
                                onClick={handleApprove}
                            >
                                Phê duyệt
                            </button>
                        </div>
                    )}
                </DialogContent>
            </Dialog>

            {/* Reject Reason Dialog */}
            <Dialog open={showRejectDialog} onOpenChange={setShowRejectDialog}>
                <DialogContent className="sm:max-w-3xl">
                    <DialogHeader>
                        <DialogTitle>Lý do từ chối</DialogTitle>
                    </DialogHeader>
                    <textarea
                        placeholder="Nhập lý do từ chối..."
                        value={rejectReason}
                        onChange={(e) => setRejectReason(e.target.value)}
                        className="mt-4"
                    />
                    <div className="flex justify-end gap-3 mt-4">
                        <Button
                            variant="outline"
                            onClick={() => setShowRejectDialog(false)}
                        >
                            Hủy
                        </Button>
                        <Button
                            variant="destructive"
                            onClick={handleReject}
                            disabled={!rejectReason.trim()}
                        >
                            Xác nhận từ chối
                        </Button>
                    </div>
                </DialogContent>
            </Dialog>
        </>
    );
};

export default ApproveAttendanceDetailPage;