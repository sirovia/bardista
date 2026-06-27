package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirovia/bardista/internal/config"
	"github.com/sirovia/bardista/internal/handler"
	"github.com/sirovia/bardista/internal/repository"
	"github.com/sirovia/bardista/internal/router"
	"github.com/sirovia/bardista/internal/service"
)

func runMigrations(db *sql.DB) {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal("Could not find migration files:", err)
	}

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			log.Fatal("Could not read migration file:", err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			log.Fatal("Could not run migration:", err)
		}

		log.Printf("Ran migration: %s", f)
	}
}

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Could not ping database:", err)
	}
	log.Println("Connected to database")

	runMigrations(db)

	// repositories
	userRepo := repository.NewUserRepository(db)

	// services
	expiryHours, _ := strconv.Atoi(cfg.JWTExpiryHours)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, time.Duration(expiryHours)*time.Hour)

	//handlers
	authHandler := handler.NewAuthHandler(authService)

	// router
	r := gin.Default()

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.Setup(r, authHandler, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
