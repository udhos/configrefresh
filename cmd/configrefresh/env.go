package main

import (
	"log"
	"os"
	"time"
)

// String extracts string from env var.
// It returns the provided defaultValue if the env var is empty.
// The string returned is also recorded in logs.
func envString(name string, defaultValue string) string {
	str := os.Getenv(name)
	if str != "" {
		log.Printf("%s=[%s] using %s=%s default=%s", name, str, name, str, defaultValue)
		return str
	}
	log.Printf("%s=[%s] using %s=%s default=%s", name, str, name, defaultValue, defaultValue)
	return defaultValue
}

// Duration extracts time.Duration value from env var.
// It returns the provided defaultValue if the env var is empty.
// The value returned is also recorded in logs.
func envDuration(name string, defaultValue time.Duration) time.Duration {
	str := os.Getenv(name)
	if str != "" {
		value, errConv := time.ParseDuration(str)
		if errConv == nil {
			log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, value, defaultValue)
			return value
		}
		log.Printf("bad %s=[%s]: error: %v", name, str, errConv)
	}
	log.Printf("%s=[%s] using %s=%v default=%v", name, str, name, defaultValue, defaultValue)
	return defaultValue
}
