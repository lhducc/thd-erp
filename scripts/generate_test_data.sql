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
ON CONFLICT (position_id) DO NOTHING;

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

-- Function to normalize Vietnamese text to ASCII
CREATE OR REPLACE FUNCTION normalize_vietnamese_text(input_text TEXT) RETURNS TEXT AS $$
DECLARE
    normalized_text TEXT;
BEGIN
    normalized_text := lower(input_text);
    
    -- Replace Vietnamese characters with ASCII equivalents
    normalized_text := replace(normalized_text, 'á', 'a');
    normalized_text := replace(normalized_text, 'à', 'a');
    normalized_text := replace(normalized_text, 'ả', 'a');
    normalized_text := replace(normalized_text, 'ã', 'a');
    normalized_text := replace(normalized_text, 'ạ', 'a');
    normalized_text := replace(normalized_text, 'ă', 'a');
    normalized_text := replace(normalized_text, 'ắ', 'a');
    normalized_text := replace(normalized_text, 'ằ', 'a');
    normalized_text := replace(normalized_text, 'ẳ', 'a');
    normalized_text := replace(normalized_text, 'ẵ', 'a');
    normalized_text := replace(normalized_text, 'ặ', 'a');
    normalized_text := replace(normalized_text, 'â', 'a');
    normalized_text := replace(normalized_text, 'ấ', 'a');
    normalized_text := replace(normalized_text, 'ầ', 'a');
    normalized_text := replace(normalized_text, 'ẩ', 'a');
    normalized_text := replace(normalized_text, 'ẫ', 'a');
    normalized_text := replace(normalized_text, 'ậ', 'a');
    
    normalized_text := replace(normalized_text, 'é', 'e');
    normalized_text := replace(normalized_text, 'è', 'e');
    normalized_text := replace(normalized_text, 'ẻ', 'e');
    normalized_text := replace(normalized_text, 'ẽ', 'e');
    normalized_text := replace(normalized_text, 'ẹ', 'e');
    normalized_text := replace(normalized_text, 'ê', 'e');
    normalized_text := replace(normalized_text, 'ế', 'e');
    normalized_text := replace(normalized_text, 'ề', 'e');
    normalized_text := replace(normalized_text, 'ể', 'e');
    normalized_text := replace(normalized_text, 'ễ', 'e');
    normalized_text := replace(normalized_text, 'ệ', 'e');
    
    normalized_text := replace(normalized_text, 'í', 'i');
    normalized_text := replace(normalized_text, 'ì', 'i');
    normalized_text := replace(normalized_text, 'ỉ', 'i');
    normalized_text := replace(normalized_text, 'ĩ', 'i');
    normalized_text := replace(normalized_text, 'ị', 'i');
    
    normalized_text := replace(normalized_text, 'ó', 'o');
    normalized_text := replace(normalized_text, 'ò', 'o');
    normalized_text := replace(normalized_text, 'ỏ', 'o');
    normalized_text := replace(normalized_text, 'õ', 'o');
    normalized_text := replace(normalized_text, 'ọ', 'o');
    normalized_text := replace(normalized_text, 'ô', 'o');
    normalized_text := replace(normalized_text, 'ố', 'o');
    normalized_text := replace(normalized_text, 'ồ', 'o');
    normalized_text := replace(normalized_text, 'ổ', 'o');
    normalized_text := replace(normalized_text, 'ỗ', 'o');
    normalized_text := replace(normalized_text, 'ộ', 'o');
    normalized_text := replace(normalized_text, 'ơ', 'o');
    normalized_text := replace(normalized_text, 'ớ', 'o');
    normalized_text := replace(normalized_text, 'ờ', 'o');
    normalized_text := replace(normalized_text, 'ở', 'o');
    normalized_text := replace(normalized_text, 'ỡ', 'o');
    normalized_text := replace(normalized_text, 'ợ', 'o');
    
    normalized_text := replace(normalized_text, 'ú', 'u');
    normalized_text := replace(normalized_text, 'ù', 'u');
    normalized_text := replace(normalized_text, 'ủ', 'u');
    normalized_text := replace(normalized_text, 'ũ', 'u');
    normalized_text := replace(normalized_text, 'ụ', 'u');
    normalized_text := replace(normalized_text, 'ư', 'u');
    normalized_text := replace(normalized_text, 'ứ', 'u');
    normalized_text := replace(normalized_text, 'ừ', 'u');
    normalized_text := replace(normalized_text, 'ử', 'u');
    normalized_text := replace(normalized_text, 'ữ', 'u');
    normalized_text := replace(normalized_text, 'ự', 'u');
    
    normalized_text := replace(normalized_text, 'ý', 'y');
    normalized_text := replace(normalized_text, 'ỳ', 'y');
    normalized_text := replace(normalized_text, 'ỷ', 'y');
    normalized_text := replace(normalized_text, 'ỹ', 'y');
    normalized_text := replace(normalized_text, 'ỵ', 'y');
    
    normalized_text := replace(normalized_text, 'đ', 'd');
    
    -- Replace spaces with dots and remove any remaining special characters
    normalized_text := replace(normalized_text, ' ', '.');
    normalized_text := regexp_replace(normalized_text, '[^a-z0-9.]', '', 'g');
    
    RETURN normalized_text;
