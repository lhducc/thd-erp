import {useRoutes} from "react-router-dom";
import PATH from "@/constants/Path";
import SignInPage from "@/pages/SignInPage";
import FirstChangePasswordPage from "@/pages/FirstChangePasswordPage";
import HomePage from "@/pages/HomePage";
import ProfilePage from "@/pages/ProfilePage";
import OfficePage from "./pages/setting/OfficePage.tsx";
import DepartmentPage from "./pages/setting/DepartmentPage.tsx";
import PositionPage from "./pages/setting/PositionPage.tsx";
import JobTitlePage from "./pages/setting/JobTitlePage.tsx";
import HierarchyLevelPage from "./pages/setting/HierarchyLevelPage.tsx";
import SettingContractPage from "./pages/setting/ContractPage.tsx";
import ContractPage from "@/pages/ContractPage.tsx";
import EmployeeDocumentPage from "@/pages/setting/EmployeeDocumentPage.tsx";

const App = () => {
    const settingRoutes = [
        {
            path: PATH.SETTING_OFFICE,
            element: <OfficePage />,
        },
        {
            path: PATH.SETTING_DEPARTMENT,
            element: <DepartmentPage />,
        },
        {
            path: PATH.SETTING_POSITION,
            element: <PositionPage />,
        },
        {
            path: PATH.SETTING_JOB_TITLE,
            element: <JobTitlePage />,
        },
        {
            path: PATH.SETTING_HIERARCHY_LEVEL,
            element: <HierarchyLevelPage />,
        },
        {
            path: PATH.SETTING_CONTRACT,
            element: <SettingContractPage />
        },
        {
            path: PATH.SETTING_EMPLOYEE_DOCUMENT,
            element: <EmployeeDocumentPage />
        }
    ];

    return useRoutes([
        {
            path: PATH.SIGN_IN,
            element: <SignInPage />,
        },
        {
            path: PATH.FIRST_CHANGE_PASSWORD,
            element: <FirstChangePasswordPage />,
        },
        {
            path: PATH.HOME,
            element: <HomePage />,
            children: [
                {
                    index: true,
                    element: <div className="h-[1000px]">Home</div>,
                },
                {
                    path: PATH.PROFILE.slice(1),
                    element: <ProfilePage />,
                },
                {
                    path: PATH.CONTRACT.slice(1),
                    element: <ContractPage />,
                },
                ...settingRoutes.map((route) => ({
                    path: route.path.replace(`${PATH.SETTING}/`, "setting/"),
                    element: route.element,
                })),
            ],
        },
    ]);
};

export default App;
