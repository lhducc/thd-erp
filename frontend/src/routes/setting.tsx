import { lazy, Suspense } from "react";
import PATH from "@/constants/Path.ts";
import Loading from "@/components/Loading.tsx";

const OfficePage = lazy(() => import("@/pages/HR/setting/OfficePage.tsx"));
const DepartmentPage = lazy(() => import("@/pages/HR/setting/DepartmentPage.tsx"));
const PositionPage = lazy(() => import("@/pages/HR/setting/PositionPage.tsx"));
const JobTitlePage = lazy(() => import("@/pages/HR/setting/JobTitlePage.tsx"));
const HierarchyLevelPage = lazy(() => import("@/pages/HR/setting/HierarchyLevelPage.tsx"));
const SettingContractPage = lazy(() => import("@/pages/HR/setting/SettingContractPage.tsx"));
const SettingDecisionPage = lazy(() => import("@/pages/HR/setting/SettingDecisionPage.tsx"));
const AllowancePage = lazy(() => import("@/pages/HR/setting/AllowancePage.tsx"));
const EmployeeDocumentPage = lazy(() => import("@/pages/HR/setting/EmployeeDocumentPage.tsx"));
const InsuranceInformationPage = lazy(() => import("@/pages/InsuranceInformationPage.tsx"));
const WorkshiftPage = lazy(() => import("@/pages/HR/checkin/workshift"));
const CreateShiftPage = lazy(() => import("@/pages/CreateShiftPage.js"));
const ShiftDetailPage = lazy(() => import("@/pages/ShiftDetailPage.js"));

export const settingRoutes = [
    {
        path: PATH.SETTING_OFFICE,
        element: (
            <Suspense fallback={<Loading />}>
                <OfficePage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_DEPARTMENT,
        element: (
            <Suspense fallback={<Loading />}>
                <DepartmentPage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_POSITION,
        element: (
            <Suspense fallback={<Loading />}>
                <PositionPage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_JOB_TITLE,
        element: (
            <Suspense fallback={<Loading />}>
                <JobTitlePage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_HIERARCHY_LEVEL,
        element: (
            <Suspense fallback={<Loading />}>
                <HierarchyLevelPage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_CONTRACT,
        element: (
            <Suspense fallback={<Loading />}>
                <SettingContractPage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_DECISION,
        element: (
            <Suspense fallback={<Loading />}>
                <SettingDecisionPage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_ALLOWANCE,
        element: (
            <Suspense fallback={<Loading />}>
                <AllowancePage />
            </Suspense>
        ),
    },
    {
        path: PATH.SETTING_EMPLOYEE_DOCUMENT,
        element: (
            <Suspense fallback={<Loading />}>
                <EmployeeDocumentPage />
            </Suspense>
        ),
    },
    {
        path: PATH.INSURANCE,
        element: (
            <Suspense fallback={<Loading />}>
                <InsuranceInformationPage />
            </Suspense>
        ),
    },
    {
        path: PATH.WORKSHIFT,
        element: (
            <Suspense fallback={<Loading />}>
                <WorkshiftPage />
            </Suspense>
        ),
    },
    {
        path: PATH.CREATESHIFT,
        element: (
            <Suspense fallback={<Loading />}>
                <CreateShiftPage />
            </Suspense>
        ),
    },
    {
        path: PATH.SHIFTDETAILPAGE,
        element: (
            <Suspense fallback={<Loading />}>
                <ShiftDetailPage />
            </Suspense>
        ),
    },
];