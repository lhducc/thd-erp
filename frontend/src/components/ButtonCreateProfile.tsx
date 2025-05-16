import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "./ui/button";
import { Plus } from "lucide-react";
import { useState } from "react";

const ButtonCreateProfile = () => {
  const [openIntern, setOpenIntern] = useState(false);
  const [openEmployee, setOpenEmployee] = useState(false);
  // const [openIntern, setOpenIntern] = useState(false);
  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger>
          <Button>
            <Plus />
            Thêm nhân viên
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem>Nhân viên chính thức</DropdownMenuItem>
          <DropdownMenuItem>Thực tập sinh</DropdownMenuItem>
          <DropdownMenuItem>Công tác viên</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  );
};

export default ButtonCreateProfile;
