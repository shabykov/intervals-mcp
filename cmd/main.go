// intervals-icu MCP server (stdio) для Claude Desktop.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/shabykov/intervals-mcp/client"
	"github.com/shabykov/intervals-mcp/server"
)

const (
	envApiURL    = "INTERVALS_API_URL"
	envApiKey    = "INTERVALS_API_KEY"
	envAthleteID = "INTERVALS_ATHLETE_ID"

	defaultApiURL = "https://intervals.icu/api/v1"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("load .env: %v", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.NewServer(cfg).Run(ctx); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func loadConfig() (client.Config, error) {
	cfg := client.Config{
		ApiURL:    getEnv(envApiURL, defaultApiURL),
		ApiKey:    os.Getenv(envApiKey),
		AthleteID: os.Getenv(envAthleteID),
	}
	if cfg.ApiKey == "" {
		return cfg, fmt.Errorf("%s is required", envApiKey)
	}
	if cfg.AthleteID == "" {
		return cfg, fmt.Errorf("%s is required", envAthleteID)
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
