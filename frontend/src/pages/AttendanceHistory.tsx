import { useAuth } from "@/context/AuthContext.tsx";
import { lazy, type ReactNode } from "react";
import ManagementAttendantHistory from "@/pages/HR/checkin/ManagementAttendantHistory.tsx";

const ManagementHistoryAttendance = lazy(
  () => import("@/pages/HR/checkin/ManagementAttendantHistory.tsx")
);
const AttendanceReport = lazy(
  () => import("@/pages/Client/HistoryAttendance.tsx")
);
const HISTORY_ATTENDANCE: Record<string, ReactNode> = {
  admin: <ManagementHistoryAttendance />,
  manager: <ManagementAttendantHistory />,
  employee: <AttendanceReport />,
};
const AttendanceHistory = () => {
  const { currentUser } = useAuth();
  console.log(currentUser?.role);
  if (!currentUser) return null;
  return <>{HISTORY_ATTENDANCE[currentUser.role]}</>;
};

export default AttendanceHistory;
