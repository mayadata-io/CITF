package docker

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// Teardown stops all the docker containers present on the machine
func (docker Docker) Teardown() error {
	// CAUTION: This function call stops all docker containers
	containersStr, err := execCommand("docker ps -q")
	if err != nil {
		return fmt.Errorf("error while getting container id. Error: %+v", err)
	}
	if containersStr != "" {
		containers := strings.Fields(containersStr)
		for _, container := range containers {
			err = runCommand("docker stop -f " + container)
			if err != nil {
				glog.Errorf("error occured while stopping docker container: %s. Error: %+v\n", container, err)
			} else {
				fmt.Printf("Stopped container: %s\n", container)
			}
		}
	}
	return nil
}
