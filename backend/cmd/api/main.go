package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!2222"})
	})

	log.Println("Server running at :8080")

	err := r.Run(":8080");
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
