package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	// Sử dụng trực tiếp chuỗi DBSource
	dsn := AppConfig.Postgres.DBSource

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối PostgreSQL:", err)
	}
	// Gọi hàm tạo enum trước khi migrate
	if err := createEnums(db); err != nil {
		log.Fatalf("Không thể tạo ENUM: %v", err)
	}
	errT := AutoMigrate(db)
	if errT != nil {
		fmt.Print(errT)
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
		//&officemodel.Office{},
		//&hrmmodel.EmployeeDocumentType{},
		//&hrmmodel.Employee{},
		//&hrmmodel.ContractType{},
		//&hrmmodel.Contract{},
		//&hrmmodel.DecisionType{},
		//&hrmmodel.Decision{},
		//&hrmmodel.DecisionEmployee{},
		//&hrmmodel.Insurance{},
		// &model.Holiday{},
		// &model.AllowedWorkingSchedule{},
		// &model.WorkShifts{},
		// &hrmmodel.Allowance{}
	)
	if err != nil {
		return fmt.Errorf("migrate thất bại: %w", err)
	}
	return nil
}
