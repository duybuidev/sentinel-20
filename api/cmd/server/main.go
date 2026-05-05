package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/duybuidev/sentinel-20/api/internal/handler"
	"github.com/duybuidev/sentinel-20/api/internal/repository"
	"github.com/duybuidev/sentinel-20/api/internal/service"
	"github.com/duybuidev/sentinel-20/api/pkg/config"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("✅ PostgreSQL connected")

	// Repositories
	containerRepo := repository.NewContainerRepository(db)
	logRepo       := repository.NewLogRepository(db)
	eventRepo     := repository.NewEventRepository(db)

	// Services
	containerSvc := service.NewContainerService(containerRepo)
	logSvc       := service.NewLogService(logRepo)
	eventSvc     := service.NewEventService(eventRepo)

	// Handlers
	containerH := handler.NewContainerHandler(containerSvc)
	logH       := handler.NewLogHandler(logSvc)
	eventH     := handler.NewEventHandler(eventSvc)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "sentinel-api"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/containers",         containerH.GetAll)
		v1.GET("/containers/summary", containerH.GetSummary)
		v1.GET("/containers/:id",     containerH.GetByID)

		v1.GET("/logs",       logH.GetAll)
		v1.GET("/logs/stats", logH.GetStats)

		v1.GET("/events", eventH.GetRecent)
	}

	log.Printf("🚀 Sentinel API running on :%s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
