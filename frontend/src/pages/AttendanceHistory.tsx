import {useAuth} from "@/context/AuthContext.tsx";
import {lazy, type ReactNode} from "react";

const ManagementHistoryAttendance = lazy(() => import("@/pages/HR/checkin/ManagementAttendantHistory.tsx"))
const AttendanceReport = lazy(() => import("@/pages/Client/attendance-report/AttendanceReport"))
const HISTORY_ATTENDANCE: Record<string, ReactNode> = {
    admin: <ManagementHistoryAttendance/>,
    manager: <ManagementHistoryAttendance/>,
    employee: <AttendanceReport/>,
};

const AttendanceHistory = () => {
    const {currentUser} = useAuth();
    if (!currentUser) return null;
    return <>{HISTORY_ATTENDANCE[currentUser.role]}</>;
};

export default AttendanceHistory;
