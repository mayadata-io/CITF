package docker

// CheckStatus should return the status of docker environment.
// It returns nil,nil for now. It has been written just to implement interface Environment
func (docker Docker) CheckStatus() (map[string]string, error) {
	return nil, nil
}
