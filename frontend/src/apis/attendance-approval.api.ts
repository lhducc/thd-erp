import type {AttendanceNeedApprovalType} from "@/types/attendance-need-approval-type.ts";

export const getAllAttendanceApprovalAPI = async (): Promise<AttendanceNeedApprovalType[]> => {
    return [
        {
            employee: {
                employee_id: "EMP001",
                full_name: "Nguyễn Văn A",
                birthday: "1990-05-15",
                gender: "male",
                work_type: "full-time",
                phone_number: "0987654321",
                email: "nguyenvana@company.com",
                address: "123 Đường ABC, Quận 1, TP.HCM",
                account_id: "ACC001",
                position_id: "POS001",
                position: {
                    position_id: "POS001",
                    position_name: "Nhân viên",
                    created_date: "2020-01-01"
                },
                job_title_id: "JT001",
                status: "active",
                manager_id: "MNG001",
                department_id: "DEP001",
                created_date: "2020-01-01",
                department: {
                    department_id: "DEP001",
                    department_name: "Phòng Kinh doanh",
                    manager: "Trần Thị B",
                    created_date: "2020-01-01",
                    office_id: "OFF001",
                    office: {
                        office_id: "OFF001",
                        office_name: "Trụ sở chính",
                        phone_number: "02838223344",
                        address: "123 Đường XYZ, Quận 1, TP.HCM",
                        latitude: 10.7722,
                        longitude: 106.6983,
                        created_date: "2020-01-01"
                    }
                },
                Decisions: []
            },
            gps: "10.7722, 106.6983",
            checkin_image_url: "https://plus.unsplash.com/premium_photo-1676977395506-2320c4d80618?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxmZWF0dXJlZC1waG90b3MtZmVlZHwxfHx8ZW58MHx8fHx8",
            checkout_image_url: "https://plus.unsplash.com/premium_photo-1676977395506-2320c4d80618?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxmZWF0dXJlZC1waG90b3MtZmVlZHwxfHx8ZW58MHx8fHx8",
            work_type: "office",
            status: "Chờ duyệt"
        },
        {
            employee: {
                employee_id: "EMP002",
                full_name: "Trần Thị B",
                birthday: "1985-08-20",
                gender: "female",
                work_type: "full-time",
                phone_number: "0912345678",
                email: "tranthib@company.com",
                address: "456 Đường XYZ, Quận 3, TP.HCM",
                account_id: "ACC002",
                position_id: "POS002",
                position: {
                    position_id: "POS002",
                    position_name: "Trưởng phòng",
                    created_date: "2019-01-01"
                },
                job_title_id: "JT002",
                status: "active",
                manager_id: "MNG002",
                department_id: "DEP001",
                created_date: "2019-01-01",
                department: {
                    department_id: "DEP001",
                    department_name: "Phòng Kinh doanh",
                    manager: "Trần Thị B",
                    created_date: "2020-01-01",
                    office_id: "OFF001",
                    office: {
                        office_id: "OFF001",
                        office_name: "Trụ sở chính",
                        phone_number: "02838223344",
                        address: "123 Đường XYZ, Quận 1, TP.HCM",
                        latitude: 10.7722,
                        longitude: 106.6983,
                        created_date: "2020-01-01"
                    }
                },
                Decisions: []
            },
            gps: "10.7722, 106.6983",
            checkin_image_url: "https://plus.unsplash.com/premium_photo-1676977395506-2320c4d80618?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxmZWF0dXJlZC1waG90b3MtZmVlZHwxfHx8ZW58MHx8fHx8",
            checkout_image_url: "https://plus.unsplash.com/premium_photo-1676977395506-2320c4d80618?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxmZWF0dXJlZC1waG90b3MtZmVlZHwxfHx8ZW58MHx8fHx8",
            work_type: "remote",
            status: "Đã duyệt"
        },
        {
            employee: {
                employee_id: "EMP003",
                full_name: "Lê Văn C",
                birthday: "1992-03-10",
                gender: "male",
                work_type: "part-time",
                phone_number: "0967890123",
                email: "levanc@company.com",
                address: "789 Đường DEF, Quận 5, TP.HCM",
                account_id: "ACC003",
                position_id: "POS003",
                position: {
                    position_id: "POS003",
                    position_name: "Nhân viên part-time",
                    created_date: "2021-01-01"
                },
                job_title_id: "JT003",
                status: "active",
                manager_id: "MNG001",
                department_id: "DEP002",
                created_date: "2021-01-01",
                department: {
                    department_id: "DEP002",
                    department_name: "Phòng Nhân sự",
                    manager: "Phạm Thị D",
                    created_date: "2020-01-01",
                    office_id: "OFF002",
                    office: {
                        office_id: "OFF002",
                        office_name: "Chi nhánh Q2",
                        phone_number: "02838223355",
                        address: "456 Đường QWER, Quận 2, TP.HCM",
                        latitude: 10.7879,
                        longitude: 106.7000,
                        created_date: "2020-01-01"
                    }
                },
                Decisions: []
            },
            gps: "10.7879, 106.7000",
            checkin_image_url: "https://plus.unsplash.com/premium_photo-1676977395506-2320c4d80618?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxmZWF0dXJlZC1waG90b3MtZmVlZHwxfHx8ZW58MHx8fHx8",
            checkout_image_url: "https://plus.unsplash.com/premium_photo-1676977395506-2320c4d80618?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxmZWF0dXJlZC1waG90b3MtZmVlZHwxfHx8ZW58MHx8fHx8",
            work_type: "hybrid",
            status: "Không duyệt"
        }
    ];
}

export const createAttendanceApprovalApi = async () => {
    // Mock implementation
    return { success: true };
}

export const updateAttendanceApprovalApi = async () => {
    // Mock implementation
    return { success: true };
}