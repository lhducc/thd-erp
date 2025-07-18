import type {ComponentProps} from "react";
import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarMenuSub,
    SidebarMenuSubButton,
    SidebarMenuSubItem, SidebarTrigger,
} from "../../ui/sidebar.tsx";
import logo from "@/assets/LogoSignInPage.svg";
import {
    House,
    Users,
    ReceiptText,
    SquareCheckBig,
    Settings,
    ChevronRight,
} from "lucide-react";
import PATH from "@/constants/Path.ts";
import {
    Collapsible,
    CollapsibleContent,
    CollapsibleTrigger,
} from "../../ui/collapsible.tsx";
import {Link} from "react-router-dom";
import checkinIcon from "@/assets/checkin-management.png";
import type {MenuItem} from "@/types";
import ItemMenu from "@/components/ui/item-menu.tsx";

const CheckinIcon = () => <img className={`w-4 h-4`} src={checkinIcon} alt="checkin icon"/>;

const HRSidebar = ({...props}: ComponentProps<typeof Sidebar>) => {
    const setting: MenuItem[] = [
        {
            name: "Trang chủ",
            url: PATH.HOME,
            icon: House,
        },
        {
            name: "Hồ sơ nhân viên",
            url: PATH.PROFILE,
            icon: Users,
        },
        {
            name: "Quản lý Check-in",
            icon: CheckinIcon,
            child: [
                {
                    name: "Quản lý ca làm việc",
                    icon: null,
                    url: PATH.WORKSHIFT,
                },
                {
                    name: "Lịch làm việc",
                    icon: null,
                    url: PATH.WORK_SCHEDULE,
                },
                {
                    name: "Lịch làm việc đăng ký",
                    icon: null,
                    url: PATH.WORK_SCHEDULE_REGISTER,
                },
                {
                    name: "Thiết lập chấm công",
                    icon: null,
                    url: PATH.ATTENDANCE_MANAGEMENT,
                },
                {
                    name: "Thiết lập đăng ký ca",
                    icon: null,
                    url: PATH.SETUP_WORKSHIFT,
                },
                {
                    name: "Quản lý chấm công",
                    icon: null,
                    url: PATH.ATTENDANCE_MANAGEMENT,
                    child: [
                        {
                            name: "Phê duyệt chấm công",
                            icon: null,
                            url: PATH.APPROVE_ATTENDANCE,
                        },
                    ]
                },
            ]
        },
        {
            name: "Hợp đồng",
            url: PATH.CONTRACT,
            icon: ReceiptText,
        },
        {
            name: "Quyết định",
            url: PATH.DECISION,
            icon: SquareCheckBig,
        },
    ];

    return (
        <Sidebar collapsible="icon" {...props}>
            <SidebarHeader className="flex items-center justify-center">
                <div className="flex items-end justify-end w-full">
                    <SidebarTrigger/>
                </div>
                <img src={logo} alt="logo" className="w-[200px]"/>
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarMenu>
                        <ItemMenu items={setting} />
                        <Collapsible asChild className="group/collapsible">
                            <SidebarMenuItem>
                                <CollapsibleTrigger asChild>
                                    <SidebarMenuButton>
                                        <Settings/>
                                        <span>Cài đặt</span>
                                        <ChevronRight
                                            className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"/>
                                    </SidebarMenuButton>
                                </CollapsibleTrigger>
                                <CollapsibleContent>
                                    <SidebarMenuSub>
                                        <Collapsible>
                                            <SidebarMenuButton>
                                                <CollapsibleTrigger className="w-full text-start">
                                                    Công ty
                                                </CollapsibleTrigger>
                                            </SidebarMenuButton>
                                            <CollapsibleContent>
                                                <SidebarMenuSub>
                                                    <SidebarMenuSubItem>
                                                        <SidebarMenuSubButton asChild>
                                                            <Link to={PATH.SETTING_OFFICE}>
                                                                <span>Văn Phòng</span>
                                                            </Link>
                                                        </SidebarMenuSubButton>
                                                    </SidebarMenuSubItem>
                                                    <SidebarMenuSubItem>
                                                        <SidebarMenuSubButton asChild>
                                                            <Link to={PATH.SETTING_DEPARTMENT}>
                                                                <span>Phòng ban</span>
                                                            </Link>
                                                        </SidebarMenuSubButton>
                                                    </SidebarMenuSubItem>
                                                    <SidebarMenuSubItem>
                                                        <SidebarMenuSubButton asChild>
                                                            <Link to={PATH.SETTING_POSITION}>
                                                                <span>Chức danh</span>
                                                            </Link>
                                                        </SidebarMenuSubButton>
                                                    </SidebarMenuSubItem>
                                                </SidebarMenuSub>
                                            </CollapsibleContent>
                                        </Collapsible>
                                        <Collapsible>
                                            <SidebarMenuButton>
                                                <CollapsibleTrigger className="w-full text-start">
                                                    Tài liệu
                                                </CollapsibleTrigger>
                                            </SidebarMenuButton>
                                            <CollapsibleContent>
                                                <SidebarMenuSub>
                                                    <SidebarMenuSubItem>
                                                        <SidebarMenuSubButton asChild>
                                                            <Link to={PATH.SETTING_CONTRACT}>
                                                                <span>Hợp đồng</span>
                                                            </Link>
                                                        </SidebarMenuSubButton>
                                                    </SidebarMenuSubItem>
                                                    <SidebarMenuSubItem>
                                                        <SidebarMenuSubButton asChild>
                                                            <Link to={PATH.SETTING_DECISION}>
                                                                <span>Quyết định</span>
                                                            </Link>
                                                        </SidebarMenuSubButton>
                                                    </SidebarMenuSubItem>
                                                    <SidebarMenuSubItem>
                                                        <SidebarMenuSubButton asChild>
                                                            <Link to={PATH.INSURANCE}>
                                                                <span>Thông tin bảo hiểm</span>
                                                            </Link>
                                                        </SidebarMenuSubButton>
                                                    </SidebarMenuSubItem>
                                                    <SidebarMenuSubItem>
                                                        <SidebarMenuSubButton asChild>
                                                            <Link to={PATH.SETTING_EMPLOYEE_DOCUMENT}>
                                                                <span>Tài liệu nhân sự</span>
                                                            </Link>
                                                        </SidebarMenuSubButton>
                                                    </SidebarMenuSubItem>
                                                </SidebarMenuSub>
                                            </CollapsibleContent>
                                        </Collapsible>
                                    </SidebarMenuSub>
                                </CollapsibleContent>
                            </SidebarMenuItem>
                        </Collapsible>
                    </SidebarMenu>
                </SidebarGroup>
            </SidebarContent>
        </Sidebar>
    );
};

export default HRSidebar;
