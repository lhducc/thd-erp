import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { X, Save } from "lucide-react";
import { useForm } from "react-hook-form";
import type { TimesheetDetail } from "@/types/timesheet";
import { useEffect } from "react";
import Loading from "@/components/Loading.tsx";

interface AdjustWorkDayForm {
    adjusted_work_day: number;
}

interface EditDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    detail: TimesheetDetail | null;
    onSubmit: (formData: AdjustWorkDayForm) => void;
    isLoading: boolean;
}

const EditDialog = ({ open, onOpenChange, detail, onSubmit, isLoading }: EditDialogProps) => {
    const { register, handleSubmit, reset, formState: { errors } } = useForm<AdjustWorkDayForm>({
        defaultValues: { adjusted_work_day: 0 }
    });

    // Reset form khi detail thay đổi
    useEffect(() => {
        if (detail) {
            reset({ adjusted_work_day: detail.work_days ?? 0 });
        } else {
            reset({ adjusted_work_day: 0 });
        }
    }, [detail, reset]);

    if (!detail) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="md:max-w-2xl">
                <DialogHeader>
                    <DialogTitle>Sửa công ngày {new Date(detail.date).toLocaleDateString("vi-VN")}</DialogTitle>
                    <DialogDescription>
                        Điều chỉnh số công cho nhân viên {detail.timesheet?.employee.full_name}
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="adjusted_work_day">Số công điều chỉnh</Label>
                        <Input
                            id="adjusted_work_day"
                            type="number"
                            step="0.5"
                            min="0"
                            max="2"
                            {...register("adjusted_work_day", {
                                required: "Vui lòng nhập số công",
                                min: { value: 0, message: "Số công không thể âm" },
                                max: { value: 2, message: "Số công tối đa là 2" },
                                valueAsNumber: true // Quan trọng: chuyển string thành number
                            })}
                        />
                        {errors.adjusted_work_day && (
                            <p className="text-sm text-red-500">{errors.adjusted_work_day.message}</p>
                        )}
                    </div>
                    <div className="flex justify-end space-x-2">
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                            <X className="h-4 w-4 mr-1" /> Hủy
                        </Button>
                        <Button type="submit">
                            <Save className="h-4 w-4 mr-1" /> {isLoading ? <Loading /> : "Lưu thay đổi"}
                        </Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
};

export default EditDialog;