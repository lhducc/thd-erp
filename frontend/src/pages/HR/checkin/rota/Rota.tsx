import { Button } from "@/components/ui/button.tsx";
import { useState, useCallback, useMemo, useRef, useEffect, memo } from "react";
import { X, GripVertical } from "lucide-react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { type Workshift } from "@/types/workshift.ts";
import {
  deleteEmployeeWorkshiftApi,
  getAllWorkshiftApi,
} from "@/apis/workshift.api.ts";
import { formatTime } from "@/lib/utils.ts";
import Loading from "@/components/Loading.tsx";
import {
  DragDropContext,
  Droppable,
  Draggable,
  type DropResult,
} from "@hello-pangea/dnd";
import { getAllShiftAllocations } from "@/apis/shift-allocation.api.ts";
import { SecondNav } from "@/pages/HR/checkin/rota/_components/SecondNav.tsx";
import { isToday } from "date-fns";
import {
  getAllEmployeeWorkshiftsApi,
  registerManyWorkshiftsApi,
} from "@/apis/employee-workshift.api";
import {
  type EmployeeWorkshift,
  type RegisterWorkshiftRequest,
} from "@/types/employee-workshift";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext.tsx";

// Memoized components để tránh re-render không cần thiết
const Nav = memo(() => {
  return (
    <div className="flex justify-between items-center">
      <h3 className="font-semibold text-3xl">Bảng phân ca</h3>
    </div>
  );
});

