-- Comprehensive ERP Test Data Generation Script
-- This script generates test data for all tables in the ERP system
-- Run this after generate_test_data.sql to add missing entities

-- Create extension for random data generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Function to generate random Vietnamese company names
CREATE OR REPLACE FUNCTION random_company_name() RETURNS TEXT AS $$
DECLARE
    company_types TEXT[] := ARRAY[
        'Công ty TNHH', 'Công ty Cổ phần', 'Tập đoàn', 'Công ty',
        'Doanh nghiệp', 'Tổng công ty'
    ];
    company_names TEXT[] := ARRAY[
        'Thăng Long', 'Hoàng Hà', 'Nam Việt', 'Đông Á', 'Bắc Thăng Long',
        'Sài Gòn', 'Mekong', 'Hà Nội', 'Đà Nẵng', 'Hải Phòng',
        'Viễn Đông', 'Kim Cương', 'Hoàng Gia', 'Thịnh Vượng', 'An Phát'
    ];
    business_areas TEXT[] := ARRAY[
        'Xây dựng', 'Công nghệ', 'Thương mại', 'Dịch vụ', 'Sản xuất',
        'Logistics', 'Tài chính', 'Bất động sản', 'Y tế', 'Giáo dục'
    ];
BEGIN
    RETURN company_types[floor(random() * array_length(company_types, 1) + 1)] || ' ' ||
           company_names[floor(random() * array_length(company_names, 1) + 1)] || ' ' ||
           business_areas[floor(random() * array_length(business_areas, 1) + 1)];
END;
$$ LANGUAGE plpgsql;

-- Function to generate random amounts (salary, allowance, etc.)
CREATE OR REPLACE FUNCTION random_amount(min_val INTEGER, max_val INTEGER) RETURNS DECIMAL AS $$
BEGIN
    RETURN (min_val + floor(random() * (max_val - min_val + 1))) * 1000;
END;
$$ LANGUAGE plpgsql;

-- Function to generate random date in range
CREATE OR REPLACE FUNCTION random_date_range(start_date DATE, end_date DATE) RETURNS DATE AS $$
BEGIN
    RETURN start_date + floor(random() * (end_date - start_date + 1))::INTEGER;
END;
$$ LANGUAGE plpgsql;

-- Function to generate random time
CREATE OR REPLACE FUNCTION random_time(start_hour INTEGER, end_hour INTEGER) RETURNS TIME AS $$
DECLARE
    hour INTEGER;
    minute INTEGER;
BEGIN
    hour := start_hour + floor(random() * (end_hour - start_hour + 1));
    minute := floor(random() * 60);
    RETURN make_time(hour, minute, 0);
END;
$$ LANGUAGE plpgsql;

-- ===== 1. EMPLOYEE DOCUMENT TYPES =====
INSERT INTO employee_document_type (id, document_type_name, created_date)
VALUES 
    ('DOC001', 'Chứng minh nhân dân', NOW()),
    ('DOC002', 'Căn cước công dân', NOW()),
    ('DOC003', 'Hộ chiếu', NOW()),
    ('DOC004', 'Bằng tốt nghiệp đại học', NOW()),
    ('DOC005', 'Bằng tốt nghiệp cao đẳng', NOW()),
    ('DOC006', 'Chứng chỉ nghề', NOW()),
    ('DOC007', 'Giấy khám sức khỏe', NOW()),
    ('DOC008', 'Sổ bảo hiểm xã hội', NOW()),
    ('DOC009', 'Hồ sơ lý lịch', NOW()),
    ('DOC010', 'Giấy xác nhận độc thân', NOW())
ON CONFLICT (id) DO NOTHING;

-- ===== 2. EMPLOYEE DOCUMENTS =====
DO $$
DECLARE
    emp RECORD;
    doc_types TEXT[] := ARRAY['DOC001', 'DOC002', 'DOC004', 'DOC007', 'DOC008'];
    doc_type TEXT;
    doc_number TEXT;
    i INTEGER;
