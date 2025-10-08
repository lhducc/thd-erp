import { lazy, Suspense } from "react";
import PATH from "@/constants/Path";
import { settingRoutes } from "@/routes/setting";
import Loading from "@/components/Loading";
import AttendanceReport from "@/pages/Client/attendance-report/AttendanceReport.tsx";
import HistoryAttendance from "@/pages/Client/HistoryAttendance.tsx";
import AttendanceHistory from "@/pages/AttendanceHistory.tsx";
import ManagementAttendantHistoryByDate from "@/pages/HR/checkin/ManagementAttendantHistoryByDate";

const HomePage = lazy(() => import("@/pages/HomePage"));
const ProfilePage = lazy(() => import("@/pages/ProfilePage"));
const ContractPage = lazy(() => import("@/pages/HR/ContractPage"));
const DetailEmployeePage = lazy(
  () => import("@/pages/HR/profile/DetailEmployeePage.tsx")
);
const WorkScheduleAuto = lazy(
  () => import("@/pages/HR/checkin/work-schedule-auto/WorkScheduleAuto")
);
const WorkScheduleRegister = lazy(
  () =>
    import("@/pages/HR/checkin/wok-schedule-register/WorkScheduleRegisterPage")
);
const SetupWorkScheduleAuto = lazy(
  () => import("@/pages/HR/checkin/work-schedule-auto/SetupWorkScheduleAuto")
);
const SetupWorkScheduleRegister = lazy(
  () =>
    import("@/pages/HR/checkin/wok-schedule-register/SetupWorkScheduleRegister")
);
const SettingWorkScheduleAuto = lazy(
  () => import("@/pages/HR/checkin/work-schedule-auto/SettingWorkScheduleAuto")
);
const SettingWorkScheduleRegister = lazy(
  () =>
    import(
      "@/pages/HR/checkin/wok-schedule-register/SettingWorkScheduleRegister"
    )
);
const ApproveAttendancePage = lazy(
  () =>
    import(
      "@/pages/HR/checkin/approve-attendance/page/ApproveAttendanceOverviewPage"
    )
);
const ApproveAttendanceDetailPage = lazy(
  () =>
    import(
      "@/pages/HR/checkin/approve-attendance/page/ApproveAttendanceDetailPage"
    )
);
const DocumentDetail = lazy(() => import("@/pages/documentDetailPage"));
const DecisionPage = lazy(() => import("@/pages/DecisionPage"));
const AttendanceManagementPage = lazy(
  () => import("@/pages/HR/checkin/AttendanceManagementPage")
);
const DecisionDetail = lazy(() => import("@/pages/DecisionDetailPage"));
const AttendancePage = lazy(
  () => import("@/pages/Client/check-in/AttendancePage")
);
const RegisterWorkshift = lazy(
  () => import("@/pages/Client/RegisterWorkshift")
);
const Rota = lazy(() => import("@/pages/HR/checkin/rota/Rota.tsx"));
// const AttendantHistory = lazy(() => import("@/pages/AttendanceHistory.tsx"));
const AttendantHistoryDetail = lazy(
  () => import("@/pages/HR/checkin/AttendantHistoryDetail")
);
const TimesheetPage = lazy(
  () => import("@/pages/HR/checkin/timesheet/TimesheetPage")
);
const TimesheetDetail = lazy(
  () => import("@/pages/HR/checkin/timesheet/TimesheetDetail")
);
const ManagerEmployee = lazy(
  () => import("@/pages/Manager/ManagerEmployeePage.tsx")
);
const ManagerEmployeeHistory = lazy(
  () => import("@/pages/Manager/ManagerEmployeeHistoryPage.tsx")
);

