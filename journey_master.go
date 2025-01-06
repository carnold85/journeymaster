package main

import (
	api "journeymaster/core"
	"log"
	"os"
)

func main() {
	// Default values
	env_defaults := map[string]string{
		"API_PORT":        "8080",
		"API_PREFIX":      "/v1",
		"API_CREDENTIALS": "test:test",
		"GIN_MODE":        "release",
		"TLS_PEM":         "conf/cert.pem",
		"TLS_KEY":         "conf/key.pem",
		"TLS_ENABLED":     "false",
	}

	// Check and set environment variables
	for key, defaultValue := range env_defaults {
		if os.Getenv(key) == "" {
			os.Setenv(key, defaultValue)
			log.Printf("Environment variable %s not set. Defaulting to: %s\n", key, defaultValue)
		} else {
			log.Printf("Environment variable %s is already set to: %s\n", key, os.Getenv(key))
		}
	}

	server, _ := api.NewJourneyMasterAPI()

	server.StartServer()
}
