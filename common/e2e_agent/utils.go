package e2e_agent

import (
	"io"
	"net/http"
)

// GetMyPublicIP fetches the public IP of the host machine
func GetMyPublicIP() (string, error) {
	url := "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}
