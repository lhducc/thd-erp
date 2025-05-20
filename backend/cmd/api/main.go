package main

import (
	"erp/backend/api/middleware"
	"erp/backend/api/router"
	"erp/backend/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	config.ConnectPostgres()

	db := config.GetDB()

	sqlDB, err := db.DB()
	if err != nil {
		panic("Không thể lấy sql.DB từ GORM: " + err.Error())
	}
	defer sqlDB.Close()
	r := gin.Default()
	r.Use(middleware.CORSMiddleware(config.AppConfig))
	router.SetupRoutes(r, db)

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!2222"})
	})

	log.Println("Server running at :8080")

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
