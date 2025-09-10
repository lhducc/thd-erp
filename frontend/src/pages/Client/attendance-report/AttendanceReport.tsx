import {useGetPersonalTimesheet} from "@/query/timesheet.query.ts";
import {useState} from "react";
import lateAttendance from "@/assets/late-attendance.svg"
import trueAttendance from "@/assets/true-attendance.svg"
import absentAttendance from "@/assets/absent-attendance.svg"
import {
    Dialog,
    DialogContent,
    DialogTitle,
} from "@/components/ui/dialog";

const AttendanceReport = () => {
    const {data: timesheetData} = useGetPersonalTimesheet(5, 2025);
    const [currentDate] = useState(new Date());
    const [selectedMonth, setSelectedMonth] = useState(currentDate.getMonth());
    const [selectedYear, setSelectedYear] = useState(currentDate.getFullYear());
    const [selectedDetail, setSelectedDetail] = useState(null);
    const [dialogOpen, setDialogOpen] = useState(false);

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


    const isDateInPast = (date: Date) => {
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        return date < today;
    };

    const getAttendanceDetail = (date: Date) => {
        if (!timesheetData?.details) return null;

        // Chuyển đổi date thành string với format YYYY-MM-DD (cân nhắc timezone)
        const dateStr = date.toLocaleDateString('en-CA'); // Format: YYYY-MM-DD

        return timesheetData.details.find(detail => {
            // Chuyển đổi detail.date thành Date object và format thành YYYY-MM-DD
            const detailDate = new Date(detail.date);
            const detailDateStr = detailDate.toLocaleDateString('en-CA');
            return detailDateStr === dateStr;
        });
    };

    // Function to format time for display
    const formatTime = (timeString: string | null) => {
        if (!timeString) return null;
        const date = new Date(timeString);
        return date.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
    };

    const getStatusAttendance = (detail: any) => {
        if (!detail) return 'Không có dữ liệu';
        if (detail.is_absent) return <img src={absentAttendance} alt="absent"/>;
        if (detail.is_late) return <img src={lateAttendance} alt="late"/>;
        return <img src={trueAttendance} alt="true"/>;
    };

    const AttendanceDialog = ({ detail }) => (
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
            <DialogContent className="bg-white p-6 rounded-lg max-w-md mx-auto">
                <DialogTitle className="text-xl font-bold mb-4">Chi tiết chấm công</DialogTitle>

                <div className="space-y-4">
                    <div>
                        <h3 className="font-medium">Trạng thái:</h3>
                        <p>
                            {detail?.is_absent ? 'Vắng mặt' :
                                detail?.is_late ? `Trễ ${detail.late_minutes} phút` : 'Đúng giờ'}
                        </p>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <h3 className="font-medium">Giờ vào:</h3>
                            <p>{formatTime(detail?.checkin_record?.timestamp) || '--:--'}</p>
                        </div>
                        <div>
                            <h3 className="font-medium">Giờ ra:</h3>
                            <p>{formatTime(detail?.checkout_record?.timestamp) || '--:--'}</p>
                        </div>
                    </div>

                    <div>
                        <h3 className="font-medium">Tổng giờ làm:</h3>
                        <p>{detail?.work_hours?.toFixed(2) || '0'} giờ</p>
                    </div>

                    <img src={detail.checkin_record.image_URL}/>
                    <img src={detail.checkout_record.image_URL}/>

                    {detail?.absent_reason && (
                        <div>
                            <h3 className="font-medium">Lý do vắng:</h3>
                            <p>{detail.absent_reason}</p>
                        </div>
                    )}
                </div>
            </DialogContent>
        </Dialog>
    );

    return (
        <div className="p-2 sm:p-4 max-w-4xl mx-auto">
            <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-gray-800">Bảng công cá nhân</h2>

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
                    <select
                        value={selectedMonth}
                        onChange={handleMonthChange}
                        className="border border-gray-300 rounded-md px-3 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        {months.map((month, index) => (
                            <option key={month} value={index}>{month}</option>
                        ))}
                    </select>
                    <select
                        value={selectedYear}
                        onChange={handleYearChange}
                        className="border border-gray-300 rounded-md px-3 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
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
                        const attendanceDetail = getAttendanceDetail(date);

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
                            >
                                <div className="text-right font-medium text-xs sm:text-sm">{day}</div>
                                <div className="mt-1 text-xs space-y-1 space-x-2">
                                    {attendanceDetail ? (
                                        <div>
                                            <div
                                                onClick={() => {
                                                setSelectedDetail(attendanceDetail);
                                                setDialogOpen(true);
                                            }} className={`p-1 rounded text-center w-full flex justify-center items-center cursor-pointer`}>
                                                {getStatusAttendance(attendanceDetail)}
                                            </div>
                                            <div className={`flex w-full justify-between items-center`}>

                                            {attendanceDetail.checkin_record && (
                                                <span className="text-sm text-red-500">
                                                    {formatTime(attendanceDetail.checkin_record.timestamp) || "--"}
                                                </span>
                                            ) || <p className={`text-red-500`}>--</p>}
                                            <span>-</span>
                                            {attendanceDetail.checkout_record && (
                                                <span className="text-sm text-red-500">
                                                    {formatTime(attendanceDetail.checkout_record.timestamp)}
                                                </span>
                                            ) || <p className={`text-red-500`}>--</p>}
                                            </div>
                                            {/*{attendanceDetail.work_hours > 0 && (*/}
                                            {/*    <div className="text-xs">*/}
                                            {/*        Giờ làm: {attendanceDetail.work_hours.toFixed(2)}h*/}
                                            {/*    </div>*/}
                                            {/*)}*/}
                                        </div>
                                    ) : (
                                        <div className="text-gray-400">Không có dữ liệu</div>
                                    )}
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>
            {selectedDetail && (
                <AttendanceDialog detail={selectedDetail} />
            )}
        </div>
    );
};

export default AttendanceReport;