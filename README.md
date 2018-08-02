# CITF

[![Build Status](https://travis-ci.org/openebs/CITF.svg?branch=master)](https://travis-ci.org/openebs/CITF)
[![Go Report](https://goreportcard.com/badge/github.com/openebs/CITF)](https://goreportcard.com/report/github.com/openebs/CITF)
[![codecov](https://codecov.io/gh/openebs/CITF/branch/master/graph/badge.svg)](https://codecov.io/gh/openebs/CITF)
[![GoDoc](https://godoc.org/github.com/openebs/CITF?status.svg)](https://godoc.org/github.com/openebs/CITF)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/openebs/CITF/blob/master/LICENSE)

**Common Integration Test Framework** is a framework that will be used organization wide for Integration Test of all OpenEBS projects.

This repository is intended to only expose generic function which will help developers in writing Integration Tests. Though it won't produce any deliverable alone.

## Directory Structure in the Project
```
OpenEBS/project
   ├── integration_test
   │   ├── project_specific_package_for_integration_test
   │   │   ├── ...
   │   │   └── files.go
   │   ├── scenario1_test.go
   │   ├── scenario2_test.go
   │   ├── ...
   │   └── scenarioN_test.go
   ├── project_specific_packages
   └── vendor
       ├── package_related_vendors
       ├── ...
       └── github.com/OpenEBS/CITF
```

> Note: Developer should keep `integration_test` completely decoupled from the rest of the project packages.

## Instantiation

Developer has to instantiate CITF using `citf.NewCITF` function, which will initialize it with all the configurations specified by `citfoptions.CreateOptions` passed to it. 

> You should not pass `K8sInclude` in `citfoptions.CreateOptions` if your environment is not yet set. otherwise it will through error.

> If you want all options except `K8sInclude` in `CreateOptions` to set to `true`; you may use `citfoptions.CreateOptionsIncludeAllButK8s` function.

> If you want all options in `CreateOptions` to set to `true`  you may use `citfoptions.CreateOptionsIncludeAll` function.

> **Note:** `citfoptions.CreateOptions.T` is compatible with golang's standard `*testing.T` as well as ginkgo's `ginkgo.GinkgoTInterface`. We can get `ginkgo.GinkgoTInterface` by `ginkgo.GinkgoT()`.

CITF struct has four fields:- 
- Environment - To Setup or TearDown the platform such as minikube, GKE, AWS etc.
- K8S - K8S will have Kubernetes ClientSet & Config.
- Docker - Docker will be used for docker related operations.
- DebugEnabled - for verbose log.

> Currently CITF environment supports minikube only.

Developer can pass environment according to their requirements.

By default it will take Minikube as environment.

## Configuration

To configure the environment of CITF, there are three ways:-
 - [Environment Variable](#environment-variable)
 - [Config File](#config-file)
 - [Default Config](#default-config)

### Environment Variable

At the time of instantiation, developer can set CITF environment using environment variable `CITF_CONF_ENVIRONMENT`.

For example:- `export CITF_CONF_ENVIRONMENT = minikube`

### Config File
If environment variable is not set then developer can pass environment using config file. The file should be in `yaml` format. 

For example:- config.yaml

```
Environment: minikube
```

### Default Config

If environment variable and config file are not present, then CITF will take default environment which is minikube.

### Platform Operations

`citf.Environment` will handle operations related to the platforms. 

In order to setup the k8s cluster, developer needs to call the `Setup()` method which will bring it up.

Developer can also check the status of the platform using `Status()` method.

Once integration test is completed, developer can delete the setup using `TearDown()` method.

> Examples for this can be found [here](https://github.com/openebs/CITF/tree/master/example).