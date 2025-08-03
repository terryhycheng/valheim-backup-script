package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	variables = []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"S3_BASE_ENDPOINT",
	}
	Envs = map[string]string{}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Error loading .env file: %v", err)
	}

	missingVariables := []string{}

	for _, variable := range variables {
		if os.Getenv(variable) == "" {
			missingVariables = append(missingVariables, variable)
		}
		Envs[variable] = os.Getenv(variable)
	}

	if len(missingVariables) > 0 {
		log.Fatalf("❌ Missing required environment variables: %v", strings.Join(missingVariables, ", "))
	}
}
