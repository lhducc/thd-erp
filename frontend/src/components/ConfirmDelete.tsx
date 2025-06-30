import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "./ui/button";
import { TriangleAlert, X } from "lucide-react";
import { useState } from "react";

type Props = {
  deleteFn: () => void;
};

const ConfirmDelete = ({ deleteFn }: Props) => {
  const [open, setOpen] = useState(false);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger>
        <Button variant="outline">
          <X />
        </Button>
      </DialogTrigger>
      <DialogContent className={`w-md`}>
        <DialogHeader>
          <DialogTitle className="text-center">Xác nhận xóa?</DialogTitle>
        </DialogHeader>
        <div className="flex w-full items-center justify-center">
          <TriangleAlert className="size-20 text-red-500" />
        </div>
        <div className="flex w-full justify-center gap-10">
          <Button variant={"outline"} onClick={() => setOpen(false)}>
            Hủy
          </Button>
          <Button onClick={() => deleteFn()}>Xóa</Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default ConfirmDelete;
