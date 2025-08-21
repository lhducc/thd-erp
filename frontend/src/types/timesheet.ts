import type {Office} from "@/types/index.ts";
import type {WorkShift} from "@/types/work-schedule.ts";
import type {AttendanceRecord} from "@/types/attendance.ts";

export interface TimesheetResponse {
    message: string;
    data: {
        data: Timesheet[];
        total: number;
        total_pages: number;
        page: number;
        limit: number;
    };
    statuscode: number;
}

export type Timesheet = {
    timesheet_list_id: string;
    time_sheet_list_name: string;
    office_id: string;
    month: number;
    year: number;
    start_date: string; // ISO datetime, e.g. "2025-05-01T00:00:00Z"
    end_date: string;
    is_locked: boolean;
    created_by: string;
    created_at: string;
    updated_by: string | null;
    updated_at: string;
    locked_by: string | null;
    locked_at: string | null;
    office: Office;
};

export type CreateTimesheet = {
    name: string;
    office_id: string;
    month: number;
    year: number;
}

export type TimesheetList = {
    timesheet_list_id: string;
    time_sheet_list_name: string;
    office_id: string;
    month: number;
    year: number;
    start_date: string;
    end_date: string;
    is_locked: boolean;
    created_by: string;
    created_at: string;
    updated_by: string | null;
    updated_at: string;
    locked_by: string | null;
    locked_at: string | null;
    office: Office;
    timesheets: TimesheetInfor[];
};

export type TimesheetInfor = {
    timesheet_id: number;
    timesheet_list_id: string;
    employee_id: string;
    month: number;
    year: number;
    office_id: string;
    department_id: string;
    total_work_days: number;
    late_shifts: number;
    total_late_minutes: number;
    annual_leave_days: number;
    personal_leave_days: number;
    business_trip_days: number;
    remote_work_days: number;
    unpaid_leave_days: number;
    other_leave_days: number;
    annual_leave_hours: number;
    remote_work_hours: number;
    unpaid_leave_hours: number;
    created_by: string;
    created_at: string;
    updated_by: string | null;
    updated_at: string;
    employee: Employee;
    department: Department;
    details?: TimesheetDetail[];
};

export type Employee = {
    employee_id: string;
    full_name: string;
    hierarchy_level: HierarchyLevel;
    position: Position;
    department: Department;
};

export type HierarchyLevel = {
    id: string;
    hierarchy_level: string;
    hierarchy_number: number;
    created_date: string;
};

export type Position = {
    position_id: string;
    position_name: string;
    created_date: string;
};

export type Department = {
    department_id: string;
    department_name: string;
    manager: string;
    created_date: string;
    office_id: string;
};


export type TimesheetDetail = {
    timesheet_detail_id: number;
    timesheet_id: number;
    date: string;
    day_of_week: number;
    work_shift_id: string;
    is_working_day: boolean;
    is_holiday: boolean;
    is_weekend: boolean;
    checkin_record_id: string;
    checkout_record_id: string;
    is_manually_adjusted: boolean;
    work_hours: number;
    work_days: number;
    is_additional_shift: boolean;
    is_late: boolean;
    late_minutes: number;
    leave_type: string | null;
    leave_hours: number;
    is_absent: boolean;
    absent_reason: string | null;
    adjustment_by: string | null;
    adjustment_at: string | null;
    work_shift: WorkShift;
    checkin_record: AttendanceRecord;
    checkout_record: AttendanceRecord;
};