END;
$$ LANGUAGE plpgsql;

-- Function to generate random email
CREATE OR REPLACE FUNCTION random_email(full_name TEXT) RETURNS TEXT AS $$
BEGIN
    RETURN normalize_vietnamese_text(full_name) || floor(random() * 1000)::TEXT || '@email.com';
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
    
    positions TEXT[] := ARRAY(SELECT p.position_id FROM position p ORDER BY p.position_id);
    job_titles TEXT[] := ARRAY(SELECT j.job_title_id FROM jobtitle j ORDER BY j.job_title_id);
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
    -- Check if we have positions and job titles before proceeding
    IF array_length(positions, 1) = 0 THEN
        RAISE EXCEPTION 'No positions found in database';
    END IF;
    IF array_length(job_titles, 1) = 0 THEN
        RAISE EXCEPTION 'No job titles found in database';
    END IF;
    
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
        login_email := normalize_vietnamese_text(emp.full_name) || '@company.vn';
        hashed_password := '$2a$10$5RKK2ig9EPLErMn/Ue2Wi.VHd156Htap9yIde6XErnTPC18GvA9ba'; -- password: "password123"
        
        INSERT INTO account (login_mail, password, first_login, role_id, employee_id, created_date)
        VALUES (login_email, hashed_password, true, 'employee', emp.employee_id, NOW())
        RETURNING id INTO new_account_id;
        
        -- Update employee with account_id
        UPDATE employee SET account_id = new_account_id WHERE employee_id = emp.employee_id;
    END LOOP;
    
    -- Create some manager accounts (avoid duplicates)
    INSERT INTO account (login_mail, password, first_login, role_id, employee_id, created_date)
    VALUES 
        ('manager@company.vn', '$2a$10$5RKK2ig9EPLErMn/Ue2Wi.VHd156Htap9yIde6XErnTPC18GvA9ba', false, 'manager', (SELECT employee_id FROM employee LIMIT 1), NOW()),
        ('admin@company.vn', '$2a$10$5RKK2ig9EPLErMn/Ue2Wi.VHd156Htap9yIde6XErnTPC18GvA9ba', false, 'admin', (SELECT employee_id FROM employee OFFSET 1 LIMIT 1), NOW())
    ON CONFLICT (login_mail) DO NOTHING;
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
DROP FUNCTION normalize_vietnamese_text(TEXT);
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
SELECT a.id, a.login_mail, a.role_id, a.employee_id, e.full_name 
FROM account a
LEFT JOIN employee e ON a.employee_id = e.employee_id
LIMIT 5;
