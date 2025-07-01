import {SidebarProvider} from "@/components/ui/sidebar";
import {Outlet} from "react-router-dom";
import AppNavbar from "@/components/layout/AppNavbar.tsx";
import AppSidebar from "@/components/layout/sidebar/AppSidebar.tsx";

const HomePage = () => {
    return (
        <SidebarProvider>
            <AppSidebar />
            <main className="w-full flex flex-col relative h-screen">
                <header className="flex items-center top-0">
                    <AppNavbar/>
                </header>
                <div className="p-8 flex-1 bg-[#F6F6F6]">
                    <Outlet/>
                </div>
            </main>
        </SidebarProvider>
    );
};

export default HomePage;
