import {useParams} from "react-router-dom";
import {
    useCalculatorTimesheet,
    useExportTimesheet,
    useGetTimesheet,
    useLockTimesheet,
    useResetTimesheet,
    useUpdateTimesheetDetail
} from "@/query/timesheet.query";
import {useCallback, useMemo, useState} from "react";
import type {TimesheetDetail, TimesheetInfor, TimesheetList} from "@/types/timesheet";
import {Dialog, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog";
import {Button} from "@/components/ui/button";
import {Edit, TimerReset} from "lucide-react";
import {toast} from "sonner";
import axios from "axios";
import Loading from "@/components/Loading";
import EditDialog from "@/pages/HR/checkin/timesheet/components/EditDialog";

interface AdjustWorkDayForm {
    adjusted_work_day: number;
}

const TimesheetDetailPage = () => {
    const {id} = useParams();
    const {data, isLoading, refetch} = useGetTimesheet(id || "");
    const {mutate: updateTimesheetDetail, isPending: isUpdatingTimesheet} = useUpdateTimesheetDetail();
    const {mutate: lockTimesheet, isPending: isLockTimesheet} = useLockTimesheet();
    const {mutate: exportTimesheet} = useExportTimesheet();
    const {mutate: calculatorTimesheet, isPending: isCalculatorTimesheet} = useCalculatorTimesheet();
    const {mutate: resetTimesheet, isPending: isResetTimesheet} = useResetTimesheet();
    const [expandedRows, setExpandedRows] = useState<Set<number>>(new Set());
    const [selectedDetail, setSelectedDetail] = useState<TimesheetDetail | null>(null);
    const [isDetailDialogOpen, setIsDetailDialogOpen] = useState(false);
    const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

    const timesheetData = data as unknown as TimesheetList;

    // Generate columns for each day in the month
    const generateDayColumns = useCallback((timesheetData: TimesheetList) => {
        const startDate = new Date(timesheetData.start_date);
        const endDate = new Date(timesheetData.end_date);
        const days = [];
        const currentDate = new Date(startDate);

        while (currentDate <= endDate) {
            const day = currentDate.getDate();
            const dayOfWeek = currentDate.getDay();
            const dateKey = currentDate.toISOString().split('T')[0];

            days.push({
                accessorKey: `day_${day}`,
                header: (
                    <div className="text-center">
                        <div>{day}</div>
                        <div className="text-xs text-gray-500">
                            {['CN', 'T2', 'T3', 'T4', 'T5', 'T6', 'T7'][dayOfWeek]}
                        </div>
                    </div>
                ),
                cell: ({row}: { row: any }) => {
                    const timesheet = row.original;
                    const detail = timesheet.details?.find((d: TimesheetDetail) =>
                        d.date.startsWith(dateKey)
                    );

                    if (!detail) {
                        return <div className="p-1 text-center">-</div>;
                    }

                    let bgColor = "bg-white";
                    let text = "";
                    let tooltip = "";

                    if (detail.is_absent) {
                        bgColor = "bg-red-100";
                        text = "V";
                        tooltip = `Vắng mặt: ${detail.absent_reason || "Không có lý do"}`;
                    } else if (detail.is_holiday) {
                        bgColor = "bg-purple-100";
                        text = "L";
                        tooltip = "Ngày lễ";
                    } else if (detail.is_weekend) {
                        bgColor = "bg-blue-100";
                        text = detail.day_of_week === 6 ? "T7" : "CN";
                        tooltip = "Cuối tuần";
                    } else if (detail.work_days > 0) {
                        bgColor = "bg-green-100";
                        text = "C";
                        tooltip = `Đi làm: ${detail.work_hours} giờ`;

                        if (detail.is_late) {
                            bgColor = "bg-orange-100";
                            tooltip += `, Trễ: ${detail.late_minutes} phút`;
                        }
                    } else if (detail.leave_type) {
                        bgColor = "bg-yellow-100";
                        text = "P";
                        tooltip = `Nghỉ phép: ${detail.leave_type}`;
                    }

                    return (
                        <div
                            className={`p-1 text-center rounded ${bgColor} cursor-pointer hover:opacity-80`}
                            title={tooltip}
                            onClick={(e) => {
                                e.stopPropagation();
                                openDetailDialog(detail);
                            }}
                        >
                            {text}
                        </div>
                    );
                }
            });

            currentDate.setDate(currentDate.getDate() + 1);
        }

        return days;
    }, []);

    // Tối ưu generateDayColumns với useMemo
    const dayColumns = useMemo(() => {
        if (!timesheetData) return [];
        return generateDayColumns(timesheetData);
    }, [timesheetData, generateDayColumns]);

    // Tối ưu các hàm xử lý sự kiện với useCallback
    const toggleRowExpansion = useCallback((timesheetId: number) => {
        setExpandedRows(prev => {
            const newExpanded = new Set(prev);
            if (newExpanded.has(timesheetId)) {
                newExpanded.delete(timesheetId);
            } else {
                newExpanded.add(timesheetId);
            }
            return newExpanded;
        });
    }, []);

    const openDetailDialog = useCallback((detail: TimesheetDetail) => {
        setSelectedDetail(detail);
        setIsDetailDialogOpen(true);
    }, []);

    const openEditDialog = useCallback((detail: TimesheetDetail) => {
        setSelectedDetail(detail);
        setIsEditDialogOpen(true);
    }, []);

    const handleReset = useCallback((timesheetDetailId: number) => {
        resetTimesheet({
            timesheetDetailId: timesheetDetailId,
        }, {
            onSuccess: async () => {
                setIsDetailDialogOpen(false);
                toast.success("Reset số công thành công");
                await refetch();
            },
            onError: (error) => {
                if (axios.isAxiosError(error)) {
                    toast.error(error.response?.data?.message || error.message);
                }
                setIsDetailDialogOpen(false);
            }
        });
    }, [resetTimesheet, refetch]);

    const handleEditSubmit = useCallback((formData: AdjustWorkDayForm) => {
        if (!selectedDetail) return;

        updateTimesheetDetail(
            {
                timesheetDetailId: selectedDetail.timesheet_detail_id,
                data: {adjusted_work_day: Number(formData.adjusted_work_day)}
            },
            {
                onSuccess: async () => {
                    setIsEditDialogOpen(false);
                    setIsDetailDialogOpen(false);
                    await refetch();
                },
            }
        );
    }, [selectedDetail, updateTimesheetDetail, refetch]);

    const handleLockTimesheet = useCallback(() => {
        if (!id) return;

        lockTimesheet({
            timesheetDetailId: id,
        }, {
            onSuccess: async () => {
                toast.success("Chốt công thành công");
                await refetch();
            },
            onError: (error) => {
                if (axios.isAxiosError(error)) {
                    toast.error(error.response?.data?.message || error.message);
                }
            }
        });
    }, [id, lockTimesheet, refetch]);

    const handleExport = useCallback(() => {
        if (!id) return;

        exportTimesheet({
            timesheetDetailId: id,
        });
    }, [id, exportTimesheet]);

    const handleCalculatorTimesheet = useCallback(() => {
        if (!id) return;

        calculatorTimesheet({
            timesheetDetailId: id,
        }, {
            onSuccess: async () => {
                toast.success("Đồng bộ bảng công thành công");
                await refetch();
            },
            onError: (error) => {
                if (axios.isAxiosError(error)) {
                    toast.error(error.response?.data?.message || error.message);
                }
            }
        });
    }, [id, calculatorTimesheet, refetch]);

    // Render expanded row with detailed daily information
    const renderRowExpansion = useCallback((timesheet: TimesheetInfor) => {
        if (!timesheet.details || timesheet.details.length === 0) return null;

        return (
            <div className="p-4 bg-gray-50 border-b">
                <h4 className="font-semibold mb-2">Chi tiết chấm công: {timesheet.employee.full_name}</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {timesheet.details?.map((detail: TimesheetDetail, idx: number) => (
                        <div
                            key={`${detail.timesheet_detail_id}-${idx}`}
                            className="bg-white p-3 rounded shadow-sm cursor-pointer hover:shadow-md transition-shadow"
                            onClick={() => openDetailDialog(detail)}
                        >
                            <div className="font-medium">
                                {new Date(detail.date).toLocaleDateString('vi-VN', {
                                    weekday: 'long',
                                    year: 'numeric',
                                    month: 'long',
                                    day: 'numeric'
                                })}
                            </div>
                            <div className="text-sm text-gray-600 mt-1">
                                Ca: {detail.work_shift?.workshift_name || "N/A"}
                            </div>
                            <div className="text-sm mt-1">
                                Trạng thái: {detail.is_absent ?
                                <span className="text-red-600">Vắng mặt ({detail.absent_reason})</span> :
                                detail.work_days > 0 ?
                                    <span className="text-green-600">Đi làm</span> :
                                    <span className="text-yellow-600">Nghỉ phép</span>
                            }
                            </div>
                            {detail.checkin_record && (
                                <div className="text-sm mt-1">
                                    Checkin: {detail.checkin_record.timestamp.slice(11, 19)}
                                    {detail.is_late && (
                                        <span className="text-orange-600 ml-2">(Trễ: {detail.late_minutes} phút)</span>
                                    )}
                                </div>
                            )}
                            {detail.checkout_record && (
                                <div className="text-sm mt-1">
                                    Checkout: {detail.checkout_record.timestamp.slice(11, 19)}
                                </div>
                            )}
                            <div className="text-sm mt-1">
                                Số giờ: {detail.work_hours} giờ ({detail.work_days} công)
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        );
    }, [openDetailDialog]);

    // Detail Dialog Component
    const DetailDialog = useCallback(() => {
        if (!selectedDetail) return null;

        return (
            <Dialog open={isDetailDialogOpen} onOpenChange={setIsDetailDialogOpen}>
                <DialogContent className="max-w-xl">
                    <DialogHeader>
                        <DialogTitle className="flex justify-between items-center gap-2">
                            <span>
                                Chi tiết ngày {selectedDetail.date && new Date(selectedDetail.date).toLocaleDateString('vi-VN')}
                            </span>
                            <div className="flex gap-2">
                                <Button size="sm" onClick={() => openEditDialog(selectedDetail)}>
                                    <Edit className="h-4 w-4 mr-1"/>
                                    Sửa công
                                </Button>
                                {selectedDetail.is_manually_adjusted && (
                                    <Button
                                        size="sm"
                                        onClick={() => handleReset(selectedDetail.timesheet_detail_id)}
                                        disabled={isResetTimesheet}
                                    >
                                        <TimerReset className="h-4 w-4 mr-1"/>
                                        {isResetTimesheet ? <Loading/> : "Reset công"}
                                    </Button>
                                )}
                            </div>
                        </DialogTitle>
                    </DialogHeader>

                    <div className="space-y-3">
                        <div>
                            <span className="font-semibold">Ca làm việc:</span> {selectedDetail.work_shift?.workshift_name || "N/A"}
                        </div>

                        <div>
                            <span className="font-semibold">Trạng thái:</span>{" "}
                            {selectedDetail.is_absent ? (
                                <span className="text-red-600">Vắng mặt ({selectedDetail.absent_reason || "Không có lý do"})</span>
                            ) : selectedDetail.work_days > 0 ? (
                                <span className="text-green-600">Đi làm</span>
                            ) : selectedDetail.leave_type ? (
                                <span className="text-yellow-600">Nghỉ phép ({selectedDetail.leave_type})</span>
                            ) : (
                                <span>Không xác định</span>
                            )}
                        </div>

                        {selectedDetail.checkin_record && (
                            <div>
                                <span className="font-semibold">Check-in:</span>{" "}
                                {selectedDetail.checkin_record.timestamp.slice(11, 19)}
                                {selectedDetail.is_late && (
                                    <span className="text-orange-600 ml-2">(Trễ: {selectedDetail.late_minutes} phút)</span>
                                )}
                            </div>
                        )}

                        {selectedDetail.checkout_record && (
                            <div>
                                <span className="font-semibold">Check-out:</span>{" "}
                                {selectedDetail.checkout_record.timestamp.slice(11, 19)}
                            </div>
                        )}

                        <div>
                            <span className="font-semibold">Số giờ làm:</span> {selectedDetail.work_hours} giờ
                        </div>

                        <div>
                            <span className="font-semibold">Số công:</span> {selectedDetail.work_days}
                        </div>

                        {selectedDetail.is_holiday && (
                            <div className="text-purple-600">
                                <span className="font-semibold">Ngày lễ</span>
                            </div>
                        )}

                        {selectedDetail.is_weekend && (
                            <div className="text-blue-600">
                                <span className="font-semibold">Cuối tuần</span>
                            </div>
                        )}
                    </div>
                </DialogContent>
            </Dialog>
        );
    }, [selectedDetail, isDetailDialogOpen, isResetTimesheet, openEditDialog, handleReset]);

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <Loading />
            </div>
        );
    }

    if (!data || !timesheetData) {
        return (
            <div className="p-4 text-center text-red-500">
                Không tìm thấy dữ liệu bảng công
            </div>
        );
    }

    return (
        <div className="container mx-auto p-4">
            {/* Header Section */}
            <div className="mb-6 p-4 bg-white rounded shadow">
                <h1 className="text-2xl font-bold mb-2">{timesheetData.time_sheet_list_name}</h1>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <p><span className="font-semibold">Tháng/Năm:</span> {timesheetData.month}/{timesheetData.year}</p>
                        <p><span className="font-semibold">Khoảng thời gian:</span> {new Date(timesheetData.start_date).toLocaleDateString('vi-VN')} - {new Date(timesheetData.end_date).toLocaleDateString('vi-VN')}</p>
                    </div>
                    <div>
                        <p><span className="font-semibold">Trạng thái:</span>
                            <span className={`px-2 py-1 rounded ml-2 ${timesheetData.is_locked ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'}`}>
                                {timesheetData.is_locked ? 'Đã khóa' : 'Chưa khóa'}
                            </span>
                        </p>
                        <p><span className="font-semibold">Người tạo:</span> {timesheetData.created_by}</p>
                        <p><span className="font-semibold">Ngày tạo:</span> {new Date(timesheetData.created_at).toLocaleDateString('vi-VN')}</p>
                    </div>
                    <div className="flex gap-4 flex-wrap">
                        <Button onClick={handleCalculatorTimesheet}>
                            {isCalculatorTimesheet ? <Loading /> : "Đồng bộ"}
                        </Button>
                        <Button onClick={handleLockTimesheet}>
                            {isLockTimesheet ? <Loading/> : "Chốt công"}
                        </Button>
                        <Button onClick={handleExport}>
                            Xuất file
                        </Button>
                    </div>
                </div>
            </div>

            {/* Timesheet Table */}
            <div className="bg-white rounded shadow overflow-x-auto">
                <table className="min-w-full">
                    <thead>
                    <tr className="bg-gray-100">
                        <th className="p-2 border">#</th>
                        <th className="p-2 border">Mã NV</th>
                        <th className="p-2 border">Họ và tên</th>
                        <th className="p-2 border">Phòng ban</th>
                        <th className="p-2 border">Chức vụ</th>
                        <th className="p-2 border">Cấp bậc</th>
                        {/* Day columns */}
                        {dayColumns?.map((col, idx) => (
                            <th key={idx} className="p-1 border text-center">{col.header}</th>
                        ))}
                        <th className="p-2 border">Phút trễ</th>
                        <th className="p-2 border">Tổng công</th>
                    </tr>
                    </thead>
                    <tbody>
                    {timesheetData.timesheets?.map((timesheet) => (
                        <>
                            <tr key={timesheet.timesheet_id} className="hover:bg-gray-50">
                                <td className="p-2 border">
                                    {timesheet.details && timesheet.details.length > 0 && (
                                        <button
                                            onClick={() => toggleRowExpansion(timesheet.timesheet_id)}
                                            className="p-1 rounded hover:bg-gray-200"
                                        >
                                            {expandedRows.has(timesheet.timesheet_id) ? "▼" : "►"}
                                        </button>
                                    )}
                                </td>
                                <td className="p-2 border">{timesheet.employee.employee_id}</td>
                                <td className="p-2 border">{timesheet.employee.full_name}</td>
                                <td className="p-2 border">{timesheet.department.department_name}</td>
                                <td className="p-2 border">{timesheet.employee.position.position_name}</td>
                                <td className="p-2 border">{timesheet.employee.hierarchy_level.hierarchy_level}</td>
                                {/* Day cells */}
                                {dayColumns?.map((col, idx) => {
                                    const cell = col.cell!({row: {original: timesheet}});
                                    return (
                                        <td key={idx} className="p-1 border">
                                            {cell}
                                        </td>
                                    );
                                })}
                                <td className="p-2 border text-center">{timesheet.total_late_minutes}</td>
                                <td className="p-2 border text-center font-semibold">{timesheet.total_work_days}</td>
                            </tr>
                            {expandedRows.has(timesheet.timesheet_id) && (
                                <tr>
                                    <td colSpan={6 + dayColumns.length + 2} className="p-0">
                                        {renderRowExpansion(timesheet)}
                                    </td>
                                </tr>
                            )}
                        </>
                    ))}
                    </tbody>
                </table>
            </div>

            {/* Legend */}
            <div className="mt-4 p-3 bg-gray-100 rounded text-sm">
                <p className="font-semibold">Chú thích:</p>
                <div className="flex flex-wrap gap-4 mt-2">
                    <div className="flex items-center">
                        <div className="w-4 h-4 bg-green-100 mr-1 border"></div>
                        C: Đi làm
                    </div>
                    <div className="flex items-center">
                        <div className="w-4 h-4 bg-orange-100 mr-1 border"></div>
                        C: Đi làm nhưng trễ
                    </div>
                    <div className="flex items-center">
                        <div className="w-4 h-4 bg-red-100 mr-1 border"></div>
                        V: Vắng mặt
                    </div>
                    <div className="flex items-center">
                        <div className="w-4 h-4 bg-yellow-100 mr-1 border"></div>
                        P: Nghỉ phép
                    </div>
                    <div className="flex items-center">
                        <div className="w-4 h-4 bg-blue-100 mr-1 border"></div>
                        T7/CN: Thứ 7/Chủ nhật
                    </div>
                    <div className="flex items-center">
                        <div className="w-4 h-4 bg-purple-100 mr-1 border"></div>
                        L: Ngày lễ
                    </div>
                </div>
            </div>

            {/* Detail Dialog */}
            <DetailDialog/>

            <EditDialog
                open={isEditDialogOpen}
                onOpenChange={setIsEditDialogOpen}
                detail={selectedDetail}
                onSubmit={handleEditSubmit}
                isLoading={isUpdatingTimesheet}
            />
        </div>
    );
};

export default TimesheetDetailPage;