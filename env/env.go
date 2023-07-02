package env

import "os"

const (
	Development = "development"
	Staging     = "staging"
	Production  = "production"
)

// Environment returns running environment
func Environment() string {
	v, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return Development
	}

	return v
}

// EVString reads environment variable
func EVString(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
