package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	serviceName    = getEnv("SERVICE_NAME", "mock-service")
	servicePort    = getEnv("SERVICE_PORT", "9000")
	sentinelAPIURL = getEnv("SENTINEL_API_URL", "http://sentinel-api:8080")
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var logMessages = map[string][]string{
	"INFO": {
		"Request processed successfully",
		"User authenticated",
		"Cache hit for key",
		"Database query completed in 12ms",
		"Health check passed",
		"Connection pool size: 10",
		"Scheduled job completed",
	},
	"WARNING": {
		"Response time exceeded 500ms",
		"Cache miss, falling back to DB",
		"Retry attempt 2/3",
		"Memory usage at 75%",
		"Connection pool near limit",
		"Rate limit threshold approaching",
	},
	"ERROR": {
		"Database connection timeout",
		"Failed to parse request body",
		"External API returned 503",
		"Queue consumer lag detected",
		"Authentication token expired",
	},
	"CRITICAL": {
		"Out of memory",
		"Disk space critically low",
		"Database unreachable",
		"Service mesh communication failure",
	},
}

// Tỉ lệ sinh log theo level
var levelWeights = []struct {
	level  string
	weight int
}{
	{"INFO", 70},
	{"WARNING", 20},
	{"ERROR", 8},
	{"CRITICAL", 2},
}

func randomLevel() string {
	n := rand.Intn(100)
	cumulative := 0
	for _, lw := range levelWeights {
		cumulative += lw.weight
		if n < cumulative {
			return lw.level
		}
	}
	return "INFO"
}

func sendLog(level, message string) {
	payload := map[string]interface{}{
		"level":   level,
		"message": message,
		"service": serviceName,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(sentinelAPIURL+"/api/v1/internal/logs", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func generateLogs() {
	ticker := time.NewTicker(time.Duration(rand.Intn(3)+1) * time.Second)
	for range ticker.C {
		level := randomLevel()
		messages := logMessages[level]
		message := messages[rand.Intn(len(messages))]
		fullMessage := fmt.Sprintf("[%s] %s", serviceName, message)
		sendLog(level, fullMessage)
		log.Printf("[%s] %s: %s", level, serviceName, message)
	}
}

func main() {
	log.Printf("🚀 %s starting on port %s", serviceName, servicePort)

	go generateLogs()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"%s"}`, serviceName)
	})

	http.HandleFunc("/crash", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[FATAL] %s received crash signal", serviceName)
		os.Exit(1)
	})

	if err := http.ListenAndServe(":"+servicePort, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
