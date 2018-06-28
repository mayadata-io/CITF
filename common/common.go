package common

import (
	"os"
	"strings"
)

// DebugEnabled governs whether to print verbose logs or not
// It can be set by Environment Variable `CITF_VERBOSE_LOG``
var DebugEnabled bool

func init() {
	// Debug-Environment detection
	debugEnv := os.Getenv("CITF_VERBOSE_LOG")

	if strings.ToLower(debugEnv) == "true" {
		DebugEnabled = true
	} else {
		DebugEnabled = false
	}
}