BEGIN
    FOR emp IN (SELECT employee_id FROM employee ORDER BY employee_id) LOOP
        -- Generate 2-4 random documents per employee
        FOR i IN 1..(2 + floor(random() * 3))::INTEGER LOOP
            doc_type := doc_types[floor(random() * array_length(doc_types, 1) + 1)];
            doc_number := 'DOC' || emp.employee_id || '-' || lpad(i::TEXT, 3, '0');
            
            INSERT INTO employee_document (
                id, employee_id, document_type_id, document_number,
                issue_date, issue_place, expiry_date, file_path, created_date
            ) VALUES (
                uuid_generate_v4()::TEXT,
                emp.employee_id,
                doc_type,
                doc_number,
                random_date_range('2020-01-01'::DATE, '2023-12-31'::DATE),
                CASE 
                    WHEN doc_type = 'DOC001' THEN 'Công an TP.HCM'
                    WHEN doc_type = 'DOC002' THEN 'Cục Cảnh sát ĐKQL cư trú và DLQG về dân cư'
                    WHEN doc_type = 'DOC004' THEN 'Đại học Bách Khoa Hà Nội'
                    WHEN doc_type = 'DOC007' THEN 'Bệnh viện Đa khoa Hà Nội'
                    ELSE 'Cơ quan có thẩm quyền'
                END,
                CASE 
                    WHEN doc_type IN ('DOC001', 'DOC002', 'DOC003') THEN random_date_range('2025-01-01'::DATE, '2035-12-31'::DATE)
                    WHEN doc_type = 'DOC007' THEN random_date_range('2024-06-01'::DATE, '2024-12-31'::DATE)
                    ELSE NULL
                END,
                '/documents/' || emp.employee_id || '/' || doc_number || '.pdf',
                NOW()
            )
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- ===== 3. CONTRACT TYPES =====
INSERT INTO contract_type (id, contract_type_name, description, created_date)
VALUES 
    ('CT001', 'Hợp đồng không thời hạn', 'Hợp đồng lao động không xác định thời hạn', NOW()),
    ('CT002', 'Hợp đồng có thời hạn', 'Hợp đồng lao động xác định thời hạn', NOW()),
    ('CT003', 'Hợp đồng thử việc', 'Hợp đồng lao động thử việc', NOW()),
    ('CT004', 'Hợp đồng theo mùa vụ', 'Hợp đồng lao động theo mùa vụ hoặc công việc nhất định', NOW()),
    ('CT005', 'Hợp đồng bán thời gian', 'Hợp đồng lao động bán thời gian', NOW())
ON CONFLICT (id) DO NOTHING;

-- ===== 4. CONTRACTS =====
DO $$
DECLARE
    emp RECORD;
    contract_types TEXT[] := ARRAY['CT001', 'CT002', 'CT003'];
    start_date DATE;
    end_date DATE;
    salary DECIMAL;
    contract_number TEXT;
    counter INTEGER := 1;
BEGIN
    FOR emp IN (SELECT employee_id FROM employee ORDER BY employee_id) LOOP
        start_date := random_date_range('2022-01-01'::DATE, '2024-06-30'::DATE);
        
        -- Determine end date based on contract type
        IF contract_types[floor(random() * array_length(contract_types, 1) + 1)] = 'CT001' THEN
            end_date := NULL; -- Indefinite contract
        ELSIF contract_types[floor(random() * array_length(contract_types, 1) + 1)] = 'CT003' THEN
            end_date := start_date + interval '2 months'; -- Trial contract
        ELSE
            end_date := start_date + interval '1 year' + (floor(random() * 24) || ' months')::interval; -- Fixed term
        END IF;
        
        salary := random_amount(8, 50); -- 8M to 50M VND
        contract_number := 'HĐLĐ-' || date_part('year', start_date)::TEXT || '-' || lpad(counter::TEXT, 4, '0');
        
        INSERT INTO contract (
            id, employee_id, contract_type_id, contract_number,
            start_date, end_date, basic_salary, position_allowance,
            other_allowance, status, created_date
        ) VALUES (
            uuid_generate_v4()::TEXT,
            emp.employee_id,
            contract_types[floor(random() * array_length(contract_types, 1) + 1)],
            contract_number,
            start_date,
            end_date,
            salary,
            salary * 0.1, -- 10% position allowance
            salary * 0.05, -- 5% other allowance
            CASE WHEN end_date IS NULL OR end_date > CURRENT_DATE THEN 'active' ELSE 'expired' END,
            NOW()
        );
        
        counter := counter + 1;
    END LOOP;
