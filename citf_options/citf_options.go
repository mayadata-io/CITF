package citfoptions

import "github.com/openebs/CITF/utils/log"

// CreateOptions specifies which fields of CITF should be included when created or reloaded
type CreateOptions struct {
	ConfigPath         string
	EnvironmentInclude bool
	K8SInclude         bool
	DockerInclude      bool
	LoggerInclude      bool
	T                  log.LoggerT
}

// CreateOptionsIncludeAll returns CreateOptions where all fields are set to `true` and ConfigPath is set to configPath
func CreateOptionsIncludeAll(configPath string, t log.LoggerT) *CreateOptions {
	var citfCreateOptions CreateOptions
	citfCreateOptions.ConfigPath = configPath
	citfCreateOptions.EnvironmentInclude = true
	citfCreateOptions.K8SInclude = true
	citfCreateOptions.DockerInclude = true
	citfCreateOptions.LoggerInclude = true
	citfCreateOptions.T = t
	return &citfCreateOptions
}

// CreateOptionsIncludeAllButEnvironment returns CreateOptions where all fields except `Environment` are set to `true` and ConfigPath is set to configPath
func CreateOptionsIncludeAllButEnvironment(configPath string, t log.LoggerT) *CreateOptions {
	citfCreateOptions := CreateOptionsIncludeAll(configPath, t)

	citfCreateOptions.EnvironmentInclude = false
	return citfCreateOptions
}

// CreateOptionsIncludeAllButK8s returns CreateOptions where all fields except `K8S` are set to `true` and ConfigPath is set to configPath
func CreateOptionsIncludeAllButK8s(configPath string, t log.LoggerT) *CreateOptions {
	citfCreateOptions := CreateOptionsIncludeAll(configPath, t)

	citfCreateOptions.K8SInclude = false
	return citfCreateOptions
}

// CreateOptionsIncludeAllButDocker returns CreateOptions where all fields except `Docker` are set to `true` and ConfigPath is set to configPath
func CreateOptionsIncludeAllButDocker(configPath string, t log.LoggerT) *CreateOptions {
	citfCreateOptions := CreateOptionsIncludeAll(configPath, t)

	citfCreateOptions.DockerInclude = false
	return citfCreateOptions
}

// CreateOptionsIncludeAllButLogger returns CreateOptions where all fields except `Logger` are set to `true` and ConfigPath is set to configPath
func CreateOptionsIncludeAllButLogger(configPath string, t log.LoggerT) *CreateOptions {
	citfCreateOptions := CreateOptionsIncludeAll(configPath, t)

	citfCreateOptions.LoggerInclude = false
	return citfCreateOptions
}
