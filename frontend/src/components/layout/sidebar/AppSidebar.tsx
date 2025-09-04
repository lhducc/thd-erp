import {useAuth} from "@/context/AuthContext.tsx";
import HRSidebar from "@/components/layout/sidebar/HRSidebar.tsx";
import ClientSidebar from "@/components/layout/sidebar/ClientSidebar.tsx";
import type {ReactNode} from "react";
import ManagerSidbar from "@/components/layout/sidebar/ManagerSidbar.tsx";

const SIDEBAR_COMPONENTS: Record<string, ReactNode> = {
    admin: <HRSidebar/>,
    hr: <HRSidebar/>,
    employee: <ClientSidebar/>,
    manager: <ManagerSidbar/>,
};

const AppSidebar = () => {
    const {currentUser} = useAuth();
    if (!currentUser) return null;
    return <>{SIDEBAR_COMPONENTS[currentUser.role]}</>;
};


export default AppSidebar;