END $$;

-- ===== 5. DECISION TYPES =====
INSERT INTO decision_type (id, decision_type_name, description, created_date)
VALUES 
    ('DT001', 'Quyết định tuyển dụng', 'Quyết định tuyển dụng nhân viên mới', NOW()),
    ('DT002', 'Quyết định thăng chức', 'Quyết định thăng chức, bổ nhiệm', NOW()),
    ('DT003', 'Quyết định tăng lương', 'Quyết định tăng lương định kỳ', NOW()),
    ('DT004', 'Quyết định khen thưởng', 'Quyết định khen thưởng nhân viên', NOW()),
    ('DT005', 'Quyết định kỷ luật', 'Quyết định kỷ luật vi phạm', NOW()),
    ('DT006', 'Quyết định nghỉ phép', 'Quyết định nghỉ phép dài ngày', NOW()),
    ('DT007', 'Quyết định chuyển công tác', 'Quyết định chuyển công tác, điều động', NOW()),
    ('DT008', 'Quyết định thôi việc', 'Quyết định chấm dứt hợp đồng lao động', NOW())
ON CONFLICT (id) DO NOTHING;

-- ===== 6. DECISIONS =====
DO $$
DECLARE
    decision_types TEXT[] := ARRAY['DT001', 'DT002', 'DT003', 'DT004'];
    decision_counter INTEGER := 1;
    decision_date DATE;
    effective_date DATE;
    decision_number TEXT;
    managers TEXT[] := ARRAY(SELECT employee_id FROM employee WHERE position_id IN ('POS005', 'POS006', 'POS007') LIMIT 5);
    decision_type TEXT;
    i INTEGER;
BEGIN
    FOR i IN 1..30 LOOP
        decision_type := decision_types[floor(random() * array_length(decision_types, 1) + 1)];
        decision_date := random_date_range('2023-01-01'::DATE, '2024-12-31'::DATE);
        effective_date := decision_date + (floor(random() * 30) || ' days')::interval;
        decision_number := 'QĐ-' || date_part('year', decision_date)::TEXT || '-' || lpad(decision_counter::TEXT, 4, '0');
        
        INSERT INTO decision (
            id, decision_type_id, decision_number, title,
            content, decision_date, effective_date, created_by, status, created_date
        ) VALUES (
            uuid_generate_v4()::TEXT,
            decision_type,
            decision_number,
            CASE decision_type
                WHEN 'DT001' THEN 'Quyết định tuyển dụng nhân viên'
                WHEN 'DT002' THEN 'Quyết định thăng chức nhân viên'
                WHEN 'DT003' THEN 'Quyết định tăng lương định kỳ'
                WHEN 'DT004' THEN 'Quyết định khen thưởng quý'
                ELSE 'Quyết định hành chính'
            END,
            'Nội dung quyết định ' || decision_number || ' về ' ||
            CASE decision_type
                WHEN 'DT001' THEN 'việc tuyển dụng nhân viên mới vào công ty'
                WHEN 'DT002' THEN 'việc thăng chức và bổ nhiệm chức vụ mới'
                WHEN 'DT003' THEN 'việc điều chỉnh mức lương cho nhân viên'
                WHEN 'DT004' THEN 'việc khen thưởng nhân viên có thành tích tốt'
                ELSE 'các vấn đề hành chính khác'
            END,
            decision_date,
            effective_date,
            managers[floor(random() * array_length(managers, 1) + 1)],
            'approved',
            NOW()
        );
        
        decision_counter := decision_counter + 1;
    END LOOP;
END $$;

-- ===== 7. DECISION EMPLOYEES (Link employees to decisions) =====
DO $$
DECLARE
    decision RECORD;
    emp_list TEXT[];
    emp_id TEXT;
    i INTEGER;
    num_employees INTEGER;
