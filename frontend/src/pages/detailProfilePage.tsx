import {Button} from "@/components/ui/button";
import {Label} from "@/components/ui/label";
import {useState} from "react";
import {Pencil, Trash2} from "lucide-react";
import DocumentEmployeeForm from "@/components/DocumentEmployeeForm"
import {useGetEmployeeById} from "@/query/employee.query.ts";
import {useParams} from "react-router-dom";
import Loading from "@/components/Loading.tsx";
import {useDocumentEmployeeById} from "@/query/useDocumentEmployee.ts";
import {formatDate} from "date-fns";

const DetailEmployeePage = () => {
    const [open, setOpen] = useState(false);
    const { id } = useParams<{ id: string }>();

    const { data: employee, isPending, isError } = useGetEmployeeById(id ?? "");

    const { data: document, isPending: isDocumentPending, isDocumentError } = useDocumentEmployeeById(id ?? "");

    if (isPending) return <Loading />;
    if (isError || !employee) return <p>Không thể tải dữ liệu. Vui lòng thử lại.</p>;

    const Header = () => (
        <div className="flex items-center justify-between bg-white p-4 border-b-1 border-gray-800">
            <Label className="text-2xl font-bold">CHI TIẾT HỒ SƠ NHÂN VIÊN</Label>
            <DocumentEmployeeForm className="w-[1261px]" open={open} setOpen={setOpen}/>
            <div className="flex space-x-2">
                <span
                    className={`px-4 py-2 border-2 rounded-full font-bold transition ${
                        employee.status === "active"
                            ? "border-green-500 text-green-500 hover:bg-green-500 hover:text-white"
                            : "border-red-500 text-red-500 hover:bg-red-500 hover:text-white"
                    }`}>
                    {employee.status === "active" ? "ACTIVE" : "INACTIVE"}
                </span>
                {["Chỉnh sửa thông tin", "Tạo hợp đồng", "Thêm tài liệu", "Xóa hồ sơ"].map((text, idx) => (
                    text === "Thêm tài liệu" ? (
                        <Button onClick={() => setOpen(true)} key={idx}
                                className="bg-red-500 hover:bg-red-600 text-white rounded-full px-3 py-2">
                            {text}
                        </Button>
                    ) : (
                        <Button key={idx} className="bg-red-500 hover:bg-red-600 text-white rounded-full px-3 py-2">
                            {text}
                        </Button>
                    )
                ))}
            </div>
        </div>
    );

    const InfoSection = () => (
        <div className="bg-white p-6 rounded-md w-[45%]">
            <h2 className="text-xl font-semibold mb-4">Thông tin nhân viên</h2>
            <div className="space-y-2">
                {[
                    {label: "Họ và tên", value: employee.full_name},
                    {label: "Mã nhân viên", value: employee.employee_id},
                    {label: "Giới tính", value: employee.gender === "male" ? "Nam" : "Nữ"},
                    {label: "Ngày sinh", value: formatDate(employee.birthday)},
                    {label: "Số điện thoại", value: employee.phone_number},
                    {label: "Email", value: employee.email},
                    {label: "Địa chỉ hiện tại", value: employee.address},
                    {label: "Chức vụ", value: employee.job_title?.job_title},
                    {label: "Vị trí", value: employee.position?.position_name},
                    {label: "Văn phòng", value: employee.department?.office?.office_name},
                    {label: "Phòng ban", value: employee.department?.department_name},
                    {label: "Ngày bắt đầu", value: formatDate(employee.created_date)},
                    {label: "Quản lý trực tiếp", value: employee.manager?.full_name},
                    {label: "Loại công việc", value: employee.work_type},
                ].map(({label, value}, index) => (
                    <div key={index} className="flex justify-between w-full border-b-1 border-gray-400 py-2">
                        <p className="text-gray-500 text-sm">{label}:</p>
                        <p className="font-medium text-left w-1/2">{value}</p>
                    </div>
                ))}
            </div>
        </div>
    );

    const DataTable = ({title, columns, data}: any) => (
        <div className="bg-white p-4 rounded-md w-full">
            <h2 className="text-lg font-bold mb-2">{title}</h2>
            <table className="w-full border text-sm">
                <thead className="bg-gray-100">
                <tr>
                    {columns.map((col: string, i: number) => (
                        <th key={i} className="border px-3 py-2 text-left">{col}</th>
                    ))}
                    <th className="border px-3 py-2">Chỉnh sửa</th>
                </tr>
                </thead>
                <tbody>
                {data.map((row: any, idx: number) => (
                    <tr key={idx}>
                        {Object.values(row).map((val: any, i: number) => (
                            <td key={i} className="border px-3 py-2">{val}</td>
                        ))}
                        <td className="border px-3 py-2 text-center">
                            <div className="flex justify-center gap-2">
                                <Pencil className="w-4 h-4 cursor-pointer"/>
                                <Trash2 className="w-4 h-4 cursor-pointer text-red-500"/>
                            </div>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
            <div className="flex justify-center mt-3 space-x-1">
                <button className="px-2 py-1 border rounded hover:bg-gray-200">
                    {"<"}
                </button>

                {[1, 2, 3, 4, 5].map((page) => (
                    <button key={page} className="px-2 py-1 border rounded hover:bg-gray-200">
                        {page}
                    </button>
                ))}

                <button className="px-2 py-1 border rounded hover:bg-gray-200">
                    {">"}
                </button>
            </div>
        </div>
    );

    // const documentData = [
    //     {type: "Chứng chỉ HSK6", status: "Hết hiệu lực", condition: "Chờ duyệt", expired: "13/05/2025"},
    //     {type: "Chứng chỉ HSK6", status: "Đang hiệu lực", condition: "Đang hiệu lực", expired: "13/05/2025"},
    //     {type: "Chứng chỉ HSK6", status: "Đang hiệu lực", condition: "Đang hiệu lực", expired: "13/05/2025"},
    //     {type: "Chứng chỉ HSK6", status: "Đang hiệu lực", condition: "Đã duyệt", expired: "13/05/2025"},
    //     {type: "Chứng chỉ HSK6", status: "Đang hiệu lực", condition: "Đang hiệu lực", expired: "13/05/2025"},
    // ];

    const contractData = [
        {id: "HD000009", name: "Hợp đồng XYZ", date: "10/07/2025", status: "Hết hiệu lực"},
        {id: "HD000001", name: "Hợp đồng lao động", date: "10/07/2025", status: "Đang hiệu lực"},
        {id: "HD000001", name: "Hợp đồng lao động", date: "10/07/2025", status: "Đang hiệu lực"},
        {id: "HD000001", name: "Hợp đồng lao động", date: "10/07/2025", status: "Đang hiệu lực"},
        {id: "HD000001", name: "Hợp đồng lao động", date: "10/07/2025", status: "Đang hiệu lực"},
    ];

    return (
        <div className="space-y-4">
            <Header/>
            <div className="flex">
                <InfoSection/>
                <div className="flex flex-col w-[55%] gap-4">
                    <DataTable
                        title="Hồ sơ nhân viên"
                        columns={["Loại tài liệu", "Trạng thái", "Tình trạng", "Ngày hết hạn"]}
                        data={document?.map(item => ({
                            ...item
                        }))}
                    />
                    <DataTable
                        title="Hợp đồng liên quan"
                        columns={["Mã Hợp đồng", "Tên hợp đồng", "Ngày tạo", "Trạng thái"]}
                        data={contractData.map(item => ({
                            ...item
                        }))}
                    />
                </div>
            </div>
        </div>
    );
};

export default DetailEmployeePage;