export const hrRoutes = [
  {
    path: PATH.HOME,
    element: (
      <Suspense fallback={<Loading />}>
        <HomePage />
      </Suspense>
    ),
    children: [
      {
        index: true,
        element: (
          <div className="h-[1000px]">
            <Suspense fallback={<Loading />}>
              <img
                src="https://thd-erp-web.thdcybersecurity.com/assets/THDHome-D95aRm-f.jpg"
                alt=""
              />
            </Suspense>
          </div>
        ),
      },
      {
        path: PATH.PROFILE,
        element: (
          <Suspense fallback={<Loading />}>
            <ProfilePage />
          </Suspense>
        ),
      },
      {
        path: PATH.CONTRACT,
        element: (
          <Suspense fallback={<Loading />}>
            <ContractPage />
          </Suspense>
        ),
      },
      {
        path: `${PATH.APPROVE}/:id`,
        element: (
          <Suspense fallback={<Loading />}>
            <ApproveAttendanceDetailPage />
          </Suspense>
        ),
      },
      {
        path: PATH.ATTENDANCE_REPORT,
        element: (
          <Suspense fallback={<Loading />}>
            <AttendanceReport />
          </Suspense>
        ),
      },
      {
        path: `${PATH.DETAIL_INFO}/:id`,
        element: (
          <Suspense fallback={<Loading />}>
            <DetailEmployeePage />
          </Suspense>
        ),
      },
      {
        path: PATH.WORK_SCHEDULE,
        element: (
          <Suspense fallback={<Loading />}>
            <WorkScheduleAuto />
          </Suspense>
        ),
      },
      {
        path: PATH.WORK_SCHEDULE_REGISTER,
        element: (
          <Suspense fallback={<Loading />}>
            <WorkScheduleRegister />
          </Suspense>
        ),
      },
      {
        path: `${PATH.WORK_SCHEDULE}/:id`,
        element: (
          <Suspense fallback={<Loading />}>
            <SetupWorkScheduleAuto />
          </Suspense>
        ),
      },
      {
        path: `${PATH.WORK_SCHEDULE_REGISTER}/:id`,
        element: (
          <Suspense fallback={<Loading />}>
            <SetupWorkScheduleRegister />
          </Suspense>
        ),
      },
      {
        path: `${PATH.WORK_SCHEDULE}/:id/setting`,
        element: (
          <Suspense fallback={<Loading />}>
            <SettingWorkScheduleAuto />
          </Suspense>
        ),
      },
      {
        path: `${PATH.WORK_SCHEDULE_REGISTER}/:id/setting`,
        element: (
          <Suspense fallback={<Loading />}>
            <SettingWorkScheduleRegister />
          </Suspense>
        ),
      },
      {
        path: PATH.APPROVE_ATTENDANCE,
        element: (
          <Suspense fallback={<Loading />}>
            <ApproveAttendancePage />
          </Suspense>
        ),
      },
      {
        path: PATH.DOCUMENTDETAIL,
        element: (
          <Suspense fallback={<Loading />}>
            <DocumentDetail />
          </Suspense>
        ),
      },
      {
        path: PATH.TIMESHEET,
        element: (
          <Suspense fallback={<Loading />}>
            <TimesheetPage />
          </Suspense>
        ),
      },
      {
        path: `${PATH.TIMESHEET}/:id`,
        element: (
          <Suspense fallback={<Loading />}>
            <TimesheetDetail />
          </Suspense>
        ),
      },
      {
        path: PATH.ROTA,
        element: (
          <Suspense fallback={<Loading />}>
            <Rota />
          </Suspense>
        ),
      },
      {
        path: PATH.ATTENDANT_HISTORY,
        element: (
          <Suspense fallback={<Loading />}>
            <AttendanceHistory />
          </Suspense>
        ),
      },
      {
        path: PATH.ATTENDANT_HISTORY_BY_DATE,
        element: (
          <Suspense fallback={<Loading />}>
            <ManagementAttendantHistoryByDate />
          </Suspense>
        ),
      },
      {
        path: `${PATH.ATTENDANT_HISTORY}/:id`,
        element: (
          <Suspense fallback={<Loading />}>
            <AttendantHistoryDetail />
          </Suspense>
        ),
      },
      {
        path: PATH.DOCUMENT_DETAIL,
        element: (
          <Suspense fallback={<Loading />}>
            <DocumentDetail />
          </Suspense>
        ),
      },
      {
        path: PATH.DECISION,
        element: (
          <Suspense fallback={<Loading />}>
            <DecisionPage />
          </Suspense>
        ),
      },
      {
        path: `${PATH.DECISION}/:id`,
        element: (
          <Suspense fallback={<Loading />}>
            <DecisionDetail />
          </Suspense>
        ),
      },
      {
        path: PATH.ATTENDANCE_MANAGEMENT,
        element: (
          <Suspense fallback={<Loading />}>
            <AttendanceManagementPage />
          </Suspense>
        ),
      },
      {
        path: PATH.DETAIL_DECISION,
        element: (
          <Suspense fallback={<Loading />}>
            <DecisionDetail />
          </Suspense>
        ),
      },
      {
        path: PATH.CHECKIN,
        element: (
          <Suspense fallback={<Loading />}>
            <AttendancePage />
          </Suspense>
        ),
      },
      {
        path: PATH.REGISTER_WORKSHIFT,
        element: (
          <Suspense fallback={<Loading />}>
            <RegisterWorkshift />
          </Suspense>
        ),
      },
      {
        path: PATH.PROFILE_MANAGER,
        element: (
          <Suspense fallback={<Loading />}>
            <ManagerEmployee />
          </Suspense>
        ),
      },
      {
        path: PATH.ATTENDANT_HISTORY_MANAGER,
        element: (
          <Suspense fallback={<Loading />}>
            <ManagerEmployeeHistory />
          </Suspense>
        ),
      },
      ...settingRoutes.map((route) => ({
        ...route,
        element: <Suspense fallback={<Loading />}>{route.element}</Suspense>,
      })),
    ],
  },
];
