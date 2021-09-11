package config

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// Config contains various configs need to run the server
type Config struct {
	Port string
}

// LoadConfig loads configs from environment variables
func LoadConfig() *Config {
	port := getPortConfig(os.Getenv(PortEnvironmentVariableKey))
	return &Config{
		Port: port,
	}
}

// getPortConfig fetch and returns port config from env variable
func getPortConfig(port string) string {
	if port != "" {
		return port
	}
	log.Info("The environment variable PORT is missing.")
	log.Info("Defaulting to 8000")
	return DefaultPort
}

const (
	PortEnvironmentVariableKey = "PORT"
	DefaultPort                = "8000"
)
