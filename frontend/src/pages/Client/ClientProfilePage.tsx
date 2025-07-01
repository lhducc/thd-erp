import { useQuery } from "@tanstack/react-query";
import { getEmployeeByIdApi } from "@/apis/profile.api";
import { useAuth } from "@/context/AuthContext";

const formatDate = (dateStr?: string) => {
    if (!dateStr) return "-";
    const date = new Date(dateStr);
    return date.toLocaleDateString("vi-VN"); // Format: DD/MM/YYYY
};

const ClientProfilePage = () => {
    const { currentUser } = useAuth();
    const employeeId = currentUser?.user_id;

    const {
        data: employeeResponse,
        isLoading,
    } = useQuery({
        queryKey: ["employee", employeeId],
        queryFn: () => getEmployeeByIdApi(employeeId!),
        enabled: !!employeeId,
    });

    const employee = employeeResponse?.data;

    return (
        <div>

            <p className={`font-semibold text-[25px]`}>Chi tiết hồ sơ nhân viên</p>
        <div className="p-6 max-w-3xl mx-auto space-y-4 text-sm">
            {isLoading ? (
                <p>Loading...</p>
            ) : employee ? (
                <>
                    <InfoRow label="Họ và tên" value={employee.full_name} />
                    <InfoRow label="Mã nhân viên" value={employee.employee_id} />
                    <InfoRow label="Giới tính" value={employee.gender} />
                    <InfoRow label="Ngày sinh" value={formatDate(employee.birthday)} />
                    <InfoRow label="Số điện thoại" value={employee.phone_number} />
                    <InfoRow label="Email" value={employee.email} />
                    <InfoRow label="Địa chỉ" value={employee.address} />
                    <InfoRow label="Ngày bắt đầu" value={formatDate(employee.created_date)} />
                    <InfoRow label="Ngày chính thức" value="-" />
                    <InfoRow label="Ngày kết thúc" value="-" />
                    <InfoRow label="Chi nhánh" value={employee.department?.office?.office_name + " - HCM"} />
                    <InfoRow label="Chức danh" value={employee.job_title?.job_title} />
                    <InfoRow label="Phòng ban" value={employee.department?.department_name} />
                    <InfoRow label="Quản lý trực tiếp" value={employee.manager?.full_name} />
                </>
            ) : (
                <p>Không có dữ liệu</p>
            )}
        </div>
        </div>
    );
};

const InfoRow = ({ label, value }: { label: string; value: string }) => (
    <div className="flex justify-between border-b py-2">
        <span className="font-semibold">{label}</span>
        <span>{value}</span>
    </div>
);

export default ClientProfilePage;
