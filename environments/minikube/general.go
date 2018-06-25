package minikube

import (
	"os"
	"strings"
	"time"

	sysutil "github.com/openebs/CITF/utils/system"
)

var (
	useSudo     = true // Default value to use sudo
	execCommand = sysutil.ExecCommandWithSudo
	runCommand  = sysutil.RunCommandWithSudo
)

// debug governs whether to print verbose logs or not
// It can be set by Environment Variable `CITF_VERBOSE_LOG``
var debug bool

func init() {
	// Debug-Environment detection
	debugEnv := os.Getenv("CITF_VERBOSE_LOG")

	if strings.ToLower(debugEnv) == "true" {
		debug = true
	} else {
		debug = false
	}

	// `sudo` use detection
	useSudoEnv := strings.ToLower(strings.TrimSpace(os.Getenv("USE_SUDO")))
	if useSudoEnv == "true" { // If it is mentioned in the environment variable to use sudo
		useSudo = true // use sudo then
	} else if useSudoEnv == "false" { // Else if it is mentioned in the environment variable not to use sudo
		useSudo = false // do not use sudo
	} // Else use default value mentioned above

	if !useSudo {
		execCommand = sysutil.ExecCommand
		runCommand = sysutil.RunCommand
	}
}

// Minikube is a struct which will be the driver for all the methods related to minikube
type Minikube struct {
	// Timeout is the timeout that will be used throughout the minikube package
	// for timeout in any operation if requires.
	Timeout time.Duration

	// WaitTimeUnit is the time duration, which will be used throughout package
	// if it needs to wait for some sub-task. (It is small timeout)
	WaitTimeUnit time.Duration
}

// NewMinikube returns a Minikube struct
func NewMinikube() Minikube {
	return Minikube{
		Timeout:      time.Minute,
		WaitTimeUnit: time.Second,
	}
}
