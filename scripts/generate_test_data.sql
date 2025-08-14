-- Script to create random test data for ERP system
-- Run this script in PostgreSQL to g-- Insert sample positions
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
ON CONFLICT (position_id) DO NOTHING;le users and related data

-- Create extension for random data generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Function to generate random Vietnamese names
CREATE OR REPLACE FUNCTION random_vietnamese_name() RETURNS TEXT AS $$
DECLARE
    first_names TEXT[] := ARRAY[
        'Nguyễn', 'Trần', 'Lê', 'Phạm', 'Hoàng', 'Huỳnh', 'Phan', 'Vũ', 'Võ', 'Đặng',
        'Bùi', 'Đỗ', 'Hồ', 'Ngô', 'Dương', 'Lý', 'Mai', 'Đinh', 'Tô', 'Lưu'
    ];
    last_names TEXT[] := ARRAY[
        'Văn Hùng', 'Thị Lan', 'Đức Minh', 'Thị Hoa', 'Văn Nam', 'Thị Mai', 'Đức Anh', 'Thị Linh',
        'Quang Minh', 'Thị Thu', 'Văn Đức', 'Thị Nga', 'Quang Huy', 'Thị Thủy', 'Văn Tuấn',
        'Thị Vân', 'Đức Long', 'Thị Hằng', 'Văn Phong', 'Thị Trang', 'Quang Đại', 'Thị Hương'
    ];
    first_name TEXT;
    last_name TEXT;
BEGIN
    first_name := first_names[floor(random() * array_length(first_names, 1) + 1)];
    last_name := last_names[floor(random() * array_length(last_names, 1) + 1)];
    RETURN first_name || ' ' || last_name;
END;
$$ LANGUAGE plpgsql;

-- Function to generate random phone numbers
CREATE OR REPLACE FUNCTION random_phone() RETURNS TEXT AS $$
BEGIN
    RETURN '0' || (floor(random() * 9) + 1)::TEXT || 
           lpad((floor(random() * 100000000))::TEXT, 8, '0');
END;
$$ LANGUAGE plpgsql;

-- Function to generate random email
CREATE OR REPLACE FUNCTION random_email(full_name TEXT) RETURNS TEXT AS $$
BEGIN
    RETURN lower(replace(replace(full_name, ' ', '.'), 'đ', 'd')) || 
           floor(random() * 1000)::TEXT || '@email.com';
END;
$$ LANGUAGE plpgsql;

-- Function to generate random birthday
CREATE OR REPLACE FUNCTION random_birthday() RETURNS DATE AS $$
BEGIN
    RETURN '1980-01-01'::DATE + (floor(random() * 15000))::INTEGER;
END;
$$ LANGUAGE plpgsql;

-- Clear existing test data (keep schema)
DELETE FROM account WHERE employee_id IS NOT NULL;
DELETE FROM decision_employees;
DELETE FROM contract;
DELETE FROM employee;
DELETE FROM jobtitle;
DELETE FROM position;
DELETE FROM department;
DELETE FROM office;

-- Insert sample offices
INSERT INTO office (office_id, office_name, phone_number, address, latitude, longitude, created_date)
VALUES 
    ('HN001', 'Văn phòng Hà Nội', '+842432345678', '123 Đường Láng, Đống Đa, Hà Nội', 21.0285, 105.8542, NOW()),
    ('HCM001', 'Văn phòng TP.HCM', '+842812345678', '456 Nguyễn Văn Cừ, Quận 5, TP.HCM', 10.7769, 106.7009, NOW()),
    ('DN001', 'Văn phòng Đà Nẵng', '+842365123456', '789 Trần Phú, Hải Châu, Đà Nẵng', 16.0471, 108.2068, NOW())
ON CONFLICT (office_id) DO NOTHING;

