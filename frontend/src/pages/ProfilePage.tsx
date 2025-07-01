import {useAuth} from "@/context/AuthContext.tsx";
import type {ReactNode} from "react";
import HRProfilePage from "@/pages/HR/HRProfilePage.tsx";
import ClientProfilePage from "@/pages/Client/ClientProfilePage.tsx";

const PROFILE_COMPONENTS: Record<string, ReactNode> = {
    admin: <HRProfilePage/>,
    HR: <HRProfilePage/>,
    Client: <ClientProfilePage/>,
};

const ProfilePage = () => {
    const {currentUser} = useAuth();
    if (!currentUser) return null;
    return <>{PROFILE_COMPONENTS[currentUser.roles]}</>;
};


export default ProfilePage;
