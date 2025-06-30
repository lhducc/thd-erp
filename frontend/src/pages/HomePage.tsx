import AppSidebar from "@/components/AppSidebar";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { Outlet } from "react-router-dom";

const HomePage = () => {
  return (
    <SidebarProvider>
      <AppSidebar />
      <main className="w-full flex flex-col relative h-screen">
        <header className="bg-rose-500 sticky top-0">
          <SidebarTrigger />
        </header>
        <div className="flex-1">
          <Outlet />
        </div>
      </main>
    </SidebarProvider>
  );
};

export default HomePage;
