package config

import (
	// checkin_model "erp/backend/internal/hrm/checkin/model"
	checkin_model "erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/hr_profile/model"
	"fmt"
	"log"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := AppConfig.Postgres.DBSource

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối PostgreSQL:", err)
	}

	// Create ENUM types first before AutoMigrate
	if err := createEnums(db); err != nil {
		log.Fatal("Không thể tạo enum types:", err)
	}

	// Run AutoMigrate to create tables
	if err := AutoMigrate(db); err != nil {
		log.Fatal("Không thể tự động migrate:", err)
	}

	// Add foreign key constraints manually after migration
	if err := addForeignKeyConstraints(db); err != nil {
		log.Fatal("Không thể tạo foreign key constraints:", err)
	}

	// Create default roles after tables are created
	if err := createDefaultRoles(db); err != nil {
		log.Fatal("Không thể tạo vai trò mặc định:", err)
	}

	fmt.Println("Đã kết nối PostgreSQL!")
	DB = db

	db.Exec("DISCARD ALL")
}

// Create ENUM types for AutoMigrate
func createEnums(db *gorm.DB) error {
	enumSQL := `
	DO $$
	BEGIN
		-- enum contract_group_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'contract_group_enum') THEN
			CREATE TYPE contract_group_enum AS ENUM (
				'Hợp đồng xác định thời hạn',
				'Hợp đồng không xác định thời hạn',
				'Hợp đồng thử việc',
				'Hợp đồng đào tạo nghề',
				'Hợp đồng dịch vụ'
			);
		END IF;

		-- enum unit_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'unit_enum') THEN
			CREATE TYPE unit_enum AS ENUM ('Năm', 'Tháng', 'Tuần', 'Ngày');
		END IF;

		-- enum working_type_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'working_type_enum') THEN
			CREATE TYPE working_type_enum AS ENUM (
				'Toàn thời gian',
				'Bán thời gian',
				'Cộng tác viên',
				'Chuyên gia',
				'Theo ca',
				'Khoán sản phẩm',
				'Khoán công việc'
			);
		END IF;

		-- enum decision_status_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'decision_status_enum') THEN
			CREATE TYPE decision_status_enum AS ENUM (
				'Đã duyệt',
				'Không duyệt',
				'Chờ duyệt'
			);
		END IF;

		-- enum decision_condition_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'decision_condition_enum') THEN
			CREATE TYPE decision_condition_enum AS ENUM (
				'Chưa hiệu lực',
				'Đang hiệu lực'
			);
		END IF;

		-- enum decision_group_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'decision_group_enum') THEN
			CREATE TYPE decision_group_enum AS ENUM (
				'Hình thức khen thưởng',
				'Hình thức kỷ luật',
				'Lý do điều chuyển',
				'Lý do tiếp nhận',
				'Lý do bổ nhiệm',
				'Lý do miễn nhiệm',
				'Lý do chấm dứt HĐLĐ'
			);
		END IF;

		-- enum document_condition_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_condition_enum') THEN
			CREATE TYPE document_condition_enum AS ENUM (
				'Hết hạn',
				'Đang hiệu lực'
			);
		END IF;

		-- enum document_status_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_status_enum') THEN
			CREATE TYPE document_status_enum AS ENUM (
				'Duyệt',
				'Không duyệt',
				'Chờ duyệt'
			);
		END IF;

		-- enum approve_status_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'approve_status_enum') THEN
			CREATE TYPE approve_status_enum AS ENUM (
				'Đã duyệt',
				'Không duyệt',
				'Chờ duyệt'
			);
		END IF;

		-- enum condition_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'condition_enum') THEN
			CREATE TYPE condition_enum AS ENUM (
				'Chưa hiệu lực',
				'Đang hiệu lực',
				'Hết hiệu lực',
				'Thanh lý'
			);
		END IF;

		-- enum document_group_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_group_enum') THEN
			CREATE TYPE document_group_enum AS ENUM (
				'Loại chứng chỉ',
				'Loại lao động',
				'Thủ tục tiếp nhận',
				'Thủ tục thôi việc'
			);
		END IF;

		-- enum work_day_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'work_day_enum') THEN
			CREATE TYPE work_day_enum AS ENUM (
				'1',
				'0.5',
				'0'
			);
		END IF;

		-- enum repeat_type_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'repeat_type_enum') THEN
			CREATE TYPE repeat_type_enum AS ENUM (
				'daily',
				'weekly',
				'monthly',
				'none'
			);
		END IF;

		-- enum status_work_schedule_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_work_schedule_enum') THEN
			CREATE TYPE status_work_schedule_enum AS ENUM (
				'expired',
				'inactive',
				'active',
				'pending'
			);
		END IF;

		-- enum weekday_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'weekday_enum') THEN
			CREATE TYPE weekday_enum AS ENUM (
				'sunday',
				'monday',
				'tuesday',
				'wednesday',
				'thursday',
				'friday',
				'saturday'
			);
		END IF;

		-- enum work_type_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'work_type_enum') THEN
			CREATE TYPE work_type_enum AS ENUM (
				'ca hành chính',
				'ca kíp'
			);
		END IF;
	END
	$$;
	`
	return db.Exec(enumSQL).Error
}

