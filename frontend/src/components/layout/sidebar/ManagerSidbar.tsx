import { useState, type ComponentProps } from "react";
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
  Calendar,
  CalendarRange,
  History,
  ChevronDown,
  UserCog,
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
      name: "Lịch làm việc",
      url: "register-workshift",
      icon: Calendar,
    },
    {
      name: "Chấm công",
      url: "attendant",
      icon: ReceiptText,
    },
    {
      name: "Bảng công cá nhân",
      url: PATH.ATTENDANCE_REPORT,
      icon: CalendarRange,
    },
    {
      name: "Lịch sử chấm công",
      url: PATH.ATTENDANT_HISTORY,
      icon: History,
    },
    {
      name: "Quản lý trực tiếp",
      icon: UserCog,
      child: [
        {
          name: "Hồ sơ",
          url: PATH.PROFILE_MANAGER,
          icon: Users,
        },
        {
          name: "Lịch làm việc",
          url: PATH.WORKSHIFT_EMPLOYEE,
          icon: Calendar,
        },
        {
          name: "Lịch sử chấm công",
          url: PATH.ATTENDANT_HISTORY_MANAGER,
          icon: History,
        },
      ],
    },
  ];

  const [openMenu, setOpenMenu] = useState<string | null>(null);

  const handleToggle = (name: string) => {
    setOpenMenu((prev) => (prev === name ? null : name));
  };

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
            {data.map((item, index) => (
              <SidebarMenuItem key={index}>
                {item.child ? (
                  <>
                    <SidebarMenuButton
                      onClick={() => handleToggle(item.name)}
                      className="flex justify-between items-center"
                    >
                      <div className="flex items-center gap-2">
                        <item.icon className="w-5 h-5" />
                        <span>{item.name}</span>
                      </div>
                      <ChevronDown
                        className={`w-4 h-4 transition-transform ${
                          openMenu === item.name ? "rotate-180" : ""
                        }`}
                      />
                    </SidebarMenuButton>

                    {openMenu === item.name && (
                      <div className="ml-6 mt-1">
                        {item.child.map((sub, i) => (
                          <SidebarMenuButton key={i} asChild>
                            <Link
                              to={sub.url ?? "#"}
                              className="flex items-center gap-2 text-sm text-gray-700 hover:text-black"
                            >
                              <sub.icon className="w-4 h-4" />
                              <span>{sub.name}</span>
                            </Link>
                          </SidebarMenuButton>
                        ))}
                      </div>
                    )}
                  </>
                ) : (
                  <SidebarMenuButton asChild>
                    <Link to={item.url}>
                      <item.icon className="w-5 h-5" />
                      <span>{item.name}</span>
                    </Link>
                  </SidebarMenuButton>
                )}
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
};

export default ClientSidebar;
