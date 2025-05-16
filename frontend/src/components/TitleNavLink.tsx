import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import PATH from "@/constants/Path";
import { Link, useLocation } from "react-router-dom";

export default function TitleNavLink() {
  function getCurrentPath() {
    const pathname = useLocation().pathname;
    [PATH.POSITION, PATH.JOB_TITLE, PATH.HIERARCHY_LEVEL].forEach((path) => {
      if (pathname.includes(path)) {
        return path;
      }
    });
    return "";
  }
  
  return (
    <Tabs defaultValue={getCurrentPath()} className="items-start">
      <TabsList className="text-foreground h-auto gap-2 rounded-none border-b bg-transparent px-0 py-1">
        <TabsTrigger
          value={PATH.POSITION}
          className="hover:bg-accent hover:text-foreground data-[state=active]:after:bg-primary data-[state=active]:hover:bg-accent relative after:absolute after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5 data-[state=active]:bg-transparent data-[state=active]:shadow-none"
        >
          <Link to={PATH.POSITION}>Vị trí công việc</Link>
        </TabsTrigger>
        <TabsTrigger
          value={PATH.JOB_TITLE}
          className="hover:bg-accent hover:text-foreground data-[state=active]:after:bg-primary data-[state=active]:hover:bg-accent relative after:absolute after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5 data-[state=active]:bg-transparent data-[state=active]:shadow-none"
        >
          <Link to={PATH.JOB_TITLE}>Chức vụ</Link>
        </TabsTrigger>
        <TabsTrigger
          value={PATH.HIERARCHY_LEVEL}
          className="hover:bg-accent hover:text-foreground data-[state=active]:after:bg-primary data-[state=active]:hover:bg-accent relative after:absolute after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5 data-[state=active]:bg-transparent data-[state=active]:shadow-none"
        >
          <Link to={PATH.HIERARCHY_LEVEL}>Cấp bậc</Link>
        </TabsTrigger>
      </TabsList>
    </Tabs>
  );
}
