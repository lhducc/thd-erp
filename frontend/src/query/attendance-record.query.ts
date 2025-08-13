import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {getAttendanceRecordByEmployeeId} from "@/apis/attendance-record.api.ts";
import api from "@/apis/api.ts";
import {toast} from "sonner";

export const useGetAttendanceRecordByEmployeeId = (id: string) =>
    useQuery({
        queryKey: ["attendanceRecord", id],
        queryFn: () => getAttendanceRecordByEmployeeId(id),
    });

interface UpdateAttendanceStatusParams {
    attendance_record_id: string;
    status: "approved" | "rejected";
    note_reject?: string;
}

export const useUpdateAttendanceStatus = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (data: UpdateAttendanceStatusParams) => {
            const response = await api.put(
                `/attendance-record/status/${data.attendance_record_id}`,
                {
                    status: data.status,
                    note_reject: data.note_reject || "",
                }
            );
            return response.data;
        },
        onSuccess: () => {
            // Invalidate and refetch the attendance records query
            queryClient.invalidateQueries({
                queryKey: ["attendance-records"]
            });
            toast.success("Cập nhật trạng thái chấm công thành công");
        },
        onError: (error) => {
            console.error("Error updating attendance status:", error);
            toast.error("Có lỗi xảy ra khi cập nhật trạng thái chấm công");
        },
    });
};