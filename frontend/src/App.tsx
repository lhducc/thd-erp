import { useRoutes } from "react-router-dom";
import PATH from "@/constants/Path";
import SignInPage from "@/pages/SignInPage";
import FirstChangePasswordPage from "@/pages/FirstChangePasswordPage";
import HomePage from "@/pages/HomePage";
import ProfilePage from "@/pages/ProfilePage";
import OfficePage from "./pages/OfficePage";
import DepartmentPage from "./pages/DepartmentPage";
import PositionPage from "./pages/PositionPage";
import JobTitlePage from "./pages/JobTitlePage";
import HierarchyLevelPage from "./pages/HierarchyLevelPage";
import ContractPage from "./pages/ContractPage";

const App = () => {
  return useRoutes([
    {
      path: PATH.HOME,
      element: <HomePage />,
      children: [
        {
          index: true,
          element: <div className="h-[1000px]">Home</div>,
        },
        {
          path: PATH.PROFILE,
          element: <ProfilePage />,
        },
        {
          path: PATH.OFFICE,
          element: <OfficePage />,
        },
        {
          path: PATH.CONTRACT,
          element: <ContractPage />,
        },
        {
          path: PATH.DEPARTMENT,
          element: <DepartmentPage />,
        },
        {
          path: PATH.POSITION,
          element: <PositionPage />,
        },
        {
          path: PATH.JOB_TITLE,
          element: <JobTitlePage />,
        },
        {
          path: PATH.HIERARCHY_LEVEL,
          element: <HierarchyLevelPage />,
        },
      ],
    },
    {
      path: PATH.SIGN_IN,
      element: <SignInPage />,
    },
    {
      path: PATH.FIRSRT_CHANGE_PASSWORD,
      element: <FirstChangePasswordPage />,
    },
  ]);
};

export default App;
