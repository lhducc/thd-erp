import type { ComponentProps } from "react";
import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarTrigger,
} from "../../ui/sidebar.tsx";
import logo from "@/assets/LogoSignInPage.svg";
import {
    House,
    Users,
    User,
    ReceiptText,
    Calendar, CalendarRange,
} from "lucide-react";
import PATH from "@/constants/Path.ts";
import { Link } from "react-router-dom";

const ClientSidebar = ({ ...props }: ComponentProps<typeof Sidebar>) => {
    const data = [
        {
            name: "Trang chủ",
            url: PATH.HOME,
            icon: House,
        },
        {
            name: "Hồ sơ nhân viên",
            url: PATH.PROFILE,
            icon: User,
        },
        {
            name: "Quản lý trực tiếp",
            url: PATH.PROFILE_MANAGER,
            icon: Users,
        },
        {
            name: "Chấm công",
            url: "attendant",
            icon: ReceiptText,
        },
        {
            name: "Lịch làm việc",
            url: "register-workshift",
            icon: Calendar,
        },
        {
            name: "Bảng công cá nhân",
            url: PATH.ATTENDANCE_REPORT,
            icon: CalendarRange,
        },
        {
            name: "Lịch sử chấm công",
            url: PATH.ATTENDANT_HISTORY,
            icon: Users,
        },
    ];

    return (
        <Sidebar collapsible="icon" {...props}>
            <SidebarHeader className="flex items-center justify-center">
                <div className="flex items-end justify-end w-full">
                    <SidebarTrigger />
                </div>
                <img src={logo} alt="logo" className="w-[200px]" />
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarMenu>
                        {data.map((item) => (
                            <SidebarMenuItem>
                                <SidebarMenuButton asChild>
                                    <Link to={item.url}>
                                        {item.icon && <item.icon />}
                                        <span>{item.name}</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        ))}
                    </SidebarMenu>
                </SidebarGroup>
            </SidebarContent>
        </Sidebar>
    );
};

export default ClientSidebar;
