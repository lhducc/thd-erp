'use client';

import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { useState, useEffect } from "react";
import { Pencil, Trash2, Download, ArrowLeft } from "lucide-react";
import { useParams, useNavigate } from "react-router-dom";
import { getDecisionByIdApi } from "@/apis/decision.api";
import { format } from "date-fns";
import { vi } from "date-fns/locale";

interface Employee {
    employee_id: string;
    fullname: string;
    position?: string;
    department?: string;
}

interface DecisionDetail {
    decision_id: string;
    decision_name: string;
    effective_date: string;
    sign_date: string;
    content: string;
    condition: string;
    attached_file: string;
    created_date: string;
    decision_type_id: string;
    decision_type_name: string;
    employees: Employee[];
}

const DecisionDetail = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [decision, setDecision] = useState<DecisionDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [downloading, setDownloading] = useState(false);

    useEffect(() => {
        const fetchDecisionDetail = async () => {
            if (!id) return;

            try {
                setLoading(true);
                const response = await getDecisionByIdApi(id);
                setDecision(response.data);
                setError(null);
            } catch (err: any) {
                setError(err.response?.data?.message || 'Không thể tải thông tin quyết định');
                console.error('Failed to fetch decision detail:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchDecisionDetail();
    }, [id]);

    const formatDate = (dateString: string) => {
        try {
            return format(new Date(dateString), 'dd/MM/yyyy', { locale: vi });
        } catch {
            return dateString;
        }
    };

    const handleExport = async () => {
        if (!decision?.attached_file) return;

        try {
            setDownloading(true);

            // Tạo một thẻ a ẩn để tải file
            const link = document.createElement('a');
            link.href = decision.attached_file;

            // Lấy tên file từ URL (phần sau cùng)
            const fileName = decision.attached_file.split('/').pop() || `quyet-dinh-${decision.decision_id}.pdf`;
            link.download = fileName;

            // Thêm vào DOM, click và xóa
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

        } catch (err) {
            console.error('Failed to download file:', err);
            alert('Không thể tải file. Vui lòng thử lại.');
        } finally {
            setDownloading(false);
        }
    };

    // Alternative: Sử dụng fetch API nếu cần xử lý CORS hoặc authentication
    const handleExportWithFetch = async () => {
        if (!decision?.attached_file) return;

        try {
            setDownloading(true);

            const response = await fetch(decision.attached_file);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);

            const link = document.createElement('a');
            link.href = url;

            // Lấy tên file từ URL hoặc sử dụng mã quyết định
            const fileName = decision.attached_file.split('/').pop() || `quyet-dinh-${decision.decision_id}.pdf`;
            link.download = fileName;

            document.body.appendChild(link);
            link.click();

            // Dọn dẹp
            window.URL.revokeObjectURL(url);
            document.body.removeChild(link);

        } catch (err) {
            console.error('Failed to download file:', err);
            alert('Không thể tải file. Vui lòng thử lại.');
        } finally {
            setDownloading(false);
        }
    };

    const handleEdit = () => {
        navigate(`/decisions/edit/${id}`);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#DB3B21]"></div>
            </div>
        );
    }

    if (error || !decision) {
        return (
            <div className="container mx-auto p-6">
                <div className="bg-red-100 text-red-700 p-4 rounded-md mb-4">
                    <p>{error || 'Không tìm thấy quyết định'}</p>
                </div>
                <Button onClick={() => navigate(-1)} className="flex items-center gap-2">
                    <ArrowLeft size={16} />
                    Quay lại
                </Button>
            </div>
        );
    }

    const Header = () => (
        <div className="flex items-center justify-between bg-white p-6 border-b-2 border-gray-800">
            <div>
                <Label className="text-3xl font-semibold">CHI TIẾT QUYẾT ĐỊNH</Label>
                <Label className="text-xl font-medium block mt-2">{decision.decision_name}</Label>
            </div>
            <div className="flex gap-3">
                <Button
                    className="p-3 bg-blue-600 text-white rounded-2xl hover:bg-blue-700 flex items-center gap-2"
                    onClick={handleEdit}
                >
                    <Pencil size={16} />
                    Chỉnh sửa
                </Button>
                <Button
                    className="p-3 bg-[#DB290D] text-white rounded-2xl hover:bg-[#b83a1a] flex items-center gap-2"
                    onClick={handleExport}
                    disabled={!decision.attached_file || downloading}
                >
                    {downloading ? (
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    ) : (
                        <Download size={16} />
                    )}
                    {downloading ? 'Đang tải...' : 'Xuất file'}
                </Button>
            </div>
        </div>
    );

    const DocumentInfo = () => (
        <div className="bg-white w-[700px] mt-[10px] rounded-md p-6 ">
            <p className="text-[22px] font-bold py-[10px] border-b border-gray-400">
                {decision.decision_name}
            </p>
            <div className="flex text-[20px]">
                <div className="space-y-4 w-[300px]">
                    <p className="font-semibold border-b border-gray-400 py-2">Mã số:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Loại quyết định:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Ngày tạo:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Ngày ký:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Hiệu lực từ ngày:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Tình trạng:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Số nhân viên áp dụng:</p>
                </div>
                <div className="space-y-4 w-[350px]">
                    <p className="text-gray-600 border-b border-gray-400 py-2">{decision.decision_id}</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">{decision.decision_type_name}</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">{formatDate(decision.created_date)}</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">{formatDate(decision.sign_date)}</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">{formatDate(decision.effective_date)}</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">{decision.condition}</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">{decision.employees?.length || 0} nhân viên</p>
                </div>
            </div>
            <div className="pt-[20px]">
                <p className="font-semibold text-lg mb-2">Nội dung</p>
                <div className="bg-gray-50 p-4 rounded-md">
                    <p className="text-gray-700 whitespace-pre-wrap">{decision.content}</p>
                </div>
            </div>

            {/* Danh sách nhân viên */}
            {decision.employees && decision.employees.length > 0 && (
                <div className="pt-[20px]">
                    <p className="font-semibold text-lg mb-2">Nhân viên áp dụng</p>
                    <div className="bg-gray-50 p-4 rounded-md max-h-60 overflow-y-auto">
                        {decision.employees.map((employee, index) => (
                            <div key={employee.employee_id} className="border-b border-gray-200 py-2 last:border-b-0">
                                <p className="font-medium">{employee.fullname}</p>
                                <p className="text-sm text-gray-600">Mã: {employee.employee_id}</p>
                                {employee.position && <p className="text-sm text-gray-500">{employee.position}</p>}
                                {employee.department && <p className="text-sm text-gray-500">{employee.department}</p>}
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );

    const DocumentAttachment = () => (
        <div className="mt-8 bg-white p-6 rounded-md mr-[30px] w-[400px]">
            <h3 className="text-xl font-semibold mb-4">Tài liệu đính kèm</h3>
            <div className="border border-gray-300 p-6 rounded-md h-80 flex flex-col items-center justify-center">
                {decision.attached_file ? (
                    <>
                        <div className="text-center mb-4">
                            <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-2">
                                <Download className="text-red-600" size={24} />
                            </div>
                            <p className="text-sm text-gray-600 truncate max-w-xs">
                                {decision.attached_file.split('/').pop()}
                            </p>
                            <p className="text-xs text-gray-400 mt-1">Click để tải xuống</p>
                        </div>
                        <Button
                            onClick={handleExport}
                            className="bg-[#DB290D] hover:bg-[#b83a1a] flex items-center gap-2"
                            disabled={downloading}
                        >
                            {downloading ? (
                                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                            ) : (
                                <Download size={16} />
                            )}
                            {downloading ? 'Đang tải...' : 'Tải xuống'}
                        </Button>
                    </>
                ) : (
                    <p className="text-gray-600 text-center">Không có tài liệu đính kèm</p>
                )}
            </div>
        </div>
    );

    return (
        <div className="min-h-screen bg-gray-50">
            <Header />
            <div className="flex justify-between p-6">
                <DocumentInfo />
                <DocumentAttachment />
            </div>
        </div>
    );
};

export default DecisionDetail;