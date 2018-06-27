package minikube

// Teardown deletes minikube
func (minikube Minikube) Teardown() error {
	// Caller of this function should have proper rights to delete minikube
	return runCommand("minikube delete")
}
