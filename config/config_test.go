/*
Copyright 2018 The OpenEBS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"os"
	"testing"

	"github.com/openebs/CITF/utils/log"
)

var logger log.Logger

// CreateFile creates yaml file for test purpose
func CreateFile() {
	fileData1 := `
environment: minikube
`
	f, err := os.Create("./test-config.yaml")

	logger.LogError(err, "unable to create config file")

	f.WriteString(fileData1)

	// Create yaml file with bad indentation
	fileData2 := `
	environment: minikube
	`
	f, err = os.Create("./test-bad-config.yaml")
	logger.LogError(err, "unable to create bad config file")

	f.WriteString(fileData2)
}

// DeleteFile deletes yaml file
func DeleteFile() {
	err := os.Remove("./test-config.yaml")
	logger.LogError(err, "unable to delete config file")

	err = os.Remove("./test-bad-config.yaml")
	logger.LogError(err, "unable to delete bad config file")
}

func TestLoadConf(t *testing.T) {
	CreateFile()
	type args struct {
		confFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "LoadConfFileNotPresent",
			args: args{
				confFilePath: "./file.yaml",
			},
			wantErr: true,
		},
		{
			name: "LoadConfSuccess",
			args: args{
				confFilePath: "./test-config.yaml",
			},
			wantErr: false,
		},
		{
			name: "LoadConfEmptyFileName",
			args: args{
				confFilePath: "",
			},
			wantErr: false,
		},
		{
			name: "LoadConfBadYamlFile",
			args: args{
				confFilePath: "./test-bad-config.yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConf(tt.args.confFilePath); (err != nil) != tt.wantErr {
				t.Errorf("LoadConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	DeleteFile()
}

func Test_getConfValueByStringField(t *testing.T) {
	type args struct {
		conf  Configuration
		field string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "KeyPresentInYaml",
			args: args{
				conf: Configuration{
					Environment: "minikube",
				},
				field: "Environment",
			},
			want: "minikube",
		},
		{
			name: "KeyNotPresentInYaml",
			args: args{
				conf: Configuration{
					Environment: "minikube",
				},
				field: "environment",
			},
			want: "<invalid reflect.Value>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfValueByStringField(tt.args.conf, tt.args.field); got != tt.want {
				t.Errorf("getConfValueByStringField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultValueByStringField(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "FieldPresentInYaml",
			args: args{
				field: "Environment",
			},
			want: "minikube",
		},
		{
			name: "FieldNotPresentInYaml",
			args: args{
				field: "environment",
			},
			want: "<invalid reflect.Value>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultValueByStringField(tt.args.field); got != tt.want {
				t.Errorf("GetDefaultValueByStringField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserConfValueByStringField(t *testing.T) {
	type args struct {
		field string
	}
	// Set Value of Conf
	Conf = Configuration{
		Environment: "minikube",
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "FieldPresentInConf",
			args: args{
				field: "Environment",
			},
			want: "minikube",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserConfValueByStringField(tt.args.field); got != tt.want {
				t.Errorf("GetUserConfValueByStringField() = %v, want %v", got, tt.want)
			}
		})
	}
	Conf = Configuration{} // Reset value of Conf to default (being a global guy)
}

func TestGetConf(t *testing.T) {
	os.Setenv("CITF_CONF_ENVIRONMENT", "minikube")

	Conf = Configuration{
		Environment: "Dear minikube",
	}

	type args struct {
		field string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		cleanup func()
	}{
		{
			name: "GetConfWithEnvSuccess",
			args: args{
				field: "Environment",
			},
			want: "minikube",
			cleanup: func() {
				os.Unsetenv("CITF_CONF_ENVIRONMENT")
			},
		},
		{
			name: "GetConfWithConfValueSuccess",
			args: args{
				field: "Environment",
			},
			want: "Dear minikube",
			cleanup: func() {
				Conf = Configuration{
					Environment: "",
				}
			},
		},
		{
			name: "GetConfWithDefaultConfValueSuccess",
			args: args{
				field: "Environment",
			},
			want:    "minikube",
			cleanup: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConf(tt.args.field); got != tt.want {
				t.Errorf("GetConf() = %v, want %v", got, tt.want)
			}
		})
		tt.cleanup()
	}
}

func TestEnvironment(t *testing.T) {
	// to store current content of `CITF_CONF_DEBUG`
	var environContent string
	// to store whether `CITF_CONF_DEBUG` was even set
	var environSet bool

	var confBak Configuration

	tests := []struct {
		name       string
		beforeFunc func()
		want       string
		afterFunc  func()
	}{
		{
			name: "`CITF_CONF_ENVIRONMENT` is set to `my-environment`",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_ENVIRONMENT")
				os.Setenv("CITF_CONF_ENVIRONMENT", "my-environment")
			},
			want: "my-environment",
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_ENVIRONMENT", environContent)
				} else {
					os.Unsetenv("CITF_CONF_ENVIRONMENT")
				}
			},
		},
		{
			name: "`CITF_CONF_ENVIRONMENT` is not set and Environment is set to `my-environment-in-conf` in `Conf`",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_ENVIRONMENT")
				if environSet {
					os.Unsetenv("CITF_CONF_ENVIRONMENT")
				}
				confBak = Conf
				Conf = Configuration{
					Environment: "my-environment-in-conf",
				}
			},
			want: "my-environment-in-conf",
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_ENVIRONMENT", environContent)
				}
				Conf = confBak
			},
		},
		{
			name: "`CITF_CONF_ENVIRONMENT` is not set and `Conf` is empty",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_ENVIRONMENT")
				if environSet {
					os.Unsetenv("CITF_CONF_ENVIRONMENT")
				}
				confBak = Conf
				Conf = Configuration{}
			},
			want: defaultConf.Environment,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_ENVIRONMENT", environContent)
				}
				Conf = confBak
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			if got := Environment(); got != tt.want {
				t.Errorf("Environment() = %v, want %v", got, tt.want)
			}
			tt.afterFunc()
		})
	}
}

func TestDebug(t *testing.T) {
	// to store current content of `CITF_CONF_DEBUG`
	var environContent string
	// to store whether `CITF_CONF_DEBUG` was even set
	var environSet bool

	var confBak Configuration

	tests := []struct {
		name       string
		beforeFunc func()
		want       bool
		afterFunc  func()
	}{
		{
			name: "`CITF_CONF_DEBUG` is set to true",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_DEBUG")
				os.Setenv("CITF_CONF_DEBUG", debugEnabledValStr)
			},
			want: debugEnabledVal,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_DEBUG", environContent)
				} else {
					os.Unsetenv("CITF_CONF_DEBUG")
				}
			},
		},
		{
			name: "`CITF_CONF_DEBUG` is set to false",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_DEBUG")
				os.Setenv("CITF_CONF_DEBUG", debugDisabledValStr)
			},
			want: debugDisabledVal,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_DEBUG", environContent)
				} else {
					os.Unsetenv("CITF_CONF_DEBUG")
				}
			},
		},
		{
			name: "`CITF_CONF_DEBUG` is not set and Debug is enabled in `Conf`",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_DEBUG")
				if environSet {
					os.Unsetenv("CITF_CONF_DEBUG")
				}
				confBak = Conf
				Conf = Configuration{
					Debug: debugEnabledVal,
				}
			},
			want: debugEnabledVal,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_DEBUG", environContent)
				}
				Conf = confBak
			},
		},
		{
			name: "`CITF_CONF_DEBUG` is not set and Debug is disabled in `Conf`",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_DEBUG")
				if environSet {
					os.Unsetenv("CITF_CONF_DEBUG")
				}
				confBak = Conf
				Conf = Configuration{
					Debug: debugDisabledVal,
				}
			},
			want: debugDisabledVal,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_DEBUG", environContent)
				}
				Conf = confBak
			},
		},
		{
			name: "`CITF_CONF_DEBUG` is not set and `Conf` is empty",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_DEBUG")
				if environSet {
					os.Unsetenv("CITF_CONF_DEBUG")
				}
				confBak = Conf
				Conf = Configuration{}
			},
			want: defaultConf.Debug,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_DEBUG", environContent)
				}
				Conf = confBak
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			if got := Debug(); got != tt.want {
				t.Errorf("Debug() = %v, want %v", got, tt.want)
			}
			tt.afterFunc()
		})
	}
}

func TestKubeMasterURL(t *testing.T) {
	// to store current content of `CITF_CONF_KUBEMASTERURL`
	var environContent string
	// to store whether `CITF_CONF_KUBEMASTERURL` was even set
	var environSet bool

	var confBak Configuration

	tests := []struct {
		name       string
		beforeFunc func()
		want       string
		afterFunc  func()
	}{
		{
			name: "`CITF_CONF_KUBEMASTERURL` is set to 'my-kube-master-url'",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_KUBEMASTERURL")
				os.Setenv("CITF_CONF_KUBEMASTERURL", "my-kube-master-url")
			},
			want: "my-kube-master-url",
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_KUBEMASTERURL", environContent)
				} else {
					os.Unsetenv("CITF_CONF_KUBEMASTERURL")
				}
			},
		},
		{
			name: "`CITF_CONF_KUBEMASTERURL` is not set and KubeMasterURL is set to `kube-master-url-in-conf` in `Conf`",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_KUBEMASTERURL")
				if environSet {
					os.Unsetenv("CITF_CONF_KUBEMASTERURL")
				}
				confBak = Conf
				Conf = Configuration{
					KubeMasterURL: "kube-master-url-in-conf",
				}
			},
			want: "kube-master-url-in-conf",
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_KUBEMASTERURL", environContent)
				}
				Conf = confBak
			},
		},
		{
			name: "`CITF_CONF_KUBEMASTERURL` is not set and `Conf` is empty",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_KUBEMASTERURL")
				if environSet {
					os.Unsetenv("CITF_CONF_KUBEMASTERURL")
				}
				confBak = Conf
				Conf = Configuration{}
			},
			want: defaultConf.KubeMasterURL,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_KUBEMASTERURL", environContent)
				}
				Conf = confBak
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			if got := KubeMasterURL(); got != tt.want {
				t.Errorf("KubeMasterURL() = %v, want %v", got, tt.want)
			}
			tt.afterFunc()
		})
	}
}

func TestKubeConfigPath(t *testing.T) {
	// to store current content of `CITF_CONF_KUBECONFIGPATH`
	var environContent string
	// to store whether `CITF_CONF_KUBECONFIGPATH` was even set
	var environSet bool

	var confBak Configuration

	tests := []struct {
		name       string
		beforeFunc func()
		want       string
		afterFunc  func()
	}{
		{
			name: "`CITF_CONF_KUBECONFIGPATH` is set to 'my-kube-config-path'",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_KUBECONFIGPATH")
				os.Setenv("CITF_CONF_KUBECONFIGPATH", "my-kube-config-path")
			},
			want: "my-kube-config-path",
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_KUBECONFIGPATH", environContent)
				} else {
					os.Unsetenv("CITF_CONF_KUBECONFIGPATH")
				}
			},
		},
		{
			name: "`CITF_CONF_KUBECONFIGPATH` is not set and KubeConfigPath is set to `kube-config-path-in-conf` in `Conf`",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_KUBECONFIGPATH")
				if environSet {
					os.Unsetenv("CITF_CONF_KUBECONFIGPATH")
				}
				confBak = Conf
				Conf = Configuration{
					KubeConfigPath: "kube-config-path-in-conf",
				}
			},
			want: "kube-config-path-in-conf",
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_KUBECONFIGPATH", environContent)
				}
				Conf = confBak
			},
		},
		{
			name: "`CITF_CONF_KUBECONFIGPATH` is not set and `Conf` is empty",
			beforeFunc: func() {
				environContent, environSet = os.LookupEnv("CITF_CONF_KUBECONFIGPATH")
				if environSet {
					os.Unsetenv("CITF_CONF_KUBECONFIGPATH")
				}
				confBak = Conf
				Conf = Configuration{}
			},
			want: defaultConf.KubeConfigPath,
			afterFunc: func() {
				if environSet {
					os.Setenv("CITF_CONF_KUBECONFIGPATH", environContent)
				}
				Conf = confBak
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			if got := KubeConfigPath(); got != tt.want {
				t.Errorf("KubeConfigPath() = %v, want %v", got, tt.want)
			}
			tt.afterFunc()
		})
	}
}
