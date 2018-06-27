package citf

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/openebs/CITF/config"
	"github.com/openebs/CITF/environments"
	"github.com/openebs/CITF/environments/docker"
	"github.com/openebs/CITF/environments/minikube"
	"github.com/openebs/CITF/utils/k8s"
)

// CITF is a struct which will be the driver for all functionalities of this framework
type CITF struct {
	Environment environments.Environment
	K8S         k8s.K8S
	Docker      docker.Docker
}

// NewCITF returns CITF struct. One need this in order to use any functionality of this framework.
func NewCITF(confFilePath string) CITF {
	var environment environments.Environment
	if err := config.LoadConf(confFilePath); err != nil {
		// Log this here
		fmt.Fprintf(GinkgoWriter, "error loading config file. Error: %+v", err)
	}

	switch config.Environment() {
	case "minikube":
		environment = minikube.NewMinikube()
	default:
		// Exit with Error using Gomega
		Expect("").ToNot(BeEmpty())
	}

	return CITF{
		K8S:         k8s.NewK8S(),
		Environment: environment,
		Docker:      docker.NewDocker(),
	}
}
