import { useState, useCallback, useMemo, useRef, memo, useEffect } from "react";
import { isToday } from "date-fns";
import { SecondNav } from "@/pages/HR/checkin/rota/_components/SecondNav.tsx";
import Loading from "@/components/Loading.tsx";
import { formatTime } from "@/lib/utils.ts";
import { getEmployeeWorkshiftsByManagerApi } from "@/apis/workshift.api";
import type { ManagerEmployeeSchedule } from "@/types/employee-workshift";

// --- Header ---
const Nav = memo(() => (
  <div className="flex justify-between items-center">
    <h3 className="font-semibold text-3xl">Lịch làm việc</h3>
  </div>
));

const DayCellView = memo(
  ({
    day,
    dayIndex,
    scheduledShift,
    actualShifts,
    columnWidths,
    onResizeColumn,
  }: {
    day: Date;
    dayIndex: number;
    scheduledShift: any | null;
    actualShifts: any[];
    columnWidths: number[];
    onResizeColumn: (index: number, width: number) => void;
  }) => {
    const resizeRef = useRef<HTMLDivElement>(null);

    const handleResizeStart = (e: React.MouseEvent) => {
      e.preventDefault();
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
        document.removeEventListener("mousemove", handleResize);
        document.removeEventListener("mouseup", handleResizeEnd);
      };

      document.addEventListener("mousemove", handleResize);
      document.addEventListener("mouseup", handleResizeEnd);
    };

    return (
      <div
        className={`relative p-2 min-h-20 border bg-gray-50 ${
          day.getDay() === 0 || day.getDay() === 6 ? "bg-red-50" : ""
        }`}
        style={{
          width: `${columnWidths[dayIndex] || 120}px`,
          minWidth: `${columnWidths[dayIndex] || 120}px`,
        }}
      >
        <div
          ref={resizeRef}
          className="absolute -right-1 top-0 bottom-0 w-2 cursor-col-resize z-10"
          onMouseDown={handleResizeStart}
        >
          <div className="absolute top-1/2 right-0 transform -translate-y-1/2 w-1 h-6 bg-gray-300 rounded"></div>
        </div>

        {/* Ca thực tế */}
        {actualShifts.map((shiftData, index) => (
          <div
            key={`actual-${shiftData.id}-${index}`}
            className="group mb-2 p-2 bg-white border border-gray-300 rounded"
          >
            <p className="text-sm font-medium">
              {shiftData.workshift?.workshift_name}
            </p>
            <p className="text-xs text-gray-600">
              {formatTime(shiftData.workshift?.start_time)} -{" "}
              {formatTime(shiftData.workshift?.end_time)}
            </p>
          </div>
        ))}

        {/* Ca dự đoán */}
        {scheduledShift && (
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
      </div>
    );
  }
);

// --- Một hàng nhân viên ---
const EmployeeRowView = memo(
  ({
    employee,
    monthDays,
    getScheduledShift,
    getActualEmployeeShift,
    columnWidths,
    onResizeColumn,
  }: {
    employee: any;
    monthDays: Date[];
    getScheduledShift: (employeeId: string, dayIndex: number) => any;
    getActualEmployeeShift: (employeeId: string, dayIndex: number) => any;
    columnWidths: number[];
    onResizeColumn: (index: number, width: number) => void;
  }) => (
    <div key={employee.id} className="flex">
      <div
        className="p-2 border-r border-b bg-white flex flex-col items-center justify-center min-w-[200px] max-w-[200px]"
        style={{
          position: "sticky",
          left: 0,
          zIndex: 40,
          backgroundColor: "white",
        }}
      >
        <p className="font-medium truncate">{employee.full_name}</p>
        <p className="text-sm text-gray-500 truncate">
          {employee.department_name}
        </p>
      </div>

      {monthDays.map((day, dayIndex) => (
        <DayCellView
          key={employee.employee_id + "-" + dayIndex}
          day={day}
          dayIndex={dayIndex}
          scheduledShift={getScheduledShift(employee.employee_id, dayIndex)}
          actualShifts={getActualEmployeeShift(employee.employee_id, dayIndex)}
          columnWidths={columnWidths}
          onResizeColumn={onResizeColumn}
        />
      ))}
    </div>
  )
);

// --- Component chính ---
const RotaManager = () => {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [columnWidths, setColumnWidths] = useState<number[]>([]);
  const [managerScheduleData, setManagerScheduleData] = useState<
    ManagerEmployeeSchedule[]
  >([]);
  const [loading, setLoading] = useState(true);

  const currentMonth = currentDate.getMonth() + 1;
  const currentYear = currentDate.getFullYear();

  const daysInMonth = useMemo(
    () => new Date(currentYear, currentMonth, 0).getDate(),
    [currentYear, currentMonth]
  );

  const monthDays = useMemo(() => {
    const days = Array.from(
      { length: daysInMonth },
      (_, i) => new Date(currentYear, currentMonth - 1, i + 1)
    );
    if (columnWidths.length !== days.length) {
      setColumnWidths(Array(days.length).fill(120));
    }
    return days;
  }, [currentYear, currentMonth, daysInMonth, columnWidths.length]);

  const handleResizeColumn = useCallback((index: number, width: number) => {
    setColumnWidths((prev) => {
      const newWidths = [...prev];
      newWidths[index] = width;
      return newWidths;
    });
  }, []);

  const handleDateChange = useCallback(
    ({ month, year }: { month: number; year: number }) => {
      const newDate = new Date(year, month - 1, 1);
      setCurrentDate(newDate);
    },
    []
  );

  const getScheduledShift = useCallback(() => null, []);
  const getActualEmployeeShift = useCallback(() => [], []);

  useEffect(() => {
    let ignore = false;

    const fetchData = async () => {
      if (
        !currentYear ||
        !currentMonth ||
        isNaN(currentYear) ||
        isNaN(currentMonth)
      ) {
        setManagerScheduleData([]);
        setLoading(false);
        return;
      }

      setLoading(true);
      try {
        const targetDate = `${currentYear}-${String(currentMonth).padStart(
          2,
          "0"
        )}-01`;
        const data = await getEmployeeWorkshiftsByManagerApi(targetDate);

        if (!ignore) {
          if (!Array.isArray(data)) {
            setManagerScheduleData([]);
            return;
          }

          // ✅ Gom theo employee_id
          const grouped = data.reduce((acc, item) => {
            if (!acc[item.employee_id]) acc[item.employee_id] = [];
            acc[item.employee_id].push(item);
            return acc;
          }, {} as Record<string, ManagerEmployeeSchedule[]>);

          // ✅ Định dạng lại dữ liệu
          const formatted = Object.entries(grouped).map(
            ([employeeId, shifts]) => ({
              employee_id: employeeId,
              full_name: shifts[0].full_name,
              shifts,
            })
          );

          setManagerScheduleData(formatted);
        }
      } catch (error) {
        console.error("Failed to load manager schedule:", error);
        if (!ignore) setManagerScheduleData([]);
      } finally {
        if (!ignore) setLoading(false);
      }
    };

    fetchData();

    return () => {
      ignore = true;
    };
  }, [currentMonth, currentYear]);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loading />
      </div>
    );
  }

  return (
    <div className="relative">
      <div className="space-y-4">
        <Nav />
        <hr />
        <SecondNav
          currentDate={currentDate}
          scheduleNames={[]}
          selectedSchedule={null}
          onDateChange={handleDateChange}
          onToggleSidebar={() => {}}
          onSelectSchedule={() => {}}
          isSidebarOpen={false}
        />
      </div>

      <div className="mt-6 overflow-auto w-full" style={{ height: "70vh" }}>
        <div className="min-w-max">
          {/* Header ngày */}
          <div className="flex bg-white sticky top-0 z-50 shadow">
            <div className="p-2 font-medium text-center border-b border-r min-w-[200px] max-w-[200px] sticky left-0 bg-white z-60">
              Nhân viên
            </div>

            {monthDays.map((day, i) => (
              <div
                key={i}
                className={`p-2 text-center font-medium border-b ${
                  day.getDay() === 0 || day.getDay() === 6 ? "bg-red-50" : ""
                } ${isToday(day) ? "border-2 border-blue-500 bg-blue-50" : ""}`}
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

          {/* Dòng nhân viên */}
          {managerScheduleData.length === 0 ? (
            <div className="flex justify-center items-center h-64 text-gray-500">
              Không có dữ liệu lịch làm việc
            </div>
          ) : (
            managerScheduleData.map((employee) => (
              <EmployeeRowView
                key={employee.employee_id}
                employee={employee}
                monthDays={monthDays}
                getScheduledShift={(id, dayIndex) =>
                  employee.shifts.find(
                    (s) => new Date(s.date).getDate() === dayIndex + 1
                  )
                }
                getActualEmployeeShift={() => []}
                columnWidths={columnWidths}
                onResizeColumn={handleResizeColumn}
              />
            ))
          )}
        </div>
      </div>
    </div>
  );
};

export default RotaManager;
