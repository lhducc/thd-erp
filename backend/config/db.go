package config

import (
	// checkin_model "erp/backend/internal/hrm/checkin/model"
	checkin_model "erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/hr_profile/model"
	"fmt"
	"log"
	"time"
	"gorm.io/gorm/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var AllModels = []interface{}{
	&model.Office{},
	&checkin_model.EmployeeWorkshift{},
	&model.Position{},
	&model.Department{},
	&model.Office{},
	&model.JobTitle{},
	&checkin_model.WorkShifts{},
	&checkin_model.EmployeeWorkshift{},
	&model.EmployeeDocumentType{},
	&model.Employee{},
	&model.ContractType{},
	&model.Contract{},
	&model.DecisionType{},
	&model.Decision{},
	&model.DecisionEmployee{},
	&model.Insurance{},
	&checkin_model.WorkShifts{},
	&model.Allowance{},
	&model.Contract{},
	&model.ContractAllowance{},
	&checkin_model.AttendanceCategory{},
	&checkin_model.AttendanceRecord{},
	&checkin_model.WorkSchedule{},
	&checkin_model.WorkScheduleShift{},
	&checkin_model.WorkScheduleManager{},
	&checkin_model.TimeSheetList{},
	&checkin_model.TimeSheet{},
	&checkin_model.TimeSheetDetail{},
}

func ConnectPostgres() {
	dsn := AppConfig.Postgres.DBSource

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Không thể kết nối PostgreSQL:", err)
	}

	// Create ENUM types first before AutoMigrate
	if err := createEnums(db); err != nil {
		log.Fatal("Không thể tạo enum types:", err)
	}
	
	errT := db.AutoMigrate(AllModels...)
	if errT != nil {
		fmt.Print(errT)
	}

	// Run AutoMigrate to create tables
	if err := AutoMigrate(db); err != nil {
		log.Fatal("Không thể tự động migrate:", err)
	}
	CreateAllContraints(db)

	// Create default roles after tables are created
	if err := createDefaultRoles(db); err != nil {
		log.Fatal("Không thể tạo vai trò mặc định:", err)
	}

	if err := CreateForeignKeysFromModels(db, AllModels); err != nil {
		log.Fatalf("Không tạo được FK cho bảng: %v", err)
	}

	fmt.Println("Đã kết nối PostgreSQL!")
	DB = db

	db.Exec("DISCARD ALL")
}

