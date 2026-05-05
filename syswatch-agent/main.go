package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

var apiBaseURL = getEnv("SENTINEL_API_URL", "http://sentinel-api:8080")

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ── HTTP helpers ────────────────────────────────────────────

func postJSON(path string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(apiBaseURL+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// ── Sync container state vào API ────────────────────────────

func syncContainer(cli *client.Client, containerID string) {
	info, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return
	}

	stats, err := cli.ContainerStatsOneShot(context.Background(), containerID)
	if err != nil {
		return
	}
	defer stats.Body.Close()

	var s types.StatsJSON
	json.NewDecoder(stats.Body).Decode(&s)

	// Tính CPU %
	cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage - s.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta  := float64(s.CPUStats.SystemUsage - s.PreCPUStats.SystemUsage)
	cpuPercent := 0.0
	if sysDelta > 0 {
		cpuPercent = (cpuDelta / sysDelta) * float64(len(s.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	payload := map[string]interface{}{
		"name":        info.Name[1:], // bỏ dấu /
		"image":       info.Config.Image,
		"status":      info.State.Status,
		"cpu_percent": cpuPercent,
		"mem_usage":   s.MemoryStats.Usage,
		"mem_limit":   s.MemoryStats.Limit,
	}

	if err := postJSON("/api/v1/internal/containers/sync", payload); err != nil {
		log.Printf("❌ Sync container %s failed: %v", info.Name, err)
	}
}

// ── Gửi event vào API ───────────────────────────────────────

func sendEvent(containerName, eventType string, exitCode int) {
	payload := map[string]interface{}{
		"container_name": containerName,
		"type":           eventType,
		"exit_code":      exitCode,
	}
	if err := postJSON("/api/v1/internal/events", payload); err != nil {
		log.Printf("❌ Send event failed: %v", err)
	}
}

// ── Điểm danh khi khởi động ─────────────────────────────────

func bootstrap(cli *client.Client) {
	fmt.Println("🔍 Bootstrapping — scanning all containers...")
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Printf("❌ Cannot list containers: %v", err)
		return
	}
	for _, c := range containers {
		syncContainer(cli, c.ID)
		fmt.Printf("   ✅ Synced: %s [%s]\n", c.Names[0][1:], c.Status)
	}
	fmt.Printf("🚀 Bootstrap done — %d containers synced\n", len(containers))
	fmt.Println("─────────────────────────────────────────────")
}

// ── Periodic sync mỗi 30s ───────────────────────────────────

func startPeriodicSync(cli *client.Client) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
			if err != nil {
				continue
			}
			for _, c := range containers {
				syncContainer(cli, c.ID)
			}
			fmt.Printf("🔄 [%s] Periodic sync — %d containers\n",
				time.Now().Format("15:04:05"), len(containers))
		}
	}()
}

// ── Main ────────────────────────────────────────────────────

func main() {
	fmt.Println("🛡️  SysWatch Agent starting...")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("❌ Cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	bootstrap(cli)
	startPeriodicSync(cli)

	fmt.Println("📡 Listening for Docker events...")
	ctx := context.Background()
	msgs, errs := cli.Events(ctx, types.EventsOptions{})

	for {
		select {
		case err := <-errs:
			log.Printf("⚠️  Docker event error: %v\n", err)

		case msg := <-msgs:
			if msg.Type != events.ContainerEventType {
				continue
			}

			name := msg.Actor.Attributes["name"]
			ts   := time.Now().Format("15:04:05")

			switch msg.Action {
			case "start":
				fmt.Printf("🟢 [%s] %s STARTED\n", ts, name)
				syncContainer(cli, msg.Actor.ID)
				sendEvent(name, "started", 0)

			case "die":
				exitCode := 0
				fmt.Sscanf(msg.Actor.Attributes["exitCode"], "%d", &exitCode)
				fmt.Printf("🚨 [%s] %s CRASHED (exit: %d)\n", ts, name, exitCode)
				sendEvent(name, "crashed", exitCode)

			case "stop":
				fmt.Printf("🟡 [%s] %s STOPPED\n", ts, name)
				sendEvent(name, "stopped", 0)

			case "restart":
				fmt.Printf("🔁 [%s] %s RESTARTED\n", ts, name)
				syncContainer(cli, msg.Actor.ID)
				sendEvent(name, "restarted", 0)
			}
		}
	}
}
