package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirovia/bardista/internal/config"
)

func main() {
	cfg := config.Load()

	r := gin.Default()

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