// Create ENUM types for AutoMigrate
func createEnums(db *gorm.DB) error {
	extensionSQL := `
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    `
	if err := db.Exec(extensionSQL).Error; err != nil {
		return fmt.Errorf("không thể tạo extension uuid-ossp: %w", err)
	}

	enumSQL := `
	DO $$
	BEGIN
		-- enum contract_group_enum
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'contract_group_enum') THEN
			CREATE TYPE contract_group_enum AS ENUM (
				'Hợp đồng thử việc',
				'Hợp đồng chính thức'
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

func GetDB() *gorm.DB {
	return DB
}

func AutoMigrateModels(db *gorm.DB, models []interface{}) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("migrate models lỗi: %w", err)
	}
	return nil
}

// AutoMigrate creates all database tables
func AutoMigrate(db *gorm.DB) error {
	// First pass - create tables without relationships to avoid circular dependencies
	err := db.Set("gorm:auto_preload", false).AutoMigrate(
		// Basic lookup tables first (no dependencies)
		&model.Role{},
		&model.Office{},
		&model.Position{},
		&model.HierarchyLevel{},
		&model.EmployeeDocumentType{},
		&model.ContractType{},
		&model.DecisionType{},
		&model.Insurance{},
		&model.Allowance{},
		&model.Department{},
		&model.JobTitle{},
		&model.Contract{},
		&model.Employee{},
		&model.Account{},
		&model.Decision{},
		&model.DecisionEmployee{},

		// Checkin models
		&checkin_model.WorkShifts{},
		&checkin_model.EmployeeWorkshift{},
		&checkin_model.AttendanceCategory{},
		&checkin_model.AttendanceRecord{},

		// Contract models
		&model.ContractAllowance{},
		&model.EmployeeDocument{},

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

func CreateForeignKeysFromModels(db *gorm.DB, models []interface{}) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, m := range models {
		stmt := &gorm.Statement{DB: tx}
		if err := stmt.Parse(m); err != nil {
			tx.Rollback()
			return fmt.Errorf("parse model lỗi: %w", err)
		}
		tableName := stmt.Schema.Table

		for _, rel := range stmt.Schema.Relationships.Relations {
			if rel.Type == schema.BelongsTo || rel.Type == schema.HasOne || rel.Type == schema.HasMany {
				for _, ref := range rel.References {
					// Bỏ qua tự tham chiếu
					if tableName == ref.PrimaryKey.Schema.Table {
						continue
					}

					// Kiểm tra cột tham chiếu có phải là PK hoặc Unique
					isReferenceValid := false
					for _, field := range ref.PrimaryKey.Schema.Fields {
						if field.DBName == ref.PrimaryKey.DBName && (field.PrimaryKey || field.Unique) {
							isReferenceValid = true
							break
						}
					}
					if !isReferenceValid {
						log.Printf("Info: Bỏ qua FK từ %s.%s vì %s.%s không phải PK/Unique",
							tableName, ref.ForeignKey.DBName,
							ref.PrimaryKey.Schema.Table, ref.PrimaryKey.DBName)
						continue
					}

					constraintName := fmt.Sprintf("fk_%s_%s", tableName, ref.ForeignKey.DBName)

					// Kiểm tra constraint đã tồn tại (phiên bản tối ưu cho single schema)
					var constraintExists bool
					checkSQL := `
                        SELECT EXISTS (
                            SELECT 1 FROM pg_constraint 
                            JOIN pg_class ON conrelid = pg_class.oid
                            WHERE conname = $1 AND pg_class.relname = $2
                        )`
					if err := tx.Raw(checkSQL, constraintName, tableName).Scan(&constraintExists).Error; err != nil {
						log.Printf("Warning: Không thể kiểm tra constraint %s: %v", constraintName, err)
						continue
					}

					if constraintExists {
						continue
					}

					// Tạo foreign key constraint
					sql := fmt.Sprintf(`
                        ALTER TABLE %s
                        ADD CONSTRAINT %s FOREIGN KEY (%s)
                        REFERENCES %s (%s)
                        ON DELETE CASCADE ON UPDATE CASCADE;
                    `,
						tableName,
						constraintName,
						ref.ForeignKey.DBName,
						ref.PrimaryKey.Schema.Table,
						ref.PrimaryKey.DBName,
					)

					if err := tx.Exec(sql).Error; err != nil {
						log.Printf("Error: Không thể tạo FK %s: %v", constraintName, err)
						continue
					}
				}
			}
		}
	}

	return tx.Commit().Error
}

func CreateAllContraints(db *gorm.DB) {
	db.Migrator().CreateConstraint(&model.Role{}, "Accounts")
	db.Migrator().CreateConstraint(&model.HierarchyLevel{}, "JobTitles")
	db.Migrator().CreateConstraint(&model.Department{}, "Office")
	db.Migrator().CreateConstraint(&model.
		JobTitle{}, "HierarchyLevel")
	db.Migrator().CreateConstraint(&model.
		Contract{}, "ContractType")
	db.Migrator().CreateConstraint(&model.
		Contract{}, "Employee")
	db.Migrator().CreateConstraint(&model.
		Contract{}, "Allowances")
	db.Migrator().CreateConstraint(&model.
		Employee{}, "Position")
	db.Migrator().CreateConstraint(&model.
		Employee{}, "JobTitle")
	db.Migrator().CreateConstraint(&model.
		Employee{}, "Manager")
	db.Migrator().CreateConstraint(&model.
		Employee{}, "Department")
	db.Migrator().CreateConstraint(&model.
		Employee{}, "Contracts")
	db.Migrator().CreateConstraint(&model.
		Employee{}, "Decisions")
	db.Migrator().CreateConstraint(&model.
		Account{}, "Role")
	db.Migrator().CreateConstraint(&model.
		Account{}, "Employee")
	db.Migrator().CreateConstraint(&model.
		Decision{}, "Employees")
	db.Migrator().CreateConstraint(&model.
		Decision{}, "DecisionType")
	db.Migrator().CreateConstraint(&checkin_model.
		WorkShifts{}, "Creator")
	db.Migrator().CreateConstraint(&checkin_model.
		EmployeeWorkshift{}, "WorkShift")
	db.Migrator().CreateConstraint(&checkin_model.
		AttendanceCategory{}, "Office")
	db.Migrator().CreateConstraint(&checkin_model.
		AttendanceRecord{}, "Employee")
	db.Migrator().CreateConstraint(&checkin_model.
		AttendanceRecord{}, "Office")
	db.Migrator().CreateConstraint(&checkin_model.
		AttendanceRecord{}, "CreateByInfo")
	db.Migrator().CreateConstraint(&checkin_model.
		AttendanceRecord{}, "AttendanceCategory")
	db.Migrator().CreateConstraint(&model.
		EmployeeDocument{}, "DocumentType")
	db.Migrator().CreateConstraint(&model.
		EmployeeDocument{}, "Employee")
	db.Migrator().CreateConstraint(&checkin_model.
		WorkSchedule{}, "Managers")
	db.Migrator().CreateConstraint(&checkin_model.
		WorkSchedule{}, "Weekdays")
	db.Migrator().CreateConstraint(&checkin_model.
		WorkSchedule{}, "Office")
	db.Migrator().CreateConstraint(&checkin_model.
		WorkScheduleShift{}, "WorkShift")
	db.Migrator().CreateConstraint(&checkin_model.
		WorkScheduleManager{}, "Employee")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetList{}, "Office")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetList{}, "Timesheets")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetList{}, "Creator")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetList{}, "Updater")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetList{}, "LockedUser")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheet{}, "Employee")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheet{}, "Office")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheet{}, "Department")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheet{}, "Details")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheet{}, "Creator")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheet{}, "Updater")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetDetail{}, "WorkShift")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetDetail{}, "CheckInRecord")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetDetail{}, "CheckOutRecord")
	db.Migrator().CreateConstraint(&checkin_model.
		TimeSheetDetail{}, "AdjustmentUser")

}
