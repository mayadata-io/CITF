package docker

import (
	"fmt"
	"strings"
)

// TearDown removes all the docker containers present on the machine
func (docker Docker) Teardown() error {
	// CAUTION: This function call deletes all docker containers
	containersStr, err := execCommand("docker ps -aq")
	if err != nil {
		return err
	}
	if containersStr != "" {
		containers := strings.Fields(containersStr)
		for _, container := range containers {
			err = runCommand("docker rm -f " + container)
			if err != nil {
				fmt.Printf("Error occured while deleting docker container: %s. Error: %+v\n", container, err)
			} else {
				fmt.Printf("Deleted container: %s\n", container)
			}
		}
	}
	return nil
}
