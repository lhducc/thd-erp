import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Loader2 } from "lucide-react";
import { exportAttendanceExcelApi } from "@/apis/attendance-record.api";

type ExportOptions = {
  byDate: Date | null;
  byMonth: Date | null;
};

type Props = {
  selectedDate: string | null;  
  onConfirm?: (options: ExportOptions) => void;
};

export default function ExportFileDialog({ selectedDate, onConfirm }: Props) {
  const [open, setOpen] = useState(false);
  const [byDateChecked, setByDateChecked] = useState(false);
  const [byMonthChecked, setByMonthChecked] = useState(false);
  const [month, setMonth] = useState<Date | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleConfirm = async () => {
    if (!byDateChecked && !byMonthChecked) {
      setError("Vui lòng chọn theo ngày hoặc theo tháng");
      return;
    }
    if (byDateChecked && !selectedDate) {
      setError("Ngày đã chọn không hợp lệ");
      return;
    }
    if (byMonthChecked && !month) {
      setError("Vui lòng chọn tháng hợp lệ");
      return;
    }

    setError("");
    setLoading(true);

    try {
      const payload: any = {};
      if (byDateChecked && selectedDate) {
        payload.date = selectedDate;
      }
      if (byMonthChecked && month) {
        const [y, m] = [month.getFullYear(), month.getMonth() + 1];
        payload.byMonth = `${y}-${String(m).padStart(2, "0")}`; // yyyy-MM
      }

      const blob = await exportAttendanceExcelApi(payload);

      const url = window.URL.createObjectURL(new Blob([blob]));
      const a = document.createElement("a");
      a.href = url;
      a.download = "attendance.xlsx";
      document.body.appendChild(a);
      a.click();
      a.remove();
    } catch (err) {
      setError("Xuất file thất bại");
    }

    setLoading(false);
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="destructive">Xuất file</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Xuất file</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {/* Theo ngày */}
          <div className="flex items-center space-x-2">
            <Checkbox
              id="byDate"
              checked={byDateChecked}
              onCheckedChange={(v) => setByDateChecked(Boolean(v))}
            />
            <Label htmlFor="byDate"> Theo ngày (ngày đã chọn:{" "}{selectedDate? new Date(selectedDate).toLocaleDateString("vi-VN")  : "chưa có"} )</Label>
          </div>

          {/* Theo tháng */}
          <div className="flex items-center space-x-2">
            <Checkbox
              id="byMonth"
              checked={byMonthChecked}
              onCheckedChange={(v) => setByMonthChecked(Boolean(v))}
            />
            <Label htmlFor="byMonth">Theo tháng</Label>
          </div>
          {byMonthChecked && (
            <Input
              type="month"
              onChange={(e) => setMonth(new Date(e.target.value))}
            />
          )}

          {/* Hiện lỗi nếu có */}
          {error && <p className="text-red-500 text-sm">{error}</p>}
        </div>

        {/* Footer */}
        <div className="flex justify-end gap-2 pt-4">
          <Button variant="outline" onClick={() => setOpen(false)}>
            Hủy
          </Button>
          <Button
            onClick={handleConfirm}
            disabled={loading}
            variant={error ? "destructive" : "default"}
          >
            {loading ? <Loader2 className="animate-spin mr-2 h-4 w-4" /> : null}
            Xác nhận
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
