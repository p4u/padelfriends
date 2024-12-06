package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds configuration values for the application.
type Config struct {
	MongoURI string
	Port     int
}

// Load reads environment variables and returns a Config struct.
// Defaults to port 7777 if not specified.
func Load() (*Config, error) {
	port := 7777
	if p := os.Getenv("PORT"); p != "" {
		pp, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT: %v", err)
		}
		port = pp
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable not set")
	}

	cfg := &Config{
		MongoURI: mongoURI,
		Port:     port,
	}
	return cfg, nil
}
