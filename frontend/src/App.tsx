import { lazy, Suspense } from "react";
import { useRoutes } from "react-router-dom";
import PATH from "@/constants/Path";
import Loading from "@/components/Loading";

const SignInPage = lazy(() => import("@/pages/HR/auth/SignInPage"));
const FirstChangePasswordPage = lazy(() => import("@/pages/FirstChangePasswordPage"));

import { hrRoutes } from "@/routes/hr";

const App = () => {
    return useRoutes([
        {
            path: PATH.SIGN_IN,
            element: (
                <Suspense fallback={<Loading />}>
                    <SignInPage />
                </Suspense>
            ),
        },
        {
            path: PATH.FIRST_CHANGE_PASSWORD,
            element: (
                <Suspense fallback={<Loading />}>
                    <FirstChangePasswordPage />
                </Suspense>
            ),
        },
        ...hrRoutes
    ]);
};

export default App;