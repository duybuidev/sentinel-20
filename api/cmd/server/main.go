package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/duybuidev/sentinel-20/api/internal/repository"
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

	containerRepo := repository.NewContainerRepository(db)
	logRepo       := repository.NewLogRepository(db)
	eventRepo     := repository.NewEventRepository(db)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "sentinel-api"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/containers", func(c *gin.Context) {
			containers, err := containerRepo.FindAll()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"data": containers, "total": len(containers)})
		})

		v1.GET("/logs", func(c *gin.Context) {
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
			logs, err := logRepo.FindAll(repository.LogFilter{
				Level:   c.Query("level"),
				Service: c.Query("service"),
				Limit:   limit,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"data": logs, "total": len(logs)})
		})

		v1.GET("/logs/stats", func(c *gin.Context) {
			stats, err := logRepo.CountByLevel()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"data": stats})
		})

		v1.GET("/events", func(c *gin.Context) {
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
			events, err := eventRepo.FindRecent(limit)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"data": events, "total": len(events)})
		})
	}

	log.Printf("🚀 Sentinel API running on :%s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