BEGIN
    FOR decision IN (SELECT id FROM decision ORDER BY created_date) LOOP
        -- Each decision affects 1-5 employees
        num_employees := 1 + floor(random() * 5)::INTEGER;
        emp_list := ARRAY(SELECT employee_id FROM employee ORDER BY random() LIMIT num_employees);
        
        FOREACH emp_id IN ARRAY emp_list LOOP
            INSERT INTO decision_employees (decision_id, employee_id, created_date)
            VALUES (decision.id, emp_id, NOW())
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- ===== 8. ALLOWANCES =====
INSERT INTO allowance (id, allowance_name, amount, allowance_type, description, created_date)
VALUES 
    ('AL001', 'Phụ cấp xăng xe', 500000, 'monthly', 'Phụ cấp xăng xe hàng tháng', NOW()),
    ('AL002', 'Phụ cấp ăn trưa', 1000000, 'monthly', 'Phụ cấp ăn trưa hàng tháng', NOW()),
    ('AL003', 'Phụ cấp điện thoại', 300000, 'monthly', 'Phụ cấp điện thoại công việc', NOW()),
    ('AL004', 'Phụ cấp đêm', 200000, 'per_shift', 'Phụ cấp làm ca đêm', NOW()),
    ('AL005', 'Phụ cấp chủ nhật', 150000, 'per_shift', 'Phụ cấp làm việc chủ nhật', NOW()),
    ('AL006', 'Phụ cấp trách nhiệm', 2000000, 'monthly', 'Phụ cấp chức vụ trách nhiệm', NOW()),
    ('AL007', 'Phụ cấp nguy hiểm', 1500000, 'monthly', 'Phụ cấp làm việc môi trường nguy hiểm', NOW()),
    ('AL008', 'Phụ cấp khu vực', 800000, 'monthly', 'Phụ cấp làm việc vùng sâu vùng xa', NOW())
ON CONFLICT (id) DO NOTHING;

-- ===== 9. CONTRACT ALLOWANCES =====
DO $$
DECLARE
    contract RECORD;
    allowances TEXT[] := ARRAY['AL001', 'AL002', 'AL003', 'AL006'];
    allowance_id TEXT;
    i INTEGER;
BEGIN
    FOR contract IN (SELECT id FROM contract) LOOP
        -- Each contract gets 1-3 random allowances
        FOR i IN 1..(1 + floor(random() * 3))::INTEGER LOOP
            allowance_id := allowances[floor(random() * array_length(allowances, 1) + 1)];
            
            INSERT INTO contract_allowance (contract_id, allowance_id, created_date)
            VALUES (contract.id, allowance_id, NOW())
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- ===== 10. INSURANCE =====
DO $$
DECLARE
    emp RECORD;
    start_date DATE;
    insurance_number TEXT;
    counter INTEGER := 1;
BEGIN
    FOR emp IN (SELECT employee_id FROM employee ORDER BY employee_id) LOOP
        start_date := random_date_range('2022-01-01'::DATE, '2024-06-30'::DATE);
        insurance_number := 'BH' || date_part('year', start_date)::TEXT || lpad(counter::TEXT, 6, '0');
        
        INSERT INTO insurance (
            id, employee_id, insurance_number, insurance_type,
            start_date, end_date, premium_amount, coverage_amount,
            provider, status, created_date
        ) VALUES (
            uuid_generate_v4()::TEXT,
            emp.employee_id,
            insurance_number,
            CASE floor(random() * 3)
                WHEN 0 THEN 'health'
                WHEN 1 THEN 'social'
                ELSE 'accident'
            END,
            start_date,
            start_date + interval '1 year',
            random_amount(1, 5), -- Premium: 1-5M VND
            random_amount(50, 200), -- Coverage: 50-200M VND
            CASE floor(random() * 3)
                WHEN 0 THEN 'Bảo Việt'
                WHEN 1 THEN 'Bảo hiểm xã hội Việt Nam'
                ELSE 'Prudential'
            END,
            'active',
            NOW()
        );
        
        counter := counter + 1;
    END LOOP;
END $$;

