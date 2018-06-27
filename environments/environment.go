package environments

// Environment is the interface which integrate all the functionalities
// that environments like minikube, docker etc should have.
type Environment interface {
	Setup() // If Setup of environment fails it is better to exit from there
	Status() (map[string]string, error)
	Teardown() error
}
