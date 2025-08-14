-- Simple script to create test data for existing tables only
-- This script works with the current database schema

-- Create extension for random data generation if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO role (id, role_name, created_date)
VALUES
    ('admin', 'Admin', NOW()),
    ('manager', 'Manager', NOW()),
    ('employee', 'Employee', NOW())

-- Insert sample offices (this table exists)
INSERT INTO office (office_id, office_name, phone_number, address, latitude, longitude, created_date)
VALUES 
    ('HN001', 'Văn phòng Hà Nội', '0243456789', '123 Đường Láng, Đống Đa, Hà Nội', 21.0285, 105.8542, NOW()),
    ('PB0001', 'Văn phòng TP.HCM', '0281234567', '456 Nguyễn Văn Cừ, Quận 5, TP.HCM', 10.7769, 106.7009, NOW()),
    ('DN001', 'Văn phòng Đà Nẵng', '0236512345', '789 Trần Phú, Hải Châu, Đà Nẵng', 16.0471, 108.2068, NOW())
ON CONFLICT (office_id) DO NOTHING;

INSERT INTO department (department_id, department_name, manager, created_date, office_id)
VALUES
    ('PB0001', 'Nghiên cứu và Phát triển', 'm@exampl', '2025-05-13 04:54:26.226035+00', 'VP0001')

INSERT INTO jobtitle (jobtitle_id, jobtitle_name, created_date, hierarchy_level_id)
VALUES
    ('CV0001', 'Thực tập sinh FE', '2025-05-16 02:30:17.893022+00', 'CB0002')

INSERT INTO hierarchylevel (id, level_name, level_order, created_date)
VALUES
    ('CB0001', 'Cán bộ', 1, NOW()),
    ('CB0002', 'Chuyên viên', 2, NOW()),
    ('CB0003', 'Trưởng phòng', 3, NOW()),
    ('CB0004', 'Giám đốc', 4, NOW());

-- Insert sample positions (this table exists)
INSERT INTO position (position_id, position_name, created_date)
VALUES 
    ('POS001', 'Nhân viên IT', NOW()),
    ('POS002', 'Nhân viên HR', NOW()),
    ('POS003', 'Nhân viên Tài chính', NOW()),
    ('POS004', 'Nhân viên Marketing', NOW()),
    ('POS005', 'Nhân viên Kinh doanh', NOW()),
    ('POS006', 'Trưởng phòng IT', NOW()),
    ('POS007', 'Trưởng phòng HR', NOW()),
    ('POS008', 'Trưởng phòng Tài chính', NOW()),
    ('POS009', 'Giám đốc', NOW())
ON CONFLICT (position_id) DO NOTHING;

-- Note: Additional roles are already created by the backend (admin, manager, employee)

-- Display what we've created
SELECT 
    'Limited Data Generation Complete!' as status,
    (SELECT COUNT(*) FROM office) as offices_created,
    (SELECT COUNT(*) FROM position) as positions_created,
    (SELECT COUNT(*) FROM role) as roles_available;

-- Display the data
SELECT 'Offices created:' as info;
SELECT office_id, office_name, phone_number FROM office ORDER BY office_id;

SELECT 'Positions created:' as info;
SELECT position_id, position_name FROM position ORDER BY position_id;

SELECT 'Available roles:' as info;
SELECT id, role_name FROM role ORDER BY id;

-- Instructions for next steps
SELECT 'NEXT STEPS:' as info;
SELECT 'The backend needs to run migrations to create all tables.' as instruction;
SELECT 'Check backend logs and ensure AutoMigrate runs with all models uncommented.' as instruction;