-- ===== 11. WORK SHIFTS =====
INSERT INTO work_shifts (id, shift_name, start_time, end_time, break_duration, work_type, created_date)
VALUES 
    ('SHIFT001', 'Ca hành chính', '08:00:00', '17:00:00', 60, 'ca hành chính', NOW()),
    ('SHIFT002', 'Ca sáng', '06:00:00', '14:00:00', 60, 'ca kíp', NOW()),
    ('SHIFT003', 'Ca chiều', '14:00:00', '22:00:00', 60, 'ca kíp', NOW()),
    ('SHIFT004', 'Ca đêm', '22:00:00', '06:00:00', 60, 'ca kíp', NOW()),
    ('SHIFT005', 'Ca linh hoạt', '09:00:00', '18:00:00', 60, 'ca hành chính', NOW()),
    ('SHIFT006', 'Ca bán thời gian sáng', '08:00:00', '12:00:00', 0, 'ca hành chính', NOW()),
    ('SHIFT007', 'Ca bán thời gian chiều', '13:00:00', '17:00:00', 0, 'ca hành chính', NOW())
ON CONFLICT (id) DO NOTHING;

-- ===== 12. EMPLOYEE WORK SHIFTS =====
DO $$
DECLARE
    emp RECORD;
    shifts TEXT[] := ARRAY['SHIFT001', 'SHIFT002', 'SHIFT003', 'SHIFT005'];
    shift_id TEXT;
    assigned_date DATE;
BEGIN
    FOR emp IN (SELECT employee_id, work_type FROM employee) LOOP
        -- Assign appropriate shift based on work type
        IF emp.work_type = 'ca hành chính' THEN
            shift_id := CASE floor(random() * 2)
                WHEN 0 THEN 'SHIFT001'
                ELSE 'SHIFT005'
            END;
        ELSE
            shift_id := CASE floor(random() * 3)
                WHEN 0 THEN 'SHIFT002'
                WHEN 1 THEN 'SHIFT003'
                ELSE 'SHIFT004'
            END;
        END IF;
        
        assigned_date := random_date_range('2024-01-01'::DATE, CURRENT_DATE);
        
        INSERT INTO employee_workshift (employee_id, work_shift_id, assigned_date, created_date)
        VALUES (emp.employee_id, shift_id, assigned_date, NOW());
    END LOOP;
END $$;

-- ===== 13. ATTENDANCE CATEGORIES =====
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

-- ===== 14. ATTENDANCE RECORDS =====
DO $$
DECLARE
    emp RECORD;
    date_iter DATE;
    end_date DATE := CURRENT_DATE;
    start_date DATE := CURRENT_DATE - interval '30 days';
    categories TEXT[] := ARRAY['ATT001', 'ATT002', 'ATT003', 'ATT004', 'ATT005', 'ATT006'];
    category TEXT;
    check_in_time TIMESTAMPTZ;
    check_out_time TIMESTAMPTZ;
    base_check_in TIME;
    base_check_out TIME;
    late_minutes INTEGER;
    early_minutes INTEGER;
