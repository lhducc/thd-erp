-- Simple script to create test data for existing tables only
-- This script works with the current database schema

-- Create extension for random data generation if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Insert roles (match schema table name: role)
INSERT INTO role (id, role_name, created_date)
VALUES
    ('admin', 'Admin', NOW()),
    ('manager', 'Manager', NOW()),
    ('employee', 'Employee', NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert hierarchy levels first (required for job titles)
INSERT INTO hierarchy_level (id, level_name, level_order, created_date)
VALUES
    ('CB0001', 'Cán bộ', 1, NOW()),
    ('CB0002', 'Chuyên viên', 2, NOW()),
    ('CB0003', 'Trưởng phòng', 3, NOW()),
    ('CB0004', 'Giám đốc', 4, NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample offices (match schema table name: office)
INSERT INTO office (office_id, office_name, phone_number, address, latitude, longitude, created_date)
VALUES 
    ('HN001', 'Văn phòng Hà Nội', '0243456789', '123 Đường Láng, Đống Đa, Hà Nội', 21.0285, 105.8542, NOW()),
    ('HCM001', 'Văn phòng TP.HCM', '0281234567', '456 Nguyễn Văn Cừ, Quận 5, TP.HCM', 10.7769, 106.7009, NOW()),
    ('DN001', 'Văn phòng Đà Nẵng', '0236512345', '789 Trần Phú, Hải Châu, Đà Nẵng', 16.0471, 108.2068, NOW())
ON CONFLICT (office_id) DO NOTHING;

-- Insert departments (depends on office)
INSERT INTO department (department_id, department_name, manager, created_date, office_id)
VALUES
    ('HR001', 'Phòng Nhân sự', '', NOW(), 'HN001'),
    ('IT001', 'Phòng Công nghệ thông tin', '', NOW(), 'HN001'),
    ('FIN001', 'Phòng Tài chính', '', NOW(), 'HN001'),
    ('MKT001', 'Phòng Marketing', '', NOW(), 'HCM001'),
    ('SALE001', 'Phòng Kinh doanh', '', NOW(), 'HCM001'),
    ('OP001', 'Phòng Vận hành', '', NOW(), 'DN001')
ON CONFLICT (department_id) DO NOTHING;

-- Insert job titles (depends on hierarchy_level, correct column name: job_title)
INSERT INTO jobtitle (job_title_id, job_title, created_date, hierarchy_level_id)
VALUES
    ('JT001', 'Lập trình viên', NOW(), 'CB0002'),
    ('JT002', 'Nhân viên nhân sự', NOW(), 'CB0002'),
    ('JT003', 'Kế toán', NOW(), 'CB0002'),
    ('JT004', 'Marketing', NOW(), 'CB0002'),
    ('JT005', 'Kinh doanh', NOW(), 'CB0002'),
    ('JT006', 'Trưởng phòng IT', NOW(), 'CB0003'),
    ('JT007', 'Trưởng phòng HR', NOW(), 'CB0003'),
    ('JT008', 'Giám đốc', NOW(), 'CB0004')
ON CONFLICT (job_title_id) DO NOTHING;

-- Insert sample positions (match schema table name: position)
INSERT INTO position (position_id, position_name, created_date)
VALUES 
    ('POS001', 'Nhân viên', NOW()),
    ('POS002', 'Nhân viên cao cấp', NOW()),
    ('POS003', 'Trưởng nhóm', NOW()),
    ('POS004', 'Phó trưởng phòng', NOW()),
    ('POS005', 'Trưởng phòng', NOW()),
    ('POS006', 'Phó giám đốc', NOW()),
    ('POS007', 'Giám đốc', NOW()),
    ('POS008', 'Chuyên viên', NOW())
ON CONFLICT (position_id) DO NOTHING;

-- Insert work shifts for attendance system
INSERT INTO work_shifts (id, shift_name, start_time, end_time, break_duration, work_type, created_date)
VALUES 
    ('SHIFT001', 'Ca hành chính', '08:00:00', '17:00:00', 60, 'ca hành chính', NOW()),
    ('SHIFT002', 'Ca sáng', '06:00:00', '14:00:00', 60, 'ca kíp', NOW()),
    ('SHIFT003', 'Ca chiều', '14:00:00', '22:00:00', 60, 'ca kíp', NOW()),
    ('SHIFT004', 'Ca đêm', '22:00:00', '06:00:00', 60, 'ca kíp', NOW()),
    ('SHIFT005', 'Ca linh hoạt', '09:00:00', '18:00:00', 60, 'ca hành chính', NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert attendance categories
INSERT INTO attendance_category (id, category_name, work_day_value, created_date)
VALUES 
    ('ATT001', 'Đi làm đúng giờ', '1', NOW()),
    ('ATT002', 'Đi muộn', '1', NOW()),
    ('ATT003', 'Về sớm', '1', NOW()),
    ('ATT004', 'Nghỉ có phép', '0', NOW()),
    ('ATT005', 'Nghỉ không phép', '0', NOW()),
    ('ATT006', 'Nghỉ ốm', '0', NOW()),
    ('ATT007', 'Nghỉ lễ', '0', NOW()),
    ('ATT008', 'Công tác', '1', NOW()),
    ('ATT009', 'Làm thêm giờ', '1.5', NOW()),
    ('ATT010', 'Làm việc cuối tuần', '2', NOW())
ON CONFLICT (id) DO NOTHING;

-- Display what we've created
SELECT 
    'Simple Data Generation Complete!' as status,
    (SELECT COUNT(*) FROM office) as offices_created,
    (SELECT COUNT(*) FROM department) as departments_created,
    (SELECT COUNT(*) FROM position) as positions_created,
    (SELECT COUNT(*) FROM jobtitle) as job_titles_created,
    (SELECT COUNT(*) FROM hierarchy_level) as hierarchy_levels_created,
    (SELECT COUNT(*) FROM role) as roles_available,
    (SELECT COUNT(*) FROM work_shifts) as work_shifts_created,
    (SELECT COUNT(*) FROM attendance_category) as attendance_categories_created;

-- Display the data
SELECT 'Offices created:' as info;
SELECT office_id, office_name, phone_number FROM office ORDER BY office_id;

SELECT 'Departments created:' as info;
SELECT department_id, department_name, office_id FROM department ORDER BY department_id;

SELECT 'Positions created:' as info;
SELECT position_id, position_name FROM position ORDER BY position_id;

SELECT 'Job titles created:' as info;
SELECT job_title_id, job_title, hierarchy_level_id FROM jobtitle ORDER BY job_title_id;

SELECT 'Work shifts created:' as info;
SELECT id, shift_name, start_time, end_time, work_type FROM work_shifts ORDER BY id;

SELECT 'Attendance categories created:' as info;
SELECT id, category_name, work_day_value FROM attendance_category ORDER BY id;

SELECT 'Available roles:' as info;
SELECT id, role_name FROM role ORDER BY id;

-- Instructions for next steps
SELECT 'NEXT STEPS:' as info;
SELECT 'Run generate_test_data.sql to create employees and accounts' as instruction;
SELECT 'Run generate_comprehensive_test_data.sql for full ERP data' as instruction;