// Component cho mỗi ô ngày - được memoized để tránh re-render không cần thiết
const DayCell = memo(
  ({
    employee,
    day,
    dayIndex,
    scheduledShift,
    assignedShifts,
    actualShifts,
    hasOverlap,
    columnWidths,
    onResizeColumn,
  }: {
    employee: any;
    day: Date;
    dayIndex: number;
    scheduledShift: any;
    assignedShifts: any[];
    actualShifts: any[];
    hasOverlap: boolean;
    columnWidths: number[];
    onResizeColumn: (index: number, width: number) => void;
  }) => {
    const { currentUser } = useAuth();
    const dayKey = `${employee.employee_id}-${dayIndex}`;
    const [isResizing, setIsResizing] = useState(false);
    const resizeRef = useRef<HTMLDivElement>(null);
    const queryClient = useQueryClient();

    const deleteMutation = useMutation({
      mutationFn: (id: number) => deleteEmployeeWorkshiftApi(id),
      onSuccess: () => {
        toast.success("Xóa ca làm việc thành công!");
        queryClient.invalidateQueries({ queryKey: ["employee-workshifts"] });
      },
      onError: (error: any) => {
        toast.error(`Xóa ca làm việc thất bại: ${error.message}`);
      },
    });

    const handleResizeStart = (e: React.MouseEvent) => {
      e.preventDefault();
      setIsResizing(true);

      const startX = e.clientX;
      const startWidth = columnWidths[dayIndex] || 120;

      const handleResize = (moveEvent: MouseEvent) => {
        const newWidth = Math.max(
          100,
          startWidth + (moveEvent.clientX - startX)
        );
        onResizeColumn(dayIndex, newWidth);
      };

      const handleResizeEnd = () => {
        setIsResizing(false);
        document.removeEventListener("mousemove", handleResize);
        document.removeEventListener("mouseup", handleResizeEnd);
      };

      document.addEventListener("mousemove", handleResize);
      document.addEventListener("mouseup", handleResizeEnd);
    };

    return (
      <Droppable droppableId={dayKey}>
        {(provided, snapshot) => (
          <div
            ref={provided.innerRef}
            {...provided.droppableProps}
            className={`relative p-2 min-h-20 border bg-gray-50 ${
              day.getDay() === 0 || day.getDay() === 6 ? "bg-red-50" : ""
            } ${
              snapshot.isDraggingOver ? "bg-blue-100 ring-2 ring-blue-400" : ""
            }`}
            style={{
              width: `${columnWidths[dayIndex] || 120}px`,
              minWidth: `${columnWidths[dayIndex] || 120}px`,
            }}
          >
            {/* Handle để resize cột */}
            <div
              ref={resizeRef}
              className="absolute -right-1 top-0 bottom-0 w-2 cursor-col-resize z-10"
              onMouseDown={handleResizeStart}
              style={{ cursor: "col-resize" }}
            >
              <div className="absolute top-1/2 right-0 transform -translate-y-1/2 w-1 h-6 bg-gray-300 rounded"></div>
            </div>

            {/* Hiển thị ca thực tế */}
            {actualShifts.map((shiftData, index) => (
              <div
                key={`actual-${shiftData.id}-${index}`}
                className="group mb-2 p-2 bg-white border border-gray-300 rounded min-w-[100px]"
              >
                <div className="flex justify-between items-center">
                  <p className="text-sm font-medium min-w-[50px]">
                    {shiftData.workshift.workshift_name}
                  </p>
                  <button
                    onClick={() => {
                      deleteMutation.mutate(shiftData.id); // Sử dụng id của employee-workshift
                    }}
                    className="hidden group-hover:inline bg-red-500/60 w-fit h-fit px-2 cursor-pointer"
                  >
                    x
                  </button>
                </div>
                <p className="text-xs text-gray-600">
                  {formatTime(shiftData.workshift.start_time)} -{" "}
                  {formatTime(shiftData.workshift.end_time)}
                </p>
              </div>
            ))}

            {/* Hiển thị ca dự đoán */}
            {scheduledShift &&
              !hasOverlap &&
              currentUser?.role !== "manager" && (
                <div className="mb-2 p-2 bg-blue-100 border border-blue-300 rounded">
                  <p className="text-sm font-medium text-blue-800">
                    {scheduledShift.workshift_name}
                  </p>
                  <p className="text-xs text-blue-600">
                    {formatTime(scheduledShift.start_time)} -{" "}
                    {formatTime(scheduledShift.end_time)}
                  </p>
                  {scheduledShift.scheduleName && (
                    <p className="text-xs mt-1 text-blue-700">
                      {scheduledShift.scheduleName}
                    </p>
                  )}
                </div>
              )}

            {/* Hiển thị ca được kéo thả */}
            {assignedShifts.map((shift, index) => (
              <Draggable
                key={shift.tempId || shift.workshift_id}
                draggableId={(shift.tempId || shift.workshift_id).toString()}
                index={index}
              >
                {(provided, snapshot) => (
                  <div
                    ref={provided.innerRef}
                    {...provided.draggableProps}
                    {...provided.dragHandleProps}
                    className={`mb-2 p-2 border rounded cursor-move flex items-start ${
                      snapshot.isDragging
                        ? "bg-green-100 border-green-400 rotate-5 shadow-lg"
                        : "bg-yellow-100 border-yellow-300"
                    }`}
                  >
                    <GripVertical className="h-3 w-3 mr-1 mt-0.5 text-yellow-600 flex-shrink-0" />
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-yellow-800 truncate">
                        Phân công: {shift.workshift_name}
                      </p>
                      <p className="text-xs text-yellow-600">
                        {formatTime(shift.start_time)} -{" "}
                        {formatTime(shift.end_time)}
                      </p>
                    </div>
                  </div>
                )}
              </Draggable>
            ))}
            {provided.placeholder}
          </div>
        )}
      </Droppable>
    );
  }
);

