import { useState, useEffect } from 'react';

interface MonthYearPickerProps {
    initialMonth?: number;
    initialYear?: number;
    onChange?: (date: { month: number; year: number }) => void;
    minYear?: number;
    maxYear?: number;
}

const MonthYearPicker: React.FC<MonthYearPickerProps> = ({
                             initialMonth,
                             initialYear,
                             onChange,
                             minYear = 1900,
                             maxYear = new Date().getFullYear() + 10
                         }) => {
    const [month, setMonth] = useState(initialMonth || new Date().getMonth() + 1);
    const [year, setYear] = useState(initialYear || new Date().getFullYear());

    const months = [
        'Tháng 1', 'Tháng 2', 'Tháng 3', 'Tháng 4', 'Tháng 5', 'Tháng 6',
        'Tháng 7', 'Tháng 8', 'Tháng 9', 'Tháng 10', 'Tháng 11', 'Tháng 12'
    ];

    useEffect(() => {
        if (onChange) {
            onChange({ month, year });
        }
    }, [month, year, onChange]);

    const handleMonthChange = (e) => {
        setMonth(parseInt(e.target.value));
    };

    const handleYearChange = (e) => {
        setYear(parseInt(e.target.value));
    };

    const generateYears = () => {
        const years = [];
        for (let y = minYear; y <= maxYear; y++) {
            years.push(y);
        }
        return years;
    };

    return (
        <div className="month-year-picker" style={{ display: 'flex', gap: '10px' }}>
            <select
                value={month}
                onChange={handleMonthChange}
                className="month-select"
                style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ccc' }}
            >
                {months.map((m, index) => (
                    <option key={m} value={index + 1}>
                        {m}
                    </option>
                ))}
            </select>

            <select
                value={year}
                onChange={handleYearChange}
                className="year-select"
                style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ccc' }}
            >
                {generateYears().map((y) => (
                    <option key={y} value={y}>
                        {y}
                    </option>
                ))}
            </select>
        </div>
    );
};

export default MonthYearPicker;