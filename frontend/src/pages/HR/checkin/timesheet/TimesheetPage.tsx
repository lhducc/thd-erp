import {useGetAllTimesheets} from "@/query/timesheet.query.ts";
import type {ColumnDef} from "@tanstack/react-table";
import {EyeIcon} from "lucide-react";
import {formatDate} from "@/lib/utils.ts";
import DataTable from "@/components/DataTable.tsx";
import {Link} from "react-router-dom";
import type {CreateTimesheet, Timesheet} from "@/types/timesheet.ts";
import {useState} from "react";
import {DialogComponent} from "@/components/DialogComponent.tsx";
import {type SubmitHandler, useForm} from "react-hook-form";
import {useMutation} from "@tanstack/react-query";
import {createTimesheetApi} from "@/apis/timesheet.api.ts";
import {toast} from "sonner";
import {Button} from "@/components/ui/button.tsx";
import {useOffice} from "@/query/useOffice.ts";
import Loading from "@/components/Loading.tsx";
import axios from "axios";

type Inputs = {
    name: string;
    month_year: string; // e.g. "2025-05"
}

const TimesheetForm = ({refresh}: { refresh: () => void }) => {
    const [open, setOpen] = useState<boolean>(false);
    const {
        register,
        handleSubmit,
    } = useForm<Inputs>()
    const [monthYearLabel, setMonthYearLabel] = useState<string>("");

    const onSubmit: SubmitHandler<Inputs> = async (data) => {
        const [year, month] = data.month_year.split("-").map(Number);
        await createTimesheet({
            name: data.name,
            month,
            year,
        });
    }

    const {mutateAsync: createTimesheet, isPending} = useMutation({
        mutationFn: (payload: CreateTimesheet) => createTimesheetApi(payload),
        onSuccess: () => {
            toast.success("Tạo bảng công mới thành công")
            setOpen(false)
            refresh()
        },
        onError: (error) => {
            if (axios.isAxiosError(error)) {
                toast.error(error.response.data.message);
            } else {
                toast.error(error.message);
            }
        }
    })

    const openButton = (
        <Button>
            Tạo bảng công
        </Button>
    )

    const content = () => {
        return (
            <form className={`flex flex-col space-y-2`} onSubmit={handleSubmit(onSubmit)}>
                <label className="text-sm font-semibold text-gray-900">
                    Tên bảng công
                </label>
                <input className="border rounded-lg p-2" defaultValue="" {...register("name", {required: true})} />
                <label className="text-sm font-semibold text-gray-900">
                    Thời gian áp dụng
                </label>
                <input
                    className="border rounded-lg p-2"
                    {...register("month_year")}
                    type="month"
                />
                <p className="text-sm text-gray-600 italic">{monthYearLabel}</p>
                <div className={`flex gap-2 items-center w-full justify-center mt-2`}>
                    <Button variant="secondary" onClick={() => setOpen(false)}>Hủy</Button>
                    <Button type="submit">{isPending ? <Loading/> : "Tạo bảng công"}</Button>
                </div>
            </form>
        )
    }

    return (
        <DialogComponent
            open={open}
            setOpen={setOpen}
            title="Tạo bảng công"
            openButton={openButton}
            content={content}
        />
    )
}

const TimesheetPage = () => {
    const {data: timesheet, isLoading: isTimesheetLoading, refetch: refreshTimesheet} = useGetAllTimesheets({
        page: 1,
        limit: 99999,
        search: ''
    });

    const columns: ColumnDef<Timesheet>[] = [
        {
            accessorKey: "timesheet_list_id",
            header: "Mã bảng công",
        },
        {
            accessorKey: "time_sheet_list_name",
            header: "Tên bảng công",
        },
        {
            accessorKey: "address",
            header: "Thời gian",
            cell: ({row}) => {
                return (
                    <>
                        {formatDate(row.original.start_date)} - {formatDate(row.original.end_date)}
                    </>
                )
            }
        },
        {
            accessorKey: "is_locked",
            header: "Trạng thái ",
            cell: ({row}) => {
                return (
                    <>{!row.original.is_locked ? "Đang áp dụng" : <p className="text-red-500">Đã chốt công</p>}</>
                )
            }
        },
        {
            id: "actions",
            header: "Xem chi tiết",
            cell: ({row}) => {
                return (
                    <Link to={`/timesheet/${row.original.timesheet_list_id}`} className="flex gap-4">
                        <EyeIcon/>
                    </Link>
                );
            },
        },
    ];
    return (
        <DataTable
            columns={columns}
            data={timesheet?.data.data || []}
            isLoading={isTimesheetLoading}
            title="Bảng công"
            buttonCreate={<TimesheetForm refresh={refreshTimesheet}/>}
            keyFilter="time_sheet_list_name"
        />
    );
};

export default TimesheetPage;
