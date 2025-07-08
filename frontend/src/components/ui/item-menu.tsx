import {
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarMenuSub,
    SidebarMenuSubButton,
    SidebarMenuSubItem
} from "@/components/ui/sidebar.tsx";
import {Collapsible, CollapsibleContent, CollapsibleTrigger} from "@/components/ui/collapsible.tsx";
import {ChevronRight} from "lucide-react";
import {Link} from "react-router-dom";
import type {MenuItem} from "@/types";

interface ItemMenuProps {
    items: MenuItem[];
}

const ItemMenu = ({items}: ItemMenuProps) => {
    return (
        <>
            {items.map((item) => (
                <SidebarMenuItem key={item.name}>
                    {item.child ? (
                        <Collapsible asChild className="group/collapsible">
                            <div>
                                <CollapsibleTrigger asChild>
                                    <SidebarMenuButton>
                                        {item.icon && <item.icon />}
                                        <span>{item.name}</span>
                                        <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                                    </SidebarMenuButton>
                                </CollapsibleTrigger>
                                <CollapsibleContent>
                                    <SidebarMenuSub>
                                        {item.child.map((subItem) => (
                                            <SidebarMenuSubItem key={subItem.name}>
                                                {subItem.child ? (
                                                    <ItemMenu items={[subItem]} />
                                                ) : (
                                                    <SidebarMenuSubButton asChild>
                                                        <Link to={subItem.url || "#"}>
                                                            {subItem.icon && <subItem.icon />}
                                                            <span>{subItem.name}</span>
                                                        </Link>
                                                    </SidebarMenuSubButton>
                                                )}
                                            </SidebarMenuSubItem>
                                        ))}
                                    </SidebarMenuSub>
                                </CollapsibleContent>
                            </div>
                        </Collapsible>
                    ) : (
                        <SidebarMenuButton asChild>
                            <Link to={item.url || "#"}>
                                {item.icon && <item.icon />}
                                <span>{item.name}</span>
                            </Link>
                        </SidebarMenuButton>
                    )}
                </SidebarMenuItem>
            ))}
        </>
    );
};

export default ItemMenu;