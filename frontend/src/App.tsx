import {useRoutes} from "react-router-dom";
import PATH from "@/constants/Path";
import SignInPage from "@/pages/SignInPage";
import FirstChangePasswordPage from "@/pages/FirstChangePasswordPage";
import ProfilePage from "@/pages/ProfilePage.tsx";
import ContractPage from "@/pages/HR/ContractPage";
import DocumentDetail from "./pages/documentDetailPage";
import DecisionPage from "./pages/DecisionPage";
import DecisionDetail from "./pages/DecisionDetailPage";
import OfficePage from "./pages/HR/setting/OfficePage.tsx";
import DepartmentPage from "./pages/HR/setting/DepartmentPage.tsx";
import PositionPage from "./pages/HR/setting/PositionPage.tsx";
import JobTitlePage from "./pages/HR/setting/JobTitlePage.tsx";
import HierarchyLevelPage from "./pages/HR/setting/HierarchyLevelPage.tsx";
import InsuranceInformationPage from "@/pages/InsuranceInformationPage.tsx";
import SettingContractPage from "@/pages/setting/SettingContractPage.tsx";
import AllowancePage from "@/pages/setting/AllowancePage.tsx";
import EmployeeDocumentPage from "@/pages/setting/EmployeeDocumentPage.tsx";
import SettingDecisionPage from "@/pages/setting/SettingDecisionPage.tsx";
import ShiftManagement from "@/pages/ShiftManagementPage.tsx";
import CreateShiftPage from "@/pages/CreateShiftPage.tsx";
import ShiftDetailPage from "@/pages/ShiftDetailPage.tsx";
import HomePage from "@/pages/HomePage.tsx";
import DetailEmployeePage from "@/pages/detailProfilePage.tsx";
import AttendancePage from "@/pages/checkin/AttendancePage.tsx";
import WorkshiftPage from "@/pages/HR/checkin/WorkshiftPage.tsx";
import AttendanceManagementPage from "@/pages/HR/checkin/AttendanceManagementPage.tsx";
import ApproveAttendancePage from "@/pages/HR/checkin/ApproveAttendancePage.tsx";
import RegisterWorkshift from "@/pages/Client/RegisterWorkshift.tsx";
import SetupWorkshift from "@/pages/HR/checkin/SetupWorkshift.tsx";
import WorkScheduleRegister from "@/pages/HR/checkin/wok-schedule-register/WorkScheduleRegister.tsx";
import SettingWorkScheduleRegister from "@/pages/HR/checkin/wok-schedule-register/SettingWorkScheduleRegister..tsx";
import WorkScheduleAuto from "@/pages/HR/checkin/work-schedule-auto/WorkScheduleAuto.tsx";
import SetupWorkScheduleAuto from "@/pages/HR/checkin/work-schedule-auto/SetupWorkScheduleAuto.tsx";
import SettingWorkScheduleAuto from "@/pages/HR/checkin/work-schedule-auto/SettingWorkScheduleAuto.tsx";


const App = () => {
    const settingRoutes = [
        {
            path: PATH.SETTING_OFFICE,
            element: <OfficePage/>,
        },
        {
            path: PATH.SETTING_DEPARTMENT,
            element: <DepartmentPage/>,
        },
        {
            path: PATH.SETTING_POSITION,
            element: <PositionPage/>,
        },
        {
            path: PATH.SETTING_JOB_TITLE,
            element: <JobTitlePage/>,
        },
        {
            path: "/schedule",
            element: <WorkScheduleRegister />
        },
        {
            path: "/schedule-setting",
            element: <SettingWorkScheduleRegister />
        },
        {
            path: PATH.SETTING_HIERARCHY_LEVEL,
            element: <HierarchyLevelPage/>,
        },
        {
            path: PATH.SETTING_CONTRACT,
            element: <SettingContractPage/>,
        },
        {
            path: PATH.SETTING_DECISION,
            element: <SettingDecisionPage/>,
        },
        {
            path: PATH.SETTING_ALLOWANCE,
            element: <AllowancePage/>,
        },
        {
            path: PATH.SETTING_EMPLOYEE_DOCUMENT,
            element: <EmployeeDocumentPage/>,
        },
        {
            path: PATH.INSURANCE,
            element: <InsuranceInformationPage/>,
        },
        {
            path: PATH.WORKSHIFT,
            element: <WorkshiftPage/>,
        },
        {
            path: PATH.SHIFT_MANAGEMENT,
            element: <ShiftManagement/>,
        },
        {
            path: PATH.CREATESHIFT,
            element: <CreateShiftPage/>,
        },
        {
            path: PATH.SHIFTDETAILPAGE,
            element: <ShiftDetailPage/>,
        },
    ];

    return useRoutes([
        {
            path: PATH.SIGN_IN,
            element: <SignInPage/>,
        },
        {
            path: PATH.FIRST_CHANGE_PASSWORD,
            element: <FirstChangePasswordPage/>,
        },
        {
            path: PATH.HOME,
            element: <HomePage/>,
            children: [
                {
                    index: true,
                    element: <div className="h-[1000px]">
                        <img src="https://thd-erp-web.thdcybersecurity.com/assets/THDHome-D95aRm-f.jpg" alt=""/>
                    </div>,
                },
                {
                    path: PATH.PROFILE.slice(1),
                    element: <ProfilePage/>,
                },
                {
                    path: PATH.CONTRACT.slice(1),
                    element: <ContractPage/>,
                },
                {
                    path: PATH.DETAIL_INFO.slice(1),
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
                    path: `${PATH.WORK_SCHEDULE}/:id/setting`,
                    element: <SettingWorkScheduleAuto />,
                },
                {
                    path: PATH.APPROVE_ATTENDANCE.slice(1),
                    element: <ApproveAttendancePage />,
                },
                {
                    path: PATH.DOCUMENTDETAIL.slice(1),
                    element: <DocumentDetail/>,
                },{
                    path: PATH.DETAIL_INFO.slice(1),
                    element: <DetailEmployeePage/>,
                },
                {
                    path: PATH.DOCUMENT_DETAIL.slice(1),
                    element: <DocumentDetail/>,
                },
                {
                    path: PATH.DECISION.slice(1),
                    element: <DecisionPage/>,
                },
                {
                    path: PATH.ATTENDANCE_MANAGEMENT.slice(1),
                    element: <AttendanceManagementPage />,
                },
                {
                    path: PATH.DETAIL_DECISION.slice(1),
                    element: <DecisionDetail/>,
                },
                {
                    path: PATH.SETUP_WORKSHIFT.slice(1),
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
                ...settingRoutes.map((route) => ({
                    path: route.path.replace(`${PATH.SETTING}/`, "setting/"),
                    element: route.element,
                })),
            ],
        },
    ])
}

export default App;
