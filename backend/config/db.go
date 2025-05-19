package config

import (
	hrmmodel "erp/backend/internal/hrm/hr_profile/model"
	// officemodel "erp/backend/internal/hrm/office/model"
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
	errT := AutoMigrate(db)
	if errT != nil {
		fmt.Print(errT)
	}

	fmt.Println("Đã kết nối PostgreSQL!")
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		// &officemodel.Office{},
		&hrmmodel.Contract{},
		&hrmmodel.Decision{},
		&hrmmodel.Employee{},
	)
	if err != nil {
		return fmt.Errorf("migrate thất bại: %w", err)
	}
	return nil
}
