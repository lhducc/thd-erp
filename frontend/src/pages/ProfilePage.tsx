import {useAuth} from "@/context/AuthContext.tsx";
import {lazy, type ReactNode} from "react";

const HRProfilePage = lazy(() => import("@/pages/HR/profile/index.tsx"))
const ClientProfilePage = lazy(() => import("@/pages/Client/ClientProfilePage.tsx"))
const PROFILE_COMPONENTS: Record<string, ReactNode> = {
    admin: <HRProfilePage/>,
    manager: <ClientProfilePage/>,
    employee: <ClientProfilePage/>,
};

const ProfilePage = () => {
    const {currentUser} = useAuth();
    if (!currentUser) return null;
    return <>{PROFILE_COMPONENTS[currentUser.role]}</>;
};


export default ProfilePage;
