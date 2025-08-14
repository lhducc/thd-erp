package config

import (
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
	if err := createEnums(db); err != nil {
		log.Fatalf("Không thể tạo ENUM: %v", err)
	}
	errT := AutoMigrate(db)
	if errT != nil {
		fmt.Print(errT)
	}

	if err := createDefaultRoles(db); err != nil {
		log.Fatalf("Không thể tạo role mặc định: %v", err)
	}

	fmt.Println("Đã kết nối PostgreSQL!")
	DB = db

	db.Exec("DISCARD ALL")

}

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
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'condition_enum') THEN
        CREATE TYPE condition_enum AS ENUM (
            'Chưa hiệu lực',
            'Đang hiệu lực',
            'Hết hiệu lực',
            'Thanh lý'
        );
    	END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'document_group_enum') THEN
		  CREATE TYPE document_group_enum AS ENUM (
			 'Loại chứng chỉ',
			 'Loại lao động',
			 'Thủ tục tiếp nhận',
			 'Thủ tục thôi việc'
		  );
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'work_day_enum') THEN
		  CREATE TYPE work_group_enum AS ENUM (
			 '1',
			 '0.5',
			 '0'
		  );
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'repeat_type_enum') THEN
			CREATE TYPE repeat_type_enum AS ENUM (
				'weekly',
				'monthly'
			);
    	END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_work_schedule_enum') THEN
			CREATE TYPE status_work_schedule_enum AS ENUM (
				'expired',
				'inactive',
				'active'
			);
		END IF;
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

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Role{},
		&model.Office{},
		&checkin_model.EmployeeWorkshift{},

		&model.Position{},
		&model.Department{},
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
	)
	fmt.Println("Migration complete")

	if err != nil {
		return fmt.Errorf("migrate thất bại: %w", err)
	}

	return nil
}

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
