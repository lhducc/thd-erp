import {Button} from "@/components/ui/button.tsx";
import {useState, useCallback, useMemo, useRef, useEffect} from "react";
import {X} from "lucide-react";
import {useQuery} from "@tanstack/react-query";
import {type Workshift} from "@/types/workshift.ts";
import {getAllWorkshiftApi} from "@/apis/workshift.api.ts";
import {formatTime} from "@/lib/utils.ts";
import Loading from "@/components/Loading.tsx";
import {DragDropContext, Droppable, Draggable, type DropResult} from "@hello-pangea/dnd";
import {getAllShiftAllocations} from "@/apis/shift-allocation.api.ts";
import {SecondNav} from "@/pages/HR/checkin/rota/_components/SecondNav.tsx";
import {isToday} from "date-fns";

const Nav = () => (
    <div className="flex justify-between items-center">
        <h3 className="font-semibold text-3xl">Bảng phân ca</h3>
        <div className="flex gap-4">
            <Button>Phân ca</Button>
            <Button>Xuất file</Button>
        </div>
    </div>
);

const Rota = () => {
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);
    const [currentDate, setCurrentDate] = useState(new Date());
    const [employeeShifts, setEmployeeShifts] = useState<Record<string, Workshift[]>>({});
    const [selectedSchedule, setSelectedSchedule] = useState<string | null>(null);
    const [isDraggingColumns, setIsDraggingColumns] = useState(false);
    const [startX, setStartX] = useState(0);
    const [scrollLeft, setScrollLeft] = useState(0);
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleMouseUp = () => setIsDraggingColumns(false);
        window.addEventListener("mouseup", handleMouseUp);
        return () => window.removeEventListener("mouseup", handleMouseUp);
    }, []);

    const currentMonth = currentDate.getMonth() + 1;
    const currentYear = currentDate.getFullYear();

    const daysInMonth = useMemo(() => {
        return new Date(currentYear, currentMonth, 0).getDate();
    }, [currentYear, currentMonth]);

    const monthDays = useMemo(() => {
        return Array.from({length: daysInMonth}, (_, i) => {
            const day = new Date(currentYear, currentMonth - 1, i + 1);
            return day;
        });
    }, [currentYear, currentMonth, daysInMonth]);

    const handleMouseDown = (e: React.MouseEvent) => {
        if (!containerRef.current) return;
        setIsDraggingColumns(true);
        setStartX(e.pageX - containerRef.current.offsetLeft);
        setScrollLeft(containerRef.current.scrollLeft);
    };

    const handleMouseMove = (e: React.MouseEvent) => {
        if (!isDraggingColumns || !containerRef.current) return;
        e.preventDefault();
        const x = e.pageX - containerRef.current.offsetLeft;
        const walk = (x - startX) * 2;
        containerRef.current.scrollLeft = scrollLeft - walk;
    };

    const handleMouseUpOrLeave = () => {
        setIsDraggingColumns(false);
    };

    const {data: shiftAllocation = [], isPending: isShiftAllocationPending} =
        useQuery({
            queryKey: ["shift-allocation", currentMonth, currentYear],
            queryFn: () => getAllShiftAllocations(currentMonth, currentYear),
        });

    const scheduleNames = useMemo(() => {
        const names = new Set<string>();
        shiftAllocation.forEach(employee => {
            employee.schedules?.forEach(schedule => {
                if (schedule.work_schedule_name) {
                    names.add(schedule.work_schedule_name);
                }
            });
        });
        return Array.from(names);
    }, [shiftAllocation]);

    const filteredEmployees = useMemo(() => {
        if (!selectedSchedule) return shiftAllocation;

        return shiftAllocation.filter(employee =>
            employee.schedules?.some(schedule =>
                schedule.work_schedule_name === selectedSchedule
            )
        );
    }, [shiftAllocation, selectedSchedule]);

    const {data: workshifts = [], isLoading, isError, error} = useQuery<Workshift[]>({
        queryKey: ['workshifts'],
        queryFn: getAllWorkshiftApi,
        enabled: isSidebarOpen,
        staleTime: 5 * 60 * 1000,
    });

    const getDayName = (date: Date): string => {
        return date.toLocaleDateString('en-US', {weekday: 'long'}).toLowerCase();
    };

    const handleDragEnd = useCallback((result: DropResult) => {
        if (!result.destination) return;

        const {source, destination} = result;

        if (source.droppableId === 'workshift-list' && destination.droppableId.includes('-')) {
            const [employeeId, dayIndex] = destination.droppableId.split('-');
            const workshift = workshifts[source.index];

            setEmployeeShifts(prev => {
                const newShifts = {...prev};
                const dayKey = `${employeeId}-${dayIndex}`;

                if (!newShifts[dayKey]) {
                    newShifts[dayKey] = [];
                }

                newShifts[dayKey] = [...newShifts[dayKey], {
                    ...workshift,
                    tempId: `${workshift.workshift_id}-${Date.now()}`
                }];

                return newShifts;
            });
        }

        if (source.droppableId.includes('-') &&
            destination.droppableId.includes('-') &&
            source.droppableId === destination.droppableId) {

            const [employeeId, dayIndex] = source.droppableId.split('-');
            const dayKey = `${employeeId}-${dayIndex}`;

            setEmployeeShifts(prev => {
                const newShifts = {...prev};
                const shiftsCopy = [...newShifts[dayKey]];
                const [removed] = shiftsCopy.splice(source.index, 1);
                shiftsCopy.splice(destination.index, 0, removed);

                newShifts[dayKey] = shiftsCopy;
                return newShifts;
            });
        }
    }, [workshifts]);

    const handleScheduleFilter = (scheduleName: string) => {
        setSelectedSchedule(prev => prev === scheduleName ? null : scheduleName);
    };

    const getEmployeeShift = (employeeId: string, dayIndex: number) => {
        const dayKey = `${employeeId}-${dayIndex}`;
        return employeeShifts[dayKey] || [];
    };

    const getScheduledShift = (employeeId: string, dayIndex: number) => {
        const employee = shiftAllocation.find(e => e.employee_id === employeeId);
        if (!employee || !employee.schedules || employee.schedules.length === 0) return null;

        const dayName = getDayName(monthDays[dayIndex]);
        for (const schedule of employee.schedules) {
            if (selectedSchedule && schedule.work_schedule_name !== selectedSchedule) {
                continue;
            }

            const weekdayShift = schedule.weekdays?.find(w => w.week_day === dayName);
            if (weekdayShift) {
                return {
                    ...weekdayShift.work_shift,
                    scheduleName: schedule.work_schedule_name
                };
            }
        }
        return null;
    };

    const handleDateChange = useCallback(
        ({ month, year }: { month: number; year: number }) => {
            const newDate = new Date(year, month - 1, 1);
            setCurrentDate(newDate);
        },
        []
    );

    return (
        <div className="relative">
            <div className="space-y-4">
                <Nav/>
                <hr/>
                <SecondNav
                    currentDate={currentDate}
                    scheduleNames={scheduleNames}
                    selectedSchedule={selectedSchedule}
                    onDateChange={handleDateChange}
                    onToggleSidebar={() => setIsSidebarOpen(!isSidebarOpen)}
                    onSelectSchedule={handleScheduleFilter}
                    isSidebarOpen={isSidebarOpen}
                />
            </div>

            <DragDropContext onDragEnd={handleDragEnd}>
                {/* Week Calendar */}
                <div
                    ref={containerRef}
                    className="mt-6 overflow-x-auto w-[80vw]"
                    onMouseDown={handleMouseDown}
                    onMouseMove={handleMouseMove}
                    onMouseUp={handleMouseUpOrLeave}
                    onMouseLeave={handleMouseUpOrLeave}
                    style={{ cursor: isDraggingColumns ? 'grabbing' : 'grab' }}
                >
                    <div className="min-w-max">
                        <div
                            className="grid gap-1"
                            style={{ gridTemplateColumns: `200px repeat(${daysInMonth}, minmax(120px, 1fr))` }}
                        >

                        <div className="p-2 font-medium text-center sticky left-0 bg-white z-10 min-w-[120px]">Nhân viên</div>
                            {monthDays.map((day, i) => (
                                <div
                                    key={i}
                                    className={`p-2 text-center font-medium border sticky top-0 min-w-[120px] ${
                                        day.getDay() === 0 || day.getDay() === 6 ? 'bg-red-50' : ''
                                    } ${
                                        isToday(day) ? 'border-2 border-blue-500 bg-blue-50' : ''
                                    }`}
                                >
                                    {day.toLocaleDateString('vi-VN', {weekday: 'short'})}
                                    <div className="text-sm">
                                        {day.getDate()}/{day.getMonth() + 1}
                                        {isToday(day) && <div className="w-2 h-2 bg-blue-500 rounded-full mx-auto mt-1"></div>}
                                    </div>
                                </div>
                            ))}
                        </div>

                        {/* Employee Rows */}
                        {isShiftAllocationPending ? (
                            <div className="flex justify-center items-center h-64">
                                <Loading/>
                            </div>
                        ) : (
                            filteredEmployees.map((employee) => (
                                <div key={employee.employee_id}
                                className="grid gap-1"
                                style={{ gridTemplateColumns: `200px repeat(${daysInMonth}, minmax(120px, 1fr))` }}
                            >
                                    <div className="p-2 border-r sticky left-0 bg-white z-10 border flex flex-col items-center justify-center">
                                        <p className="font-medium">{employee.full_name}</p>
                                        <p className="text-sm text-gray-500">{employee.department_name}</p>
                                    </div>

                                    {monthDays.map((day, dayIndex) => {
                                        const dayKey = `${employee.employee_id}-${dayIndex}`;
                                        const scheduledShift = getScheduledShift(employee.employee_id, dayIndex);
                                        const assignedShifts = getEmployeeShift(employee.employee_id, dayIndex);

                                        return (
                                            <Droppable key={employee.employee_id + '-' + dayIndex} droppableId={dayKey}>
                                                {(provided) => (
                                                    <div
                                                        ref={provided.innerRef}
                                                        {...provided.droppableProps}
                                                        className={`p-2 min-h-20 border bg-gray-50 min-w-[120px] ${
                                                            day.getDay() === 0 || day.getDay() === 6 ? 'bg-red-50' : ''
                                                        }`}
                                                    >
                                                        {scheduledShift && (
                                                            <div className="mb-2 p-2 bg-blue-100 rounded">
                                                                <p className="text-sm font-medium">{scheduledShift.workshift_name}</p>
                                                                <p className="text-xs">
                                                                    {formatTime(scheduledShift.start_time)} - {formatTime(scheduledShift.end_time)}
                                                                </p>
                                                                {scheduledShift.scheduleName && (
                                                                    <p className="text-xs mt-1 text-blue-700">
                                                                        {scheduledShift.scheduleName}
                                                                    </p>
                                                                )}
                                                            </div>
                                                        )}

                                                        {assignedShifts.map((shift, index) => (
                                                            <Draggable
                                                                key={shift.tempId || shift.workshift_id}
                                                                draggableId={(shift.tempId || shift.workshift_id).toString()}
                                                                index={index}
                                                            >
                                                            {(provided) => (
                                                                    <div
                                                                        ref={provided.innerRef}
                                                                        {...provided.draggableProps}
                                                                        {...provided.dragHandleProps}
                                                                        className="mb-2 p-2 bg-green-100 rounded cursor-move"
                                                                    >
                                                                        <p className="text-sm font-medium">{shift.workshift_name}</p>
                                                                        <p className="text-xs">
                                                                            {formatTime(shift.start_time)} - {formatTime(shift.end_time)}
                                                                        </p>
                                                                    </div>
                                                                )}
                                                            </Draggable>
                                                        ))}
                                                        {provided.placeholder}
                                                    </div>
                                                )}
                                            </Droppable>
                                        );
                                    })}
                                </div>
                            ))
                        )}
                    </div>
                </div>

                {/* Sidebar */}
                <div className={`
                    fixed top-0 right-0 h-full w-96 bg-white shadow-lg z-50
                    transform transition-transform duration-300 ease-in-out
                    ${isSidebarOpen ? 'translate-x-0' : 'translate-x-full'}
                `}>
                    <div className="p-4 h-full flex flex-col">
                        <div className="flex justify-between items-center mb-6">
                            <h2 className="text-xl font-bold">Danh sách ca làm việc</h2>
                            <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setIsSidebarOpen(false)}
                            >
                                <X className="h-5 w-5"/>
                            </Button>
                        </div>

                        {isLoading ? (
                            <div className="flex-1 flex items-center justify-center">
                                <Loading/>
                            </div>
                        ) : isError ? (
                            <div className="flex-1 flex items-center justify-center">
                                <p className="text-red-500">Lỗi khi tải dữ liệu: {(error as Error).message}</p>
                            </div>
                        ) : (
                            <Droppable droppableId="workshift-list" isDropDisabled={true}>
                                {(provided) => (
                                    <div
                                        ref={provided.innerRef}
                                        {...provided.droppableProps}
                                        className="flex-1 overflow-y-auto"
                                    >
                                        {workshifts.map((shift, index) => (
                                            <Draggable
                                                key={shift.tempId || shift.workshift_id}
                                                draggableId={(shift.tempId || shift.workshift_id).toString()}
                                                index={index}
                                            >

                                            {(provided) => (
                                                    <div
                                                        ref={provided.innerRef}
                                                        {...provided.draggableProps}
                                                        {...provided.dragHandleProps}
                                                        className="mb-2 p-2 border rounded-lg cursor-move"
                                                    >
                                                        <h3 className="font-medium">{shift.workshift_name}: {formatTime(shift.start_time)} - {formatTime(shift.end_time)}</h3>
                                                    </div>
                                                )}
                                            </Draggable>
                                        ))}
                                        {provided.placeholder}
                                    </div>
                                )}
                            </Droppable>
                        )}
                    </div>
                </div>
            </DragDropContext>

            {/* Overlay */}
            {isSidebarOpen && (
                <div
                    className="fixed inset-0 z-40"
                    onClick={() => setIsSidebarOpen(false)}
                />
            )}
        </div>
    );
}

export default Rota;