import {Link, useParams} from "react-router-dom";
import {useAllEmployee, useGetEmployeeById} from "@/query/useEmployee.ts";
import Loading from "@/components/Loading.tsx";
import DataTable from "@/components/DataTable.tsx";
import type {ColumnDef} from "@tanstack/react-table";
import type {Employee} from "@/types/employee.ts";
import more from "@/assets/more.svg";
import {formatDate} from "date-fns";
import {Button} from "@/components/ui/button.tsx";

const AttendantHistoryDetail = () => {
    const { id } = useParams<{ id: string }>();

    const { data: employee, isPending, isError } = useGetEmployeeById(id ?? "");
    console.log(employee);
    const { data: attendant, isPending: isPendingAttendant, isError: isErrorAttendant } = useAllEmployee();

    if (isPending) return <Loading />;
    if (isError || !employee) return <p>Không thể tải dữ liệu. Vui lòng thử lại.</p>;

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

    const columns: ColumnDef<Employee>[] = [
        {
            accessorKey: "employee_id",
            header: "Mã NV",
        },
        {
            accessorKey: "full_name",
            header: "Tên nhân viên",
        },
        {
            accessorKey: "department.office.office_name",
            header: "Văn phòng",
        },{
            accessorKey: "department.department_name",
            header: "Phòng ban",
        },{
            accessorKey: "job_title.job_title",
            header: "Chức danh",
        },{
            accessorKey: "position.position_name",
            header: "Cấp bậc",
        },
        {
            accessorKey: "is_gps",
            header: "Xem chi tiết",
            cell: ({row}) => {
                const data = row.original
                return(
                    <div className={`flex justify-center items-center`}>
                        <Link to={`${data.employee_id}`}>
                            <img className={`w-[25px] h-[25px]`} src={more} alt=""/>
                        </Link>
                    </div>
                )
            },
        },
    ];

    const createAttendant = (
        <Button value={`Tạo chấm công`} />
    )

    return (
        <>
            <DataTable
                columns={columns}
                data={attendant || []}
                isLoading={isPendingAttendant}
                navLink={InfoSection}
                title="Bảng lịch sử chấm công chi tiết"
                buttonCreate={createAttendant}
                // keyFilter="full_name"
            />
        </>
    );
};

export default AttendantHistoryDetail;
