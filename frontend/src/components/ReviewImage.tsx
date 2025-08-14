import {
    Dialog,
    DialogTrigger,
    DialogContent,
    DialogOverlay,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { useState } from "react";

const ReviewImage = ({ url }: { url: string }) => {
    const [open, setOpen] = useState(false);

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button variant="outline" onClick={() => setOpen(true)}>
                    Review
                </Button>
            </DialogTrigger>

            <DialogContent
                className="fixed left-1/2 top-1/2 z-50 grid w-full max-w-2xl translate-x-[-50%] translate-y-[-50%] gap-4 border border-gray-200 bg-white p-6 shadow-lg duration-200 rounded-lg"
            >
                <img
                    src={url}
                    alt="Attendance preview"
                    className="w-full h-auto rounded-md"
                />
            </DialogContent>

            <DialogOverlay className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm" />
        </Dialog>
    );
};

export default ReviewImage;
