-- -- ERP Database Schema Creation Script
-- -- This script creates all tables from scratch without using GORM AutoMigrate

-- -- Create extension for UUID generation if not exists
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -- Drop tables if they exist (in reverse dependency order)
-- DROP TABLE IF EXISTS timesheet_detail CASCADE;
-- DROP TABLE IF EXISTS timesheets CASCADE;
-- DROP TABLE IF EXISTS timesheet_list CASCADE;
-- DROP TABLE IF EXISTS work_schedule_manager CASCADE;
-- DROP TABLE IF EXISTS work_schedule_shift CASCADE;
-- DROP TABLE IF EXISTS work_schedule CASCADE;
-- DROP TABLE IF EXISTS attendance_record CASCADE;
-- DROP TABLE IF EXISTS attendance_category CASCADE;
-- DROP TABLE IF EXISTS contract_allowance CASCADE;
-- DROP TABLE IF EXISTS contract CASCADE;
-- DROP TABLE IF EXISTS allowance CASCADE;
-- DROP TABLE IF EXISTS insurance CASCADE;
-- DROP TABLE IF EXISTS decision_employees CASCADE;
-- DROP TABLE IF EXISTS decision CASCADE;
-- DROP TABLE IF EXISTS decision_type CASCADE;
-- DROP TABLE IF EXISTS employee_document CASCADE;
-- DROP TABLE IF EXISTS employee_document_type CASCADE;
-- DROP TABLE IF EXISTS employee_workshift CASCADE;
-- DROP TABLE IF EXISTS work_shifts CASCADE;
-- DROP TABLE IF EXISTS account CASCADE;
-- DROP TABLE IF EXISTS employee CASCADE;
-- DROP TABLE IF EXISTS jobtitle CASCADE;
-- DROP TABLE IF EXISTS hierarchy_level CASCADE;
-- DROP TABLE IF EXISTS contract_type CASCADE;
-- DROP TABLE IF EXISTS department CASCADE;
-- DROP TABLE IF EXISTS position CASCADE;
-- DROP TABLE IF EXISTS office CASCADE;
-- DROP TABLE IF EXISTS role CASCADE;

-- -- Create ENUM types
-- DO $$
-- BEGIN
--     -- enum contract_group_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'contract_group_enum') THEN
--         CREATE TYPE contract_group_enum AS ENUM (
--             'Hợp đồng xác định thời hạn',
--             'Hợp đồng không xác định thời hạn',
--             'Hợp đồng thử việc',
--             'Hợp đồng đào tạo nghề',
--             'Hợp đồng dịch vụ'
--         );
--     END IF;

--     -- enum unit_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'unit_enum') THEN
--         CREATE TYPE unit_enum AS ENUM ('Năm', 'Tháng', 'Tuần', 'Ngày');
--     END IF;

--     -- enum working_type_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'working_type_enum') THEN
--         CREATE TYPE working_type_enum AS ENUM (
--             'Toàn thời gian',
--             'Bán thời gian',
--             'Cộng tác viên',
--             'Chuyên gia',
--             'Theo ca',
--             'Khoán sản phẩm',
--             'Khoán công việc'
--         );
--     END IF;

--     -- enum decision_status_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'decision_status_enum') THEN
--         CREATE TYPE decision_status_enum AS ENUM (
--             'Đã duyệt',
--             'Không duyệt',
--             'Chờ duyệt'
--         );
--     END IF;

--     -- enum decision_condition_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'decision_condition_enum') THEN
--         CREATE TYPE decision_condition_enum AS ENUM (
--             'Chưa hiệu lực',
--             'Đang hiệu lực'
--         );
--     END IF;

--     -- enum decision_group_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'decision_group_enum') THEN
--         CREATE TYPE decision_group_enum AS ENUM (
--             'Hình thức khen thưởng',
--             'Hình thức kỷ luật',
--             'Lý do điều chuyển',
--             'Lý do tiếp nhận',
--             'Lý do bổ nhiệm',
--             'Lý do miễn nhiệm',
--             'Lý do chấm dứt HĐLĐ'
--         );
--     END IF;