// Component cho mỗi hàng nhân viên - được memoized
const EmployeeRow = memo(
  ({
    employee,
    monthDays,
    selectedSchedule,
    shiftAllocation,
    processedEmployeeWorkshifts,
    employeeShifts,
    getScheduledShift,
    getActualEmployeeShift,
    isShiftOverlap,
    columnWidths,
    onResizeColumn,
  }: {
    employee: any;
    monthDays: Date[];
    selectedSchedule: string | null;
    shiftAllocation: any[];
    processedEmployeeWorkshifts: Record<string, EmployeeWorkshift[]>;
    employeeShifts: Record<string, Workshift[]>;
    getScheduledShift: (employeeId: string, dayIndex: number) => any;
    getActualEmployeeShift: (employeeId: string, dayIndex: number) => any;
    isShiftOverlap: (actualShift: Workshift, scheduledShift: any) => boolean;
    columnWidths: number[];
    onResizeColumn: (index: number, width: number) => void;
  }) => {
    return (
      <div key={employee.employee_id} className="flex">
        <div
          className="p-2 border-r border-b bg-white flex flex-col items-center justify-center min-w-[200px] max-w-[200px]"
          style={{
            position: "sticky",
            left: 0,
            zIndex: 40,
            backgroundColor: "white",
          }}
        >
          <p className="font-medium truncate max-w-full">
            {employee.full_name}
          </p>
          <p className="text-sm text-gray-500 truncate max-w-full">
            {employee.department_name}
          </p>
        </div>

        {monthDays.map((day, dayIndex) => {
          const scheduledShift = getScheduledShift(
            employee.employee_id,
            dayIndex
          );
          const assignedShifts =
            employeeShifts[`${employee.employee_id}-${dayIndex}`] || [];
          const actualShifts = getActualEmployeeShift(
            employee.employee_id,
            dayIndex
          );

          // Kiểm tra xem có ca thật và ca dự đoán trùng nhau không
          const hasOverlap =
            actualShifts.length > 0 &&
            scheduledShift &&
            actualShifts.some((actualShift) =>
              isShiftOverlap(actualShift, scheduledShift)
            );

          return (
            <DayCell
              key={employee.employee_id + "-" + dayIndex}
              employee={employee}
              day={day}
              dayIndex={dayIndex}
              scheduledShift={scheduledShift}
              assignedShifts={assignedShifts}
              actualShifts={actualShifts}
              hasOverlap={hasOverlap}
              columnWidths={columnWidths}
              onResizeColumn={onResizeColumn}
            />
          );
        })}
      </div>
    );
  }
);

