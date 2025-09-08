import { useState } from 'react';
import {useQuery, useMutation, useQueryClient} from "@tanstack/react-query";
import {
    deleteEmployeeWorkshiftApi, getAllWorkshiftApi,
    getEmployeeWorkshiftsApi, getRegisterWorkshiftApi,
    registerEmployeeWorkshiftApi, updateEmployeeWorkshiftApi
} from "@/apis/workshift.api.ts";
import type { EmployeeWorkshift } from '@/types/employee-workshift';
import {toast} from "sonner";
import {useAuth} from "@/context/AuthContext.tsx";

const RegisterWorkshift = () => {
    const [currentDate] = useState(new Date());
    const [selectedMonth, setSelectedMonth] = useState(currentDate.getMonth());
    const [selectedYear, setSelectedYear] = useState(currentDate.getFullYear());
    const [selectedDate, setSelectedDate] = useState<Date | null>(null);
    const [selectedWorkshift, setSelectedWorkshift] = useState<string>('');
    const [showModal, setShowModal] = useState(false);
    const [editingWorkshiftId, setEditingWorkshiftId] = useState<string | null>(null);
    const queryClient = useQueryClient();
    // Vietnamese month names
    const months = [
        'Tháng 1', 'Tháng 2', 'Tháng 3', 'Tháng 4',
        'Tháng 5', 'Tháng 6', 'Tháng 7', 'Tháng 8',
        'Tháng 9', 'Tháng 10', 'Tháng 11', 'Tháng 12'
    ];

    // Generate years range (current year -5 to +5)
    const years = Array.from({ length: 11 }, (_, i) => currentDate.getFullYear() - 5 + i);

    // Get days in month and create calendar grid
    const daysInMonth = new Date(selectedYear, selectedMonth + 1, 0).getDate();
    const firstDayOfMonth = new Date(selectedYear, selectedMonth, 1).getDay();

    // Create calendar days array
    const calendarDays = [];

    // Add empty slots for days before the first day of month
    for (let i = 0; i < firstDayOfMonth; i++) {
        calendarDays.push(null);
    }

    // Add actual days of the month
    for (let day = 1; day <= daysInMonth; day++) {
        calendarDays.push(new Date(selectedYear, selectedMonth, day));
    }

    const handleMonthChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setSelectedMonth(parseInt(e.target.value));
    };

    const handleYearChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        setSelectedYear(parseInt(e.target.value));
    };

    const handlePrev = () => {
        if (selectedMonth === 0) {
            setSelectedMonth(11);
            setSelectedYear(selectedYear - 1);
        } else {
            setSelectedMonth(selectedMonth - 1);
        }
    };

    const handleNext = () => {
        if (selectedMonth === 11) {
            setSelectedMonth(0);
            setSelectedYear(selectedYear + 1);
        } else {
            setSelectedMonth(selectedMonth + 1);
        }
    };

    const { currentUser } = useAuth();

    const {
        data: employeeWorkshifts,
    } = useQuery({
        queryKey: ["employeeWorkshifts", currentUser?.user_id, selectedMonth, selectedYear],
        queryFn: () => getEmployeeWorkshiftsApi(selectedMonth + 1, selectedYear),
    });

    const {
        data: allWorkshifts,
    } = useQuery({
        queryKey: ["register-workshifts"],
        queryFn: getRegisterWorkshiftApi,
    });


    const filteredWorkshifts = allWorkshifts
        ?.filter(ws => ws.work_shift !== null) // lọc ca hợp lệ
        ?.map(ws => ({
            id: ws.workshift_id,
            name: ws.work_shift.workshift_name,
            start: ws.work_shift.start_time,
            end: ws.work_shift.end_time
        }));


    // Mutations
    const registerMutation = useMutation({
        mutationFn: ({ employee_id, work_shift_id, date }:
                     { employee_id: string; work_shift_id: string; date: string }) =>
            registerEmployeeWorkshiftApi(employee_id, work_shift_id, date),
        onSuccess: () => {
            toast.success('Đăng ký ca làm việc thành công!');
            queryClient.invalidateQueries({ queryKey: ["employeeWorkshifts"] });
            setShowModal(false);
        },
        onError: (error) => {
            toast.error(`Đăng ký ca làm việc thất bại: ${error.message}`);
        }
    });

    const deleteMutation = useMutation({
        mutationFn: (id: number) => deleteEmployeeWorkshiftApi(id),
        onSuccess: () => {
            toast.success('Xóa ca làm việc thành công!');
            queryClient.invalidateQueries({ queryKey: ["employeeWorkshifts"] });
        },
        onError: (error) => {
            toast.error(`Xóa ca làm việc thất bại: ${error.message}`);
        }
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, work_shift_id }: { id: string; work_shift_id: string }) =>
            updateEmployeeWorkshiftApi(id, work_shift_id),
        onSuccess: () => {
            toast.success('Cập nhật ca làm việc thành công!');
            queryClient.invalidateQueries({ queryKey: ["employeeWorkshifts"] });
            setShowModal(false);
        },
        onError: (error) => {
            toast.error(`Cập nhật ca làm việc thất bại: ${error.message}`);
        }
    });

    // Handlers
    const handleDateClick = (date: Date) => {
        const today = new Date();
        today.setHours(0, 0, 0, 0);

        if (date >= today) {
            setSelectedDate(date);
            const existingShift = getWorkshiftForDate(date);
            if (existingShift) {
                setEditingWorkshiftId(existingShift.ID);
                setSelectedWorkshift(existingShift.work_shift_id);
            } else {
                setEditingWorkshiftId(null);
                setSelectedWorkshift('');
            }
            setShowModal(true);
        }
    };

    const handleRegister = () => {
        if (!selectedDate || !selectedWorkshift || !currentUser?.user_id) return;

        const dateStr = selectedDate.toISOString().split('T')[0];

        if (editingWorkshiftId) {
            updateMutation.mutate({
                id: editingWorkshiftId.toString(),
                work_shift_id: selectedWorkshift
            });
        } else {
            registerMutation.mutate({
                employee_id: currentUser.user_id,
                work_shift_id: selectedWorkshift,
                date: dateStr
            });
        }
    };

    const handleDelete = (workshiftId: number) => {
        if (window.confirm('Bạn có chắc chắn muốn xóa ca làm việc này?')) {
            deleteMutation.mutate(workshiftId);
        }
    };

    // Helper functions
    const getWorkshiftForDate = (date: Date): EmployeeWorkshift | undefined => {
        if (!employeeWorkshifts) return undefined;
        const dateStr = date.toISOString().split('T')[0];
        return employeeWorkshifts.find(ws =>
            new Date(ws.date).toISOString().split('T')[0] === dateStr
        );
    };

    const isDateInPast = (date: Date) => {
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        return date < today;
    };

    const getWorkshiftName = (workshiftId: string) => {
        return allWorkshifts?.find(ws => ws.workshift_id === workshiftId)?.workshift_name || 'Unknown';
    };

    return (
        <div className="p-2 sm:p-4 max-w-4xl mx-auto">
            <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-gray-800">Lịch làm việc</h2>

            {/* Month/Year Selector */}
            <div className="flex flex-col sm:flex-row justify-between items-center bg-gray-100 p-2 sm:p-3 rounded-lg mb-4 sm:mb-6 shadow-sm gap-2 sm:gap-0">
                <button
                    onClick={handlePrev}
                    className="w-full sm:w-auto bg-gray-800 text-white rounded-md sm:rounded-l-md sm:rounded-r-none border-r border-gray-100 py-2 hover:bg-red-700 hover:text-white px-3"
                >
                    <div className="flex flex-row align-middle justify-center sm:justify-start">
                        <svg className="w-5 mr-2" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                            <path fillRule="evenodd" d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z" clipRule="evenodd"></path>
                        </svg>
                        <p className="ml-2">Trước</p>
                    </div>
                </button>

                {/* Month/Year selectors */}
                <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
                    <select value={selectedMonth} onChange={handleMonthChange} className="...">
                        {months.map((month, index) => (
                            <option key={month} value={index}>{month}</option>
                        ))}
                    </select>
                    <select value={selectedYear} onChange={handleYearChange} className="...">
                        {years.map(year => (
                            <option key={year} value={year}>{year}</option>
                        ))}
                    </select>
                </div>

                <button
                    onClick={handleNext}
                    className="w-full sm:w-auto bg-gray-800 text-white rounded-md sm:rounded-r-md sm:rounded-l-none py-2 border-l border-gray-200 hover:bg-red-700 hover:text-white px-3"
                >
                    <div className="flex flex-row align-middle justify-center sm:justify-start">
                        <span className="mr-2">Sau</span>
                        <svg className="w-5 ml-2" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                            <path fillRule="evenodd" d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z" clipRule="evenodd"></path>
                        </svg>
                    </div>
                </button>
            </div>

            {/* Calendar Grid */}
            <div className="bg-white rounded-lg shadow overflow-hidden border">
                {/* Weekday Headers */}
                <div className="grid grid-cols-7 bg-gray-100">
                    {['CN', 'T2', 'T3', 'T4', 'T5', 'T6', 'T7'].map(day => (
                        <div key={day} className="py-1 sm:py-2 text-center font-medium text-gray-600 text-xs sm:text-sm">
                            {day}
                        </div>
                    ))}
                </div>

                {/* Calendar Days */}
                <div className="grid grid-cols-7 gap-px bg-gray-200">
                    {calendarDays.map((date, index) => {
                        if (!date) {
                            return <div key={index} className="min-h-[60px] sm:min-h-[80px] md:min-h-[100px] bg-gray-100" />;
                        }

                        const day = date.getDate();
                        const isToday = date.toDateString() === new Date().toDateString();
                        const isPast = isDateInPast(date);
                        const workshift = getWorkshiftForDate(date);
                        return (
                            <div
                                key={index}
                                className={`min-h-[60px] sm:min-h-[80px] md:min-h-[100px] p-1 sm:p-2 ${
                                    isToday
                                        ? 'border-2 border-red-500 bg-gray-300'
                                        : isPast
                                            ? 'bg-gray-100'
                                            : 'bg-white hover:bg-gray-50 cursor-pointer'
                                }`}
                                onClick={() => !isPast && handleDateClick(date)}
                            >
                                <div className="text-right font-medium text-xs sm:text-sm">{day}</div>
                                <div className="mt-1 text-xs">
                                    {workshift ? (
                                        <div className="bg-blue-100 text-blue-800 p-1 rounded text-center">
                                            {getWorkshiftName(workshift.workshift_id)}
                                            {!isPast && (
                                                <button
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        handleDelete(workshift?.id);
                                                    }}
                                                    className="text-red-500 text-xs ml-1"
                                                >
                                                    ×
                                                </button>
                                            )}
                                        </div>
                                    ) : isPast ? (
                                        <div className="text-gray-400">Không có ca</div>
                                    ) : (
                                        <div className="text-gray-500">Chọn ca</div>
                                    )}
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            {/* Registration/Edit Modal */}
            {showModal && selectedDate && (
                <div className="fixed inset-0 bg-black/70 bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 w-full max-w-md">
                        <h3 className="text-lg font-semibold mb-4">
                            {editingWorkshiftId ? 'Cập nhật' : 'Đăng ký'} ca làm việc ngày {selectedDate.toLocaleDateString('vi-VN')}
                        </h3>

                        <div className="mb-4">
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Chọn ca làm việc
                            </label>
                            <select
                                value={selectedWorkshift}
                                onChange={(e) => setSelectedWorkshift(e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                                disabled={registerMutation.isPending || updateMutation.isPending}
                            >
                                <option value="">-- Chọn ca --</option>
                                {filteredWorkshifts?.map((ws, index) => (
                                    <option key={`${ws.code}-${index}`} value={ws.code}>
                                        {ws.name}
                                    </option>
                                ))}

                            </select>
                        </div>

                        <div className="flex justify-end gap-2">
                            <button
                                onClick={() => setShowModal(false)}
                                className="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400"
                                disabled={registerMutation.isPending || updateMutation.isPending}
                            >
                                Hủy
                            </button>
                            <button
                                onClick={handleRegister}
                                disabled={!selectedWorkshift || registerMutation.isPending || updateMutation.isPending}
                                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-blue-300"
                            >
                                {registerMutation.isPending || updateMutation.isPending
                                    ? 'Đang xử lý...'
                                    : editingWorkshiftId ? 'Cập nhật' : 'Đăng ký'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default RegisterWorkshift;