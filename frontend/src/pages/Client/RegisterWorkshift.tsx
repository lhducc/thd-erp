import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  deleteEmployeeWorkshiftByEmployeeApi,
  getAllWorkshiftApi,
  getEmployeeWorkshiftsApi,
  getRegisterWorkshiftApi,
  getWorkshiftInfoApi,
  registerEmployeeWorkshiftApi,
  updateEmployeeWorkshiftApi,
} from "@/apis/workshift.api.ts";
import type { EmployeeWorkshift } from "@/types/employee-workshift";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext.tsx";

const RegisterWorkshift = () => {
  const [currentDate] = useState(new Date());
  const [selectedMonth, setSelectedMonth] = useState(currentDate.getMonth());
  const [selectedYear, setSelectedYear] = useState(currentDate.getFullYear());
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [selectedWorkshift, setSelectedWorkshift] = useState<string>("");
  const [showModal, setShowModal] = useState(false);
  const [editingWorkshiftId, setEditingWorkshiftId] = useState<string | null>(
    null
  );
  const queryClient = useQueryClient();
  // Vietnamese month names
  const months = [
    "Tháng 1",
    "Tháng 2",
    "Tháng 3",
    "Tháng 4",
    "Tháng 5",
    "Tháng 6",
    "Tháng 7",
    "Tháng 8",
    "Tháng 9",
    "Tháng 10",
    "Tháng 11",
    "Tháng 12",
  ];

  // Generate years range (current year -5 to +5)
  const years = Array.from(
    { length: 11 },
    (_, i) => currentDate.getFullYear() - 5 + i
  );

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

  const { data: employeeWorkshifts } = useQuery({
    queryKey: [
      "employeeWorkshifts",
      currentUser?.user_id,
      selectedMonth,
      selectedYear,
    ],
    queryFn: () => getEmployeeWorkshiftsApi(selectedMonth + 1, selectedYear),
  });

  const { data: registerWorkshifts } = useQuery({
    queryKey: ["register-workshifts"],
    queryFn: getRegisterWorkshiftApi,
  });

  const { data: allWorkshifts } = useQuery({
    queryKey: ["all-workshifts"],
    queryFn: getAllWorkshiftApi,
  });

  const { data: workshiftInfo, isLoading: isLoadingWorkshiftInfo } = useQuery({
    queryKey: ["workshift-info"],
    queryFn: getWorkshiftInfoApi,
  });

  const isScheduleAuto = workshiftInfo?.is_schedule_auto === true;

  const getWorkshiftsForSelectedDay = () => {
    if (!selectedDate || !registerWorkshifts) return [];

    // Chuyển đổi thứ từ Date object sang dạng string giống trong data
    const daysOfWeek = [
      "sunday",
      "monday",
      "tuesday",
      "wednesday",
      "thursday",
      "friday",
      "saturday",
    ];
    const dayIndex = selectedDate.getDay();
    const currentDay = daysOfWeek[dayIndex];

    // Lọc các ca làm việc cho thứ hiện tại
    return registerWorkshifts
      .filter((ws) => ws.week_day === currentDay && ws.work_shift !== null)
      .map((ws) => ({
        id: ws.workshift_id,
        name: ws.work_shift.workshift_name,
        start: ws.work_shift.start_time,
        end: ws.work_shift.end_time,
      }));
  };

  // Mutations
  const registerMutation = useMutation({
    mutationFn: ({
      employee_id,
      work_shift_id,
      date,
    }: {
      employee_id: string;
      work_shift_id: string;
      date: string;
    }) => registerEmployeeWorkshiftApi(employee_id, work_shift_id, date),
    onSuccess: () => {
      toast.success("Đăng ký ca làm việc thành công!");
      queryClient.invalidateQueries({ queryKey: ["employeeWorkshifts"] });
      setShowModal(false);
    },
    onError: (error) => {
      toast.error(`Đăng ký ca làm việc thất bại: ${error.message}`);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteEmployeeWorkshiftByEmployeeApi(id),
    onSuccess: () => {
      toast.success("Xóa ca làm việc thành công!");
      queryClient.invalidateQueries({ queryKey: ["employeeWorkshifts"] });
    },
    onError: (error) => {
      toast.error(`Xóa ca làm việc thất bại: ${error.message}`);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({
      id,
      work_shift_id,
    }: {
      id: string;
      work_shift_id: string;
    }) => updateEmployeeWorkshiftApi(id, work_shift_id),
    onSuccess: () => {
      toast.success("Cập nhật ca làm việc thành công!");
      queryClient.invalidateQueries({ queryKey: ["employeeWorkshifts"] });
      setShowModal(false);
    },
    onError: (error) => {
      toast.error(`Cập nhật ca làm việc thất bại: ${error.message}`);
    },
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
        setSelectedWorkshift("");
      }
      setShowModal(true);
    }
  };

  const handleRegister = () => {
    if (!selectedDate || !selectedWorkshift || !currentUser?.user_id) return;

    // Sửa đổi: Tạo chuỗi ngày tháng chính xác theo múi giờ địa phương
    const year = selectedDate.getFullYear();
    const month = String(selectedDate.getMonth() + 1).padStart(2, "0");
    const day = String(selectedDate.getDate()).padStart(2, "0");
    const dateStr = `${year}-${month}-${day}`;

    if (editingWorkshiftId) {
      updateMutation.mutate({
        id: editingWorkshiftId.toString(),
        work_shift_id: selectedWorkshift,
      });
    } else {
      registerMutation.mutate({
        employee_id: currentUser.user_id,
        work_shift_id: selectedWorkshift,
        date: dateStr,
      });
    }
  };

  const handleDelete = (workshiftId: number) => {
    if (window.confirm("Bạn có chắc chắn muốn xóa ca làm việc này?")) {
      deleteMutation.mutate(workshiftId);
    }
  };

  const getWorkshiftForDate = (date: Date): EmployeeWorkshift | undefined => {
    if (!employeeWorkshifts) return undefined;

    const dateStr = date.toLocaleDateString("en-CA");

    return employeeWorkshifts.find((ws) => {
      const wsDate = new Date(ws.date);
      const wsDateStr = wsDate.toLocaleDateString("en-CA");
      return wsDateStr === dateStr;
    });
  };

  const getWorkshiftsForDate = (date: Date): EmployeeWorkshift[] => {
    if (!employeeWorkshifts) return [];

    // 2. Lọc các ca làm việc
    return employeeWorkshifts.filter((ws) => {
      // Parse ngày từ API, NHƯNG chỉ lấy YYYY-MM-DD, bỏ qua timezone.
      // Giả sử `ws.date` là chuỗi có format "2025-09-11T17:00:00Z"
      const apiDateStr = ws.date;
      // Tách phần date (YYYY-MM-DD) ra khỏi chuỗi
      const [datePart] = apiDateStr.split("T"); // -> "2025-09-11"
      // Tách datePart thành các thành phần
      const [apiYear, apiMonth, apiDay] = datePart.split("-").map(Number); // -> [2025, 9, 11]

      // 3. So sánh năm, tháng, ngày một cách thủ công
      // Lưu ý: Tháng trong JS Date là 0-indexed, nên apiMonth - 1
      const apiDateObj = new Date(apiYear, apiMonth - 1, apiDay);

      // So sánh timestamp của hai ngày (đã được set cùng 1 múi giờ xác định)
      return apiDateObj.getTime() === date.getTime();
    });
  };

  const isDateInPast = (date: Date) => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    return date < today;
  };

  const getWorkshiftName = (workshiftId: string) => {
    // Ưu tiên tìm trong registerWorkshifts trước
    const workshiftFromRegister = registerWorkshifts?.find(
      (ws) => ws.workshift_id === workshiftId && ws.work_shift !== null
    );
    if (workshiftFromRegister) {
      return workshiftFromRegister.work_shift.workshift_name;
    }

    // Nếu không tìm thấy trong registerWorkshifts, thử tìm trong allWorkshifts
    const workshiftFromAll = allWorkshifts?.find(
      (ws) => ws.workshift_id === workshiftId
    );
    if (workshiftFromAll) {
      return workshiftFromAll.workshift_name;
    }

    // Nếu vẫn không tìm thấy, trả về 'Unknown'
    return "Unknown";
  };

  function getHourAndMinutesFromTime(timeString) {
    if (!/^\d{2}:\d{2}(:\d{2})?$/.test(timeString)) {
      throw new Error("Invalid time string");
    }
    return timeString.slice(0, 5); // keeps only HH:MM
  }

  return (
    <div className="p-2 sm:p-4 max-w-4xl mx-auto">
      <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-gray-800">
        Lịch làm việc
      </h2>

      {/* Month/Year Selector */}
      <div className="flex flex-col sm:flex-row justify-between items-center bg-gray-100 p-2 sm:p-3 rounded-lg mb-4 sm:mb-6 shadow-sm gap-2 sm:gap-0">
        <button
          onClick={handlePrev}
          className="w-full sm:w-auto bg-gray-800 text-white rounded-md sm:rounded-l-md sm:rounded-r-none border-r border-gray-100 py-2 hover:bg-red-700 hover:text-white px-3"
        >
          <div className="flex flex-row align-middle justify-center sm:justify-start">
            <svg
              className="w-5 mr-2"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fillRule="evenodd"
                d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z"
                clipRule="evenodd"
              ></path>
            </svg>
            <p className="ml-2">Trước</p>
          </div>
        </button>

        {/* Month/Year selectors */}
        <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
          <select
            value={selectedMonth}
            onChange={handleMonthChange}
            className="..."
          >
            {months.map((month, index) => (
              <option key={month} value={index}>
                {month}
              </option>
            ))}
          </select>
          <select
            value={selectedYear}
            onChange={handleYearChange}
            className="..."
          >
            {years.map((year) => (
              <option key={year} value={year}>
                {year}
              </option>
            ))}
          </select>
        </div>

        <button
          onClick={handleNext}
          className="w-full sm:w-auto bg-gray-800 text-white rounded-md sm:rounded-r-md sm:rounded-l-none py-2 border-l border-gray-200 hover:bg-red-700 hover:text-white px-3"
        >
          <div className="flex flex-row align-middle justify-center sm:justify-start">
            <span className="mr-2">Sau</span>
            <svg
              className="w-5 ml-2"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fillRule="evenodd"
                d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z"
                clipRule="evenodd"
              ></path>
            </svg>
          </div>
        </button>
      </div>

      {/* Calendar Grid */}
      <div className="bg-white rounded-lg shadow overflow-hidden border">
        {/* Weekday Headers */}
        <div className="grid grid-cols-7 bg-gray-100">
          {["CN", "T2", "T3", "T4", "T5", "T6", "T7"].map((day) => (
            <div
              key={day}
              className="py-1 sm:py-2 text-center font-medium text-gray-600 text-xs sm:text-sm"
            >
              {day}
            </div>
          ))}
        </div>

        {/* Calendar Days */}
        <div className="grid grid-cols-7 gap-px bg-gray-200">
          {calendarDays.map((date, index) => {
            if (!date) {
              return (
                <div
                  key={index}
                  className="min-h-[60px] sm:min-h-[80px] md:min-h-[100px] bg-gray-100"
                />
              );
            }

            const day = date.getDate();
            const isToday = date.toDateString() === new Date().toDateString();
            const isPast = isDateInPast(date);
            const workshifts = getWorkshiftsForDate(date); // Đổi thành số nhiều
            console.log(workshifts);
            return (
              <div
                key={index}
                className={`min-h-[60px] sm:min-h-[80px] md:min-h-[100px] p-1 sm:p-2 ${
                  isToday
                    ? "border-2 border-red-500 bg-gray-300"
                    : isPast
                    ? "bg-gray-100"
                    : "bg-white hover:bg-gray-50 cursor-pointer"
                }`}
                onClick={() => {
                  if (!isPast && !isScheduleAuto) {
                    handleDateClick(date);
                  }
                }}
              >
                <div className="text-right font-medium text-xs sm:text-sm">
                  {day}
                </div>
                <div className="mt-1 text-xs">
                  {workshifts.length > 0 ? (
                    <div className="space-y-1">
                      {workshifts.map((workshift) => (
                        <div
                          key={workshift.id}
                          className="bg-blue-100 text-blue-800 p-1 rounded text-center"
                        >
                          <div className={`flex flex-col gap-3`}>
                            <p>{workshift.workshift.workshift_name}</p>
                            <p className={`text-black mb-2`}>
                              {getHourAndMinutesFromTime(
                                workshift.workshift.start_time
                              )}{" "}
                              -{" "}
                              {getHourAndMinutesFromTime(
                                workshift.workshift.end_time
                              )}
                            </p>
                          </div>
                          {!isPast && !isScheduleAuto && (
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                handleDelete(workshift.id);
                              }}
                              className="text-red-500 text-xs ml-1"
                            >
                              ×
                            </button>
                          )}
                        </div>
                      ))}
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
              {editingWorkshiftId ? "Cập nhật" : "Đăng ký"} ca làm việc ngày{" "}
              {selectedDate.toLocaleDateString("vi-VN")}
            </h3>

            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Chọn ca làm việc
              </label>
              <select
                value={selectedWorkshift}
                onChange={(e) => setSelectedWorkshift(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                disabled={
                  registerMutation.isPending || updateMutation.isPending
                }
              >
                <option value="">-- Chọn ca --</option>
                {getWorkshiftsForSelectedDay().map((ws) => (
                  <option key={ws.id} value={ws.id}>
                    {ws.name} ({ws.start} - {ws.end})
                  </option>
                ))}
              </select>
            </div>

            <div className="flex justify-end gap-2">
              <button
                onClick={() => setShowModal(false)}
                className="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400"
                disabled={
                  registerMutation.isPending || updateMutation.isPending
                }
              >
                Hủy
              </button>
              <button
                onClick={handleRegister}
                disabled={
                  isScheduleAuto ||
                  !selectedWorkshift ||
                  registerMutation.isPending ||
                  updateMutation.isPending
                }
                className={`px-4 py-2 rounded-md text-white ${
                  isScheduleAuto
                    ? "bg-gray-400 cursor-not-allowed"
                    : "bg-blue-600 hover:bg-blue-700"
                }`}
              >
                {registerMutation.isPending || updateMutation.isPending
                  ? "Đang xử lý..."
                  : editingWorkshiftId
                  ? "Cập nhật"
                  : "Đăng ký"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default RegisterWorkshift;
