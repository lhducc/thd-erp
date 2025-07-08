import React, { useState } from 'react';
import SelectDay from "@/components/attendance/SelectDay.tsx";

const RegisterWorkshift = () => {
    const [currentDate] = useState(new Date());
    const [selectedMonth, setSelectedMonth] = useState(currentDate.getMonth());
    const [selectedYear, setSelectedYear] = useState(currentDate.getFullYear());

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
        calendarDays.push(day);
    }

    const handleMonthChange = (e) => {
        setSelectedMonth(parseInt(e.target.value));
    };

    const handleYearChange = (e) => {
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

    return (
        <div className="p-2 sm:p-4 max-w-4xl mx-auto">
            <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-gray-800">Lịch làm việc</h2>

            {/* Month/Year Selector - Responsive */}
            <div className="flex flex-col sm:flex-row justify-between items-center bg-gray-100 p-2 sm:p-3 rounded-lg mb-4 sm:mb-6 shadow-sm gap-2 sm:gap-0">
                {/* Prev Button - Full width on mobile, normal on desktop */}
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

                {/* Month/Year Selectors - Stack on mobile, row on desktop */}
                <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
                    <select
                        value={selectedMonth}
                        onChange={handleMonthChange}
                        className="px-3 py-2 border border-gray-300 bg-white rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm sm:text-base"
                    >
                        {months.map((month, index) => (
                            <option key={month} value={index}>{month}</option>
                        ))}
                    </select>

                    <select
                        value={selectedYear}
                        onChange={handleYearChange}
                        className="px-3 py-2 border border-gray-300 bg-white rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm sm:text-base"
                    >
                        {years.map(year => (
                            <option key={year} value={year}>{year}</option>
                        ))}
                    </select>
                </div>

                {/* Next Button - Full width on mobile, normal on desktop */}
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

            {/* Calendar Grid - Responsive */}
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
                    {calendarDays.map((day, index) => (
                        <div
                            key={index}
                            className={`min-h-[60px] sm:min-h-[80px] md:min-h-[100px] p-1 sm:p-2 ${
                                day === currentDate.getDate() &&
                                selectedMonth === currentDate.getMonth() &&
                                selectedYear === currentDate.getFullYear()
                                    ? 'border-2 border-red-500 bg-gray-300'
                                    : 'bg-white'
                            }`}
                        >
                            {day && (
                                <>
                                    <div className="text-right font-medium text-xs sm:text-sm">{day}</div>
                                    <div className="mt-1 text-xs text-gray-500">
                                        <SelectDay />
                                    </div>
                                </>
                            )}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default RegisterWorkshift;