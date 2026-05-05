package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/duybuidev/sentinel-20/api/pkg/config"
)

func main() {
	cfg := config.Load()

	// Connect PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("✅ PostgreSQL connected")

	// Gin router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "sentinel-api"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/containers", func(c *gin.Context) {
			var containers []map[string]interface{}
			rows, err := db.Queryx("SELECT * FROM containers ORDER BY updated_at DESC")
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			for rows.Next() {
				row := make(map[string]interface{})
				rows.MapScan(row)
				containers = append(containers, row)
			}
			c.JSON(200, gin.H{"data": containers, "total": len(containers)})
		})

		v1.GET("/logs", func(c *gin.Context) {
			level   := c.DefaultQuery("level", "")
			service := c.DefaultQuery("service", "")
			limit   := c.DefaultQuery("limit", "100")

			query := "SELECT * FROM logs WHERE 1=1"
			if level   != "" { query += fmt.Sprintf(" AND level='%s'", level) }
			if service != "" { query += fmt.Sprintf(" AND service='%s'", service) }
			query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %s", limit)

			var logs []map[string]interface{}
			rows, err := db.Queryx(query)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			for rows.Next() {
				row := make(map[string]interface{})
				rows.MapScan(row)
				logs = append(logs, row)
			}
			c.JSON(200, gin.H{"data": logs, "total": len(logs)})
		})

		v1.GET("/events", func(c *gin.Context) {
			var events []map[string]interface{}
			rows, err := db.Queryx("SELECT * FROM events ORDER BY created_at DESC LIMIT 50")
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			for rows.Next() {
				row := make(map[string]interface{})
				rows.MapScan(row)
				events = append(events, row)
			}
			c.JSON(200, gin.H{"data": events})
		})
	}

	log.Printf("🚀 Sentinel API running on :%s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
