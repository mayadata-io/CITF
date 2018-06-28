package citf

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openebs/CITF/common"
	"github.com/openebs/CITF/config"
	"github.com/openebs/CITF/environments"
	"github.com/openebs/CITF/environments/docker"
	"github.com/openebs/CITF/environments/minikube"
	"github.com/openebs/CITF/utils/k8s"
)

// CITF is a struct which will be the driver for all functionalities of this framework
type CITF struct {
	Environment  environments.Environment
	K8S          k8s.K8S
	Docker       docker.Docker
	DebugEnabled bool
}

// NewCITF returns CITF struct. One need this in order to use any functionality of this framework.
func NewCITF(confFilePath string) (CITF, error) {
	var environment environments.Environment
	if err := config.LoadConf(confFilePath); err != nil {
		// Log this here
		// Here, we don't want to return fatal error since we want to continue
		// executing the function with default configuration even if it fails
		glog.Errorf("error loading config file. Error: %+v", err)
	}

	switch config.Environment() {
	case "minikube":
		environment = minikube.NewMinikube()
	default:
		// Exit with Error
		return CITF{}, fmt.Errorf("platform: %q is not suppported by CITF", config.Environment())
	}

	k8sInstance, err := k8s.NewK8S()
	if err != nil {
		return CITF{}, err
	}

	return CITF{
		K8S:          k8sInstance,
		Environment:  environment,
		Docker:       docker.NewDocker(),
		DebugEnabled: common.DebugEnabled,
	}, nil
}
