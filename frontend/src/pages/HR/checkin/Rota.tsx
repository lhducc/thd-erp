import {Button} from "@/components/ui/button";
import MonthYearPicker from "@/components/MonthYearPicker";
import {useState, useCallback, useMemo} from "react";
import {Menu, X} from "lucide-react";
import {useQuery} from "@tanstack/react-query";
import {type WorkShift} from "@/types/Workshift";
import {getAllWorkshiftApi} from "@/apis/workshift.api";
import {formatTime} from "@/lib/utils";
import Loading from "@/components/Loading";
import {DragDropContext, Droppable, Draggable, type DropResult} from "@hello-pangea/dnd";

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
    const [tags] = useState(["1", "2", "3"]);
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);
    const [currentDate, setCurrentDate] = useState(new Date());
    const [employeeShifts, setEmployeeShifts] = useState<Record<string, WorkShift[]>>({
        'emp1': [],
        'emp2': [],
        'emp3': [],
    });

    // Memoize weekDays calculation
    const weekDays = useMemo(() => {
        const startDate = new Date(currentDate);
        startDate.setDate(currentDate.getDate() - currentDate.getDay());
        return Array.from({length: 7}).map((_, i) => {
            const day = new Date(startDate);
            day.setDate(startDate.getDate() + i);
            return day;
        });
    }, [currentDate]);

    const {data: workshifts = [], isLoading, isError, error} = useQuery<WorkShift[]>({
        queryKey: ['workshifts'],
        queryFn: getAllWorkshiftApi,
        enabled: isSidebarOpen,
        staleTime: 5 * 60 * 1000,
    });

    const handleDragEnd = useCallback((result: DropResult) => {
        if (!result.destination) return;

        const {source, destination} = result;

        if (source.droppableId === 'workshift-list' && destination.droppableId.startsWith('employee-')) {
            const employeeId = destination.droppableId.replace('employee-', '');
            const workshift = workshifts[source.index];

            setEmployeeShifts(prev => ({
                ...prev,
                [employeeId]: [...(prev[employeeId] || []), {
                    ...workshift,
                    tempId: `${workshift.workshift_id}-${Date.now()}`
                }]
            }));
        }

        if (source.droppableId.startsWith('employee-') &&
            destination.droppableId.startsWith('employee-') &&
            source.droppableId === destination.droppableId) {
            const employeeId = source.droppableId.replace('employee-', '');
            const newShifts = Array.from(employeeShifts[employeeId]);
            const [removed] = newShifts.splice(source.index, 1);
            newShifts.splice(destination.index, 0, removed);

            setEmployeeShifts(prev => ({
                ...prev,
                [employeeId]: newShifts
            }));
        }
    }, [workshifts, employeeShifts]);

    const SecondNav = useCallback(() => {
        const handleDateChange = useCallback(({month, year}: { month: number; year: number }) => {
            const newDate = new Date(currentDate);
            newDate.setFullYear(year, month - 1);

            if (newDate.getTime() !== currentDate.getTime()) {
                setCurrentDate(newDate);
            }
        }, [currentDate]);

        return (
            <div className="flex justify-between items-center py-2">
                <MonthYearPicker
                    initialMonth={currentDate.getMonth() + 1}
                    initialYear={currentDate.getFullYear()}
                    onChange={handleDateChange}
                />
                <div className="flex gap-2">
                    {tags.map((data, index) => (
                        <p key={index} className="bg-green-100 rounded-sm px-3 py-1">
                            X {data}
                        </p>
                    ))}
                </div>
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setIsSidebarOpen(!isSidebarOpen)}
                >
                    <Menu className="h-5 w-5"/>
                </Button>
            </div>
        );
    }, [currentDate, tags, isSidebarOpen]);

    return (
        <div className="relative">
            <div className="space-y-4">
                <Nav/>
                <hr/>
                <SecondNav/>
            </div>

            <DragDropContext onDragEnd={handleDragEnd}>
                {/* Week Calendar */}
                <div className="mt-6 overflow-x-scroll w-[1600px]">
                    <div className="min-w-max">
                        <div className="grid grid-cols-8 gap-1">
                            <div className="p-2 font-medium">Nhân viên</div>
                            {weekDays.map((day, i) => (
                                <div key={i} className="p-2 text-center font-medium">
                                    {day.toLocaleDateString('vi-VN', { weekday: 'short' })}
                                    <div className="text-sm">
                                        {day.getDate()}/{day.getMonth() + 1}
                                    </div>
                                </div>
                            ))}
                        </div>

                        {Object.keys(employeeShifts).map((employeeId) => (
                            <div key={employeeId} className="grid grid-cols-8 gap-1 border-t">
                                <div className="p-2 border-r">NV {employeeId.replace('emp', '')}</div>
                                {weekDays.map((day, dayIndex) => (
                                    <Droppable
                                        key={dayIndex}
                                        droppableId={`employee-${employeeId}-day-${dayIndex}`}
                                    >
                                        {(provided) => (
                                            <div
                                                ref={provided.innerRef}
                                                {...provided.droppableProps}
                                                className="min-h-20 p-2 bg-gray-50"
                                            >
                                                {employeeShifts[employeeId]
                                                    .filter(shift => true) // Add your filter logic here
                                                    .map((shift, shiftIndex) => (
                                                        <Draggable
                                                            key={shift.tempId || shift.workshift_id}
                                                            draggableId={shift.tempId || shift.workshift_id}
                                                            index={shiftIndex}
                                                        >
                                                            {(provided) => (
                                                                <div
                                                                    ref={provided.innerRef}
                                                                    {...provided.draggableProps}
                                                                    {...provided.dragHandleProps}
                                                                    className="mb-1 p-1 bg-white border rounded text-sm"
                                                                >
                                                                    {shift.workshift_name}: {formatTime(shift.start_time)}-{formatTime(shift.end_time)}
                                                                </div>
                                                            )}
                                                        </Draggable>
                                                    ))}
                                                {provided.placeholder}
                                            </div>
                                        )}
                                    </Droppable>
                                ))}
                            </div>
                        ))}
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
                                                key={shift.workshift_id}
                                                draggableId={shift.workshift_id}
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
                    className="fixed inset-0 bg-opacity-50 z-40"
                    onClick={() => setIsSidebarOpen(false)}
                />
            )}
        </div>
    );
};

export default Rota;