import {useParams} from "react-router-dom";
import {useGetEmployeeById} from "@/query/employee.query.ts";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import type {ColumnDef} from "@tanstack/react-table";
import {useGetAttendanceRecordByEmployeeId} from "@/query/attendance-record.query.ts";
import type {AttendanceRecord} from "@/types/attendance.ts";
import {formatDate, getTimeFromTimestamp} from "@/lib/utils.ts";
import ReviewImage from "@/components/ReviewImage.tsx";
import {CreateAttendant} from "@/components/attendance/CreateAttendance.tsx";

const AttendantHistoryDetail = () => {
    const { id } = useParams<{ id: string }>();

    const { data: employee, isPending, isError } = useGetEmployeeById(id ?? "");
    const { data: attendant, isPending: isPendingAttendant, isError: isErrorAttendant, refetch: refreshAttendance } = useGetAttendanceRecordByEmployeeId(id || "");

    if (isPending) return <Loading />;
    if (isError || isErrorAttendant) return <p>Không thể tải dữ liệu. Vui lòng thử lại.</p>;

    const InfoSection = (
        <div className="p-6 rounded-md">
            <h2 className="text-xl font-semibold mb-4">Thông tin nhân viên</h2>
            <div className={`flex gap-5`}>
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
        </div>
    );

    const columns: ColumnDef<AttendanceRecord>[] = [
        {
            accessorKey: "timestamp",
            header: "Ngày",
            cell: ({row}) => formatDate(row.original.timestamp)
        },
        {
            accessorKey: "timestamp",
            header: "Giờ chấm công",
            cell: ({row}) => getTimeFromTimestamp(row.original.timestamp)
        },
        {
            accessorKey: "AttendanceCategory.attendance_category_name",
            header: "Hình thức chấm công",
        },
        {
            accessorKey: "department.department_name",
            header: "Vị trí",
        },
        {
            accessorKey: "image_URL",
            header: "xem hình ảnh chấm công",
            cell: ({row}) => {
                const url = row.original.image_URL
                return(
                    <>
                        {   url !== ""  ?
                            <ReviewImage url={url} /> :
                            "Trống"
                        }
                    </>
                )
            }
        },
        {
            accessorKey: "create_by_info.full_name",
            header: "Người tạo",
        },
    ];

    return (
        <>
            <DataTable
                columns={columns}
                data={attendant || []}
                isLoading={isPendingAttendant}
                navLink={InfoSection}
                title="Bảng lịch sử chấm công chi tiết"
                buttonCreate={<CreateAttendant employeeId={id || ""} refresh={refreshAttendance} />}
                keyFilter="timestamp"
            />
        </>
    );
};

export default AttendantHistoryDetail;
