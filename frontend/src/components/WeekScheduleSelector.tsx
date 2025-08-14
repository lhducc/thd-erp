import { useState, useEffect } from "react";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from "@/components/ui/card";
import { Clock, CalendarDays } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { formatTime } from "@/lib/utils.ts";
import type { WeekdaySelection } from "@/types/work-schedule.ts";
import type { Workshift } from "@/types/workshift.ts";

type WorkshiftSchedulerProps = {
    workshifts: Workshift[];
    onSelectionChange: (selections: WeekdaySelection[]) => void;
    initialSelections?: WeekdaySelection[];
};

export const WorkshiftScheduler = ({
                                       workshifts,
                                       onSelectionChange,
                                       initialSelections = []
                                   }: WorkshiftSchedulerProps) => {
    const days = [
        { id: 'mon', name: 'Thứ 2', week_day: 'monday' },
        { id: 'tue', name: 'Thứ 3', week_day: 'tuesday' },
        { id: 'wed', name: 'Thứ 4', week_day: 'wednesday' },
        { id: 'thu', name: 'Thứ 5', week_day: 'thursday' },
        { id: 'fri', name: 'Thứ 6', week_day: 'friday' },
        { id: 'sat', name: 'Thứ 7', week_day: 'saturday' },
        { id: 'sun', name: 'Chủ nhật', week_day: 'sunday' }
    ];

    const [selectedDays, setSelectedDays] = useState<string[]>([]);
    const [selectedShifts, setSelectedShifts] = useState<string[]>([]);

    // Initialize with default values
    useEffect(() => {
        if (initialSelections.length > 0) {
            const initialDays = [...new Set(initialSelections.map(sel => sel.week_day))];
            setSelectedDays(initialDays);

            const initialShifts = [...new Set(initialSelections.map(sel => sel.workshift_id))];
            setSelectedShifts(initialShifts);
        }
    }, [initialSelections]);

    // Handle day selection
    const handleDayToggle = (weekDay: string) => {
        const newSelectedDays = selectedDays.includes(weekDay)
            ? selectedDays.filter(day => day !== weekDay)
            : [...selectedDays, weekDay];

        setSelectedDays(newSelectedDays);
        updateSelections(newSelectedDays, selectedShifts);
    };

    // Handle shift selection
    const handleShiftSelect = (shiftId: string) => {
        const newSelectedShifts = selectedShifts.includes(shiftId)
            ? selectedShifts.filter(id => id !== shiftId)
            : [...selectedShifts, shiftId];

        setSelectedShifts(newSelectedShifts);
        updateSelections(selectedDays, newSelectedShifts);
    };

    // Update parent component with selections
    const updateSelections = (days: string[], shifts: string[]) => {
        const selections: WeekdaySelection[] = days.flatMap(day =>
            shifts.map(shiftId => ({
                week_day: day,
                workshift_id: shiftId
            }))
        );

        onSelectionChange(selections);
    };

    return (
        <div className="flex flex-col lg:flex-row gap-6 p-4 bg-gray-50 rounded-lg">
            {/* Days Column */}
            <Card className="lg:w-1/4">
                <CardHeader>
                    <div className="flex items-center gap-2">
                        <CalendarDays className="w-5 h-5 text-blue-600" />
                        <CardTitle>Ngày làm việc</CardTitle>
                    </div>
                    <CardDescription>Chọn ngày cần áp dụng lịch</CardDescription>
                </CardHeader>
                <CardContent className="overflow-y-auto">
                    <div className="space-y-3 p-1">
                        {days.map(day => (
                            <div
                                key={day.id}
                                className="flex items-center space-x-3 p-3 hover:bg-gray-100 rounded-lg transition-colors"
                            >
                                <Checkbox
                                    id={day.id}
                                    checked={selectedDays.includes(day.week_day)}
                                    onCheckedChange={() => handleDayToggle(day.week_day)}
                                />
                                <Label htmlFor={day.id} className="text-base cursor-pointer">
                                    {day.name}
                                </Label>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>

            {/* Workshifts Column */}
            <div className="lg:w-3/4 space-y-4">
                <Card>
                    <CardHeader>
                        <div className="flex items-center gap-2">
                            <Clock className="w-5 h-5 text-blue-600" />
                            <CardTitle>Ca làm việc</CardTitle>
                        </div>
                        <CardDescription>
                            {selectedDays.length > 0
                                ? "Chọn ca làm việc sẽ áp dụng cho tất cả ngày đã chọn"
                                : "Vui lòng chọn ít nhất một ngày làm việc"}
                        </CardDescription>
                    </CardHeader>
                </Card>

                <div className="overflow-y-auto border rounded-lg" style={{ maxHeight: '400px' }}>
                    {selectedDays.length > 0 ? (
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4">
                            {workshifts?.map(shift => (
                                <Card
                                    key={shift.workshift_id}
                                    className={`hover:shadow-md transition-shadow ${
                                        selectedShifts.includes(shift.workshift_id)
                                            ? "border-primary border-2"
                                            : ""
                                    }`}
                                >
                                    <CardContent className="p-4">
                                        <div className="flex items-start space-x-3">
                                            <Checkbox
                                                id={`shift-${shift.workshift_id}`}
                                                checked={selectedShifts.includes(shift.workshift_id)}
                                                onCheckedChange={() => handleShiftSelect(shift.workshift_id)}
                                                className="mt-1"
                                            />
                                            <div className="flex-1">
                                                <div className="flex items-center justify-between">
                                                    <Label htmlFor={`shift-${shift.workshift_id}`} className="text-base font-medium">
                                                        {shift.workshift_name}
                                                    </Label>
                                                    <Badge variant="outline">{shift.time_of_day}</Badge>
                                                </div>
                                                <div className="mt-2 space-y-1 text-sm text-gray-600">
                                                    <div className="flex">
                                                        <span className="w-24">Thời gian:</span>
                                                        <span>
                                                            {formatTime(shift.start_time)} - {formatTime(shift.end_time)}
                                                        </span>
                                                    </div>
                                                    <div className="flex">
                                                        <span className="w-24">Check-in:</span>
                                                        <span>
                                                            {formatTime(shift.checkin_from)} - {formatTime(shift.checkin_to)}
                                                        </span>
                                                    </div>
                                                    <div className="flex">
                                                        <span className="w-24">Check-out:</span>
                                                        <span>
                                                            {formatTime(shift.checkout_from)} - {formatTime(shift.checkout_to)}
                                                        </span>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))}
                        </div>
                    ) : (
                        <div className="flex items-center justify-center h-40 text-gray-500">
                            Vui lòng chọn ít nhất một ngày làm việc
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default WorkshiftScheduler;