const Rota = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [currentDate, setCurrentDate] = useState(new Date());
  const [employeeShifts, setEmployeeShifts] = useState<
    Record<string, Workshift[]>
  >({});
  const [selectedSchedule, setSelectedSchedule] = useState<string | null>(null);
  const [isDraggingColumns, setIsDraggingColumns] = useState(false);
  const [startX, setStartX] = useState(0);
  const [scrollLeft, setScrollLeft] = useState(0);
  const [columnWidths, setColumnWidths] = useState<number[]>([]);
  const containerRef = useRef<HTMLDivElement>(null);
  const [draggingItem, setDraggingItem] = useState<string | null>(null);
  const { currentUser } = useAuth();
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedDepartment, setSelectedDepartment] = useState("");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");

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
    const days = Array.from({ length: daysInMonth }, (_, i) => {
      const day = new Date(currentYear, currentMonth - 1, i + 1);
      return day;
    });

    // Khởi tạo chiều rộng cột nếu chưa có
    if (columnWidths.length !== days.length) {
      setColumnWidths(Array(days.length).fill(120));
    }

    return days;
  }, [currentYear, currentMonth, daysInMonth, columnWidths.length]);

  // Thêm query để lấy dữ liệu ca làm việc thực tế
  const {
    data: employeeWorkshiftsData = {},
    isPending: isEmployeeWorkshiftsPending,
    refetch: refetchEmployeeWorkshifts,
  } = useQuery({
    queryKey: ["employee-workshifts", currentMonth, currentYear],
    queryFn: () => getAllEmployeeWorkshiftsApi(currentMonth, currentYear),
  });

  // Mutation để đăng ký nhiều ca làm việc
  const registerManyMutation = useMutation({
    mutationFn: registerManyWorkshiftsApi,
    onSuccess: () => {
      toast.success("Phân ca thành công!");
      // Refetch dữ liệu ca thật sau khi phân ca
      refetchEmployeeWorkshifts();
      // Xóa ca ảo sau khi phân ca thành công
      setEmployeeShifts({});
    },
    onError: (error) => {
      toast.error("Phân ca thất bại: " + error.message);
    },
  });

  // Xử lý dữ liệu employee workshifts
  const processedEmployeeWorkshifts = useMemo(() => {
    if (Array.isArray(employeeWorkshiftsData)) {
      const result: Record<string, EmployeeWorkshift[]> = {};
      employeeWorkshiftsData.forEach((shift: EmployeeWorkshift) => {
        if (!result[shift.employee_id]) {
          result[shift.employee_id] = [];
        }
        result[shift.employee_id].push(shift);
      });
      return result;
    }
    return employeeWorkshiftsData;
  }, [employeeWorkshiftsData]);

  const handleMouseDown = (e: React.MouseEvent) => {
    if (!containerRef.current) return;
    setIsDraggingColumns(true);
    setStartX(e.pageX - containerRef.current.offsetLeft);
    setScrollLeft(containerRef.current.scrollLeft);
  };

  const handleMouseMove = useCallback(
    (e: React.MouseEvent) => {
      if (!isDraggingColumns || !containerRef.current) return;
      e.preventDefault();
      const x = e.pageX - containerRef.current.offsetLeft;
      const walk = (x - startX) * 2;
      containerRef.current.scrollLeft = scrollLeft - walk;
    },
    [isDraggingColumns, startX, scrollLeft]
  );

  const handleMouseUpOrLeave = useCallback(() => {
    setIsDraggingColumns(false);
  }, []);

  const { data: shiftAllocation = [], isPending: isShiftAllocationPending } =
    useQuery({
      queryKey: ["shift-allocation", currentMonth, currentYear],
      queryFn: () => getAllShiftAllocations(currentMonth, currentYear),
    });

  const scheduleNames = useMemo(() => {
    const names = new Set<string>();
    shiftAllocation?.forEach((employee) => {
      employee.schedules?.forEach((schedule) => {
        if (schedule.work_schedule_name) {
          names.add(schedule.work_schedule_name);
        }
      });
    });
    return Array.from(names);
  }, [shiftAllocation]);

  const filteredEmployees = useMemo(() => {
    let employees = shiftAllocation;

    // Lọc theo schedule
    if (selectedSchedule) {
      employees = employees.filter((employee) =>
        employee.schedules?.some(
          (schedule) => schedule.work_schedule_name === selectedSchedule
        )
      );
    }

    // Lọc theo phòng ban
    if (selectedDepartment) {
      employees = employees.filter(
        (e) =>
          e.department_name &&
          e.department_name.toLowerCase() === selectedDepartment.toLowerCase()
      );
    }

    // Lọc theo tên nhân viên
    if (searchTerm.trim()) {
      employees = employees.filter((e) =>
        e.full_name.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Sort theo tên phòng ban
    employees = [...employees].sort((a, b) => {
      const depA = a.department_name?.toLowerCase() || "";
      const depB = b.department_name?.toLowerCase() || "";
      return sortOrder === "asc"
        ? depA.localeCompare(depB)
        : depB.localeCompare(depA);
    });

    return employees;
  }, [
    shiftAllocation,
    selectedSchedule,
    selectedDepartment,
    searchTerm,
    sortOrder,
  ]);

  const {
    data: workshifts = [],
    isLoading,
    isError,
    error,
  } = useQuery<Workshift[]>({
    queryKey: ["workshifts"],
    queryFn: getAllWorkshiftApi,
    enabled: isSidebarOpen,
    staleTime: 5 * 60 * 1000,
  });

  const getDayName = useCallback((date: Date): string => {
    return date?.toLocaleDateString("en-US", { weekday: "long" }).toLowerCase();
  }, []);

  const handleDragEnd = useCallback(
    (result: DropResult) => {
      const { source, destination } = result;
      setDraggingItem(null);

      // Nếu không có destination (kéo ra ngoài vùng drop) thì không làm gì
      if (!destination) return;

      // Kéo từ sidebar vào lịch
      if (
        source.droppableId === "workshift-list" &&
        destination.droppableId.includes("-")
      ) {
        const [employeeId, dayIndex] = destination.droppableId.split("-");
        const workshift = workshifts[source.index];

        setEmployeeShifts((prev) => {
          const newShifts = { ...prev };
          const dayKey = `${employeeId}-${dayIndex}`;

          if (!newShifts[dayKey]) {
            newShifts[dayKey] = [];
          }

          // Kiểm tra xem ca đã tồn tại chưa để tránh trùng lặp
          const shiftExists = newShifts[dayKey].some(
            (s) => s.workshift_id === workshift.workshift_id
          );

          if (!shiftExists) {
            newShifts[dayKey] = [
              ...newShifts[dayKey],
              {
                ...workshift,
                tempId: `${workshift.workshift_id}-${Date.now()}`,
              },
            ];
          }

          return newShifts;
        });
        return; // Kết thúc xử lý sau khi thêm ca mới
      }

      // Kéo thả giữa các ô khác nhau
      if (
        source.droppableId.includes("-") &&
        destination.droppableId.includes("-") &&
        source.droppableId !== destination.droppableId
      ) {
        const [sourceEmployeeId, sourceDayIndex] =
          source.droppableId.split("-");
        const [destEmployeeId, destDayIndex] =
          destination.droppableId.split("-");

        setEmployeeShifts((prev) => {
          const newShifts = { ...prev };
          const sourceKey = `${sourceEmployeeId}-${sourceDayIndex}`;
          const destKey = `${destEmployeeId}-${destDayIndex}`;

          // Lấy ca từ ô nguồn
          const sourceShifts = [...(newShifts[sourceKey] || [])];

          // Kiểm tra index có hợp lệ không
          if (source.index >= sourceShifts.length) return newShifts;

          const [movedShift] = sourceShifts.splice(source.index, 1);

          // Thêm vào ô đích
          if (!newShifts[destKey]) {
            newShifts[destKey] = [];
          }

          // Kiểm tra xem ca đã tồn tại chưa
          const shiftExists = newShifts[destKey].some(
            (s) => s?.workshift_id === movedShift?.workshift_id
          );

          if (!shiftExists && movedShift) {
            // Đảm bảo index không vượt quá độ dài mảng
            const insertIndex = Math.min(
              destination.index,
              newShifts[destKey].length
            );

            newShifts[destKey] = [
              ...newShifts[destKey].slice(0, insertIndex),
              {
                ...movedShift,
                tempId: `${movedShift.workshift_id}-${Date.now()}`,
              },
              ...newShifts[destKey].slice(insertIndex),
            ];
          }

          // Cập nhật ô nguồn
          if (sourceShifts.length === 0) {
            delete newShifts[sourceKey];
          } else {
            newShifts[sourceKey] = sourceShifts;
          }

          return newShifts;
        });
        return; // Kết thúc xử lý sau khi di chuyển ca
      }

      // Kéo thả trong cùng một ô
      if (
        source.droppableId.includes("-") &&
        destination.droppableId.includes("-") &&
        source.droppableId === destination.droppableId
      ) {
        const [employeeId, dayIndex] = source.droppableId.split("-");
        const dayKey = `${employeeId}-${dayIndex}`;

        setEmployeeShifts((prev) => {
          const newShifts = { ...prev };
          const shiftsCopy = [...(newShifts[dayKey] || [])];

          // Kiểm tra index có hợp lệ không
          if (
            source.index >= shiftsCopy.length ||
            destination.index > shiftsCopy.length
          )
            return newShifts;

          const [removed] = shiftsCopy.splice(source.index, 1);
          shiftsCopy.splice(destination.index, 0, removed);

          newShifts[dayKey] = shiftsCopy;
          return newShifts;
        });
      }
    },
    [workshifts]
  );

  const handleDragStart = useCallback((result: any) => {
    setDraggingItem(result.draggableId);
  }, []);

  const handleScheduleFilter = useCallback((scheduleName: string) => {
    setSelectedSchedule((prev) =>
      prev === scheduleName ? null : scheduleName
    );
  }, []);

  const getActualEmployeeShift = useCallback(
    (employeeId: string, dayIndex: number) => {
      const day = monthDays[dayIndex];

      // Lấy ngày theo định dạng YYYY-MM-DD
      const dayString = `${day.getFullYear()}-${(day.getMonth() + 1)
        .toString()
        .padStart(2, "0")}-${day.getDate().toString().padStart(2, "0")}`;

      const employeeShifts = processedEmployeeWorkshifts[employeeId] || [];

      return employeeShifts
        .filter((shift: EmployeeWorkshift) => {
          try {
            // Nếu shift.date đã là ISO thì cắt phần ngày trực tiếp
            const shiftDateString = shift.date.substring(0, 10); // "YYYY-MM-DD"
            return shiftDateString === dayString;
          } catch (error) {
            console.error("Error parsing date:", shift.date, error);
            return false;
          }
        })
        .map((shift: EmployeeWorkshift) => ({
          id: shift.id, // Thêm id của employee-workshift
          workshift: shift.workshift,
          date: shift.date,
        }));
    },
    [monthDays, processedEmployeeWorkshifts]
  );

  // Hàm kiểm tra xem ca thật và ca dự đoán có trùng nhau không
  const isShiftOverlap = useCallback(
    (actualShift: Workshift, scheduledShift: any): boolean => {
      if (!scheduledShift) return false;
      return (
        actualShift.workshift_id === scheduledShift.workshift_id ||
        (actualShift.start_time === scheduledShift.start_time &&
          actualShift.end_time === scheduledShift.end_time)
      );
    },
    []
  );

  const getScheduledShift = useCallback(
    (employeeId: string, dayIndex: number) => {
      const employee = shiftAllocation.find(
        (e) => e.employee_id === employeeId
      );
      if (!employee || !employee.schedules || employee.schedules.length === 0)
        return null;

      const day = monthDays[dayIndex];
      if (!day) return null;

      const dayName = getDayName(day);
      for (const schedule of employee.schedules) {
        if (schedule.is_schedule_auto === false) {
          continue; // Bỏ qua schedule nếu is_schedule_auto là false
        }
        if (
          selectedSchedule &&
          schedule.work_schedule_name !== selectedSchedule
        ) {
          continue;
        }

        const weekdayShift = schedule.weekdays?.find(
          (w) => w.week_day === dayName
        );
        if (weekdayShift) {
          return {
            ...weekdayShift.work_shift,
            scheduleName: schedule.work_schedule_name,
          };
        }
      }
      return null;
    },
    [shiftAllocation, selectedSchedule, getDayName, monthDays]
  );

  // Hàm xử lý khi nhấn nút phân ca - CHỈ lấy ca được kéo thả
  const handleAssignShifts = useCallback(() => {
    // Tạo dữ liệu để gửi API - chỉ lấy ca từ employeeShifts (kéo thả)
    const shiftsToRegister: RegisterWorkshiftRequest[] = [];

    // Duyệt qua tất cả employeeShifts (ca được kéo thả)
    Object.entries(employeeShifts).forEach(([dayKey, shifts]) => {
      const [employeeId, dayIndexStr] = dayKey.split("-");
      const dayIndex = parseInt(dayIndexStr);

      if (dayIndex >= 0 && dayIndex < monthDays.length) {
        const day = monthDays[dayIndex + 1];

        shifts.forEach((shift) => {
          shiftsToRegister.push({
            employee_id: employeeId,
            workshift_id: shift.workshift_id,
            date: day.toISOString(),
          });
        });
      }
    });

    if (shiftsToRegister.length === 0) {
      toast.warning(
        "Không có ca nào để phân. Hãy kéo ca từ danh sách bên phải vào lịch."
      );
      return;
    }

    // Gọi API
    registerManyMutation.mutate(shiftsToRegister);
  }, [employeeShifts, monthDays, registerManyMutation]);
  // Hàm phân ca tự động từ schedule
  const handleAssignFromSchedule = useCallback(() => {
    const shiftsToRegister: RegisterWorkshiftRequest[] = [];

    filteredEmployees.forEach((employee) => {
      monthDays.forEach((day, dayIndex) => {
        const scheduledShift = getScheduledShift(
          employee.employee_id,
          dayIndex
        );
        const actualShifts = getActualEmployeeShift(
          employee.employee_id,
          dayIndex
        );

        // Chỉ thêm ca dự đoán nếu KHÔNG có ca thật nào
        if (scheduledShift && actualShifts.length === 0) {
          shiftsToRegister.push({
            employee_id: employee.employee_id,
            workshift_id: scheduledShift.workshift_id,
            date: day.toISOString(),
          });
        }
      });
    });

    if (shiftsToRegister.length === 0) {
      toast.warning("Không có ca dự đoán nào để phân.");
      return;
    }

    registerManyMutation.mutate(shiftsToRegister);
  }, [
    filteredEmployees,
    monthDays,
    getScheduledShift,
    getActualEmployeeShift,
    registerManyMutation,
  ]);

  const handleDateChange = useCallback(
    ({ month, year }: { month: number; year: number }) => {
      const newDate = new Date(year, month - 1, 1);
      setCurrentDate(newDate);
      // Reset ca phân công khi đổi tháng
      setEmployeeShifts({});
    },
    []
  );

  // Xử lý resize cột
  const handleResizeColumn = useCallback((index: number, width: number) => {
    setColumnWidths((prev) => {
      const newWidths = [...prev];
      newWidths[index] = width;
      return newWidths;
    });
  }, []);

  // Đếm tổng số ca đã được kéo thả
  const totalAssignedShifts = useMemo(() => {
    return Object.values(employeeShifts).reduce(
      (total, shifts) => total + shifts.length,
      0
    );
  }, [employeeShifts]);

  return (
    <div className="relative">
      <div className="space-y-4">
        <Nav />

        <hr />
        <SecondNav
          currentDate={currentDate}
          scheduleNames={scheduleNames}
          selectedSchedule={selectedSchedule}
          onDateChange={handleDateChange}
          onToggleSidebar={() => setIsSidebarOpen(!isSidebarOpen)}
          onSelectSchedule={handleScheduleFilter}
          isSidebarOpen={isSidebarOpen}
        />

        <div className="flex items-center gap-3 my-3">
          <input
            type="text"
            placeholder="Tìm theo tên nhân viên..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="border rounded px-3 py-2 w-64"
          />

          <select
            value={selectedDepartment}
            onChange={(e) => setSelectedDepartment(e.target.value)}
            className="border rounded px-3 py-2"
          >
            <option value="">Tất cả phòng ban</option>
            {Array.from(
              new Set(shiftAllocation.map((e) => e.department_name))
            ).map((dept) => (
              <option key={dept} value={dept}>
                {dept}
              </option>
            ))}
          </select>

          <Button
            variant="outline"
            onClick={() => setSortOrder(sortOrder === "asc" ? "desc" : "asc")}
          >
            Sắp xếp theo phòng ban {sortOrder === "asc" ? "↑" : "↓"}
          </Button>
        </div>
      </div>

      <DragDropContext onDragEnd={handleDragEnd} onDragStart={handleDragStart}>
        <div
          ref={containerRef}
          className="mt-6 overflow-auto w-[80vw]"
          style={{
            cursor: isDraggingColumns ? "grabbing" : "grab",
            height: "70vh",
          }}
          onMouseDown={handleMouseDown}
          onMouseMove={handleMouseMove}
          onMouseUp={handleMouseUpOrLeave}
          onMouseLeave={handleMouseUpOrLeave}
        >
          <div className="min-w-max">
            {/* Header với sticky top */}
            <div
              className="flex bg-white"
              style={{
                position: "sticky",
                top: 0,
                zIndex: 50,
                boxShadow: "0 2px 4px rgba(0,0,0,0.1)",
              }}
            >
              <div
                className="p-2 font-medium text-center border-b border-r min-w-[200px] max-w-[200px]"
                style={{
                  position: "sticky",
                  left: 0,
                  zIndex: 60,
                  backgroundColor: "white",
                }}
              >
                Nhân viên
              </div>
              {monthDays.map((day, i) => (
                <div
                  key={i}
                  className={`p-2 text-center font-medium border-b ${
                    day.getDay() === 0 || day.getDay() === 6 ? "bg-red-50" : ""
                  } ${
                    isToday(day) ? "border-2 border-blue-500 bg-blue-50" : ""
                  }`}
                  style={{
                    width: `${columnWidths[i] || 120}px`,
                    minWidth: `${columnWidths[i] || 120}px`,
                    backgroundColor: "white",
                  }}
                >
                  {day.toLocaleDateString("vi-VN", { weekday: "short" })}
                  <div className="text-sm">
                    {day.getDate()}/{day.getMonth() + 1}
                    {isToday(day) && (
                      <div className="w-2 h-2 bg-blue-500 rounded-full mx-auto mt-1"></div>
                    )}
                  </div>
                </div>
              ))}
            </div>

            {isShiftAllocationPending || isEmployeeWorkshiftsPending ? (
              <div className="flex justify-center items-center h-64">
                <Loading />
              </div>
            ) : (
              filteredEmployees?.map((employee) => (
                <EmployeeRow
                  key={employee.employee_id}
                  employee={employee}
                  monthDays={monthDays}
                  selectedSchedule={selectedSchedule}
                  shiftAllocation={shiftAllocation}
                  processedEmployeeWorkshifts={processedEmployeeWorkshifts}
                  employeeShifts={employeeShifts}
                  getScheduledShift={getScheduledShift}
                  getActualEmployeeShift={getActualEmployeeShift}
                  isShiftOverlap={isShiftOverlap}
                  columnWidths={columnWidths}
                  onResizeColumn={handleResizeColumn}
                />
              ))
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div
          className={`
          fixed top-0 right-0 h-full w-96 bg-white shadow-lg z-50
          transform transition-transform duration-300 ease-in-out
          ${isSidebarOpen ? "translate-x-0" : "translate-x-full"}
        `}
        >
          <div className="p-4 h-full flex flex-col">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-bold">Danh sách ca làm việc</h2>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setIsSidebarOpen(false)}
              >
                <X className="h-5 w-5" />
              </Button>
            </div>

            {isLoading ? (
              <div className="flex-1 flex items-center justify-center">
                <Loading />
              </div>
            ) : isError ? (
              <div className="flex-1 flex items-center justify-center">
                <p className="text-red-500">
                  Lỗi khi tải dữ liệu: {(error as Error).message}
                </p>
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
                        draggableId={(
                          shift.tempId || shift.workshift_id
                        ).toString()}
                        index={index}
                      >
                        {(provided, snapshot) => (
                          <div
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            {...provided.dragHandleProps}
                            className={`mb-2 p-2 border rounded-lg cursor-move flex items-start ${
                              snapshot.isDragging
                                ? "bg-green-100 border-green-400 shadow-md"
                                : ""
                            }`}
                          >
                            <GripVertical className="h-4 w-4 mr-2 mt-0.5 text-gray-500 flex-shrink-0" />
                            <div className="flex-1">
                              <h3 className="font-medium">
                                {shift.workshift_name}:{" "}
                                {formatTime(shift.start_time)} -{" "}
                                {formatTime(shift.end_time)}
                              </h3>
                            </div>
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

      {/* Nút phân ca cố định ở góc */}
      <div className={`fixed bottom-6 right-6 z-50 flex flex-col gap-2`}>
        <Button
          onClick={handleAssignShifts}
          disabled={registerManyMutation.isPending || totalAssignedShifts === 0}
          size="lg"
          className="shadow-lg"
        >
          {registerManyMutation.isPending ? (
            <>
              <Loading className="mr-2 h-4 w-4" />
              Đang phân ca...
            </>
          ) : (
            `Phân ca (${totalAssignedShifts})`
          )}
        </Button>

        <Button
          onClick={handleAssignFromSchedule}
          disabled={registerManyMutation.isPending}
          variant="outline"
          size="sm"
          className={`shadow-lg ${
            currentUser?.role === "manager" ? "hidden" : ""
          } `}
        >
          Phân ca tự động từ lịch
        </Button>
      </div>

      {isSidebarOpen && (
        <div
          className="fixed inset-0 z-40"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}
    </div>
  );
};

export default Rota;