BEGIN
    FOR emp IN (SELECT e.employee_id, ew.work_shift_id 
                FROM employee e 
                LEFT JOIN employee_workshift ew ON e.employee_id = ew.employee_id 
                ORDER BY e.employee_id) LOOP
        
        -- Get shift times (default to office hours if no shift assigned)
        SELECT start_time, end_time INTO base_check_in, base_check_out 
        FROM work_shifts 
        WHERE id = COALESCE(emp.work_shift_id, 'SHIFT001');
        
        date_iter := start_date;
        WHILE date_iter <= end_date LOOP
            -- Skip weekends for office shifts (about 70% of the time)
            IF EXTRACT(dow FROM date_iter) NOT IN (0, 6) OR random() < 0.3 THEN
                
                -- Choose attendance category (mostly normal attendance)
                category := CASE 
                    WHEN random() < 0.7 THEN 'ATT001' -- Normal attendance
                    WHEN random() < 0.85 THEN 'ATT002' -- Late
                    WHEN random() < 0.95 THEN 'ATT003' -- Early leave
                    ELSE categories[floor(random() * array_length(categories, 1) + 1)]
                END;
                
                -- Calculate actual times based on category
                IF category = 'ATT001' THEN -- Normal
                    check_in_time := date_iter + base_check_in + (floor(random() * 10 - 5) || ' minutes')::interval;
                    check_out_time := date_iter + base_check_out + (floor(random() * 10 - 5) || ' minutes')::interval;
                    late_minutes := 0;
                    early_minutes := 0;
                ELSIF category = 'ATT002' THEN -- Late
                    late_minutes := 5 + floor(random() * 55)::INTEGER; -- 5-60 minutes late
                    check_in_time := date_iter + base_check_in + (late_minutes || ' minutes')::interval;
                    check_out_time := date_iter + base_check_out;
                    early_minutes := 0;
                ELSIF category = 'ATT003' THEN -- Early leave
                    early_minutes := 10 + floor(random() * 50)::INTEGER; -- 10-60 minutes early
                    check_in_time := date_iter + base_check_in;
                    check_out_time := date_iter + base_check_out - (early_minutes || ' minutes')::interval;
                    late_minutes := 0;
                ELSE -- Absent categories
                    check_in_time := NULL;
                    check_out_time := NULL;
                    late_minutes := 0;
                    early_minutes := 0;
                END IF;
                
                INSERT INTO attendance_record (
                    employee_id, attendance_date, check_in_time, check_out_time,
                    work_shift_id, attendance_category_id, late_minutes, early_leave_minutes,
                    notes, created_date
                ) VALUES (
                    emp.employee_id,
                    date_iter,
                    check_in_time,
                    check_out_time,
                    emp.work_shift_id,
                    category,
                    late_minutes,
                    early_minutes,
                    CASE 
                        WHEN category = 'ATT002' THEN 'Đến muộn do kẹt xe'
                        WHEN category = 'ATT003' THEN 'Về sớm có việc cá nhân'
                        WHEN category = 'ATT004' THEN 'Nghỉ phép đã đăng ký'
                        WHEN category = 'ATT005' THEN 'Nghỉ không báo trước'
                        WHEN category = 'ATT006' THEN 'Nghỉ ốm có giấy y tế'
                        ELSE NULL
                    END,
                    NOW()
                )
                ON CONFLICT DO NOTHING;
            END IF;
            
            date_iter := date_iter + interval '1 day';
        END LOOP;
    END LOOP;
END $$;

-- ===== 15. WORK SCHEDULES =====
DO $$
DECLARE
    offices TEXT[] := ARRAY(SELECT office_id FROM office);
    managers TEXT[] := ARRAY(SELECT employee_id FROM employee WHERE position_id IN ('POS005', 'POS006', 'POS007'));
    office_id TEXT;
    manager_id TEXT;
    start_date DATE;
    end_date DATE;
    schedule_name TEXT;
    i INTEGER;
BEGIN
    FOR i IN 1..10 LOOP
        office_id := offices[floor(random() * array_length(offices, 1) + 1)];
        manager_id := managers[floor(random() * array_length(managers, 1) + 1)];
        start_date := random_date_range('2024-01-01'::DATE, '2024-06-30'::DATE);
        end_date := start_date + interval '3 months' + (floor(random() * 6) || ' months')::interval;
        schedule_name := 'Lịch làm việc ' || 
                        CASE floor(random() * 4)
                            WHEN 0 THEN 'Quý ' || EXTRACT(quarter FROM start_date)::TEXT
                            WHEN 1 THEN 'Tháng ' || EXTRACT(month FROM start_date)::TEXT
                            WHEN 2 THEN 'Dự án đặc biệt'
                            ELSE 'Thường xuyên'
                        END || ' - ' || office_id;
        
        INSERT INTO work_schedule (
            schedule_name, start_date, end_date, repeat_type, 
            status, office_id, created_by, created_date
        ) VALUES (
            schedule_name,
            start_date,
            end_date,
            CASE floor(random() * 4)
                WHEN 0 THEN 'daily'
                WHEN 1 THEN 'weekly'
                WHEN 2 THEN 'monthly'
                ELSE 'none'
            END,
            CASE floor(random() * 3)
                WHEN 0 THEN 'active'
                WHEN 1 THEN 'inactive'
                ELSE 'pending'
            END,
            office_id,
            manager_id,
            NOW()
        );
    END LOOP;
END $$;

-- ===== 16. WORK SCHEDULE SHIFTS =====
DO $$
DECLARE
    schedule RECORD;
    shifts TEXT[] := ARRAY['SHIFT001', 'SHIFT002', 'SHIFT003', 'SHIFT005'];
    weekdays TEXT[] := ARRAY['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday'];
    shift_id TEXT;
    weekday TEXT;
    i INTEGER;
