import PATH from "@/constants/Path.ts";
import HomePage from "@/pages/HomePage.tsx";
import ProfilePage from "@/pages/ProfilePage.tsx";
import ContractPage from "@/pages/HR/ContractPage.tsx";
import DetailEmployeePage from "@/pages/detailProfilePage.tsx";
import WorkScheduleAuto from "@/pages/HR/checkin/work-schedule-auto/WorkScheduleAuto.tsx";
import SetupWorkScheduleRegister from "@/pages/HR/checkin/wok-schedule-register/SetupWorkScheduleRegister.tsx";
import SetupWorkScheduleAuto from "@/pages/HR/checkin/work-schedule-auto/SetupWorkScheduleAuto.tsx";
import SettingWorkScheduleAuto from "@/pages/HR/checkin/work-schedule-auto/SettingWorkScheduleAuto.tsx";
import ApproveAttendancePage from "@/pages/HR/checkin/ApproveAttendancePage.tsx";
import DocumentDetail from "@/pages/documentDetailPage.tsx";
import DecisionPage from "@/pages/DecisionPage.tsx";
import AttendanceManagementPage from "@/pages/HR/checkin/AttendanceManagementPage.tsx";
import DecisionDetail from "@/pages/DecisionDetailPage.tsx";
import SetupWorkshift from "@/pages/HR/checkin/SetupWorkshift.tsx";
import AttendancePage from "@/pages/Client/check-in/AttendancePage.tsx";
import RegisterWorkshift from "@/pages/Client/RegisterWorkshift.tsx";
import {settingRoutes} from "@/routes/setting.tsx";
import {Suspense} from "react";
import Loading from "@/components/Loading.tsx";
import WorkScheduleRegister from "@/pages/HR/checkin/wok-schedule-register/WorkScheduleRegister.tsx";
import Rota from "@/pages/HR/checkin/Rota.tsx";
import SettingWorkScheduleRegister from "@/pages/HR/checkin/wok-schedule-register/SettingWorkScheduleRegister..tsx";
import AttendantHistory from "@/pages/HR/checkin/AttendantHistory.tsx";
import AttendantHistoryDetail from "@/pages/HR/checkin/AttendantHistoryDetail.tsx";

export const hrRoutes = [
    {
        path: PATH.HOME,
        element: <HomePage/>,
        children: [
            {
                index: true,
                element: <div className="h-[1000px]">
                    <Suspense fallback={<Loading />}>
                        <img src="https://thd-erp-web.thdcybersecurity.com/assets/THDHome-D95aRm-f.jpg" alt=""/>
                    </Suspense>
                </div>,
            },
            {
                path: PATH.PROFILE,
                element: <ProfilePage/>,
            },
            {
                path: PATH.CONTRACT,
                element: <ContractPage/>,
            },
            {
                path: `${PATH.DETAIL_INFO}/:id`,
                element: <DetailEmployeePage/>,
            },
            {
                path: PATH.WORK_SCHEDULE,
                element: <WorkScheduleAuto />,
            },
            {
                path: PATH.WORK_SCHEDULE_REGISTER,
                element: <WorkScheduleRegister />,
            },
            {
                path: `${PATH.WORK_SCHEDULE}/:id`,
                element: <SetupWorkScheduleAuto />,
            },
            {
                path: `${PATH.WORK_SCHEDULE_REGISTER}/:id`,
                element: <SetupWorkScheduleRegister />,
            },
            {
                path: `${PATH.WORK_SCHEDULE}/:id/setting`,
                element: <SettingWorkScheduleAuto />,
            },
            {
                path: `${PATH.WORK_SCHEDULE_REGISTER}/:id/setting`,
                element: <SettingWorkScheduleRegister />,
            },
            {
                path: PATH.APPROVE_ATTENDANCE,
                element: <ApproveAttendancePage />,
            },
            {
                path: PATH.DOCUMENTDETAIL,
                element: <DocumentDetail/>,
            },
            {
                path: PATH.ROTA,
                element: <Rota/>,
            },
            {
                path: PATH.ATTENDANT_HISTORY,
                element: <AttendantHistory />,
            },
            {
                path: `${PATH.ATTENDANT_HISTORY}/:id`,
                element: <AttendantHistoryDetail />,
            },
            {
                path: PATH.DOCUMENT_DETAIL,
                element: <DocumentDetail/>,
            },
            {
                path: PATH.DECISION,
                element: <DecisionPage/>,
            },
            {
                path: PATH.ATTENDANCE_MANAGEMENT,
                element: <AttendanceManagementPage />,
            },
            {
                path: PATH.DETAIL_DECISION,
                element: <DecisionDetail/>,
            },
            {
                path: PATH.SETUP_WORKSHIFT,
                element: <SetupWorkshift />,
            },
            {
                path: PATH.CHECKIN,
                element: <AttendancePage/>,
            },
            {
                path: PATH.REGISTER_WORKSHIFT,
                element: <RegisterWorkshift/>,
            },
            ...settingRoutes,
        ],
    },
]