import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";

interface ImageDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    imageUrl: string;
    title?: string;
}

const ImageDialog = ({ open, onOpenChange, imageUrl, title = "Hình ảnh chấm công" }: ImageDialogProps) => {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-4xl">
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                </DialogHeader>
                <div className="flex justify-center">
                    <img
                        src={imageUrl}
                        alt="Hình ảnh chấm công"
                        className="max-h-[70vh] max-w-full object-contain"
                    />
                </div>
            </DialogContent>
        </Dialog>
    );
};

export default ImageDialog;