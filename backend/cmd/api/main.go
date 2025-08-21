package main

import (
	"erp/backend/api/middleware"
	"erp/backend/api/router"
	"erp/backend/config"
	"erp/backend/pkg/job"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	config.LoadConfig()
	config.LoadMailConfig()
	config.ConnectPostgres()
	config.LoadMinIOConfig()
	job.InitWorkerPool(10)

	db := config.GetDB()
	if db == nil {
		panic("Database connection is nil - failed to connect to PostgreSQL")
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		panic("Không thể lấy sql.DB từ GORM: " + err.Error())
	}
	defer sqlDB.Close()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(config.AppConfig))
	router.SetupRoutes(r, db)

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	// Luôn chạy HTTP
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}()

	certPath := "./ssl/server.crt"
	keyPath := "./ssl/server.key"

	if config.AppConfig.Server.SSL && fileExists(certPath) && fileExists(keyPath) {
		server := &http.Server{
			Addr:         ":8443",
			Handler:      r,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		log.Println("HTTPS Server running at :8443")
		if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
			log.Fatalf("Failed to run HTTPS server: %v", err)
		}
	} else {
		log.Println("⚠️ SSL disabled or certificate files not found — running HTTP only.")
		select {}
	}
	// Start HTTP server on port 8080
	go func() {
		log.Println("HTTP Server running at :8080")
		if err := r.Run(":8080"); err != nil {
			log.Printf("HTTP server failed: %v", err)
		}
	}()

	// Start HTTPS server on port 8443 if SSL is enabled
	if config.AppConfig.Server.SSL {
		log.Println("HTTPS Server running at :8443")
		err = r.RunTLS(":8443", "/app/ssl/server.crt", "/app/ssl/server.key")
		if err != nil {
			log.Fatalf("HTTPS server failed: %v", err)
		}
	} else {
		// If SSL is disabled, just run HTTP server and wait
		select {}
	}

}