--     -- enum document_condition_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_condition_enum') THEN
--         CREATE TYPE document_condition_enum AS ENUM (
--             'Hết hạn',
--             'Đang hiệu lực'
--         );
--     END IF;

--     -- enum document_status_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_status_enum') THEN
--         CREATE TYPE document_status_enum AS ENUM (
--             'Duyệt',
--             'Không duyệt',
--             'Chờ duyệt'
--         );
--     END IF;

--     -- enum approve_status_enum
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'approve_status_enum') THEN
--         CREATE TYPE approve_status_enum AS ENUM (
--             'Đã duyệt',
--             'Không duyệt',
--             'Chờ duyệt'
--         );
--     END IF;

--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'condition_enum') THEN
--         CREATE TYPE condition_enum AS ENUM (
--             'Chưa hiệu lực',
--             'Đang hiệu lực',
--             'Hết hiệu lực',
--             'Thanh lý'
--         );
--     END IF;

--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_group_enum') THEN
--         CREATE TYPE document_group_enum AS ENUM (
--             'Loại chứng chỉ',
--             'Loại lao động',
--             'Thủ tục tiếp nhận',
--             'Thủ tục thôi việc'
--         );
--     END IF;

--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'work_day_enum') THEN
--         CREATE TYPE work_day_enum AS ENUM (
--             '1',
--             '0.5',
--             '0'
--         );
--     END IF;

--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'repeat_type_enum') THEN
--         CREATE TYPE repeat_type_enum AS ENUM (
--             'weekly',
--             'monthly'
--         );
--     END IF;

--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_work_schedule_enum') THEN
--         CREATE TYPE status_work_schedule_enum AS ENUM (
--             'expired',
--             'inactive',
--             'active'
--         );
--     END IF;

--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'weekday_enum') THEN
--         CREATE TYPE weekday_enum AS ENUM (
--             'sunday',
--             'monday',
--             'tuesday',
--             'wednesday',
--             'thursday',
--             'friday',
--             'saturday'
--         );
--     END IF;

--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'work_type_enum') THEN
--         CREATE TYPE work_type_enum AS ENUM (
--             'ca hành chính',
--             'ca kíp'
--         );
--     END IF;
-- END
-- $$;

-- -- Create tables in proper dependency order

