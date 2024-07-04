package client

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_config"
)

const (
	powerOn       = "power_on"
	powerOff      = "power_off"
	powerStateOn  = "on"
	powerStateOff = "off"
)

var timeoutSecs = 90

type Machine struct {
	SystemID    string   `json:"system_id,omitempty"`
	IPAddresses []net.IP `json:"ip_addresses,omitempty"`
	PowerState  string   `json:"power_state,omitempty"`
}

func maasApiCall(method string, url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to make new %s request to  %s, error %v", method, url, err)
	}
	parameters := strings.Split(e2e_config.GetConfig().MaasOauthApiToken, ":")
	oauthNonce := generateRandomOauthNonce(11)
	oauthConsumerKey := parameters[0]
	oauthToken := parameters[1]
	oauthSignature := "%26" + parameters[2]
	oauthSignatureMethod := "PLAINTEXT"
	oauthTimestamp := fmt.Sprintf("%v", time.Now().UTC().Unix())
	oauthVersion := "2.0"

	header := fmt.Sprintf("OAuth oauth_consumer_key=%s,oauth_token=%s,oauth_signature_method=%s,oauth_timestamp=%s,oauth_nonce=%s,oauth_version=%s,oauth_signature=%s",
		oauthConsumerKey,
		oauthToken,
		oauthSignatureMethod,
		oauthTimestamp,
		oauthNonce,
		oauthVersion,
		oauthSignature,
	)
	req.Header.Add("Authorization", header)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make %s api call to  %s, error %v", method, url, err)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("mas api reponse code is not 200, actual response code  is %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read api response, error %v", err)
	}
	return body, err
}

func listMaasMachines() ([]Machine, error) {
	url := fmt.Sprintf("http://%s/MAAS/api/2.0/machines/", e2e_config.GetConfig().MaasEndpoint)
	body, err := maasApiCall("GET", url)
	var response []Machine
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error while parsing response, %v", err)
	}
	return response, err
}

func getMaasMachineId(ip string) (string, error) {
	machines, err := listMaasMachines()
	if err != nil {
		return "", fmt.Errorf("failed to get maas machines, error %v", err)
	}
	var machineId string
	for _, machine := range machines {
		ipArray := transformIPArray(machine.IPAddresses)
		for _, machineIp := range ipArray {
			if machineIp == ip {
				machineId = machine.SystemID
			}
		}
	}
	if machineId == "" {
		return "", fmt.Errorf("failed to find system id for ip %s", ip)
	}
	return machineId, err
}

func getMaasMachineStatus(ip string) (string, error) {
	machines, err := listMaasMachines()
	if err != nil {
		return "", fmt.Errorf("failed to get maas machines, error %v", err)
	}

	for _, machine := range machines {
		if strings.Contains(fmt.Sprintf("%s", machine.IPAddresses), ip) {
			return machine.PowerState, err
		}
	}
	return "", err
}

func maasMachinePowerOnOff(operation string, systemId string) error {
	url := fmt.Sprintf("http://%s/MAAS/api/2.0/machines/%s/?op=%s",
		e2e_config.GetConfig().MaasEndpoint,
		systemId,
		operation)
	_, err := maasApiCall("POST", url)
	if err != nil {
		return err
	}
	return err
}

func generateRandomOauthNonce(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func transformIPArray(ipArray []net.IP) []string {
	s := make([]string, 0)
	for _, ip := range ipArray {
		s = append(s, ip.String())
	}
	return s
}
