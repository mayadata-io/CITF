package config

import (
	"fmt"
	"io/ioutil"

	"os"
	"reflect"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Configuration is struct to hold the configurations of CITF
type Configuration struct {
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// Conf will contain configurations for CITF
var Conf Configuration
var defaultConf Configuration

func init() {
	defaultConf = Configuration{
		Environment: "minikube",
	}
}

// LoadConf loads the configuration from the file which path is supplied
func LoadConf(confFilePath string) error {
	if len(confFilePath) == 0 {
		return nil
	}

	yamlBytes, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return fmt.Errorf("error reading file: %q. Error: %+v", confFilePath, err)
	}

	err = yaml.Unmarshal(yamlBytes, Conf)
	if err != nil {
		return fmt.Errorf("error parsing file: %q. Error: %+v", confFilePath, err)
	}
	return nil
}

func getConfValueByStringField(conf Configuration, field string) string {
	r := reflect.ValueOf(conf)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

// fields should be in exact case as the field is present in struct Configuration
func GetDefaultValueByStringField(field string) string {
	return getConfValueByStringField(defaultConf, field)
}

// fields should be in exact case as the field is present in struct Configuration
func GetUserConfValueByStringField(field string) string {
	return getConfValueByStringField(Conf, field)
}

func GetConf(field string) string {
	if value, ok := os.LookupEnv("CITF_CONF_" + strings.ToUpper(field)); ok {
		return value
	}
	if value := GetUserConfValueByStringField(field); len(value) != 0 {
		return value
	}
	return GetDefaultValueByStringField(field)
}

func Environment() string {
	return GetConf("Environment")
}
