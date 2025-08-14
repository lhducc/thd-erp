import { type JSX } from "react";
import {
    Dialog,
    DialogContent,
    DialogOverlay,
    DialogTrigger
} from "@/components/ui/dialog.tsx";
import {DialogTitle} from "@radix-ui/react-dialog";

export const DialogComponent = ({
                                    open,
                                    setOpen,
                                    openButton,
                                    title,
                                    content
                                }: {
    open: boolean;
    setOpen: (value: boolean) => void;
    openButton: JSX.Element;
    title: string;
    content: () => JSX.Element;
}) => {
    return (
        <>
            <Dialog open={open} onOpenChange={setOpen}>

                <DialogTrigger asChild>{openButton}</DialogTrigger>
                {open && (
                    <DialogContent
                        className="fixed left-1/2 top-1/2 z-50 grid w-full max-w-2xl translate-x-[-50%] translate-y-[-50%] gap-4 border border-gray-200 bg-white p-6 shadow-lg duration-200 rounded-lg"
                    >
                        <DialogTitle className="text-xl font-bold">
                            {title}
                        </DialogTitle>
                        {content()}
                    </DialogContent>
                )}

                <DialogOverlay className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm" />
            </Dialog>
        </>
    );
};
