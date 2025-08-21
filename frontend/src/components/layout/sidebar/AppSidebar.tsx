import {useAuth} from "@/context/AuthContext.tsx";
import HRSidebar from "@/components/layout/sidebar/HRSidebar.tsx";
import ClientSidebar from "@/components/layout/sidebar/ClientSidebar.tsx";
import type {ReactNode} from "react";

const SIDEBAR_COMPONENTS: Record<string, ReactNode> = {
    admin: <HRSidebar/>,
    manager: <HRSidebar/>,
    employee: <ClientSidebar/>,
};

const AppSidebar = () => {
    const {currentUser} = useAuth();
    console.log(currentUser)
    if (!currentUser) return null;
    return <>{SIDEBAR_COMPONENTS[currentUser.role]}</>;
};


export default AppSidebar;
