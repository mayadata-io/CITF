package minikube

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

// runPostStartCommandsForMinikube runs the commands required when run minikube as --vm-driver=none
// Assumption: Environment variables `USER` and `HOME` is well defined.
func (minikube Minikube) runPostStartCommandsForMinikubeNoneDriver() {
	userName := os.Getenv("USER")
	homeDir := os.Getenv("HOME")
	commands := []string{
		"mv /root/.kube " + homeDir + "/.kube",
		"chown -R " + userName + " " + homeDir + "/.kube",
		"chgrp -R " + userName + " " + homeDir + "/.kube",
		"mv /root/.minikube " + homeDir + "/.minikube",
		"chown -R " + userName + " " + homeDir + "/.minikube",
		"chgrp -R " + userName + " " + homeDir + "/.minikube",
	}

	for _, command := range commands {
		fmt.Printf("Running %q\n", command)
		output, err := execCommand(command)
		if err != nil {
			fmt.Printf("Running %q failed. Error: %+v\n", command, err)
		} else {
			fmt.Printf("Run %q successfully. Output: %s\n", command, output)
		}
	}
}

// StartMinikube method starts minikube with `--vm-driver=none` option.
func (minikube Minikube) StartMinikube() {
	err := runCommand("minikube start --vm-driver=none")
	// We can also use following:
	// "minikube start --vm-driver=none --feature-gates=MountPropagation=true --cpus=1 --memory=1024 --v=3 --alsologtostderr"
	if err != nil {
		glog.Fatal(err)
	}

	envChangeMinikubeNoneUser := os.Getenv("CHANGE_MINIKUBE_NONE_USER")
	if debug {
		fmt.Printf("Environ CHANGE_MINIKUBE_NONE_USER = %q\n", envChangeMinikubeNoneUser)
	}
	if envChangeMinikubeNoneUser == "true" {
		// Below commands shall automatically run in this case.
		if debug {
			fmt.Println("Returning from setup.")
		}
		return
	}

	minikube.waitForDotKubeDirToBeCreated()

	minikube.waitForDotMinikubeDirToBeCreated()

	minikube.runPostStartCommandsForMinikubeNoneDriver()
}

// Setup checks if a teardown is required before minikube start
// if so it does that and then start the minikube.
// It does nothing when minikube is already running.
// it prints status too.
func (minikube Minikube) Setup() {
	minikubeStatus, err := minikube.Status()

	if debug {
		if err != nil {
			fmt.Printf("Error occured while checking minikube status. Error: %+v\n", err)
		} else {
			fmt.Printf("minikube status: %q\n", minikubeStatus)
		}
	}

	teardownRequired := false
	startRequired := false

	status, ok := minikubeStatus["minikube"]
	if !ok {
		fmt.Println("\"minikube\" not present in status. May be minikube is not accessible. Aborting...")
		os.Exit(1)
	}
	if status == "" { // This means cluster itself is not there
		fmt.Println("cluster is not up. will start the machine")
		startRequired = true // So, Start the minikube
	} else if status == "Stopped" { // Cluster is there but it is stopped
		fmt.Println("minikube cluster is present but not \"Running\", so will tearing down the machine then start again.")
		teardownRequired = true // We need to teardown it first
		startRequired = true    // Then also we need to start the machine
	} else if status != "Running" { // If cluster is there and machine is not in "Stopped" or "Running" state
		// Then there is a problem
		fmt.Printf("minikube is in unknown state. State: %q. Aborting...", status)
		os.Exit(1)
	} else { // Else minikube is Running so we need not do anything.
		fmt.Println("minikube is already Running.")
	}

	// If we figured out that a teardown is needed then do so
	if teardownRequired {
		err = minikube.Teardown()
		if err != nil {
			fmt.Printf("Error while deleting machine. Error: %+v\n", err)
		} else {
			fmt.Println("minikube deleted.")
		}
	}

	// If we figured out that a start is needed then do so
	if startRequired {
		minikube.StartMinikube()
	}
}
