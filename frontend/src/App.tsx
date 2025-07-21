import {useRoutes} from "react-router-dom";
import PATH from "@/constants/Path";
import SignInPage from "@/pages/HR/auth/SignInPage.tsx";
import FirstChangePasswordPage from "@/pages/FirstChangePasswordPage";
import {hrRoutes} from "@/routes/hr.tsx";

const App = () => {
    return useRoutes([
        {
            path: PATH.SIGN_IN,
            element: <SignInPage/>,
        },
        {
            path: PATH.FIRST_CHANGE_PASSWORD,
            element: <FirstChangePasswordPage/>,
        },
        ...hrRoutes
    ])
}

export default App;
