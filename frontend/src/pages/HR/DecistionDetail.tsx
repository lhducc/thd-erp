import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { getDecisionByIdApi } from '@/apis/decision.api';
import { format } from 'date-fns';
import { vi } from 'date-fns/locale';

// Interface cho dữ liệu quyết định
interface DecisionDetail {
    decision_id: string;
    decision_name: string;
    effective_date: string;
    sign_date: string;
    content: string;
    condition: string;
    attached_file: string;
    created_date: string;
    decision_type: string;
    employees: Array<{
        employee_id: string;
        full_name: string;
        position: string;
        department: string;
    }>;
}

const DecisionDetail = () => {
    const { id } = useParams<{ id: string }>();
    const [decision, setDecision] = useState<DecisionDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

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
            return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: vi });
        } catch {
            return dateString;
        }
    };

    const getConditionBadgeVariant = (condition: string) => {
        switch (condition) {
            case 'Đang hiệu lực':
                return 'default';
            case 'Hết hiệu lực':
                return 'destructive';
            case 'Chưa hiệu lực':
                return 'secondary';
            default:
                return 'outline';
        }
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#DB3B21]"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="container mx-auto p-6">
                <div className="bg-destructive/15 text-destructive p-4 rounded-md mb-4">
                    <p>{error}</p>
                </div>
                <Button asChild>
                    <Link to="/decisions">Quay lại danh sách</Link>
                </Button>
            </div>
        );
    }

    if (!decision) {
        return (
            <div className="container mx-auto p-6">
                <div className="bg-destructive/15 text-destructive p-4 rounded-md mb-4">
                    <p>Không tìm thấy quyết định</p>
                </div>
                <Button asChild>
                    <Link to="/decisions">Quay lại danh sách</Link>
                </Button>
            </div>
        );
    }

    return (
        <div className="container mx-auto p-6 space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold">Chi tiết quyết định</h1>
                    <p className="text-muted-foreground">Mã quyết định: {decision.decision_id}</p>
                </div>
                <Button asChild>
                    <Link to="/decisions">Quay lại danh sách</Link>
                </Button>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Thông tin chính */}
                <div className="lg:col-span-2 space-y-6">
                    <Card>
                        <CardHeader>
                            <CardTitle>Thông tin quyết định</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div>
                                    <label className="text-sm font-medium text-muted-foreground">Tên quyết định</label>
                                    <p className="text-lg font-semibold">{decision.decision_name}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-muted-foreground">Loại quyết định</label>
                                    <p className="text-lg">{decision.decision_type}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-muted-foreground">Ngày ký</label>
                                    <p className="text-lg">{formatDate(decision.sign_date)}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-muted-foreground">Ngày hiệu lực</label>
                                    <p className="text-lg">{formatDate(decision.effective_date)}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-muted-foreground">Ngày tạo</label>
                                    <p className="text-lg">{formatDate(decision.created_date)}</p>
                                </div>
                                <div>
                                    <label className="text-sm font-medium text-muted-foreground">Tình trạng</label>
                                    <Badge variant={getConditionBadgeVariant(decision.condition)}>
                                        {decision.condition}
                                    </Badge>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Nội dung quyết định */}
                    <Card>
                        <CardHeader>
                            <CardTitle>Nội dung quyết định</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="prose max-w-none">
                                <p className="whitespace-pre-wrap">{decision.content}</p>
                            </div>
                        </CardContent>
                    </Card>

                    {/* File đính kèm */}
                    {decision.attached_file && (
                        <Card>
                            <CardHeader>
                                <CardTitle>File đính kèm</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="flex items-center justify-between p-4 border rounded-lg">
                                    <div className="flex items-center space-x-3">
                                        <div className="p-2 bg-muted rounded">
                                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                            </svg>
                                        </div>
                                        <div>
                                            <p className="font-medium">{decision.attached_file}</p>
                                            <p className="text-sm text-muted-foreground">Tài liệu đính kèm</p>
                                        </div>
                                    </div>
                                    <Button variant="outline" asChild>
                                        <a
                                            href={`/api/files/${decision.attached_file}`}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                        >
                                            Tải xuống
                                        </a>
                                    </Button>
                                </div>
                            </CardContent>
                        </Card>
                    )}
                </div>

                {/* Danh sách nhân viên */}
                <div className="space-y-6">
                    <Card>
                        <CardHeader>
                            <CardTitle>Nhân viên áp dụng</CardTitle>
                            <p className="text-sm text-muted-foreground">
                                {decision.employees.length} nhân viên
                            </p>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-3">
                                {decision.employees.map((employee) => (
                                    <div key={employee.employee_id} className="p-3 border rounded-lg">
                                        <div className="flex justify-between items-start">
                                            <div>
                                                <p className="font-medium">{employee.full_name}</p>
                                                <p className="text-sm text-muted-foreground">ID: {employee.employee_id}</p>
                                                <p className="text-sm">{employee.position}</p>
                                                <p className="text-sm text-muted-foreground">{employee.department}</p>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>

                    {/* Actions */}
                    <Card>
                        <CardHeader>
                            <CardTitle>Thao tác</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-3">
                            <Button className="w-full" variant="outline" asChild>
                                <Link to={`/decisions/edit/${decision.decision_id}`}>
                                    Chỉnh sửa quyết định
                                </Link>
                            </Button>
                            <Button className="w-full" variant="outline" asChild>
                                <a
                                    href={`/api/decisions/${decision.decision_id}/export`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                >
                                    Xuất PDF
                                </a>
                            </Button>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    );
};

export default DecisionDetail;