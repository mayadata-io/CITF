package system

import (
	"os"
)

// GetenvFallback search for the key in environment, if it is there then
// this function returns the value otherwise it returns the fallback value
func GetenvFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