-- Insert sample departments
INSERT INTO department (department_id, department_name, manager, created_date, office_id)
VALUES 
    ('HR001', 'Phòng Nhân sự', '', NOW(), 'HN001'),
    ('IT001', 'Phòng Công nghệ thông tin', '', NOW(), 'HN001'),
    ('FIN001', 'Phòng Tài chính', '', NOW(), 'HN001'),
    ('MKT001', 'Phòng Marketing', '', NOW(), 'HCM001'),
    ('SALE001', 'Phòng Kinh doanh', '', NOW(), 'HCM001'),
    ('OP001', 'Phòng Vận hành', '', NOW(), 'DN001')
ON CONFLICT (department_id) DO NOTHING;

-- Insert sample hierarchy levels
INSERT INTO hierarchy_level (id, level_name, level_order, created_date)
VALUES 
    ('LV001', 'Nhân viên', 1, NOW()),
    ('LV002', 'Trưởng nhóm', 2, NOW()),
    ('LV003', 'Trưởng phòng', 3, NOW()),
    ('LV004', 'Phó giám đốc', 4, NOW()),
    ('LV005', 'Giám đốc', 5, NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample positions
INSERT INTO position (position_id, position_name, description, created_date, hierarchy_level_id)
VALUES 
    ('POS001', 'Nhân viên', 'Nhân viên cơ bản', NOW(), 'LV001'),
    ('POS002', 'Nhân viên Senior', 'Nhân viên có kinh nghiệm', NOW(), 'LV001'),
    ('POS003', 'Trưởng nhóm', 'Quản lý nhóm nhỏ', NOW(), 'LV002'),
    ('POS004', 'Phó trưởng phòng', 'Phụ trách một phần công việc phòng', NOW(), 'LV002'),
    ('POS005', 'Trưởng phòng', 'Quản lý toàn bộ phòng ban', NOW(), 'LV003'),
    ('POS006', 'Phó giám đốc', 'Hỗ trợ giám đốc', NOW(), 'LV004'),
    ('POS007', 'Giám đốc', 'Quản lý toàn công ty', NOW(), 'LV005'),
    ('POS008', 'Chuyên viên', 'Chuyên gia trong lĩnh vực', NOW(), 'LV001')
ON CONFLICT (position_id) DO NOTHING;

-- Insert sample job titles
INSERT INTO jobtitle (job_title_id, job_title, created_date)
VALUES 
    ('JT001', 'Lập trình viên', NOW()),
    ('JT002', 'Nhân viên nhân sự', NOW()),
    ('JT003', 'Kế toán', NOW()),
    ('JT004', 'Marketing', NOW()),
    ('JT005', 'Kinh doanh', NOW()),
    ('JT006', 'Vận hành', NOW()),
    ('JT007', 'Thiết kế', NOW()),
    ('JT008', 'Tester', NOW())
ON CONFLICT (job_title_id) DO NOTHING;

-- Generate 50 random employees
DO $$
DECLARE
    i INTEGER;
    emp_id TEXT;
    full_name TEXT;
    birthday DATE;
    gender TEXT;
    work_type TEXT;
    phone TEXT;
    email TEXT;
    address TEXT;
    position_id TEXT;
    job_title_id TEXT;
    department_id TEXT;
    
    positions TEXT[] := ARRAY['POS001', 'POS002', 'POS003', 'POS004', 'POS005', 'POS006', 'POS007', 'POS008'];
    job_titles TEXT[] := ARRAY['JT001', 'JT002', 'JT003', 'JT004', 'JT005', 'JT006', 'JT007', 'JT008'];
    departments TEXT[] := ARRAY['HR001', 'IT001', 'FIN001', 'MKT001', 'SALE001', 'OP001'];
    work_types TEXT[] := ARRAY['ca hành chính', 'ca kíp'];
    genders TEXT[] := ARRAY['Nam', 'Nữ'];
    addresses TEXT[] := ARRAY[
        'Quận Ba Đình, Hà Nội',
        'Quận Hoàn Kiếm, Hà Nội', 
        'Quận Đống Đa, Hà Nội',
        'Quận 1, TP.HCM',
        'Quận 3, TP.HCM',
        'Quận 5, TP.HCM',
        'Quận Hải Châu, Đà Nẵng',
        'Quận Thanh Khê, Đà Nẵng'
    ];
BEGIN
    FOR i IN 1..50 LOOP
        emp_id := 'EMP' || lpad(i::TEXT, 5, '0');
        full_name := random_vietnamese_name();
        birthday := random_birthday();
        gender := genders[floor(random() * array_length(genders, 1) + 1)];
        work_type := work_types[floor(random() * array_length(work_types, 1) + 1)];
        phone := random_phone();
        email := random_email(full_name);
        address := addresses[floor(random() * array_length(addresses, 1) + 1)];
        position_id := positions[floor(random() * array_length(positions, 1) + 1)];
        job_title_id := job_titles[floor(random() * array_length(job_titles, 1) + 1)];
        department_id := departments[floor(random() * array_length(departments, 1) + 1)];
        
        INSERT INTO employee (
            employee_id, full_name, birthday, gender, work_type, 
            phone_number, email, address, position_id, job_title_id, 
            status, department_id, created_date
        ) VALUES (
            emp_id, full_name, birthday, gender, work_type::work_type_enum,
            phone, email, address, position_id, job_title_id,
            'active', department_id, NOW()
        );
    END LOOP;
END $$;

-- Generate accounts for some employees
DO $$
DECLARE
    emp RECORD;
    login_email TEXT;
    hashed_password TEXT;
    new_account_id BIGINT;
BEGIN
    -- Create accounts for first 20 employees
    FOR emp IN (SELECT employee_id, full_name, email FROM employee LIMIT 20) LOOP
        login_email := lower(replace(replace(emp.full_name, ' ', '.'), 'đ', 'd')) || '@company.vn';
        hashed_password := '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'; -- password: "password123"
        
        INSERT INTO account (login_mail, password, first_login, role_id, employee_id, created_date)
        VALUES (login_email, hashed_password, true, 'employee', emp.employee_id, NOW())
        RETURNING id INTO new_account_id;
        
        -- Update employee with account_id
        UPDATE employee SET account_id = new_account_id WHERE employee_id = emp.employee_id;
    END LOOP;
    
    -- Create some manager accounts
    INSERT INTO account (login_mail, password, first_login, role_id, employee_id, created_date)
    VALUES 
        ('manager@company.vn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'manager', (SELECT employee_id FROM employee LIMIT 1), NOW()),
        ('admin@company.vn', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'admin', (SELECT employee_id FROM employee OFFSET 1 LIMIT 1), NOW());
END $$;

-- Set some employees as managers of departments
DO $$
DECLARE
    dept RECORD;
    manager_emp_id TEXT;
BEGIN
    FOR dept IN (SELECT department_id FROM department) LOOP
        SELECT employee_id INTO manager_emp_id 
        FROM employee 
        WHERE department_id = dept.department_id 
        ORDER BY random() 
        LIMIT 1;
        
        IF manager_emp_id IS NOT NULL THEN
            UPDATE department 
            SET manager = manager_emp_id 
            WHERE department_id = dept.department_id;
        END IF;
    END LOOP;
END $$;

-- Clean up functions
DROP FUNCTION random_vietnamese_name();
DROP FUNCTION random_phone();
DROP FUNCTION random_email(TEXT);
DROP FUNCTION random_birthday();

-- Show summary
SELECT 
    'Data Generation Complete!' as status,
    (SELECT COUNT(*) FROM office) as offices_created,
    (SELECT COUNT(*) FROM department) as departments_created,
    (SELECT COUNT(*) FROM position) as positions_created,
    (SELECT COUNT(*) FROM jobtitle) as job_titles_created,
    (SELECT COUNT(*) FROM employee) as employees_created,
    (SELECT COUNT(*) FROM account) as accounts_created;

-- Show sample data
SELECT 'Sample Employees:' as info;
SELECT employee_id, full_name, gender, work_type, department_id, status 
FROM employee 
LIMIT 5;

SELECT 'Sample Accounts:' as info;
SELECT a.account_id, a.login_mail, a.role_id, a.employee_id, e.full_name 
FROM account a
LEFT JOIN employee e ON a.employee_id = e.employee_id
LIMIT 5;
