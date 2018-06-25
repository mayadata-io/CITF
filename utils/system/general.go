package system

import (
	"os"
	"strings"
)

// debug governs whether to print verbose logs or not
// It can be set by Environment Variable `CITF_VERBOSE_LOG``
var debug bool

func init() {
	debugEnv := os.Getenv("CITF_VERBOSE_LOG")

	if strings.ToLower(debugEnv) == "true" {
		debug = true
	} else {
		debug = false
	}
}

// GetenvFallback search for the key in environment, if it is there then
// this function returns the value otherwise it returns the fallback value
func GetenvFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
