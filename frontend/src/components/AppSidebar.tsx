import type { ComponentProps } from "react";
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
  SidebarMenuSubItem,
} from "./ui/sidebar";
import logo from "@/assets/LogoSignInPage.svg";
import {
  House,
  Users,
  ReceiptText,
  SquareCheckBig,
  Calendar,
  Settings,
  ChevronRight,
} from "lucide-react";
import PATH from "@/constants/Path";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "./ui/collapsible";
import { Link } from "react-router-dom";

const AppSidebar = ({ ...props }: ComponentProps<typeof Sidebar>) => {
  const data = [
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
      name: "Hợp đồng",
      url: PATH.CONTRACT,
      icon: ReceiptText,
    },
    {
      name: "Quyết định",
      url: PATH.DECISION,
      icon: SquareCheckBig,
    },
    {
      name: "Lịch thực tập",
      url: "#",
      icon: Calendar,
    },
  ];

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="flex items-center justify-center">
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
            <Collapsible asChild className="group/collapsible">
              <SidebarMenuItem>
                <CollapsibleTrigger asChild>
                  <SidebarMenuButton>
                    <Settings />
                    <span>Cài đặt</span>
                    <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
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

export default AppSidebar;
