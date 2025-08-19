import {useAuth} from "@/context/AuthContext.tsx";
import {lazy, type ReactNode} from "react";

const ManagementHistoryAttendance = lazy(() => import("@/pages/HR/checkin/ManagementAttendantHistory.tsx"))
const HistoryAttendance = lazy(() => import("@/pages/Client/history-attendance/HistoryAttendance.tsx"))
const HISTORY_ATTENDANCE: Record<string, ReactNode> = {
    admin: <ManagementHistoryAttendance/>,
    manager: <ManagementHistoryAttendance/>,
    employee: <HistoryAttendance/>,
};

const AttendanceHistory = () => {
    const {currentUser} = useAuth();
    if (!currentUser) return null;
    return <>{HISTORY_ATTENDANCE[currentUser.role]}</>;
};

export default AttendanceHistory;
