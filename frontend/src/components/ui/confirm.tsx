import {useState} from "react";
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import {TriangleAlert, X} from "lucide-react";

type Props = {
    message: string;
    btn: React.ReactNode;
    onConfirm: () => void;
}

const Confirm = ({message, btn, onConfirm} : Props) => {
    const [open, setOpen] = useState(false);

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger>
                {btn}
            </DialogTrigger>
            <DialogContent className={`w-md`}>
                <DialogHeader>
                    <DialogTitle className="text-center">{message}</DialogTitle>
                </DialogHeader>
                <div className="flex w-full justify-center gap-10">
                    <Button variant={"outline"} onClick={() => setOpen(false)}>
                        Hủy
                    </Button>
                    <Button onClick={() => onConfirm()}>Xác nhận</Button>
                </div>
            </DialogContent>
        </Dialog>
    );
};

export default Confirm;
