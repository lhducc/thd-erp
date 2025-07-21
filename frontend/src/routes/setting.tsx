import PATH from "@/constants/Path.ts";
import OfficePage from "@/pages/HR/setting/OfficePage.tsx";
import DepartmentPage from "@/pages/HR/setting/DepartmentPage.tsx";
import PositionPage from "@/pages/HR/setting/PositionPage.tsx";
import JobTitlePage from "@/pages/HR/setting/JobTitlePage.tsx";
import SetupWorkScheduleRegister from "@/pages/HR/checkin/wok-schedule-register/SetupWorkScheduleRegister.tsx";
import SettingWorkScheduleRegister from "@/pages/HR/checkin/wok-schedule-register/SettingWorkScheduleRegister..tsx";
import HierarchyLevelPage from "@/pages/HR/setting/HierarchyLevelPage.tsx";
import SettingContractPage from "@/pages/HR/setting/SettingContractPage.tsx";
import SettingDecisionPage from "@/pages/HR/setting/SettingDecisionPage.tsx";
import AllowancePage from "@/pages/HR/setting/AllowancePage.tsx";
import EmployeeDocumentPage from "@/pages/HR/setting/EmployeeDocumentPage.tsx";
import InsuranceInformationPage from "@/pages/InsuranceInformationPage.tsx";
import WorkshiftPage from "@/pages/HR/checkin/WorkshiftPage.tsx";
import CreateShiftPage from "@/pages/CreateShiftPage.js";
import ShiftDetailPage from "@/pages/ShiftDetailPage.js";

export const settingRoutes = [
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
    // {
    //     path: PATH.SHIFT_MANAGEMENT,
    //     element: <ShiftManagement/>,
    // },
    {
        path: PATH.CREATESHIFT,
        element: <CreateShiftPage/>,
    },
    {
        path: PATH.SHIFTDETAILPAGE,
        element: <ShiftDetailPage/>,
    },
];
