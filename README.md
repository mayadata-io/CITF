# CITF
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