-- 1. Role table (no dependencies)
CREATE TABLE IF NOT EXISTS role (
    id VARCHAR(50) PRIMARY KEY,
    role_name VARCHAR(100) NOT NULL,
    created_date TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Office table (no dependencies)
CREATE TABLE office (
    office_id VARCHAR(8) PRIMARY KEY,
    office_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    address TEXT,
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Position table (no dependencies)
CREATE TABLE IF NOT EXISTS position (
    position_id VARCHAR(50) PRIMARY KEY,
    position_name VARCHAR(100) NOT NULL,
    created_date TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Hierarchy level table (no dependencies)
CREATE TABLE IF NOT EXISTS hierarchy_level (
    id VARCHAR(50) PRIMARY KEY,
    level_name VARCHAR(100) NOT NULL,
    level_order INTEGER,
    created_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employeedocumenttype
(
    id VARCHAR(50) PRIMARY KEY,
    type_name VARCHAR(255) NOT NULL,
    group_type document_group_enum,
    created_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contracttype
(
    id VARCHAR(50) PRIMARY KEY,
    type_name VARCHAR(255) NOT NULL,
    group_type contract_group_enum,
    created_date TIMESTAMPTZ DEFAULT NOW()
);

-- -- 4. Department table (depends on office)
-- CREATE TABLE department (
--     department_id VARCHAR(10) PRIMARY KEY,
--     department_name VARCHAR(100) NOT NULL,
--     description TEXT,
--     manager VARCHAR(100),
--     created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     office_id VARCHAR(8),
--     FOREIGN KEY (office_id) REFERENCES office(office_id)
-- );

-- -- 6. Job title table (depends on hierarchy_level)
-- CREATE TABLE IF NOT EXISTS jobtitle (
--     job_title_id VARCHAR(50) PRIMARY KEY,
--     job_title VARCHAR(255) NOT NULL,
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     hierarchy_level_id VARCHAR(50),
--     FOREIGN KEY (hierarchy_level_id) REFERENCES hierarchy_level(id)
-- );

-- -- 7. Employee table (depends on position, jobtitle, department)
-- CREATE TABLE IF NOT EXISTS employee (
--     employee_id VARCHAR(50) PRIMARY KEY,
--     full_name VARCHAR(255) NOT NULL,
--     birthday DATE,
--     gender VARCHAR(10) CHECK (gender IN ('Nam', 'Nữ', 'Khác')),
--     work_type work_type_enum,
--     phone_number VARCHAR(20),
--     email VARCHAR(255),
--     address TEXT,
--     account_id BIGINT,
--     position_id VARCHAR(50),
--     job_title_id VARCHAR(50),
--     status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
--     manager VARCHAR(50),
--     department_id VARCHAR(10),
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     schedule_id INTEGER,
--     FOREIGN KEY (position_id) REFERENCES position(position_id),
--     FOREIGN KEY (job_title_id) REFERENCES jobtitle(job_title_id),
--     FOREIGN KEY (department_id) REFERENCES department(department_id),
--     FOREIGN KEY (manager) REFERENCES employee(employee_id)
-- );

-- -- 8. Account table (depends on role and employee)
-- CREATE TABLE IF NOT EXISTS account (
--     id BIGSERIAL PRIMARY KEY,
--     login_mail VARCHAR(255) UNIQUE NOT NULL,
--     password VARCHAR(255) NOT NULL,
--     first_login BOOLEAN DEFAULT true,
--     role_id VARCHAR(50),
--     employee_id VARCHAR(50),
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (role_id) REFERENCES role(id),
--     FOREIGN KEY (employee_id) REFERENCES employee(employee_id)
-- );

-- -- Update employee table to add foreign key to account
-- ALTER TABLE employee ADD CONSTRAINT fk_employee_account 
--     FOREIGN KEY (account_id) REFERENCES account(id);

-- -- 9. Employee document type table
-- CREATE TABLE IF NOT EXISTS employee_document_type (
--     id VARCHAR(50) PRIMARY KEY,
--     type_name VARCHAR(255) NOT NULL,
--     group_type document_group_enum,
--     created_date TIMESTAMPTZ DEFAULT NOW()
-- );

-- -- 10. Employee document table
-- CREATE TABLE IF NOT EXISTS employee_document (
--     id BIGSERIAL PRIMARY KEY,
--     employee_id VARCHAR(50) NOT NULL,
--     document_type_id VARCHAR(50),
--     document_name VARCHAR(255),
--     document_number VARCHAR(100),
--     issue_date DATE,
--     expiry_date DATE,
--     issuer VARCHAR(255),
--     file_path VARCHAR(500),
--     notes TEXT,
--     status document_status_enum DEFAULT 'Chờ duyệt',
--     condition document_condition_enum DEFAULT 'Đang hiệu lực',
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (employee_id) REFERENCES employee(employee_id),
--     FOREIGN KEY (document_type_id) REFERENCES employee_document_type(id)
-- );

-- -- 11. Contract type table
-- CREATE TABLE IF NOT EXISTS contract_type (
--     id VARCHAR(50) PRIMARY KEY,
--     type_name VARCHAR(255) NOT NULL,
--     group_type contract_group_enum,
--     created_date TIMESTAMPTZ DEFAULT NOW()
-- );

-- -- 12. Contract table
-- CREATE TABLE IF NOT EXISTS contract (
--     id BIGSERIAL PRIMARY KEY,
--     contract_number VARCHAR(100) UNIQUE,
--     employee_id VARCHAR(50) NOT NULL,
--     contract_type_id VARCHAR(50),
--     start_date DATE,
--     end_date DATE,
--     salary DECIMAL(15,2),
--     status VARCHAR(20) DEFAULT 'active',
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (employee_id) REFERENCES employee(employee_id),
--     FOREIGN KEY (contract_type_id) REFERENCES contract_type(id)
-- );

-- -- 13. Decision type table
-- CREATE TABLE IF NOT EXISTS decision_type (
--     id VARCHAR(50) PRIMARY KEY,
--     type_name VARCHAR(255) NOT NULL,
--     group_type decision_group_enum,
--     created_date TIMESTAMPTZ DEFAULT NOW()
-- );

-- -- 14. Decision table
-- CREATE TABLE IF NOT EXISTS decision (
--     id BIGSERIAL PRIMARY KEY,
--     decision_number VARCHAR(100) UNIQUE,
--     decision_type_id VARCHAR(50),
--     title VARCHAR(500),
--     content TEXT,
--     effective_date DATE,
--     status decision_status_enum DEFAULT 'Chờ duyệt',
--     condition decision_condition_enum DEFAULT 'Chưa hiệu lực',
--     created_by VARCHAR(50),
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (decision_type_id) REFERENCES decision_type(id),
--     FOREIGN KEY (created_by) REFERENCES employee(employee_id)
-- );

-- -- 15. Decision employees junction table
-- CREATE TABLE IF NOT EXISTS decision_employees (
--     decision_id BIGINT,
--     employee_id VARCHAR(50),
--     PRIMARY KEY (decision_id, employee_id),
--     FOREIGN KEY (decision_id) REFERENCES decision(id) ON DELETE CASCADE,
--     FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
-- );

-- -- 16. Allowance table
-- CREATE TABLE IF NOT EXISTS allowance (
--     id BIGSERIAL PRIMARY KEY,
--     allowance_name VARCHAR(255) NOT NULL,
--     amount DECIMAL(15,2),
--     unit unit_enum,
--     created_date TIMESTAMPTZ DEFAULT NOW()
-- );

-- -- 17. Contract allowance junction table
-- CREATE TABLE IF NOT EXISTS contract_allowance (
--     contract_id BIGINT,
--     allowance_id BIGINT,
--     amount DECIMAL(15,2),
--     PRIMARY KEY (contract_id, allowance_id),
--     FOREIGN KEY (contract_id) REFERENCES contract(id) ON DELETE CASCADE,
--     FOREIGN KEY (allowance_id) REFERENCES allowance(id) ON DELETE CASCADE
-- );

-- -- 18. Insurance table
-- CREATE TABLE IF NOT EXISTS insurance (
--     id BIGSERIAL PRIMARY KEY,
--     employee_id VARCHAR(50) NOT NULL,
--     insurance_number VARCHAR(100),
--     start_date DATE,
--     end_date DATE,
--     premium DECIMAL(15,2),
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (employee_id) REFERENCES employee(employee_id)
-- );

-- -- 19. Work shifts table
-- CREATE TABLE IF NOT EXISTS work_shifts (
--     id VARCHAR(20) PRIMARY KEY,
--     shift_name VARCHAR(100) NOT NULL,
--     start_time TIME,
--     end_time TIME,
--     break_duration INTEGER DEFAULT 0,
--     work_type work_type_enum,
--     created_date TIMESTAMPTZ DEFAULT NOW()
-- );

-- -- 20. Employee workshift table
-- CREATE TABLE IF NOT EXISTS employee_workshift (
--     id BIGSERIAL PRIMARY KEY,
--     employee_id VARCHAR(50) NOT NULL,
--     work_shift_id VARCHAR(20),
--     assigned_date DATE,
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (employee_id) REFERENCES employee(employee_id),
--     FOREIGN KEY (work_shift_id) REFERENCES work_shifts(id)
-- );

-- -- 21. Attendance category table
-- CREATE TABLE IF NOT EXISTS attendance_category (
--     id VARCHAR(50) PRIMARY KEY,
--     category_name VARCHAR(100) NOT NULL,
--     work_day_value work_day_enum DEFAULT '1',
--     created_date TIMESTAMPTZ DEFAULT NOW()
-- );

-- -- 22. Attendance record table
-- CREATE TABLE IF NOT EXISTS attendance_record (
--     id BIGSERIAL PRIMARY KEY,
--     employee_id VARCHAR(50) NOT NULL,
--     attendance_date DATE NOT NULL,
--     check_in_time TIMESTAMPTZ,
--     check_out_time TIMESTAMPTZ,
--     work_shift_id VARCHAR(20),
--     attendance_category_id VARCHAR(50),
--     late_minutes INTEGER DEFAULT 0,
--     early_leave_minutes INTEGER DEFAULT 0,
--     notes TEXT,
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (employee_id) REFERENCES employee(employee_id),
--     FOREIGN KEY (work_shift_id) REFERENCES work_shifts(id),
--     FOREIGN KEY (attendance_category_id) REFERENCES attendance_category(id)
-- );

-- -- 23. Work schedule table
-- CREATE TABLE IF NOT EXISTS work_schedule (
--     id SERIAL PRIMARY KEY,
--     schedule_name VARCHAR(255) NOT NULL,
--     start_date DATE,
--     end_date DATE,
--     repeat_type repeat_type_enum,
--     status status_work_schedule_enum DEFAULT 'inactive',
--     office_id VARCHAR(8),
--     created_by VARCHAR(50),
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (office_id) REFERENCES office(office_id),
--     FOREIGN KEY (created_by) REFERENCES employee(employee_id)
-- );

-- -- 24. Work schedule shift table
-- CREATE TABLE IF NOT EXISTS work_schedule_shift (
--     id BIGSERIAL PRIMARY KEY,
--     work_schedule_id INTEGER,
--     work_shift_id VARCHAR(20),
--     weekday weekday_enum,
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (work_schedule_id) REFERENCES work_schedule(id) ON DELETE CASCADE,
--     FOREIGN KEY (work_shift_id) REFERENCES work_shifts(id)
-- );

-- -- 25. Work schedule manager table
-- CREATE TABLE IF NOT EXISTS work_schedule_manager (
--     id BIGSERIAL PRIMARY KEY,
--     work_schedule_id INTEGER,
--     manager_id VARCHAR(50),
--     created_date TIMESTAMPTZ DEFAULT NOW(),
--     FOREIGN KEY (work_schedule_id) REFERENCES work_schedule(id) ON DELETE CASCADE,
--     FOREIGN KEY (manager_id) REFERENCES employee(employee_id)
-- );

-- -- Create indexes for better performance
-- CREATE INDEX IF NOT EXISTS idx_employee_department ON employee(department_id);
-- CREATE INDEX IF NOT EXISTS idx_employee_position ON employee(position_id);
-- CREATE INDEX IF NOT EXISTS idx_employee_job_title ON employee(job_title_id);
-- CREATE INDEX IF NOT EXISTS idx_employee_manager ON employee(manager);
-- CREATE INDEX IF NOT EXISTS idx_account_login_mail ON account(login_mail);
-- CREATE INDEX IF NOT EXISTS idx_attendance_employee_date ON attendance_record(employee_id, attendance_date);
-- CREATE INDEX IF NOT EXISTS idx_contract_employee ON contract(employee_id);

-- -- Insert default roles
-- INSERT INTO role (id, role_name, created_date) VALUES 
--     ('admin', 'admin', NOW()),
--     ('manager', 'manager', NOW()),
--     ('employee', 'employee', NOW())
-- ON CONFLICT (id) DO NOTHING;

-- SELECT 'Database schema created successfully!' AS status;
-- SELECT 'Tables created: ' || COUNT(*) AS tables_count 
-- FROM information_schema.tables 
-- WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