BEGIN
    FOR schedule IN (SELECT id FROM work_schedule) LOOP
        -- Each schedule has shifts for different weekdays
        FOR i IN 1..5 LOOP -- Monday to Friday
            shift_id := shifts[floor(random() * array_length(shifts, 1) + 1)];
            weekday := weekdays[i];
            
            INSERT INTO work_schedule_shift (
                work_schedule_id, work_shift_id, weekday, created_date
            ) VALUES (
                schedule.id, shift_id, weekday::weekday_enum, NOW()
            );
        END LOOP;
    END LOOP;
END $$;

-- ===== 17. WORK SCHEDULE MANAGERS =====
DO $$
DECLARE
    schedule RECORD;
    managers TEXT[] := ARRAY(SELECT employee_id FROM employee WHERE position_id IN ('POS004', 'POS005'));
    manager_id TEXT;
    i INTEGER;
BEGIN
    FOR schedule IN (SELECT id FROM work_schedule) LOOP
        -- Each schedule has 1-2 managers
        FOR i IN 1..(1 + floor(random() * 2))::INTEGER LOOP
            manager_id := managers[floor(random() * array_length(managers, 1) + 1)];
            
            INSERT INTO work_schedule_manager (
                work_schedule_id, manager_id, created_date
            ) VALUES (
                schedule.id, manager_id, NOW()
            )
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- Clean up helper functions
DROP FUNCTION random_company_name();
DROP FUNCTION random_amount(INTEGER, INTEGER);
DROP FUNCTION random_date_range(DATE, DATE);
DROP FUNCTION random_time(INTEGER, INTEGER);

-- Show comprehensive summary
SELECT 
    'Comprehensive Data Generation Complete!' as status,
    (SELECT COUNT(*) FROM employee_document_type) as document_types,
    (SELECT COUNT(*) FROM employee_document) as employee_documents,
    (SELECT COUNT(*) FROM contract_type) as contract_types,
    (SELECT COUNT(*) FROM contract) as contracts,
    (SELECT COUNT(*) FROM decision_type) as decision_types,
    (SELECT COUNT(*) FROM decision) as decisions,
    (SELECT COUNT(*) FROM decision_employees) as decision_employees,
    (SELECT COUNT(*) FROM allowance) as allowances,
    (SELECT COUNT(*) FROM contract_allowance) as contract_allowances,
    (SELECT COUNT(*) FROM insurance) as insurances,
    (SELECT COUNT(*) FROM work_shifts) as work_shifts,
    (SELECT COUNT(*) FROM employee_workshift) as employee_workshifts,
    (SELECT COUNT(*) FROM attendance_category) as attendance_categories,
    (SELECT COUNT(*) FROM attendance_record) as attendance_records,
    (SELECT COUNT(*) FROM work_schedule) as work_schedules,
    (SELECT COUNT(*) FROM work_schedule_shift) as work_schedule_shifts,
    (SELECT COUNT(*) FROM work_schedule_manager) as work_schedule_managers;

-- Show sample data from key tables
SELECT 'Sample Attendance Records (Recent):' as info;
SELECT ar.employee_id, e.full_name, ar.attendance_date, ar.check_in_time, ar.check_out_time, 
       ac.category_name, ar.late_minutes, ar.early_leave_minutes
FROM attendance_record ar
JOIN employee e ON ar.employee_id = e.employee_id
JOIN attendance_category ac ON ar.attendance_category_id = ac.id
ORDER BY ar.attendance_date DESC, ar.employee_id
LIMIT 10;

SELECT 'Sample Contracts:' as info;
SELECT c.employee_id, e.full_name, ct.contract_type_name, c.start_date, c.end_date, 
       c.basic_salary, c.status
FROM contract c
JOIN employee e ON c.employee_id = e.employee_id
JOIN contract_type ct ON c.contract_type_id = ct.id
ORDER BY c.start_date DESC
LIMIT 5;

SELECT 'Sample Work Schedules:' as info;
SELECT ws.schedule_name, ws.start_date, ws.end_date, ws.status, o.office_name
FROM work_schedule ws
JOIN office o ON ws.office_id = o.office_id
ORDER BY ws.created_date DESC
LIMIT 5;