// AutoMigrate creates all database tables
func AutoMigrate(db *gorm.DB) error {
	// Migrate models in dependency order to avoid foreign key issues
	err := db.AutoMigrate(
		// Basic lookup tables first (no dependencies)
		&model.Role{},
		&model.Office{},
		&model.Position{},
		&model.HierarchyLevel{},
		&model.EmployeeDocumentType{},
		&model.ContractType{},
		// DecisionType has no dependencies
		&model.DecisionType{},
		&model.Insurance{},
		&model.Allowance{},

		// Department depends on Office
		&model.Department{},

		// JobTitle depends on HierarchyLevel
		&model.JobTitle{},

		// Employee depends on Position, JobTitle, Department
		&model.Employee{},

		// Account depends on Employee and Role
		&model.Account{},

		// Checkin models
		&checkin_model.WorkShifts{},
		&checkin_model.EmployeeWorkshift{},
		&checkin_model.AttendanceCategory{},
		&checkin_model.AttendanceRecord{},

		// Contract models (after Employee)
		&model.Contract{},
		&model.ContractAllowance{},
		&model.EmployeeDocument{},

		// Decision depends on DecisionType and Employee (both already migrated)
		&model.Decision{},
		&model.DecisionEmployee{},

		// Work schedule models
		&checkin_model.WorkSchedule{},
		&checkin_model.WorkScheduleShift{},
		&checkin_model.WorkScheduleManager{},
		&checkin_model.TimeSheetList{},
		&checkin_model.TimeSheet{},
		&checkin_model.TimeSheetDetail{},
	)

	fmt.Println("Migration complete")

	if err != nil {
		return fmt.Errorf("migrate thất bại: %w", err)
	}

	return nil
}

// Add foreign key constraints manually after migration
func addForeignKeyConstraints(db *gorm.DB) error {
	// Add foreign key constraint from decision.decision_type_id to decisiontype.decision_type_id
	constraintSQL := `
	DO $$
	BEGIN
		-- Add foreign key constraint for decision.decision_type_id -> decisiontype.decision_type_id
		IF NOT EXISTS (
			SELECT 1 FROM information_schema.table_constraints 
			WHERE constraint_name = 'fk_decision_decision_type_id'
		) THEN
			ALTER TABLE decision 
			ADD CONSTRAINT fk_decision_decision_type_id 
			FOREIGN KEY (decision_type_id) REFERENCES decisiontype(decision_type_id) 
			ON UPDATE CASCADE ON DELETE SET NULL;
		END IF;
	END
	$$;
	`
	return db.Exec(constraintSQL).Error
}

func GetDB() *gorm.DB {
	return DB
}

// Create default roles after migration
func createDefaultRoles(db *gorm.DB) error {
	defaultRoles := []string{"admin", "manager", "employee"}

	for _, name := range defaultRoles {
		var count int64
		if err := db.Model(&model.Role{}).Where("role_name = ?", name).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			role := model.Role{
				ID:          name,
				RoleName:    name,
				CreatedDate: time.Now(),
			}
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
