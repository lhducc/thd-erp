import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import PATH from "@/constants/Path";
import { Link, useLocation } from "react-router-dom";

export default function NavLinkAttendantHistory() {
  const location = useLocation();

  function getCurrentPath() {
    return (
      [PATH.ATTENDANT_HISTORY_BY_DATE, PATH.ATTENDANT_HISTORY].find((path) =>
        location.pathname.includes(path)
      ) || ""
    );
  }

  return (
    <Tabs defaultValue={getCurrentPath()} className="items-start">
      <TabsList className="text-foreground h-auto gap-2 rounded-none border-b bg-transparent px-0 py-1">
        <TabsTrigger
          value={PATH.ATTENDANT_HISTORY}
          className="hover:bg-accent hover:text-foreground data-[state=active]:after:bg-primary data-[state=active]:hover:bg-accent relative after:absolute after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5 data-[state=active]:bg-transparent data-[state=active]:shadow-none"
        >
          <Link to={PATH.ATTENDANT_HISTORY}>Danh sách nhân viên</Link>
        </TabsTrigger>
        <TabsTrigger
          value={PATH.ATTENDANT_HISTORY_BY_DATE}
          className="hover:bg-accent hover:text-foreground data-[state=active]:after:bg-primary data-[state=active]:hover:bg-accent relative after:absolute after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5 data-[state=active]:bg-transparent data-[state=active]:shadow-none"
        >
          <Link to={PATH.ATTENDANT_HISTORY_BY_DATE}>Dữ liệu theo ngày</Link>
        </TabsTrigger>
      </TabsList>
    </Tabs>
  );
}
