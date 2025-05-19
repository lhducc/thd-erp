import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "./ui/button";
import { Plus } from "lucide-react";
import CreateEmployeeForm from "./CreateEmployeeForm";
import { useState } from "react";

const ButtonCreateProfile = () => {
  const [open, setOpen] = useState(false);
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
          <DropdownMenuItem onClick={() => setOpen(true)}>
            Nhân viên chính thức
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setOpen(true)}>
            Thực tập sinh
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setOpen(true)}>
            Công tác viên
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <CreateEmployeeForm open={open} setOpen={setOpen} />
    </>
  );
};

export default ButtonCreateProfile;
