import { useMemo, useState, useEffect } from "react";
import { Button } from "@/components/ui/button.tsx";
import { DialogComponent } from "@/components/DialogComponent.tsx";
import { Label } from "@radix-ui/react-label";
import { Input } from "@/components/ui/input.tsx";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { createAttendanceRecordByAdminId } from "@/apis/attendance-record.api.ts";
import type { CreateManualRecord } from "@/types/attendance.ts";
import {toRFC3339} from "@/lib/utils.ts";

type Props = {
    employeeId: string;
    refresh: () => void;
};

export const CreateAttendant = ({ employeeId, refresh }: Props) => {
    const [timestamp, setTimestamp] = useState(() => toRFC3339(new Date()));
    const [open, setOpen] = useState(false);

    const [form, setForm] = useState<CreateManualRecord>({
        employee_id: employeeId,
        timestamp: toRFC3339(new Date()),
    });

    useEffect(() => {
        setForm(prev => ({
            ...prev,
            timestamp: timestamp
        }));
    }, [timestamp]);

    const { mutateAsync: createAttendance, isPending  } = useMutation({
        mutationFn: (data: CreateManualRecord) => createAttendanceRecordByAdminId(data),
        onSuccess: () => {
            toast.success("Thêm chấm công thành công");
            refresh();
            setOpen(false);
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        await createAttendance(form);
    };

    const openButton = useMemo(
        () => <Button>Tạo chấm công</Button>,
        []
    );

    const content = () => {
        return (
            <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
                <div>
                    <Label htmlFor="employeeId">Mã nhân viên</Label>
                    <Input
                        id="employeeId"
                        value={employeeId}
                        readOnly
                        className="bg-gray-100 cursor-not-allowed"
                    />
                </div>

                <div>
                    <Label htmlFor="timestamp">Thời gian chấm công</Label>
                    <Input
                        id="timestamp"
                        type="datetime-local"
                        value={timestamp.slice(0, 16)}
                        onChange={(e) => {
                            const newTimestamp = toRFC3339(new Date(e.target.value));
                            setTimestamp(newTimestamp);
                        }}
                        required
                    />
                </div>

                <Button type="submit" disabled={isPending}>
                    {isPending ? "Đang xử lý..." : "Xác nhận"}
                </Button>
            </form>
        )
    }

    return (
        <DialogComponent
            open={open}
            setOpen={setOpen}
            openButton={openButton}
            title="Tạo chấm công"
            content={content}
        />
    );
};
